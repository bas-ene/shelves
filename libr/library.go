package libr

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type ISBN string

type Book struct {
	Isbn       ISBN     `json:"isbn"`
	GoogleId   string   `json:"googleId"`
	Title      string   `json:"title"`
	Authors    []string `json:"authors"`
	Pages      int      `json:"pages"`
	Categories []string `json:"categories"`
	Language   string   `json:"language"`
	Price      float64  `json:"price"`
	Currency   string   `json:"currency"`
	PagesRead  int      `json:"pagesRead"`
	ImgHref    string   `json:"imgHref"`
}

type Library []Book

func (lib *Library) Save(filename string) error {
	data, err := json.MarshalIndent(*lib, "", "   ")

	if err != nil {
		return fmt.Errorf("could not marshal books: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func (lib *Library) RemoveBook(gId string) *Book {
	var delBook *Book = nil

	*lib = slices.DeleteFunc(*lib, func(b Book) bool {
		if b.GoogleId == gId {
			delBook = &b
			return true
		}
		return false
	})

	return delBook
}

func (lib *Library) ReadBook(gId string, nPages int) *Book {
	for _, book := range *lib {
		if book.GoogleId == gId {
			if nPages < 0 {
				book.PagesRead = max(0, book.PagesRead+nPages)
			} else {
				book.PagesRead = min(book.Pages, book.Pages+nPages)
			}
			return &book
		}
	}
	return nil
}

func (lib *Library) AddBook(book *Book) {
	if book == nil {
		return
	}
	*lib = append(*lib, *book)
}

func LoadLibrary(filename string) (Library, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return Library{}, nil
	}

	var books Library
	err = json.Unmarshal(data, &books)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal books: %w", err)
	}
	return books, nil
}

func (book Book) String() string {
	str := "Title: %s\n" +
		"Authors: %v\n" +
		"Pages: %d\n" +
		"Categories: %v\n" +
		"Language: %s\n" +
		"Price: %.2f %s\n"
	return fmt.Sprintf(str, book.Title, book.Authors, book.Pages, book.Categories, book.Language, book.Price, book.Currency)
}
