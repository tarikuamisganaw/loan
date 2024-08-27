package repositories

import (
	"context"
	"errors"
	"loan/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoanRepository struct {
	collection *mongo.Collection
}

// Initialize a new LoanRepository with a MongoDB collection
func NewLoanRepository(db *mongo.Database) *LoanRepository {
	return &LoanRepository{
		collection: db.Collection("loans"),
	}
}

// Save a new loan in the database
func (r *LoanRepository) Save(loan *domain.Loan) error {
	_, err := r.collection.InsertOne(context.TODO(), loan)
	return err
}

// Find a loan by its ID and User ID
func (r *LoanRepository) FindByIDAndUser(loanID string, userID string) (*domain.Loan, error) {
	var loan domain.Loan
	filter := bson.M{"_id": loanID, "user_id": userID}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&loan)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &loan, err
}

// Find a loan by its ID
func (r *LoanRepository) FindByID(loanID string) (*domain.Loan, error) {
	var loan domain.Loan
	filter := bson.M{"_id": loanID}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&loan)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &loan, err
}

// Find all loans, with optional status and order filters
func (r *LoanRepository) FindAll(status string, order string) ([]domain.Loan, error) {
	var loans []domain.Loan
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	findOptions := options.Find()
	if order == "desc" {
		findOptions.SetSort(bson.D{{"_id", -1}})
	} else {
		findOptions.SetSort(bson.D{{"_id", 1}})
	}

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &loans)
	return loans, err
}

// Update a loan in the database
func (r *LoanRepository) Update(loan *domain.Loan) error {
	filter := bson.M{"_id": loan.ID}
	update := bson.M{"$set": bson.M{
		"status": loan.Status,
	}}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// Delete a loan by its ID
func (r *LoanRepository) Delete(loanID string) error {
	filter := bson.M{"_id": loanID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("loan not found")
	}
	return nil
}
