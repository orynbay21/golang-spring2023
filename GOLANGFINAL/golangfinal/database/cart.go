package database

import (
	"context"
	"errors"
	"log"
	"time"

	"golangfinal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	//defining custom errors
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find product")
	ErrUserIDIsNotValid   = errors.New("user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add product to cart")
	ErrCantRemoveItem     = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from cart ")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	//b4 u add the project to the cart
	//get the product u want from the db(get by id)

	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	//now that u got it from db check for err
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productcart []models.ProductUser 

	//.All() is putting a product from searchfromdb to the productCart
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	//to update smth u need id(user id)
	
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	//if everything is well u return nil instead of the error
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	//USER - is a collection
	//UserCart -  is a field in that collection
	//therefore we just want to update this field by removind one item from it

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	//what is $pull in mongodb??
	//pull is mongodb method to delete that specific product
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	
	//upd has to have some way to get to THIS users data
	//and go to that users cart
	//and delete that product with that particular id

	//UpdateMany function from mongodb
	//method to update 1 or more documents that matches with the specified filter criteria in a collection.
	//try replacing it with UpdateOne
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	var getcartitems models.User
	var ordercart models.Order
	//checkout:

	//initializing the fields from the Order struct
	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Orderered_At = time.Now()
	//make an empty slice of the user's product cart that is ALL an order now

	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true
	
	//fetch the cart of the user
	//find the cart total price
	//create an order with the items
	//added order to the user Collection
	//added items in the cart to order list
	//empty up the cart

	//find the cart total price: 

	//UserCart of the user contains multiple instances of Product Model
	//we need to flatte this data in the form {usercart}-apple,{usercart}-grapes,{usercart}-banana etc.
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	//here The path parameter is used by the Atlas Search operators to specify the field or fields to be searched

	//group the flattened data by the !id of the user, find out the !total  by !summing the !price field 
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getusercart []bson.M
		
	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}

	var total_price int32
	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
		//int's size is platform-dependent, whereas int32's size is always 32 bits.
	}
	//finished creating the fields of the order
	ordercart.Price = int(total_price)

	//create an order itself
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}
	//find the _id field which is =="id" value
	//using MongoDB query syntax
	//"decode" method decodes the result of the query
	//and maps it to the "getcartitems" variable
	//if there is no match for the query, the 'err' msg will appear

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}
	//adding an orderlist
	//order_list is the alias of orderCart
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	
	//creating an update operation that will add items to a MongoDB array field called order_list
	//the "$each" operator allows multiple items to be added to the array at once
	//The items being added to the array are from the "getcartitems.UserCart" variable.
	//The "$[]" operator is used to target all elements in the "orders" array
	
	//orders is an alias for the orderCart
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}

	//_ stands for the number of the documents modified, and we ignore it
	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		log.Println(err)
	}
	//emptying the cart after buying everything
	usercart_empty := make([]models.ProductUser, 0)

	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filtered, updated)

	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {
	//instant buy - taking a product and not putting it into the cart but buying it instantly instead
	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	//taking the structure from the Product Cart
	var product_details models.ProductUser
	//even though u dont have to put a product in the cart, u still have to creat an order for it

	//taking the structure from the Order
	var orders_detail models.Order

	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Orderered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true
	

	//finding the product u want to instantly buy, decode it and put it in the cart
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)

	if err != nil {
		log.Println(err)
	}
	//checkout 
	//the total price ==the price of the product
	orders_detail.Price = product_details.Price
	
	//id of the user that wanted to buy that product is used in the filter
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}
	//orders is an alies for the order cart
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
