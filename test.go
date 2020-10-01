package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type login_user struct {
	User_id  string
	Name     string
	Password string
	Status   int
}

var mongo_client *mongo.Database

func init_mongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

func main() {
	// Rest of the code will go here
	client := init_mongo()
	mongo_client = client.Database("testdb")

	// insert_mongo("taysk", "sen kwan", "password1", 1)
	http.HandleFunc("/add_user", add_user)
	http.HandleFunc("/get_users", get_users)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func insert_mongo(user_id string, name string, password string, status int) {

	collection := mongo_client.Collection("login_users")
	new_user := login_user{user_id, name, password, status}

	insertResult, err := collection.InsertOne(context.TODO(), new_user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a user data: ", user_id, " , name ", name, ", password", password, ", status,", status)
	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)

}

func add_user(w http.ResponseWriter, r *http.Request) {

	var data login_user

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	log.Println(data)

	insert_mongo(data.User_id, data.Name, data.Password, data.Status)

}

func get_users(w http.ResponseWriter, r *http.Request) {

	collection := mongo_client.Collection("login_users")

	// Call the collection's Find() method to return Cursor obj
	// with all of the col's documents
	cursor, err := collection.Find(context.TODO(), bson.D{})

	// Find() method raised an error
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(context.TODO())

		// If the API call was a success
	} else {
		// iterate over docs using Next()
		w.Header().Set("Content-Type", "application/json")
		for cursor.Next(context.TODO()) {

			// declare a result BSON object
			var result bson.M
			err := cursor.Decode(&result)

			// If there is a cursor.Decode error
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				os.Exit(1)

				// If there are no cursor.Decode errors
			} else {
				// w.Write(byte[] ("\nresult type:", reflect.TypeOf(result)))
				fmt.Println("result:", result)
				json.NewEncoder(w).Encode(result)
			}
		}
	}

}
