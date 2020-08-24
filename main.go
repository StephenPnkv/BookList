package main


import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Book struct {
	ID int `json:id`
	Title string `json:title`
	Author string `json:author`
	Year string `json:year`
}

var books []Book 

func main(){
	router := mux.NewRouter() 
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books",updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	port := fmt.Sprintf(":%d",8000)
	http.ListenAndServe(port,router)
}

func getBook(w http.ResponseWriter, r *http.Request){
	log.Println("getbook is called.")
}
func addBook(w http.ResponseWriter, r *http.Request){
	log.Println("addbook is called.")
}

func updateBook(w http.ResponseWriter, r *http.Request){
	log.Println("updatebook is called.")
}
func removeBook(w http.ResponseWriter, r *http.Request){
	log.Println("removebook is called.")
}
func getBooks(w http.ResponseWriter, r *http.Request){
	log.Println("getbooks is called.")
}



