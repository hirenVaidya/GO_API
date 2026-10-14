func RequireQuota(b *billing.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := auth.TenantFrom(r.Context())
			ok, err := b.WithinQuota(r.Context(), t.ID, t.Plan)
			if err != nil || !ok {
				http.Error(w, "plan limit reached, please upgrade", http.StatusPaymentRequired)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}