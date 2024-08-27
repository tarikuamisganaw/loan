package repositories

import (
	"context"
	"errors"
	"loan/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

type UserRepository interface {
	Create(user domain.User) error
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) Create(user domain.User) error {
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

func (r *mongoUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	filter := bson.D{{"email", email}}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *mongoUserRepository) Update(user *domain.User) error {
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"is_verified", user.IsVerified},
		}},
	}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}
func (r *mongoUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	var user domain.User
	// fmt.Printf("Searching for user with username: %s\n", username)
	err := r.collection.FindOne(context.TODO(), bson.M{"name": username}).Decode(&user)
	// fmt.Printf("User found: %+v\n", user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *mongoUserRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &users)

	return users, err
}

func (r *mongoUserRepository) DeleteUser(id primitive.ObjectID) error {

	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
