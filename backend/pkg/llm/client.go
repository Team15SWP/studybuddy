package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"study_buddy/internal/config"
)

const dailyQuota = 50

type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type OpenRouterClient struct {
	APIKey     string
	OpenAi     *config.OpenAI
	DisabledAt []*time.Time
	Count      []int
}

func NewOpenRouterClient(openAi *config.OpenAI) *OpenRouterClient {
	return &OpenRouterClient{
		APIKey:     openAi.ApiKeys[openAi.Ind],
		OpenAi:     openAi,
		DisabledAt: make([]*time.Time, len(openAi.ApiKeys), len(openAi.ApiKeys)),
		Count:      make([]int, len(openAi.ApiKeys), len(openAi.ApiKeys)),
	}
}

func (c *OpenRouterClient) Complete(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": "mistralai/mistral-7b-instruct:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyJSON, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	c.APIKey, err = c.NextValidKey()
	if err != nil {
		return "", fmt.Errorf("c.NextValidKey: %w", err)
	}

	c.Count[c.OpenAi.Ind]++
	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if c.Count[c.OpenAi.Ind] >= dailyQuota {
		c.InvalidKey()
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm error: %s", string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("llm response had no choices")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *OpenRouterClient) NextValidKey() (string, error) {
	for i := 0; i < len(c.DisabledAt); i++ {
		c.OpenAi.Ind = (c.OpenAi.Ind + 1) % int32(len(c.OpenAi.ApiKeys))
		if c.DisabledAt[c.OpenAi.Ind] == nil || time.Since(*c.DisabledAt[c.OpenAi.Ind]) > 24*time.Hour {
			c.DisabledAt[c.OpenAi.Ind] = nil
			return c.OpenAi.ApiKeys[c.OpenAi.Ind], nil
		}
	}
	return "", fmt.Errorf("no valid key")
}

func (c *OpenRouterClient) InvalidKey() {
	tt := time.Now()
	c.DisabledAt[c.OpenAi.Ind] = &tt
}
