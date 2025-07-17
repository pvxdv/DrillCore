package debtStorage

import (
	"context"
	"drillCore/internal/model"
	"errors"
)

type DebtStorage interface {
	Save(ctx context.Context, debt *model.Debt) (int64, error)
	Debts(ctx context.Context, userID int64) ([]*model.Debt, error)
	Update(ctx context.Context, debt model.Debt) error
	Delete(ctx context.Context, id int64) error
}

var (
	ErrDebtNotFound = errors.New("debt not found")
)
