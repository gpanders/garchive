package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Link struct {
	Title  string
	Path   string
	Domain string
}

var links string
var protocolRe = regexp.MustCompile("^(?:http|https|ftp)://")

func serveIndex(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	search := r.URL.Query().Get("search")
	p, err := getLinks(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, p)
}

func getLinks(query string) ([]Link, error) {
	f, err := os.Open(links)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 2

	var links []Link
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(query)

	for _, row := range rows {
		title := row[0]
		path := protocolRe.ReplaceAllLiteralString(row[1], "")
		if !re.MatchString(title) && !re.MatchString(path) {
			continue
		}

		domain := strings.SplitN(path, "/", 2)[0]
		links = append(links, Link{Title: title, Path: path, Domain: domain})
	}

	return links, nil
}

func main() {
	addr := flag.String("addr", "0.0.0.0", "bind address")
	port := flag.Int("port", 8080, "bind port")
	archive := flag.String("archive", "data", "absolute path to archive directory")
	flag.StringVar(&links, "links", "links.csv", "path to links CSV file")

	flag.Parse()

	http.Handle("/ar/", http.StripPrefix("/ar/", http.FileServer(http.Dir(*archive))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveIndex)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), nil))
}
