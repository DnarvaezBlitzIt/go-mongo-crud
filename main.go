package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Persona struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nombre   string             `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido string             `json:"apellido,omitempty" bson:"apellido,omitempty"`
}

var client *mongo.Client

func renderPersonas(res http.ResponseWriter, req *http.Request) {

	log.Println("Mostrando data...")

	res.Header().Add("content-type", "application/json")
	collection := client.Database("go-mongo-crud").Collection("personas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)

	var personas []bson.M

	for cursor.Next(ctx) {

		var persona bson.M
		if err = cursor.Decode(&persona); err != nil {
			log.Fatal(err)
		}

		personas = append(personas, persona)

	}

	json.NewEncoder(res).Encode(personas)
	log.Println(personas)
	res.WriteHeader(http.StatusOK)

}

func crearPersona(res http.ResponseWriter, req *http.Request) {

	log.Println("Persistiendo data...")

	res.Header().Add("content-type", "application/json")
	var persona Persona
	json.NewDecoder(req.Body).Decode(&persona)
	collection := client.Database("go-mongo-crud").Collection("personas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, persona)
	json.NewEncoder(res).Encode(result)

}

func buscarPersona(res http.ResponseWriter, req *http.Request) {

	log.Println("Enviando data...")

	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var personaBusqueda Persona

	collection := client.Database("go-mongo-crud").Collection("personas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Persona{ID: id}).Decode(&personaBusqueda)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Persona no encontrada")
		res.Write([]byte("Persona no encontrada."))
		return
	}
	log.Printf("Persona encontrada: %s", personaBusqueda)
	json.NewEncoder(res).Encode(personaBusqueda)
	res.WriteHeader(http.StatusCreated)

}

func modificarPersona(res http.ResponseWriter, req *http.Request) {

	log.Println("Enviando data...")

	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var personaBusqueda Persona

	collection := client.Database("go-mongo-crud").Collection("personas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Persona{ID: id}).Decode(&personaBusqueda)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Persona no encontrada")
		res.Write([]byte("Persona no encontrada."))
		return
	}

	var personaModificada Persona

	json.NewDecoder(req.Body).Decode(&personaModificada)

	update, error := collection.UpdateOne(ctx, bson.M{"_id": personaBusqueda.ID}, bson.M{"$set": bson.M{"nombre": personaModificada.Nombre, "apellido": personaModificada.Apellido}})
	if error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Persona no modificada."))
		return
	}

	log.Printf("Persona modificada: %v", update)
	json.NewEncoder(res).Encode(personaModificada)
	res.WriteHeader(http.StatusOK)

}

func eliminarPersona(res http.ResponseWriter, req *http.Request) {

	log.Println("Enviando data...")

	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var personaBusqueda Persona

	collection := client.Database("go-mongo-crud").Collection("personas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Persona{ID: id}).Decode(&personaBusqueda)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Persona no encontrada")
		res.Write([]byte("Persona no encontrada."))
		return
	}

	update, error := collection.DeleteOne(ctx, personaBusqueda)
	if error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Persona no eliminada."))
		return
	}

	log.Printf("Persona modificada: %v", update)
	res.WriteHeader(http.StatusOK)

}

func main() {

	log.Println("Entrando en el servicio...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/crearPersona", crearPersona).Methods("POST")
	router.HandleFunc("/buscarPersona/{id}", buscarPersona).Methods("GET")
	router.HandleFunc("/modificarPersona/{id}", modificarPersona).Methods("PUT")
	router.HandleFunc("/renderPersonas", renderPersonas).Methods("GET")
	router.HandleFunc("/eliminarPersona/{id}", eliminarPersona).Methods("DELETE")
	http.ListenAndServe(":8081", router)

}
