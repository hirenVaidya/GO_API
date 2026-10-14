package internal
package llm

import "context"

type Message struct{ Role, Content string }

type Provider interface {
    // Streams tokens to out; returns usage when done.
    Stream(ctx context.Context, system string, msgs []Message, out chan<- string) (Usage, error)
    Embed(ctx context.Context, texts []string) ([][]float32, error)
}

type Usage struct{ InputTokens, OutputTokens int }