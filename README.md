# groq-go
The Groq Go Client is a Go library that provides easy access to the [Groq](https://groq.com/) API for creating chat completions and managing models. This client simplifies interactions with the API, including handling streaming responses.

## Installation
```bash
go get -u github.com/magicx-ai/groq-go/groq
```

## Usage
### Creating Chat Completions
```go
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
    fmt.Println(fmt.Errorf("error occurred: %v", err))
    return
}

fmt.Println(resp.Choices[0].Message.Content)
```

### Listing Models
```go
cli := groq.NewClient(apiKey, &http.Client{})

modelsResponse, err := cli.ListModels()
if err != nil {
    fmt.Println(fmt.Errorf("error occurred: %v", err))
    return
}

for _, model := range modelsResponse.Data {
    fmt.Printf("Model ID: %s, Owned By: %s, Active: %t, Context Window: %d\n", model.ID, model.OwnedBy, model.Active, model.ContextWindow)
}
```

### Retrieving a Model
```go
cli := groq.NewClient(apiKey, &http.Client{})

modelID := groq.ModelIDMIXTRAL
modelResponse, err := cli.RetrieveModel(modelID)
if err != nil {
    fmt.Println(fmt.Errorf("error occurred: %v", err))
    return
}

fmt.Printf("Model ID: %s, Owned By: %s, Active: %t, Context Window: %d\n", modelResponse.ID, modelResponse.OwnedBy, modelResponse.Active, modelResponse.ContextWindow)
```

### Streaming Chat Completions
```go
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
    fmt.Println(fmt.Errorf("error occurred: %v", err))
    return
}

for res := range respCh {
    if res.Error != nil {
        fmt.Println(fmt.Errorf("error occurred: %v", res.Error))
        break
    }
    fmt.Printf("Response: %+v\n", res.Response.Choices[0].Delta)
}
```

## Testing
Mock groq.Client
```bash
mockgen github.com/magicx-ai/groq-go/groq Client
```

You can use go generate
```go
//go:generate mockgen -destination {as_you_want} github.com/magicx-ai/groq-go/groq Client
```
