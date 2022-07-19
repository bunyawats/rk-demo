package service

import (
	"context"
	"fmt"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	customer_collection = "customer"
)

type (
	MongoService struct {
		mongoDB func() *mongo.Database
		ctx     context.Context
		logger  *rkentry.LoggerEntry
	}

	CustomerDoc struct {
		CusId primitive.ObjectID `bson:"_id"`
		Fname string             `bson:"fname"`
		Lname string             `bson:"lname"`
		Age   int                `bson:"age"`
	}
)

func NewMongoService(conFun func() *mongo.Database, ctx context.Context) *MongoService {

	logger := rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	service := &MongoService{
		conFun,
		ctx,
		logger,
	}

	return service
}

func (s *MongoService) getCustomerCollection() *mongo.Collection {

	opts := options.CreateCollection()
	err := s.mongoDB().CreateCollection(s.ctx, customer_collection, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
	return s.mongoDB().Collection(customer_collection)

}

func (s *MongoService) GetCustomer(id string) *CustomerDoc {

	s.logger.Info("Call MongoService.GetUser")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create objectId: %v", err.Error()))
	}

	res := s.getCustomerCollection().FindOne(s.ctx, bson.M{"_id": objectId})
	if res.Err() != nil {
		s.logger.Error(fmt.Sprintf("Error while get customer from MongoDB: %v", res.Err().Error()))
	}
	customer := &CustomerDoc{}
	err = res.Decode(customer)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while get customer from MongoDB: %v", err.Error()))
	}

	return customer
}

func (s *MongoService) CreateCustomer(c *CustomerDoc) (string, error) {

	s.logger.Info("Call MongoService.CreateUser")

	c.CusId = primitive.NewObjectID()
	_, err := s.getCustomerCollection().InsertOne(s.ctx, c)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create c in MongoDB: %v", err.Error()))
		return "", err
	}

	return c.CusId.Hex(), nil
}

func (s *MongoService) ListCustomer() (*[]*CustomerDoc, error) {

	s.logger.Info("Call MongoService.ListCustomer")

	userList := make([]*CustomerDoc, 0)

	cursor, err := s.getCustomerCollection().Find(s.ctx, bson.D{})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while find customer list in MongoDB: %v", err.Error()))
		return nil, err
	}

	if err = cursor.All(s.ctx, &userList); err != nil {
		s.logger.Error(fmt.Sprintf("Error while decode customer list in MongoDB cursor: %v", err.Error()))
		return nil, err
	}

	return &userList, nil
}

func (s *MongoService) UpdateCustomer(c *CustomerDoc) (int64, error) {

	s.logger.Info("Call MongoService.UpdateCustomer")

	res, err := s.getCustomerCollection().UpdateOne(
		s.ctx,
		bson.M{
			"_id": c.CusId,
		},
		bson.D{
			{
				"$set", bson.D{
					{"fname", c.Fname},
					{"lname", c.Lname},
					{"age", c.Age},
				},
			},
		},
	)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while update user from MangoD : %v", err.Error()))
		return 0, err
	}
	return res.ModifiedCount, nil

}

func (s *MongoService) DeleteCustomer(id string) (int64, error) {

	s.logger.Info("Call MongoService.DeleteCustomer")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create objectId: %v", err.Error()))
		return 0, err
	}

	res, err := s.getCustomerCollection().DeleteOne(s.ctx, bson.M{"_id": objectId})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while delete user from MangoD : %v", err.Error()))
		return 0, err
	}
	return res.DeletedCount, nil
}
