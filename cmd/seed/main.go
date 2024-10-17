package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/masatrio/bookstore-api/config"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := seed(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}

func seed(db *sql.DB) error {
	users := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"Satrio", "satrio@test.test", "hashed_password"},
		{"Satmoko", "satmoko@test.test", "hashed_password"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
		if err != nil {
			return err
		}
	}

	books := []struct {
		Title  string
		Author string
		Price  float64
	}{
		{"The Hobbit", "J.R.R. Tolkien", 150000},
		{"1984", "George Orwell", 120000},
		{"The Catcher in the Rye", "J.D. Salinger", 130000},
		{"To Kill a Mockingbird", "Harper Lee", 140000},
		{"The Great Gatsby", "F. Scott Fitzgerald", 110000},
		{"Pride and Prejudice", "Jane Austen", 160000},
		{"Moby Dick", "Herman Melville", 180000},
		{"War and Peace", "Leo Tolstoy", 250000},
		{"The Divine Comedy", "Dante Alighieri", 300000},
		{"Crime and Punishment", "Fyodor Dostoevsky", 190000},
		{"The Odyssey", "Homer", 220000},
		{"Brave New World", "Aldous Huxley", 140000},
		{"The Iliad", "Homer", 230000},
		{"The Brothers Karamazov", "Fyodor Dostoevsky", 210000},
		{"The Alchemist", "Paulo Coelho", 170000},
		{"One Hundred Years of Solitude", "Gabriel Garcia Marquez", 200000},
		{"Don Quixote", "Miguel de Cervantes", 240000},
		{"Ulysses", "James Joyce", 260000},
		{"The Sound and the Fury", "William Faulkner", 180000},
		{"Madame Bovary", "Gustave Flaubert", 190000},
		{"The Lord of the Rings", "J.R.R. Tolkien", 270000},
		{"Jane Eyre", "Charlotte Bronte", 150000},
		{"The Old Man and the Sea", "Ernest Hemingway", 120000},
		{"Wuthering Heights", "Emily Bronte", 140000},
		{"The Stranger", "Albert Camus", 130000},
		{"Les Miserables", "Victor Hugo", 280000},
		{"Anna Karenina", "Leo Tolstoy", 250000},
		{"The Count of Monte Cristo", "Alexandre Dumas", 260000},
	}

	for _, book := range books {
		_, err := db.Exec("INSERT INTO books (title, author, price) VALUES ($1, $2, $3)", book.Title, book.Author, book.Price)
		if err != nil {
			return err
		}
	}

	return nil
}
