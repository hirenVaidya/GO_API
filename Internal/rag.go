func (r *Retriever) TopK(ctx context.Context, tenantID uuid.UUID, queryVec []float32, k int) ([]string, error) {
	rows, err := r.db.Query(ctx, `
        SELECT content FROM chunks
        WHERE tenant_id = $1
        ORDER BY embedding <=> $2
        LIMIT $3`, tenantID, pgvector.NewVector(queryVec), k)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}