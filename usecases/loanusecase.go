package usecaseses

import (
	"errors"
	"loan/domain"
	"loan/repositories"
)

// Apply for Loan Use Case
type LoanUseCases struct {
	LoanRepository repositories.LoanRepository
	UserRepository repositories.UserRepository
}

// Apply for Loan
func (u *LoanUseCases) ApplyForLoan(userID string, amount float64, interestRate float64, durationMonths int) (string, error) {
	user, err := u.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		return "", errors.New("user does not exist")
	}

	loan := domain.Loan{
		ID:             generateLoanID(), // Assume a utility function for generating loan IDs
		UserID:         userID,
		Amount:         amount,
		InterestRate:   interestRate,
		DurationMonths: durationMonths,
		Status:         "pending",
	}

	err = u.LoanRepository.Save(&loan)
	if err != nil {
		return "", err
	}

	return "Loan application submitted successfully", nil
}

// View Loan Status
func (u *LoanUseCases) ViewLoanStatus(loanID string, userID string) (*domain.Loan, error) {
	loan, err := u.LoanRepository.FindByIDAndUser(loanID, userID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, errors.New("loan not found or you don't have access")
	}
	return loan, nil
}

// Admin View All Loans
func (u *LoanUseCases) ViewAllLoans(status, order string) ([]domain.Loan, error) {
	loans, err := u.LoanRepository.FindAll(status, order)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

// Admin Approve/Reject Loan
func (u *LoanUseCases) UpdateLoanStatus(loanID string, status string) (string, error) {
	loan, err := u.LoanRepository.FindByID(loanID)
	if err != nil || loan == nil {
		return "", errors.New("loan not found")
	}

	if status != "approved" && status != "rejected" {
		return "", errors.New("invalid status")
	}

	loan.Status = status
	err = u.LoanRepository.Update(loan)
	if err != nil {
		return "", err
	}

	return "Loan " + status + " successfully", nil
}

// Admin Delete Loan
func (u *LoanUseCases) DeleteLoan(loanID string) (string, error) {
	loan, err := u.LoanRepository.FindByID(loanID)
	if err != nil || loan == nil {
		return "", errors.New("loan not found")
	}

	err = u.LoanRepository.Delete(loanID)
	if err != nil {
		return "", err
	}

	return "Loan deleted successfully", nil
}
