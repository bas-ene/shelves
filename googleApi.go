package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/books/v1"
	"google.golang.org/api/option"
	"shelves/libr"
)

type GoogleAPI struct {
	client *books.Service
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
	return &b
}

func (api *GoogleAPI) GetBook(googleId string) (*libr.Book, error) {
	call := api.client.Volumes.Get(googleId)
	item, _ := call.Do()
	b := buildBook(item)
	return b, nil
}

func (api *GoogleAPI) GetBookOptions(isbn libr.ISBN) ([]*libr.Book, error) {
	call := api.client.Volumes.List(fmt.Sprintf("isbn:%s", isbn))
	resp, err := call.Do()
	if err != nil {
		log.Printf("Errore nell'esecuzione della chiamata: %v", err)
		return nil, err
	}
	var bookList []*libr.Book
	for _, item := range resp.Items {
		b := buildBook(item)
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
	return &GoogleAPI{client: client}
}
