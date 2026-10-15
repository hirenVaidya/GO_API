package api
authSvc := auth.NewService(dbPool)
billingSvc := billing.NewService(dbPool)

r := chi.NewRouter()
r.Route("/v1", func(r chi.Router) {
	r.Use(authSvc.Middleware)                  // 1. who is calling
	r.Use(billing.RequireQuota(billingSvc))    // 2. are they within plan limits
	r.Post("/chat", chatHandler.Stream)        // 3. the handler
})