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
	user_collection = "user"
)

type (
	MongoService struct {
		mongoDB func() *mongo.Database
		ctx     context.Context
		logger  *rkentry.LoggerEntry
	}

	User struct {
		Id   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
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

func (s *MongoService) getUserCollection() *mongo.Collection {

	opts := options.CreateCollection()
	err := s.mongoDB().CreateCollection(s.ctx, user_collection, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
	return s.mongoDB().Collection(user_collection)

}

func (s *MongoService) GetUser(id string) *User {

	s.logger.Info("Call MongoService.GetUser")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create objectId: %v", err.Error()))
	}

	res := s.getUserCollection().FindOne(s.ctx, bson.M{"_id": objectId})
	if res.Err() != nil {
		s.logger.Error(fmt.Sprintf("Error while get user from MongoDB: %v", res.Err().Error()))
	}
	user := &User{}
	err = res.Decode(user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while get user from MongoDB: %v", err.Error()))
	}

	return user
}

func (s *MongoService) CreateUser(name string) *User {

	s.logger.Info("Call MongoService.CreateUser")

	user := &User{
		Id:   primitive.NewObjectID(),
		Name: name,
	}
	_, err := s.getUserCollection().InsertOne(s.ctx, user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create user in MongoDB: %v", err.Error()))
	}

	return user
}

func (s *MongoService) ListUser() *[]*User {

	s.logger.Info("Call MongoService.ListUser")

	userList := make([]*User, 0)

	cursor, err := s.getUserCollection().Find(s.ctx, bson.D{})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while find user list in MongoDB: %v", err.Error()))
	}

	if err = cursor.All(s.ctx, &userList); err != nil {
		s.logger.Error(fmt.Sprintf("Error while decode user list in MongoDB cursor: %v", err.Error()))
	}

	return &userList
}

func (s *MongoService) UpdateUser(id string, name string) int64 {

	s.logger.Info("Call MongoService.UpdateUser")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create objectId: %v", err.Error()))
	}
	u := &User{
		Id:   objectId,
		Name: name,
	}
	res, err := s.getUserCollection().UpdateOne(
		s.ctx,
		bson.M{
			"_id": objectId,
		},
		bson.D{
			{
				"$set", bson.D{
					{"name", u.Name},
				},
			},
		},
	)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while update user from MangoD : %v", err.Error()))
	}
	return res.ModifiedCount

}

func (s *MongoService) DeleteUser(id string) int64 {

	s.logger.Info("Call MongoService.DeleteUser")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while create objectId: %v", err.Error()))
	}

	res, err := s.getUserCollection().DeleteOne(s.ctx, bson.M{"_id": objectId})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while delete user from MangoD : %v", err.Error()))
	}
	return res.DeletedCount
}
