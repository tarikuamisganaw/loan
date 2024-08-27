package controller

import (
	"loan/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	ApplyForLoanUseCase      usecase.ApplyForLoanUseCase
	ViewLoanStatusUseCase    usecase.ViewLoanStatusUseCase
	ViewAllLoansAdminUseCase usecase.ViewAllLoansAdminUseCase
	UpdateLoanStatusUseCase  usecase.UpdateLoanStatusUseCase
	DeleteLoanUseCase        usecase.DeleteLoanUseCase
}

func (c *LoanController) ApplyForLoan(ctx *gin.Context) {
	var input struct {
		Amount         float64 `json:"amount"`
		InterestRate   float64 `json:"interest_rate"`
		DurationMonths int     `json:"duration_months"`
	}
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := ctx.GetString("user_id") // Assume we get user ID from session/token
	message, err := c.ApplyForLoanUseCase.Execute(userID, input.Amount, input.InterestRate, input.DurationMonths)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (c *LoanController) ViewLoanStatus(ctx *gin.Context) {
	loanID := ctx.Param("id")
	userID := ctx.GetString("user_id")

	loan, err := c.ViewLoanStatusUseCase.Execute(loanID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loan)
}

func (c *LoanController) ViewAllLoans(ctx *gin.Context) {
	status := ctx.Query("status")
	order := ctx.Query("order")

	loans, err := c.ViewAllLoansAdminUseCase.Execute(status, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loans)
}

func (c *LoanController) UpdateLoanStatus(ctx *gin.Context) {
	loanID := ctx.Param("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	message, err := c.UpdateLoanStatusUseCase.Execute(loanID, input.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	loanID := ctx.Param("id")

	message, err := c.DeleteLoanUseCase.Execute(loanID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
