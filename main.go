package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"shelves/libr"
	"shelves/templs"
	"strconv"
)

var SAVEFILE string = "./library.json"

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	// fmt.Printf("i: %d, err: %s", i, err)
	return err == nil
}

func checkIsbn(isbn string) bool {
	switch len(isbn) {
	case 10:
		if isNumeric(isbn) || (isNumeric(isbn[:9]) && isbn[9] == 'x') {
			return true
		}
		return false
	case 13:
		return isNumeric(isbn)
	default:
		return false
	}
}

func main() {
	myLib, err := libr.LoadLibrary(SAVEFILE)
	if err != nil {
		log.Fatal("Could not retrieve library")
	}
	bookApi := NewGoogleBooksApi(GOOGLEAPIKEY)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
		return
	})

	http.HandleFunc("GET /library/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templs.LibraryComp(myLib).Render(r.Context(), w)
		return
	})

	http.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Search: ")
		isbn := r.URL.Query().Get("isbn")
		log.Printf("ISBN: %s\n", isbn)
		isbn = strings.TrimSpace(strings.ReplaceAll(isbn, "-", ""))
		if checkIsbn(isbn) {
			log.Println("Valid isbn, searching")
			searched, _ := bookApi.SearchByIsbn(libr.ISBN(isbn))
			log.Printf("Found: %v", searched)
			templs.SearchedBooksComp(searched).Render(r.Context(), w)
		} else {
			log.Println("Invalid isbn")
			w.Write([]byte{'f'})
		}
		return
	})

	http.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		defer myLib.Save(SAVEFILE)
		r.ParseForm()
		gId := r.PostFormValue("gId")
		b, _ := bookApi.GetBook(gId)
		myLib.AddBook(b)
		templs.BookInLib(*b).Render(r.Context(), w)
		return
	})

	http.HandleFunc("DELETE /remove/{gId}", func(w http.ResponseWriter, r *http.Request) {
		defer myLib.Save(SAVEFILE)
		gId := r.PathValue("gId")
		if myLib.RemoveBook(gId) != nil {
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
		//TODO:ROAR TORNA NUOVA LIBRERIA? o riesco a togliere solo il libro?
		return
	})

	fmt.Println("Listening on :9090")
	http.ListenAndServe(":9090", nil)
}
