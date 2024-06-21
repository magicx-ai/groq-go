package groq

import (
	"bytes"
	"context"
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

func (c *client) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan *ChatCompletionStreamResponse, func(), error) {
	if !req.Stream {
		return nil, nil, fmt.Errorf("stream must be set to true")
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)
	httpReq, err := http.NewRequestWithContext(ctxWithCancel, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		cancel()

		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	responseCh := make(chan *ChatCompletionStreamResponse)

	cli := &sse.Client{
		HTTPClient:        c.client,
		ResponseValidator: sse.DefaultValidator,
		Backoff:           sse.DefaultClient.Backoff,
	}

	conn := cli.NewConnection(httpReq)
	go func() {
		defer close(responseCh)

		err := conn.Connect()
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, context.Canceled) {
			responseCh <- &ChatCompletionStreamResponse{
				Error: errors.Wrap(err, "failed to connect to the server"),
			}
		}
	}()

	remover := conn.SubscribeToAll(func(event sse.Event) {
		if strings.Contains(event.Data, "DONE") {
			cancel()
			return
		}

		var chatResp ChatCompletionResponse
		if err := json.Unmarshal([]byte(event.Data), &chatResp); err != nil {
			responseCh <- &ChatCompletionStreamResponse{Error: errors.Wrap(err, "failed to unmarshal response")}
			return
		}

		responseCh <- &ChatCompletionStreamResponse{Response: chatResp}
	})

	return responseCh, func() {
		cancel()
		remover()
	}, nil
}
