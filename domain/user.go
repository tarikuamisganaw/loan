package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email,omitempty"`
	Password   string             `bson:"password,omitempty"`
	Profile    string             `bson:"profile,omitempty"`
	IsVerified bool               `bson:"is_verified,omitempty"`
}

type AuthUser struct {
	Username string `json:"name"`
	Password string `json:"password"`
}
type UserUsecase interface {
	Register(req RegisterRequest) error
	VerifyEmail(email, token string) error
	Login(*AuthUser) (string, string, error)
	GetAllUsers() ([]User, error)
	DeleteUser(objectID primitive.ObjectID) error
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Profile  string `json:"profile" binding:"required"`
}
