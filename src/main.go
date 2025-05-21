package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	echoLambda     *echoadapter.EchoLambda
	client         *mongo.Client
	userCollection *mongo.Collection
)

// var validate = validator.New()
var uri = "add the uri of yor MONGODB ATLAS here"

func initMongoDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}

	// Ping to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	userCollection = client.Database("serverless-crud").Collection("users") // Adjust DB/collection name
	fmt.Println("✅ Connected to MongoDB")
}

type Users struct {
	Name  string `json:"name"  bson:"name"`
	Email string `json:"email"  bson:"email"`
}

func init() {
	fmt.Println("echo cold start")

	initMongoDB()
	e := echo.New()

	e.POST("/create/user", func(c echo.Context) error {

		//you are creating new user
		var user Users
		//binding the request with the model
		if err := c.Bind(&user); err != nil {
			log.Printf("failed to decode request: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid request format",
			})
		}

		//created contex for the mongoConnection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//geting the collection
		userCollection = client.Database("serverless-crud").Collection("users")

		//store the data into the db
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			log.Println("Error in inserting the data")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "failed to insert user",
			})
		}

		fmt.Printf("Inserted the _id: %v", result)

		//seding back the response to user
		return c.JSON(http.StatusCreated, echo.Map{
			"message": "user create successfully",
		})
	})

	e.GET("/get/users", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("❌ Failed to retrieve users: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to retrieve users",
			})
		}
		defer cursor.Close(ctx)

		var users []Users
		if err := cursor.All(ctx, &users); err != nil {
			log.Printf("❌ Failed to decode users: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "failed to decode users",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "users retrieved successfully",
			"data":    users,
		})
	})
	e.PUT("/update/user/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			log.Printf("invalid id format: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid is a format",
			})
		}

		var updatedUser Users
		if err := c.Bind(&updatedUser); err != nil {
			log.Printf("failed to decode request: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid request format",
			})
		}

		//UPDATE THE USER TO THE DB
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"name":  updatedUser.Name,
				"email": updatedUser.Email,
			},
		}
		result, err := userCollection.UpdateByID(ctx, objectID, update)
		if err != nil {
			log.Printf("failed to update the user")
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "failed to update the user",
			})
		}

		if result.MatchedCount == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "user not found",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "user updated successfully",
		})
	})

	e.DELETE("/delete/user/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			log.Printf("invalid id format: %v", err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "invalid id format",
			})
		}

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			log.Println("error with deleting the user")

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "erro with delteing the user",
			})
		}

		if result.DeletedCount == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "user not found",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "user deleted successfully",
		})

	})

	echoLambda = echoadapter.New(e)

}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return echoLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(handler)
}
