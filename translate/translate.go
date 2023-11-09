package translate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Define a struct that matches the structure of the JSON data
type Translation struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func Translate(sourceText string) string {

	apiKey := goDotEnvVariable("APIKEY")
	baseURL := "https://translation.googleapis.com/language/translate/v2"

	// sourceText := "Hello, world!"
	sourceText = url.QueryEscape(sourceText)
	targetLanguage := "zh-CN" // Mandarin

	// return "你好世界"

	// Construct the API URL with query parameters
	url := fmt.Sprintf("%s?q=%s&target=%s&key=%s", baseURL, sourceText, targetLanguage, apiKey)

	// Send a GET request to the API
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return ""
	}

	var translationResponse Translation
	err = json.Unmarshal([]byte(responseBody), &translationResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON: ", err)
	}

	// Access and print the value you're interested in (e.g., "translatedText")
	if len(translationResponse.Data.Translations) > 0 {
		translatedText := translationResponse.Data.Translations[0].TranslatedText
		// fmt.Println("Translation:", translatedText)
		return translatedText
	}

	return ""
}
