package stress

import (
	"context"
	"github.com/magicx-ai/groq-go/groq"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func CompletionFunc(apiKey string) func() error {
	return func() error {
		httpCli := &http.Client{
			Timeout: 3 * time.Second,
		}

		c := groq.NewClient(apiKey, httpCli)

		completion, err := c.CreateChatCompletion(groq.ChatCompletionRequest{
			Messages: []groq.Message{
				{
					Role:    groq.MessageRoleSystem,
					Content: "You are a developer.",
				},
				{
					Role:    groq.MessageRoleUser,
					Content: "How do I write a function in Golang?",
				},
			},
			Model:      groq.ModelIDLLAMA370B,
			MaxTokens:  1,
			NumChoices: 1,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create chat completion")
		}

		if completion.Model != string(groq.ModelIDLLAMA370B) {
			return errors.New("model mismatch")
		}

		if len(completion.Choices) != 1 {
			return errors.New("unexpected number of choices")
		}

		return nil
	}
}

func CompletionStreamFunc(apiKey string) func() error {
	return func() error {
		httpCli := &http.Client{
			Timeout: 3 * time.Second,
		}

		c := groq.NewClient(apiKey, httpCli)

		ctx := context.Background()

		stream, closer, err := c.CreateChatCompletionStream(ctx, groq.ChatCompletionRequest{
			Messages: []groq.Message{
				{
					Role:    groq.MessageRoleSystem,
					Content: "You are a developer.",
				},
				{
					Role:    groq.MessageRoleUser,
					Content: "How do I write a function in Golang?",
				},
			},
			Model:      groq.ModelIDLLAMA370B,
			MaxTokens:  1,
			NumChoices: 1,
			Stream:     true,
		})
		if closer != nil {
			defer closer()
		}
		if err != nil {
			return errors.Wrap(err, "failed to create chat completion")
		}

		for r := range stream {
			if r.Error != nil {
				return errors.Wrap(r.Error, "failed to get response")
			}

			if r.Response.Model != string(groq.ModelIDLLAMA370B) {
				return errors.New("model mismatch")
			}

			if len(r.Response.Choices) != 1 {
				return errors.New("unexpected number of choices")
			}
		}

		if _, isOpen := <-stream; isOpen {
			return errors.New("stream should be closed")
		}

		return nil
	}
}

func ListModelsFunc(apiKey string) func() error {
	return func() error {
		httpCli := &http.Client{
			Timeout: 3 * time.Second,
		}

		c := groq.NewClient(apiKey, httpCli)
		models, err := c.ListModels()
		if err != nil {
			return errors.Wrap(err, "failed to list models")
		}
		if len(models.Data) == 0 {
			return errors.New("models list is empty")
		}

		return nil
	}
}

func RetrieveModelFunc(apiKey string) func() error {
	return func() error {
		httpCli := &http.Client{
			Timeout: 3 * time.Second,
		}

		c := groq.NewClient(apiKey, httpCli)

		model, err := c.RetrieveModel(groq.ModelIDLLAMA370B)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve model")
		}

		if model.ID != groq.ModelIDLLAMA370B {
			return errors.New("model mismatch")
		}

		return nil
	}
}
