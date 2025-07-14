# 🚀 Task App Backend

This is a **Task App backend**, a web service built with **Go**, **Gin**, **GORM**, and **MySQL/PostgreSQL**.  
It analyzes web pages, extracting details like:

- HTML version
- Title
- Number of `<h1>` and `<h2>` tags
- Counts of internal, external, and broken links
- Detects if a login form exists

It then stores this in a relational database and provides a REST API to view or manage it.

---

## ⚙️ Tech Stack

- **Go** (Gin framework)
- **GORM** (ORM for Go)
- **MySQL**  database

---

## 🚀 Quickstart

### ✅ Prerequisites

- Go >= 1.18 installed
- MySQL running (e.g. on `localhost:3306` )
- A database created, e.g `taskdb`

---

### 🔧 Configuration

Edit your DB connection string in `database/db.go` .

For MySQL:
dsn := "username:password@tcp(127.0.0.1:3306)/taskdb?charset=utf8mb4&parseTime=True&loc=Local"

### 🏗️ Build and run
Install dependencies & run the server:
- go mod tidy
- go run main.go  or you can run ( go run . ) to start with some dummy data . 
