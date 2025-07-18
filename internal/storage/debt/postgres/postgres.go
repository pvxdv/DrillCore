package postgres

import (
	"context"
	"database/sql"
	"drillCore/internal/config"
	"drillCore/internal/model"
	debtStorage "drillCore/internal/storage/debt"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type DebtStorage struct {
	db     *pgx.Conn
	logger *zap.SugaredLogger
}

func New(ctx context.Context, cfg *config.DbEnvs, logger *zap.SugaredLogger) (*DebtStorage, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		defer func() { _ = conn.Close(ctx) }()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DebtStorage{db: conn, logger: logger}, nil
}

func (s *DebtStorage) Close(ctx context.Context) error {
	if err := s.db.Close(ctx); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

func (s *DebtStorage) Save(ctx context.Context, debt *model.Debt) (int64, error) {
	q := `INSERT INTO debt (user_id, description, amount, return_date)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`

	var debtID int64
	err := s.db.QueryRow(ctx, q,
		debt.UserID,
		debt.Description,
		debt.Amount,
		debt.ReturnDate,
	).Scan(&debtID)
	if err != nil {
		return -1, fmt.Errorf("failed to insert debt: %w", err)
	}

	s.logger.Debugf("successfully added debt (ID: %d) for user %d", debtID, debt.UserID)
	return debtID, nil
}

func (s *DebtStorage) Get(ctx context.Context, id int64) (*model.Debt, error) {
	q := `SELECT id, user_id, description, amount, return_date FROM debt WHERE id = $1`

	var d model.Debt
	var date sql.NullTime
	err := s.db.QueryRow(ctx, q, id).Scan(&d.ID, &d.UserID, &d.Description, &d.Amount, &date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, debtStorage.ErrDebtNotFound
		}
		return nil, fmt.Errorf("failed to get debt: %w", err)
	}

	if date.Valid {
		d.ReturnDate = &date.Time
	}
	return &d, nil
}

func (s *DebtStorage) Debts(ctx context.Context, userID int64) ([]*model.Debt, error) {
	q := `SELECT id, user_id, description, amount, return_date
		 FROM debt WHERE user_id = $1`

	row, err := s.db.Query(ctx, q, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to find debt: %w", debtStorage.ErrDebtNotFound)
		}
		return nil, fmt.Errorf("failed to get debt: %w", err)
	}

	res := make([]*model.Debt, 0)
	for row.Next() {
		var debt model.Debt
		var date sql.NullTime

		err = row.Scan(
			&debt.ID,
			&debt.UserID,
			&debt.Description,
			&debt.Amount,
			&date,
		)
		if err != nil {
			s.logger.Warnw("failed to scan debt", zap.Error(err))
		}

		if date.Valid {
			debt.ReturnDate = &date.Time
		}

		res = append(res, &debt)
	}

	return res, nil
}

func (s *DebtStorage) Update(ctx context.Context, debt *model.Debt) error {
	q := `UPDATE debt 
		 SET user_id = $1, 
		     description = $2, 
		     amount = $3, 
		     return_date = $4 
		 WHERE id = $5`

	result, err := s.db.Exec(ctx, q,
		debt.UserID,
		debt.Description,
		debt.Amount,
		debt.ReturnDate,
		debt.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update debt: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("failed to update debt: %w", debtStorage.ErrDebtNotFound)
	}

	s.logger.Debugf("successfully updated debt (ID: %d)", debt.ID)
	return nil
}

func (s *DebtStorage) Delete(ctx context.Context, id int64) error {
	q := "DELETE FROM debt WHERE id = $1"

	result, err := s.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete debt: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("failed to update debt: %w", debtStorage.ErrDebtNotFound)
	}

	s.logger.Debugf("successfully deleted debt (ID: %d)", id)
	return nil
}
