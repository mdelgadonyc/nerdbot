package delphi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func Delphi(question string) string {

	apiKey := goDotEnvVariable("CHATKEY")
	apiEndpoint := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	// Stripping out the trailing newline is necessary for the JSON payload
	// question = strings.TrimRight(question, "\n")

	payloadString := fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "%s"}],
		"temperature": 0.7
		}`, question)

	payload := strings.NewReader(payloadString)

	client := &http.Client{}
	req, err := http.NewRequest(method, apiEndpoint, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var response OpenAIResponse

	// Unmarshal the JSON response into the OpenAIResponse struct
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Access the parased data
	if len(response.Choices) > 0 {
		responseText := response.Choices[0].Message
		fmt.Println("API Response: ", responseText)
		return responseText.Content
	} else {
		fmt.Println("No choices in the API response.")
	}
	return ""

	// Print the response
	//fmt.Println(string(body))
}
