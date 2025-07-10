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
	fmt.Println("Starting analysis for:", urlEntry.URL)
	urlEntry.Status = "running"
	db.Save(urlEntry)
	fmt.Println("Fetching URL:", urlEntry.URL)	
	

	resp, err := http.Get(urlEntry.URL)
	if err != nil {
		urlEntry.Status = "error"
		db.Save(urlEntry)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		urlEntry.Status = "error"
		db.Save(urlEntry)
		return
	}

	parsedBase, err := url.Parse(urlEntry.URL)
	if err != nil {
		urlEntry.Status = "error"
		db.Save(urlEntry)
		return
	}

	urlEntry.HtmlVersion = detectHTMLVersion(doc)
	urlEntry.Title = extractTitle(doc)
	urlEntry.H1Count, urlEntry.H2Count = countHeadings(doc)
	urlEntry.InternalLinks, urlEntry.ExternalLinks, urlEntry.BrokenLinks = countLinks(parsedBase, doc)
	urlEntry.HasLoginForm = detectLoginForm(doc)
	urlEntry.Status = "done"

	db.Save(urlEntry)
	fmt.Println("Finished analysis for:", urlEntry.URL)
}

func detectHTMLVersion(n *html.Node) string {
	// Look for <!DOCTYPE html>
	for c := n; c != nil; c = c.NextSibling {
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

func countLinks(base *url.URL, n *html.Node) (int, int, int) {
	var internal, external, broken int
	var f func(*html.Node)
	client := &http.Client{Timeout: 3 * time.Second}

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
					if statusCode >= 400 {
						broken++
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return internal, external, broken
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
