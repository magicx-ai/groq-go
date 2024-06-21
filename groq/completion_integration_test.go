//go:build integration

package groq

import (
	"context"
	"fmt"
	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

var apiKey = func() string {
	return os.Getenv("GROQ_API_KEY")
}()

func TestCreateChatCompletion(t *testing.T) {
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set")
	}

	httpCli := &http.Client{
		Timeout: 3 * time.Second,
	}

	c := NewClient(apiKey, httpCli)

	completion, err := c.CreateChatCompletion(ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    MessageRoleSystem,
				Content: "You are a developer.",
			},
			{
				Role:    MessageRoleUser,
				Content: "How do I write a function in Golang?",
			},
		},
		Model:      ModelIDLLAMA370B,
		MaxTokens:  1,
		NumChoices: 1,
	})
	require.NoError(t, err, "failed to create chat completion")

	assert.EqualValues(t, ModelIDLLAMA370B, completion.Model)
	assert.Len(t, completion.Choices, 1)
	fmt.Println("Role:", completion.Choices[0].Message.Role)
	fmt.Println("Content:", completion.Choices[0].Message.Content)
}

func TestCreateChatCompletionStream(t *testing.T) {
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set")
	}

	defer leaktest.Check(t)()

	httpCli := &http.Client{
		Timeout: 3 * time.Second,
	}

	c := NewClient(apiKey, httpCli)
	ctx := context.Background()
	stream, closer, err := c.CreateChatCompletionStream(ctx, ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    MessageRoleSystem,
				Content: "You are a developer.",
			},
			{
				Role:    MessageRoleUser,
				Content: "How do I write a function in Golang?",
			},
		},
		Model:      ModelIDLLAMA370B,
		MaxTokens:  1,
		NumChoices: 1,
		Stream:     true,
	})
	if closer != nil {
		defer closer()
	}
	require.NoError(t, err, "failed to create chat completion")

	for r := range stream {
		require.NoError(t, r.Error, "failed to get response")
		assert.EqualValues(t, ModelIDLLAMA370B, r.Response.Model)
		assert.Len(t, r.Response.Choices, 1)
	}
	assert.Empty(t, stream, "stream should be empty")
}

func TestListModels(t *testing.T) {
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set")
	}

	httpCli := &http.Client{
		Timeout: 3 * time.Second,
	}

	c := NewClient(apiKey, httpCli)
	models, err := c.ListModels()
	require.NoError(t, err, "failed to list models")
	require.NotEmpty(t, models, "models list is empty")

	currentModels := map[ModelID]bool{
		ModelIDLLAMA38B:  false,
		ModelIDLLAMA370B: false,
		ModelIDMIXTRAL:   false,
		ModelIDGEMMA:     false,
	}

	for _, model := range models.Data {
		// check if the model is in the current models
		if _, ok := currentModels[model.ID]; !ok {
			fmt.Printf("model %s is not in the current models\n", model.ID)
			continue
		}
	}
}

func TestRetrieveModel(t *testing.T) {
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set")
	}

	testcases := []struct {
		modelID ModelID
	}{
		{
			modelID: ModelIDLLAMA38B,
		},
		{
			modelID: ModelIDLLAMA370B,
		},
		{
			modelID: ModelIDMIXTRAL,
		},
		{
			modelID: ModelIDGEMMA,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(string(tc.modelID), func(t *testing.T) {
			t.Parallel()

			httpCli := &http.Client{
				Timeout: 3 * time.Second,
			}

			c := NewClient(apiKey, httpCli)

			model, err := c.RetrieveModel(tc.modelID)
			require.NoError(t, err, "failed to retrieve model")

			assert.Equal(t, tc.modelID, model.ID)
			assert.True(t, model.Active, "model is in-active")
		})
	}
}
