package main


import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

type Book struct {
	ID int `json:id`
	Title string `json:title`
	Author string `json:author`
	Year string `json:year`
}

var books []Book 

func main(){

	books = append(books,Book{ID: 1, Title: "Pointers", Author: "Mr. Pointer", Year: "2012"},
		Book{ID: 2, Title: "Recursion", Author: "Mr. Recursion", Year: "2019"},
		Book{ID: 3, Title: "Concurrency", Author: "Mr. Concurrency", Year: "2014"},
		Book{ID: 4, Title: "Go: Data Structures & Algorithms", Author: "Mr. Data", Year: "2015"},
		Book{ID: 5, Title: "Hacking with Go", Author: "Mr. Hacker", Year: "2019"},
	)

	router := mux.NewRouter() 
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books",updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	port := fmt.Sprintf(":%d",8000)
	http.ListenAndServe(port,router)
}

func getBooks(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(books)
}
func addBook(w http.ResponseWriter, r *http.Request){
	var book Book 
	_ = json.NewDecoder(r.Body).Decode(&book)

	books = append(books, book)

	json.NewEncoder(w).Encode(books)
}
func getBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	i, _ := strconv.Atoi(params["id"])

	for _, book := range books{
		if book.ID == i{
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request){
	log.Println("updatebook is called.") 
}
func removeBook(w http.ResponseWriter, r *http.Request){
	log.Println("removebook is called.")
}



