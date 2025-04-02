package handlers

import (
	"context"
	"net/http"
	"recipe-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection,
		ctx,
	}
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)

	if err = cur.All(handler.ctx, &recipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &recipes)

}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
//	'404':
//	    description: Invalid recipe ID
func (handler *RecipesHandler) UpdateRecipeHandler(context *gin.Context) {
	id := context.Param("id")
	var recipe models.Recipe

	objectId, _ := primitive.ObjectIDFromHex(id)

	if err := context.ShouldBindJSON(&recipe); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{
		{"$set", bson.M{
			"name":         recipe.ID,
			"instructions": recipe.Instructions,
			"ingredients":  recipe.Ingredients,
			"tags":         recipe.Tags,
		}},
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been updated",
	})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func (handler *RecipesHandler) DeleteRecipeHandler(context *gin.Context) {
	id := context.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *RecipesHandler) SearchRecipesHandlerByTag(context *gin.Context) {
	tag := context.Query("tag")

	recipes := make([]any, 0)

	cur, err := handler.collection.Find(handler.ctx, bson.M{
		"tags": tag,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)

	if err = cur.All(handler.ctx, &recipes); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &recipes)
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
func (handler *RecipesHandler) NewRecipeHandler(context *gin.Context) {
	var recipe models.Recipe

	if err := context.ShouldBindJSON(&recipe); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	if _, err := handler.collection.InsertOne(handler.ctx, &recipe); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, &recipe)

}
