package books

import (
	"errors"
	"go-mongo/config"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32 `bson:"Price,truncate"`
}

func AllBooks(r *http.Request) ([]Book, error) {
	ctx := r.Context()
	c, err := config.Books.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var bks []Book
	if err = c.All(ctx, &bks); err != nil {
		return nil, err
	}

	return bks, nil
}

func OneBook(r *http.Request) (Book, error) {
	ctx := r.Context()
	var bk Book
	isbn := r.FormValue("isbn")
	if isbn == "" {
		return bk, errors.New("400. Bad Request.")
	}
	c, err := config.Books.Find(ctx, bson.M{"Isbn": isbn})
	if err != nil {
		return bk, err
	}
	var bks []Book
	if err = c.All(ctx, &bks); err != nil {
		return bk, err
	}
	if len(bks) < 1 {
		return bk, nil
	}
	return bks[0], nil
}

func PutBook(r *http.Request) (Book, error) {
	ctx := r.Context()
	// get form values
	bk := Book{}
	bk.Isbn = r.FormValue("isbn")
	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")

	// validate form values
	if bk.Isbn == "" || bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Price must be a number.")
	}
	bk.Price = float32(f64)

	// insert values
	_, err = config.Books.InsertOne(ctx, bson.M{"Isbn": bk.Isbn, "Title": bk.Title, "Author": bk.Author, "Price": bk.Price})
	if err != nil {
		return bk, errors.New("500. Internal Server Error." + err.Error())
	}
	return bk, nil
}

func UpdateBook(r *http.Request) (Book, error) {
	ctx := r.Context()
	// get form values
	bk := Book{}
	bk.Isbn = r.FormValue("isbn")
	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")

	if bk.Isbn == "" || bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad Request. Fields can't be empty.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Enter number for price.")
	}
	bk.Price = float32(f64)

	// insert values
	filter := bson.D{{"Isbn", bk.Isbn}}
	update := bson.D{{"$set", bson.M{"Title": bk.Title, "Author": bk.Author, "Price": bk.Price}}}
	_, err = config.Books.UpdateOne(ctx, filter, update)
	if err != nil {
		return bk, err
	}
	return bk, nil
}

func DeleteBook(r *http.Request) error {
	ctx := r.Context()
	isbn := r.FormValue("isbn")
	if isbn == "" {
		return errors.New("400. Bad Request.")
	}

	_, err := config.Books.DeleteOne(ctx, bson.M{"Isbn": isbn})
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
