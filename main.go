package main

import (
	"fmt"
	"net/http"
	"os"

	"shelves/libr"
	"shelves/templs"
	"strconv"
)

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	// fmt.Printf("i: %d, err: %s", i, err)
	return err == nil
}

func checkIsbn(text string, lastChar rune) bool {
	if lastChar == 0 {
		return true
	}
	if lastChar == 'x' {
		return len(text) == 10
	}
	if len(text) > 13 {
		return false
	}
	return isNumeric(text)
}

func main() {
	myLib, err := libr.LoadLibrary("./library.json")
	if err != nil {
		fmt.Println("Error while loading library: %w", err)
		os.Exit(1)
	}
	// bookApi := NewGoogleBooksApi(GOOGLEAPIKEY)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templs.BookComp(myLib[0]).Render(r.Context(), w)
	})

	fmt.Println("Listening on :9090")
	http.ListenAndServe(":9090", nil)
}

// http.HandleFunc("GET /api/search/{isbn}", func(w http.ResponseWriter, r *http.Request) {
// 	isbn := r.PathValue("isbn")
// 	res, _ := bookApi.GetBookOptions(ISBN(isbn))
// 	j, _ := json.MarshalIndent(res, "", "   ")
// 	w.Write(j)
// })
//
// //TODO: dovra diventare un POST
// http.HandleFunc("GET /api/add/{googleId}", func(w http.ResponseWriter, r *http.Request) {
// 	defer myLib.Save("./library.json")
// 	gId := r.PathValue("googleId")
// 	res, _ := bookApi.GetBook(gId)
// 	myLib = append(myLib, *res)
// 	j, _ := json.MarshalIndent(res, "", "   ")
// 	w.Write(j)
// })
//
// //TODO: dovra diventare un DELETE
// http.HandleFunc("GET /api/remove/{googleId}", func(w http.ResponseWriter, r *http.Request) {
// 	defer myLib.Save("./library.json")
// 	gId := r.PathValue("googleId")
// 	deletedBook := myLib.RemoveBook(gId)
// 	j, _ := json.MarshalIndent(deletedBook, "", "   ")
// 	w.Write(j)
// })
//
// //TODO: DIVENTARE PUT
// http.HandleFunc("GET /api/read/{googleId}", func(w http.ResponseWriter, r *http.Request) {
// 	defer myLib.Save("./library.json")
// 	gId := r.PathValue("googleId")
// 	nPages, err := strconv.Atoi(r.URL.Query().Get("nPages"))
// 	if err != nil {
// 		str := "Parametro nPages deve essere un numero intero"
// 		j, _ := json.Marshal(str)
// 		w.Write(j)
// 		return
// 	}
// 	book := myLib.ReadBook(gId, nPages)
// 	j, _ := json.MarshalIndent(book, "", "   ")
//
// 	w.Write(j)
// })
//
// log.Fatal(http.ListenAndServe(":9090", nil))
