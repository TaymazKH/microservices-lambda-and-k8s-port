package main

import (
    "html/template"
    "log"
    "net/http"
)

var templates = template.Must(template.ParseFiles("base.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received at: %v", r.URL)

    data := struct {
        Title   string
        Message string
    }{
        Title:   "Welcome",
        Message: "This is the home page.",
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := templates.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        log.Println("Template execution error:", err)
    }
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received at: %v", r.URL)

    data := struct {
        Title   string
        Message string
    }{
        Title:   "About Us",
        Message: "This page is the about us.",
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := templates.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        log.Println("Template execution error:", err)
    }
}
