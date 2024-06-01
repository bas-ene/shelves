package googleapi

import (
	"context"
	"fmt"
	"log"

	"shelves/libr"

	"google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

type GoogleAPI struct {
	client *books.Service
	//TODO: avoid it growing without control
	cache map[string]*libr.Book
}

func buildBook(item *books.Volume) *libr.Book {

	b := libr.Book{
		//TODO: handle the isbn situation better
		Isbn:       libr.ISBN(item.VolumeInfo.IndustryIdentifiers[0].Identifier),
		Title:      item.VolumeInfo.Title,
		GoogleId:   item.Id,
		Authors:    item.VolumeInfo.Authors,
		Pages:      int(item.VolumeInfo.PageCount),
		Categories: item.VolumeInfo.Categories,
		Language:   item.VolumeInfo.Language,
		Price:      0.0,
		Currency:   "EUR",
		PagesRead:  0,
	}
	if item.SaleInfo.Saleability != "NOT_FOR_SALE" {
		b.Price = item.SaleInfo.ListPrice.Amount
		b.Currency = item.SaleInfo.ListPrice.CurrencyCode
	}
	if item.VolumeInfo.ImageLinks != nil {
		b.ImgHref = item.VolumeInfo.ImageLinks.Thumbnail
	} else {
		b.ImgHref = "https://books.google.it/googlebooks/images/no_cover_thumb.gif"
	}

	return &b
}

func (api *GoogleAPI) GetBook(googleId string) (*libr.Book, error) {
	if _, ok := api.cache[googleId]; ok {
		log.Printf("Book in cache")
		return api.cache[googleId], nil
	}

	call := api.client.Volumes.Get(googleId)
	log.Printf("Searching for GoggleId: %s", googleId)
	item, err := call.Do()
	if err != nil {
		log.Printf("Errore nell'esecuzione della chiamata: %v", err)
		return nil, err
	}
	b := buildBook(item)
	api.cache[googleId] = b
	return b, nil
}

func (api *GoogleAPI) SearchByIsbn(isbn libr.ISBN) ([]*libr.Book, error) {
	call := api.client.Volumes.List(fmt.Sprintf("isbn:%s", isbn))
	resp, err := call.Do()

	if err != nil {
		log.Printf("Errore nell'esecuzione della chiamata: %v", err)
		return nil, err
	}
	var bookList []*libr.Book
	for _, item := range resp.Items {
		b, _ := api.GetBook(item.Id)
		bookList = append(bookList, b)
	}
	return bookList, nil
}

func NewGoogleBooksApi(apiKey string) *GoogleAPI {
	client, err := books.NewService(context.TODO(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Errore nella creazione del servizio: %v", err)
		return nil
	}
	return &GoogleAPI{client: client, cache: make(map[string]*libr.Book)}
}
