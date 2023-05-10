package controllers

import (
  "context"
  "fmt"
  "net/http"
  "time"

  "golangfinal/models"
	//"log"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson" //pkg bson for reading,writing,manipulating BSON

  //BSON - binary serialization format used to store docs

  "go.mongodb.org/mongo-driver/bson/primitive" //pkg primitive has types similar to GO primitivies for BSON types
  //that dont have direct Go primitive representations

 // "go.mongodb.org/mongo-driver/mongo" //pkg mongo -MongoDB driver API for GO
)
/*
type Product struct {
  Product_ID       primitive.ObjectID   bson:"_id"
  Product_Name     *string              json:"product_name"
  Price            *int               json:"price"
  Comment             []Comment           json:"comment" bson:"comment"
}
type Comment struct{
  Comment_id     primitive.ObjectID    bson:"_id"
  Comment        *string                 json:"comment" bson:"comment"
  Rating          *int                    json:"rating" bson:"rating"
}

*/
func AddComment() gin.HandlerFunc {
  return func(c *gin.Context) {
    product_id:= c.Query("id") //returns the value of the key if it exists
    if product_id == "" {
      //setting the response header - metadata that can be included with JSON data to provide additional context and information about the data being transmitted.
      c.Header("Content-Type", "application/json")

      c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
      //any remaining handlers that are set up to execute in the current HTTP request will not be executed.
      c.Abort()
      return
    }
    comment,err:=primitive.ObjectIDFromHex(product_id)
    //address, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
      c.IndentedJSON(500, "Internal Server Error")
    }
    var comments models.Comment
    comments.Comment_id=primitive.NewObjectID()
    if err = c.BindJSON(&comments); err != nil {
      c.IndentedJSON(http.StatusNotAcceptable, err.Error())
    }
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
    filter := bson.D{primitive.E{Key: "_id", Value: comment}}
    update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "comment", Value: comments}}}}
    _, err =ProductCollection.UpdateOne(ctx, filter, update)
    if err != nil {
      fmt.Println(err)
    }
    c.IndentedJSON(200,"Successfully added your comment!")
    defer cancel()
    ctx.Done()
  }
}