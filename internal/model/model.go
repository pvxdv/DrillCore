package model

import "time"

// Debt
// @Description Represents a debt position, such as a credit or loan to friends/family.
// @Description Contains basic information about the debt obligation.
type Debt struct {
	ID          int64      `json:"id,omitempty" example:"1"`
	UserID      int64      `json:"user_id" example:"1"`
	Description string     `json:"description" example:"Loan for car purchase"`
	Amount      int64      `json:"amount" example:"1000000"`
	ReturnDate  *time.Time `json:"return_date,omitempty" example:"2025-01-02T15:04:05Z"`
}
