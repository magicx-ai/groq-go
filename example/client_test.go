package example

import (
	"fmt"
	"net/http"

	"github.com/magicx-ai/groq-go/groq"
)

const (
	apiKey = "your-api-key"
)

func ExampleClient_CreateChatCompletion() {
	cli := groq.NewClient(apiKey, &http.Client{})

	req := groq.ChatCompletionRequest{
		Messages: []groq.Message{
			{
				Role:    "user",
				Content: "Explain the importance of fast language models",
			},
		},
		Model:       groq.ModelIDLLAMA370B,
		MaxTokens:   150,
		Temperature: 0.7,
		TopP:        0.9,
		NumChoices:  1,
		Stream:      false,
	}

	resp, err := cli.CreateChatCompletion(req)
	if err != nil {
		fmt.Println(fmt.Errorf("error is occurred: %v", err))
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func ExampleClient_ListModels() {
	client := groq.NewClient(apiKey, &http.Client{})

	modelsResponse, err := client.ListModels()
	if err != nil {
		fmt.Println(fmt.Errorf("error is occurred: %v", err))
		return
	}

	if len(modelsResponse.Data) == 0 {
		fmt.Println(fmt.Errorf("expected at least one model, got none"))
		return
	}

	for _, model := range modelsResponse.Data {
		fmt.Printf("Model ID: %s, Owned By: %s, Active: %t, Context Window: %d\n", model.ID, model.OwnedBy, model.Active, model.ContextWindow)
	}
}

func ExampleClient_RetrieveModel() {
	client := groq.NewClient(apiKey, &http.Client{})

	// You can use the predefined model IDs
	modelID := groq.ModelIDMIXTRAL
	// Or you can find the model ID from the ListModels API
	// modelID = groq.ModelID("mixtral-8x7b-32768")
	modelResponse, err := client.RetrieveModel(modelID)
	if err != nil {
		fmt.Println(fmt.Errorf("error is occurred: %v", err))
		return
	}

	fmt.Printf("Model ID: %s, Owned By: %s, Active: %t, Context Window: %d\n", modelResponse.ID, modelResponse.OwnedBy, modelResponse.Active, modelResponse.ContextWindow)
}

func ExampleClient_CreateChatCompletionStream() {
	cli := groq.NewClient(apiKey, &http.Client{})

	req := groq.ChatCompletionRequest{
		Messages: []groq.Message{
			{
				Role:    "user",
				Content: "Explain the importance of fast language models",
			},
		},
		Model:       groq.ModelIDLLAMA370B,
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
		NumChoices:  1,
		Stream:      true,
	}

	respCh, err := cli.CreateChatCompletionStream(req)
	if err != nil {
		fmt.Println(fmt.Errorf("error is occurred: %v", err))
		return
	}

	for res := range respCh {
		if res.Error != nil {
			fmt.Println(fmt.Errorf("error is occurred: %v", res.Error))
			break
		}
		fmt.Printf("Response: %+v\n", res.Response.Choices[0].Delta)
	}
}
