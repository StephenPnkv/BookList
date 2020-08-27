//Stephen Penkov 8/26/20
//A basic REST API connected to a mongodb database that handles CRUD operations.
//Tested using Postman

package main


import (
	"fmt"
	"time"
	"log"
	"context"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID int 					`bson:"id,omitempty" json:"id"`
	Title string 			`bson:"title,omitempty" json:"title"`
	Author string 			`bson:"author,omitempty" json:"author"`
	Year string 			`bson:"year,omitempty" json:"year"`
}


var books []Book 

func logFatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}

//Load env variables
var mongoURI, dbName, colName, port string 

func init(){
	gotenv.Load()
	mongoURI = os.Getenv("MONGO_URI")
	dbName = os.Getenv("DBNAME")
	colName = os.Getenv("COLLECTION_NAME")
	port = os.Getenv("PORT")
}


func getMongoDBConnection()(*mongo.Client,error){
	//Helper function connects a client to mongodb, returns 
	//client and error.

	//Connect to client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	
	if err != nil{
		log.Fatal(err)
	}
	//Test connection
	err = client.Ping(context.TODO(),nil)
	if err != nil{
		log.Fatal(err)
	}
	//Return client and err
	return client, err 
}

func getMongoDBCollection()(*mongo.Collection, error){
	//Function establishes a connection to mongodb and returns the collection and error.

	client, err := getMongoDBConnection()
	if err != nil{
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(colName)
	return collection, err 
}

func main(){

	fmt.Println("Starting server. . .")
	//Create router instance
	fmt.Printf("Database name: %s, Collection name: %s\nMongoURI: %s\n",
			dbName,colName,mongoURI)

	router := mux.NewRouter() 

	//Create routes
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}",updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	//Create server and bind to port
	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(port,router)
}

func getBooks(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")
	collection, err := getMongoDBCollection()
	var allBooks []Book 
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	cursor, e := collection.Find(ctx,bson.M{})
	if e != nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}` ))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var book Book 
		cursor.Decode(&book)
		allBooks = append(allBooks,book)
	}
	if err = cursor.Err(); err != nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}` ))
		return
	}
	json.NewEncoder(res).Encode(allBooks)
}

func addBook(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")

	var book Book 
	json.NewDecoder(req.Body).Decode(&book)
	fmt.Printf("Added Book - ID: %d | title: %s | author: %s | year: %s\n", book.ID, book.Title, book.Author, book.Year)

	collection, err := getMongoDBCollection()
	if err != nil{
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx,book)
	json.NewEncoder(res).Encode(result)
}

func getBook(res http.ResponseWriter, req *http.Request){
	
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	var book Book
	
	id, _ := strconv.Atoi(params["id"])

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	collection, err := getMongoDBCollection()
	if err != nil{
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx,filter).Decode(&book)
	if err != nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}` ))
		return
	}
	json.NewEncoder(res).Encode(&book)
}

func updateBook(res http.ResponseWriter, req *http.Request){
	
	res.Header().Add("content-type", "application/json")
	var book Book
	
	json.NewDecoder(req.Body).Decode(&book)
	fmt.Printf("Updated Book - ID: %d | title: %s | author: %s | year: %s\n", book.ID, book.Title, book.Author, book.Year)

	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])

	filter := bson.M{"id": id}
	updatedBook := bson.D{
			{"$set",bson.D{
				{"id",id},
				{"title",book.Title},
				{"author",book.Author},
				{"year",book.Year},
			}},
	}
	
	collection, err := getMongoDBCollection()
	if err != nil{
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, e := collection.UpdateOne(
		ctx,
		filter,
		updatedBook,
	)
	if e != nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}` ))
		return
	}

	json.NewEncoder(res).Encode(&result)
}

func removeBook(res http.ResponseWriter, req *http.Request){
	params := mux.Vars(req) 
	id, _ := strconv.Atoi(params["id"])
	filter := bson.D{{"id", id}}
	fmt.Printf("Removed Book - ID: %d\n", id)

	collection, err := getMongoDBCollection()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	result, e := collection.DeleteOne(ctx,filter)
	if e != nil {
		log.Fatal(e)
	}
	json.NewEncoder(res).Encode(&result)
}



