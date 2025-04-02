// Package main Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/sanjbh/go-recipe-api
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Sanjay Bhattacharya <sanjbh@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"os"
	"recipe-api/handlers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// var client *mongo.Client
// var collection *mongo.Collection
// var ctx context.Context
// var recipes []models.Recipe

var recipesHandler *handlers.RecipesHandler

func init() {
	var mongoURL, dbName string
	ctx := context.Background()

	if u, ok := os.LookupEnv("MONGO_URL"); !ok {
		mongoURL = "mongodb://admin:password@localhost:27017/test?authSource=admin"
	} else {
		mongoURL = u
	}

	if d, ok := os.LookupEnv("MONGO_DB"); !ok {
		dbName = "demo"
	} else {
		dbName = d
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.
		Database(dbName).
		Collection("recipes")

	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {

	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandlerByTag)
	router.Run()
}
