package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DigiEgg struct
type DigiEgg struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Stage string             `json:"stage" bson:"stage"` // Incubation stage
}

// Global MongoDB variable
var collection *mongo.Collection

// Connect to MongoDB
func connectDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Get MongoDB configuration
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Failed to connect MongoDB:", err)
	}

	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ MongoDB connection failed:", err)
	}

	// Connect to the database and collection
	collection = client.Database(mongoDB).Collection("eggs")
	fmt.Println("✅ Connected to MongoDB!")
}

func main() {
	// Connect to the database
	connectDB()

	// Initialize Gin
	r := gin.Default()

	// Define routes
	r.POST("/eggs", createEgg)   // Create DigiEgg
	r.GET("/eggs", getAllEggs)   // Get all DigiEggs
	r.GET("/eggs/:id", getEgg)   // Get DigiEgg by ID
	r.PUT("/eggs/:id", updateEgg) // Update DigiEgg
	r.DELETE("/eggs/:id", deleteEgg) // Delete DigiEgg

	// Run the server
	r.Run(":8080")
}

// Create DigiEgg
func createEgg(c *gin.Context) {
	var newEgg DigiEgg

	// Bind JSON request body
	if err := c.ShouldBindJSON(&newEgg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set MongoDB ID
	newEgg.ID = primitive.NewObjectID()

	// Insert into database
	_, err := collection.InsertOne(context.TODO(), newEgg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create egg"})
		return
	}

	c.JSON(http.StatusCreated, newEgg)
}

// Get all DigiEggs
func getAllEggs(c *gin.Context) {
	var eggs []DigiEgg

	// Find all records
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get eggs"})
		return
	}
	defer cursor.Close(context.TODO())

	// Decode the results
	for cursor.Next(context.TODO()) {
		var egg DigiEgg
		cursor.Decode(&egg)
		eggs = append(eggs, egg)
	}

	c.JSON(http.StatusOK, eggs)
}

// Get DigiEgg by ID
func getEgg(c *gin.Context) {
	id := c.Param("id")

	// Convert ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Find the egg by ID
	var egg DigiEgg
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&egg)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Egg not found"})
		return
	}

	c.JSON(http.StatusOK, egg)
}

// Update DigiEgg
func updateEgg(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedEgg DigiEgg
	if err := c.ShouldBindJSON(&updatedEgg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the egg in the database
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedEgg},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update egg"})
		return
	}

	c.JSON(http.StatusOK, updatedEgg)
}

// Delete DigiEgg
func deleteEgg(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Delete the egg from the database
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete egg"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Egg deleted"})
}
