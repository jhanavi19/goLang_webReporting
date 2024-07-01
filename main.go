package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Report struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Occupation string             `json:"occupation,omitempty" bson:"occupation,omitempty"`
	Hobby      string             `json:"hobby,omitempty" bson:"hobby,omitempty"`
}

func initMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)

	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

func generateReport(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var report Report
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("mydatabase").Collection("reportDetails")
	result, err := collection.InsertOne(context.Background(), report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Report saved successfully",
		"id":      result.InsertedID,
	})
	w.WriteHeader(http.StatusOK)
}

func getReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("mydatabase").Collection("reportDetails")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var reports []Report
	for cur.Next(context.Background()) {
		var report Report
		err := cur.Decode(&report)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reports = append(reports, report)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reports)
}

func main() {

	initMongo()

	router := mux.NewRouter()

	router.HandleFunc("/generate-report", generateReport).Methods("POST")
	router.HandleFunc("/api/reports", getReports).Methods("GET")

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
