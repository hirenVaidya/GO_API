func (h *ChatHandler) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	tenant := auth.TenantFrom(r.Context())
	var req struct {
		Question string `json:"question"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	qvec, _ := h.llm.Embed(r.Context(), []string{req.Question})
	ctxChunks, _ := h.retriever.TopK(r.Context(), tenant.ID, qvec[0], 5)

	system := "Answer ONLY using the context below. If unsure, say you don't know.\n\n" +
		strings.Join(ctxChunks, "\n---\n")

	out := make(chan string)
	var usage llm.Usage
	go func() {
		defer close(out)
		usage, _ = h.llm.Stream(r.Context(), system,
			[]llm.Message{{Role: "user", Content: req.Question}}, out)
	}()

	for tok := range out {
		fmt.Fprintf(w, "data: %s\n\n", tok)
		flusher.Flush()
	}
	h.billing.Record(r.Context(), tenant.ID, "chat_message", usage)
}