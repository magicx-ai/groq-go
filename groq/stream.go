package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/tmaxmax/go-sse"
)

type ChatCompletionStreamResponse struct {
	Response ChatCompletionResponse
	Error    error
}

func (c *client) CreateChatCompletionStream(req ChatCompletionRequest) (<-chan *ChatCompletionStreamResponse, error) {
	if !req.Stream {
		return nil, fmt.Errorf("stream must be set to true")
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	responseCh := make(chan *ChatCompletionStreamResponse)

	conn := sse.DefaultClient.NewConnection(httpReq)

	var remover sse.EventCallbackRemover
	remover = conn.SubscribeToAll(func(event sse.Event) {
		// Stream is terminated when the server sends "[DONE]"
		if strings.Contains(event.Data, "DONE") {
			// Close the response channel
			close(responseCh)
			// Remove the event subscriber itself
			remover()
		}

		var chatResp ChatCompletionResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResp); err != nil {
			responseCh <- &ChatCompletionStreamResponse{Error: errors.Wrap(err, "failed to unmarshal response")}
			return
		}

		responseCh <- &ChatCompletionStreamResponse{Response: chatResp}
	})

	go func() {
		err := conn.Connect()
		if err != nil && !errors.Is(err, io.EOF) {
			responseCh <- &ChatCompletionStreamResponse{
				Error: errors.Wrap(err, "failed to connect to the server"),
			}
		}
	}()

	return responseCh, nil
}
