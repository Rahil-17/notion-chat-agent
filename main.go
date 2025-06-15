package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// loadEnv loads environment variables from .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

// fetchNotionContent retrieves content from a Notion page using the Notion API
// Returns the page content as a string and any error that occurred
func fetchNotionContent(pageID, token string) (string, error) {
	// Construct the Notion API URL for fetching page blocks
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", pageID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the JSON response into a structured format
	var result struct {
		Results []map[string]interface{} `json:"results"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	// Extract text content from each block
	var lines []string
	for _, blk := range result.Results {
		// Every block has a "type" and an object keyed by that type
		blkType, ok := blk["type"].(string)
		if !ok {
			continue
		}

		if blkObj, ok := blk[blkType].(map[string]interface{}); ok {
			// Most content blocks expose "rich_text"
			if rtArr, ok := blkObj["rich_text"].([]interface{}); ok {
				for _, rtRaw := range rtArr {
					if rt, ok := rtRaw.(map[string]interface{}); ok {
						if plain, ok := rt["plain_text"].(string); ok {
							lines = append(lines, plain)
						}
					}
				}
			}
		}
	}

	return strings.Join(lines, "\n"), nil
}

// askOpenAI sends a question to OpenAI's API along with the Notion context
// Returns the AI's response and any error that occurred
func askOpenAI(question, context, apiKey string) (string, error) {
	// Construct the prompt with the Notion content and user's question
	prompt := fmt.Sprintf(`
You are a helpful assistant. Answer the user's question based on the following Notion content.

--- Notion Content ---
%s
-----------------------

Question: %s
Answer:
`, context, question)

	// Prepare the request payload for OpenAI API
	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	// Make the API request to OpenAI
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the OpenAI response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		return strings.TrimSpace(result.Choices[0].Message.Content), nil
	}
	return "", fmt.Errorf("no response from OpenAI")
}

func main() {
	// Load environment variables
	loadEnv()
	openaiKey := os.Getenv("OPENAI_API_KEY")
	notionKey := os.Getenv("NOTION_API_KEY")
	pageID := os.Getenv("NOTION_PAGE_ID")

	// Fetch content from Notion
	fmt.Println("📄 Fetching Notion content...")
	context, err := fetchNotionContent(pageID, notionKey)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}
	if len(context) == 0 {
		fmt.Println("❌ No content fetched.")
		return
	}
	fmt.Println("✅ Content fetched.")

	// Get user's question
	fmt.Print("\n❓ Enter your question: ")
	reader := bufio.NewReader(os.Stdin)
	question, _ := reader.ReadString('\n')
	question = strings.TrimSpace(question)

	// Get AI's response
	answer, err := askOpenAI(question, context, openaiKey)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("\n🤖 Answer:\n" + answer)
}
