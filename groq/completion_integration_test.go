//go:build integration

package groq

import (
	"fmt"
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
		Stream:     false,
	})
	require.NoError(t, err, "failed to create chat completion")

	assert.EqualValues(t, ModelIDLLAMA370B, completion.Model)
	assert.Len(t, completion.Choices, 1)
	fmt.Println("Role:", completion.Choices[0].Message.Role)
	fmt.Println("Content:", completion.Choices[0].Message.Content)
}
