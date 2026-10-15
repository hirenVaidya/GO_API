package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func hashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// CreateAPIKey returns the full key ONCE. Only the hash is stored.
func (s *Service) CreateAPIKey(ctx context.Context, tenantID uuid.UUID) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := "dc_" + hex.EncodeToString(b)

	_, err := s.db.Exec(ctx, `
		INSERT INTO api_keys (tenant_id, key_hash, prefix)
		VALUES ($1, $2, $3)`, tenantID, hashKey(key), key[:8])
	if err != nil {
		return "", err
	}
	return key, nil
}

// Middleware authenticates "Authorization: Bearer dc_xxx" and loads the tenant.
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		key := strings.TrimPrefix(h, "Bearer ")
		if key == "" || key == h {
			http.Error(w, "missing API key", http.StatusUnauthorized)
			return
		}

		var t Tenant
		err := s.db.QueryRow(r.Context(), `
			SELECT t.id, t.plan
			FROM api_keys k
			JOIN tenants t ON t.id = k.tenant_id
			WHERE k.key_hash = $1 AND k.revoked_at IS NULL`,
			hashKey(key)).Scan(&t.ID, &t.Plan)
		if err != nil {
			http.Error(w, "invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithTenant(r.Context(), t)))
	})
}
