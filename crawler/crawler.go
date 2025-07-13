package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"github.com/dalibouali/task-app-backend/models"
)


func AnalyzeUrl(db *gorm.DB, urlEntry *models.Url) {
	// Always reload the latest from DB
	var current models.Url
	if err := db.First(&current, urlEntry.ID).Error; err != nil {
		fmt.Println("Error fetching latest URL record:", err)
		return
	}

	// Skip if not queued anymore (stopped or processed already)
	if current.Status != "queued" {
		fmt.Println("Skipping", current.URL, "status is now", current.Status)
		return
	}

	fmt.Println("Starting analysis for:", current.URL)
	current.Status = "running"
	db.Save(&current)

	fmt.Println("Fetching URL:", current.URL)
	resp, err := http.Get(current.URL)
	if err != nil {
		current.Status = "error"
		db.Save(&current)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		current.Status = "error"
		db.Save(&current)
		return
	}

	parsedBase, err := url.Parse(current.URL)
	if err != nil {
		current.Status = "error"
		db.Save(&current)
		return
	}

	current.HtmlVersion = detectHTMLVersion(doc)
	current.Title = extractTitle(doc)
	current.H1Count, current.H2Count = countHeadings(doc)
	internalLinks, externalLinks, brokenLinks := countLinks(parsedBase, doc)
	current.InternalLinks = internalLinks
	current.ExternalLinks = externalLinks
	current.BrokenLinksCount = len(brokenLinks)
	current.BrokenLinksList = brokenLinks
	current.HasLoginForm = detectLoginForm(doc)
	current.Status = "done"

	db.Save(&current)
	// Clear existing broken links
	db.Unscoped().Where("url_id = ?", current.ID).Delete(&models.BrokenLink{})

	// Save new broken links
	for _, b := range brokenLinks {
		b.UrlID = current.ID
		db.Create(&b)
	}

	fmt.Println("Finished analysis for:", current.URL)
}


func detectHTMLVersion(n *html.Node) string {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type == html.DoctypeNode {
            if strings.EqualFold(c.Data, "html") {
                return "HTML5"
            }
            return "Older HTML"
        }
    }
    return "Unknown"
}

func extractTitle(n *html.Node) string {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return title
}

func countHeadings(n *html.Node) (int, int) {
	var h1, h2 int
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				h1++
			case "h2":
				h2++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return h1, h2
}

func countLinks(base *url.URL, n *html.Node) (int, int, []models.BrokenLink) {
	var internal, external int
	var brokenLinks []models.BrokenLink
	client := &http.Client{Timeout: 3 * time.Second}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := strings.TrimSpace(attr.Val)
					if link == "" || strings.HasPrefix(link, "#") {
						continue
					}
					parsedLink, err := url.Parse(link)
					if err != nil {
						continue
					}
					if parsedLink.Host == "" || parsedLink.Host == base.Host {
						internal++
					} else {
						external++
					}

					// Check link status code
					fullUrl := base.ResolveReference(parsedLink).String()
					statusCode := getStatusCode(client, fullUrl)
					if statusCode >= 400 || statusCode == 0 {
						brokenLinks = append(brokenLinks, models.BrokenLink{
							URL: fullUrl,
							StatusCode: statusCode,
						})
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return internal, external, brokenLinks
}


func getStatusCode(client *http.Client, link string) int {
	resp, err := client.Head(link)
	if err != nil {
		return 0 // treat as no error
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func detectLoginForm(n *html.Node) bool {
	var found bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && strings.ToLower(attr.Val) == "password" {
					found = true
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return found
}
