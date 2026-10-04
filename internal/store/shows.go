package store

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// dedupe trims labels and removes empties and duplicates, keeping order.
func dedupe(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// CreateShow inserts the show and all its seats in ONE transaction,
// so a show never exists half-created.
func (s *Store) CreateShow(ctx context.Context, name string, seats []string, price int64, limit int) (*Show, error) {
	seats = dedupe(seats)
	id := uuid.NewString()

	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`INSERT INTO shows(id, name, price_paise, per_user_limit) VALUES($1,$2,$3,$4)`,
		id, name, price, limit); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO seats(show_id, label) SELECT $1, unnest($2::text[])`,
		id, seats); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetShow(ctx, id)

}

func (s *Store) GetShow(ctx context.Context, id string) (*Show, error) {
	var sh Show
	err := s.pool.QueryRow(ctx, `SELECT id::text, name, price_paise, per_user_limit FROM shows WHERE id=$1`, id).Scan(&sh.ID, &sh.Name, &sh.PricePaise, &sh.PerUserLimit)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `SELECT label, status FROM seats WHERE show_id=$1 ORDER BY label`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sh.Seats = []SeatInfo{}
	for rows.Next() {
		var label, status string
		if err := rows.Scan(&label, &status); err != nil {
			return nil, err
		}
		sh.Seats = append(sh.Seats, SeatInfo{Label: label, Status: status})
		switch status {
		case "available":
			sh.Available++
		case "held":
			sh.Held++
		case "confirmed":
			sh.Confirmed++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sh.TotalSeats = len(sh.Seats)
	return &sh, nil
}
