package domain

type Loan struct {
	ID             string  `json:"id" bson:"_id"`
	UserID         string  `json:"user_id" bson:"user_id"`
	Amount         float64 `json:"amount" bson:"amount"`
	InterestRate   float64 `json:"interest_rate" bson:"interest_rate"`
	DurationMonths int     `json:"duration_months" bson:"duration_months"`
	Status         string  `json:"status" bson:"status"` // pending, approved, rejected
}
