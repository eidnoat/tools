package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	apiURL = "https://api.cognitive.microsofttranslator.com/translate"
	region = "eastasia"
)

type translateRequest struct {
	Text string `json:"text"`
}

type detectedLanguage struct {
	Language string `json:"language"`
	Score    float64 `json:"score"`
}

type translation struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type translateResponse struct {
	DetectedLanguage detectedLanguage `json:"detectedLanguage"`
	Translations     []translation    `json:"translations"`
}

func translateText(apiKey, content string) (string, error) {
	body := []translateRequest{{Text: content}}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// Query parameters
	q := req.URL.Query()
	q.Add("api-version", "3.0")
	q.Add("to", "zh-Hans")
	q.Add("to", "en")
	req.URL.RawQuery = q.Encode()

	// Headers
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	req.Header.Set("Ocp-Apim-Subscription-Region", region)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ClientTraceId", uuid.New().String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no translation result")
	}

	respData := result[0]
	targetLang := "zh-Hans"
	if respData.DetectedLanguage.Language == "zh-Hans" {
		targetLang = "en"
	}

	for _, trans := range respData.Translations {
		if trans.To == targetLang {
			return trans.Text, nil
		}
	}

	return "", fmt.Errorf("target language translation not found")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: translate <API_KEY> <TEXT>")
		os.Exit(1)
	}

	apiKey := os.Args[1]
	text := os.Args[2]

	result, err := translateText(apiKey, text)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
