package billing

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"yourmodule/internal/llm"
)

// Monthly message limits per plan
var planLimits = map[string]int{
	"free":     50,
	"starter":  1000,
	"pro":      10000,
	"business": 50000,
}

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// WithinQuota checks if the tenant still has messages left this month.
func (s *Service) WithinQuota(ctx context.Context, tenantID uuid.UUID, plan string) (bool, error) {
	limit, ok := planLimits[plan]
	if !ok {
		limit = planLimits["free"]
	}

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var used int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM usage_events
		WHERE tenant_id = $1
		  AND kind = 'chat_message'
		  AND created_at >= $2`, tenantID, monthStart).Scan(&used)
	if err != nil {
		return false, err
	}
	return used < limit, nil
}

// Record saves one usage event after a chat finishes.
func (s *Service) Record(ctx context.Context, tenantID uuid.UUID, kind string, u llm.Usage) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO usage_events (tenant_id, kind, input_tokens, output_tokens)
		VALUES ($1, $2, $3, $4)`,
		tenantID, kind, u.InputTokens, u.OutputTokens)
	return err
}
