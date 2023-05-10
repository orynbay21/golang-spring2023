package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golangfinal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson" //pkg bson for reading,writing,manipulating BSON

	//BSON - binary serialization format used to store docs

	"go.mongodb.org/mongo-driver/bson/primitive" //pkg primitive has types similar to GO primitivies for BSON types
	//that dont have direct Go primitive representations

	"go.mongodb.org/mongo-driver/mongo" //pkg mongo -MongoDB driver API for GO
)

/*
Context package defines the Context type, which carries deadlines, cancellation signals and other request-scoped values across API boundaries and between processes
Context is an object whose first purpose is to cancel an operation with potentially big latency
Context is a great way to store and transfer data between the methods of our program
context.Background()         context.TODO()
Context.WithCancel() – явный сигнал отмены
	Context.WithTimeout() – сигнал с таймаутом например 100 сек
	Context.WithDeadline() – сигнал с таймстемпом
Context.WithValue() –передача данных через контекст
Принимает родительский контекст, ключ и значение
Возращает только контекст без функции отмены
Используется только в крайних случаях, передача данных через контекст не рекомендуется
*/

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id") //returns the value of the key if it exists
		if user_id == "" {
			//setting the response header - metadata that can be included with JSON data to provide additional context and information about the data being transmitted.
			c.Header("Content-Type", "application/json")

			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			//any remaining handlers that are set up to execute in the current HTTP request will not be executed.
			c.Abort()
			return
		}
		//ObjectIDFromHex creates a new ObjectID from a hexadecimal string. It returns an error if the hex string is not a valid ObjectID.
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		//accesing the addresses from the models.Address slice
		var addresses models.Address

		//NewObjectID() function generates a new object id
		addresses.Address_id = primitive.NewObjectID()
		//parse the JSON request body and bind it to the addresses variable.
		// If there is an error during this process, the error message is returned.
		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		//дочерний            //родительский контекст
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		/*
			what is bson?
			it is a binary representation of the stored documents used in MongoDb
			easy+flexible data processing
			D-ordered representation of BSON document - slice
			M-unordered -map
			A-ordered -array
			E-single element inside a slice
		*/
		//stages of agregation pipeline:(agregation-скопление)
		// Each stage takes input from the previous stage and produces output that is passed to the next stage in the pipeline.

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		//filters a BSON(binary json) document
		//bson.D constructor creates on ordered list of key-value pairs
		//$match key operator filters docs based on a condition that is set in Value
		//Value is another BSON document that specifies the _id field to match against
		//in a primitive.E (single element inside a slice)

		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		//unwind operator is used to deconstructing an array field in a document and create
		//separate output documents for each item in the array."flattening"
		// the path is "address", which suggests that there is an array field in the document with this name that needs to be unwound.

		//purpose here - find out how many addressess this user has bcs we have home address+work address

		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		// for each group of documents with the same $address_id value, the $sum operator will count the number
		//of documents in the group and return a value of 1 for each document.
		//Finally, the $group operator will count all of these 1s and return the total count of documents in each group.

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		//var pointcursor *mongo.Cursor - A pointer to the result set of a query

		//The 'ctx' context is passed to the Aggregate function in order to provide contextual information
		//about the context in which the operation is being performed.(info abt timouts,cancellation signals)
		//ctx is being passed to the aggregate function to ensure that the operation doesnt block indefinitely

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addressinfo []bson.M
		//bson.M -map  alias to []primitive.M

		//converting json to golang un
		///START FROM HERE
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			//All() function is used to retrieve all of the documents in a MongoDB collection.
			panic(err)
		}

		var size int32
		/*
			for index,value:=range arrayName{
				}
			_ to ignore the index
		*/
		for _, address_no := range addressinfo {
			//
			count := address_no["count"]
			size = count.(int32) //state that the value stored in count is of type int32
			//assign that value to the size variable
			//size ==number of addresses
		}

		//if there is only work address u can add home address
		//but more than 2 is not allowed
		if size < 2 {
			//filter query
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			//use $push when u expect that some values will be there
			//and ull need to add more
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
			c.IndentedJSON(200,"Successfully added your address!")
		} else {
			c.IndentedJSON(400, "Not Allowed ")
		}
		//Canceling this context releases resources associated with it
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editaddress models.Address //the new edited address
		//is of the same type as the old one
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		//filtering by the id of the specifical user that want to change it's home address
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		//the 0s in the variable names indicate that we are editing the Home address, or 1st
		//address associated with the current id and not the Work address
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}
		// to update at most one document in the collection.
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated the Home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "something Went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the Work Address")
	}
}

// deletes both home and work address for simplicity
func DeleteAddress() gin.HandlerFunc {
	//gin context helps to get access to things like query
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		//creating an empty slice
		addresses := make([]models.Address, 0)

		//to convert a user id string in hexdecimal
		//format to an ObjectID type in the MongoDB go driver
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		//check for error during  the conversion process
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		//why d u have timeout?
		//whenever a server is working w a database, u cant have it endlessly waiting for the
		//result(fr ex if the server is going down)

		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		//The address field of the document is updated with a new value,
		//which is stored in a variable called addresses.
		//addresses variable is filled with 0s as we mentioned b4
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted!")
	}
}
