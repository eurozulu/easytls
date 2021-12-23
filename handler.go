package main

import (
	"easytls/templates"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

const contentFolder = "content"
const templateFolder = "templates"

type Handler struct {
	Verbose bool
	temps   map[string]*template.Template
}

var templatePages = []string{"basepage.html"}

var pages = findPages()

func (h Handler) handleRequest(w http.ResponseWriter, req *http.Request) {
	if h.Verbose {
		log.Printf("request to %s: %v", req.URL, req.Header)
	}
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	switch req.Method {
	case http.MethodGet:
		h.handleGetRequest(w, req)

	case http.MethodPost:
		h.handlePostRequest(w, req)

	default:
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (h Handler) handleGetRequest(w http.ResponseWriter, req *http.Request) {
	p, err := h.loadPage(strings.TrimLeft(req.URL.Path, "/"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t, ok := h.temps[p.TemplateName]
	if !ok {
		log.Printf("%s is an invalid page, requests unknown template '%s'", req.URL.Path, p.TemplateName)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, p); err != nil {
		log.Printf("Problem rendering %s  %w", req.URL.Path, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h Handler) handlePostRequest(w http.ResponseWriter, req *http.Request) {
}

func (h Handler) loadPage(name string) (*templates.Page, error) {
	pn, ok := pages[name]
	if !ok {
		return nil, fmt.Errorf("unknown page")
	}

	fn := path.Join(contentFolder, pn)
	by, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	p := &templates.Page{}
	if err = json.Unmarshal(by, p); err != nil {
		return nil, err
	}
	return p, nil
}

func NewHandler() *Handler {
	temps := map[string]*template.Template{}
	for _, tn := range templatePages {
		tp := path.Join(templateFolder, tn)
		t, err := template.ParseFiles(tp)
		if err != nil {
			log.Fatalf("template file %v is invalid,  %w", tn, err)
		}
		temps[tn] = t
	}
	return &Handler{temps: temps}
}

func findPages() map[string]string {
	m := map[string]string{}
	fis, err := ioutil.ReadDir(contentFolder)
	if err != nil {
		log.Fatalf("Invalid content folder  %w", err)
	}
	for _, fi := range fis {
		if fi.IsDir() || !fi.Mode().IsRegular() {
			continue
		}
		ext := path.Ext(fi.Name())
		if !strings.EqualFold(ext, ".json") {
			continue
		}
		n := fi.Name()[:len(fi.Name())-len(ext)]
		m[n] = fi.Name()
	}
	return m
}
