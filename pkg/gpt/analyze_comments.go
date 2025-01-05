// Structy pro JSON vstup a výstup
package gpt

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type GPT struct {
	GPTKey   string
	GPTUrl   string
	LogLevel string
}

type AnalyzeRequest struct {
	TaskText string   `json:"task_text"`
	Comments []string `json:"comments"`
}

type Response struct {
	Data         map[string]interface{}
	Usage        map[string]interface{}
	FinishReason string `json:"finish_reason"`
}

func New(gptUrl, gptKey, logLevel string) *GPT {
	return &GPT{
		GPTKey:   gptKey,
		GPTUrl:   gptUrl,
		LogLevel: logLevel,
	}
}

// Funkce pro analýzu komentářů pomocí OpenAI GPT
func (gpt GPT) AnalyzeCommentsWithPrompt(taskText string, comments []string) (response Response, err error) {
	client := resty.New()

	// Sestavení promptu
	prompt := "You are an assistant tasked with analyzing whether comments modify the scope of a given task. " +
		"Here is the task description: \"" + taskText + "\"\n" +
		"For each comment, determine:\n" +
		"- If it changes the task scope (yes/no).\n" +
		"- If yes, specify what should be updated in the new_description field.\n\n"

	for i, comment := range comments {
		prompt += fmt.Sprintf("Comment %d: \"%s\"\n", i+1, comment)
	}
	prompt += "\nOutput must strictly be a valid JSON object. Do not include any explanations or additional text outside of the JSON format. Provide the output as a JSON array in this valid json format:\n" +
		"{ comments_analysis: [ {comment_id: int(), changes_scope: bool() }, ...], new_description: \"none or text of the new suggested description.\"}"

	// Odeslání promptu do GPT-4
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+gpt.GPTKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-4",
			"messages": []map[string]string{
				{"role": "system", "content": "You are an assistant tasked with analyzing task comments."},
				{"role": "user", "content": prompt},
			},
			"max_tokens": 500,
		}).
		Post(gpt.GPTUrl + "v1/chat/completions")

	if err != nil {
		if gpt.LogLevel == "debug" {
			log.Println("DEBUG [gpt.AnalyzeCommentsWithPrompt]: failed to contact OpenAI API:" + err.Error())
		}
		return response, fmt.Errorf("failed to contact OpenAI API: %w", err)
	}

	// Debugging: Výpis odpovědi
	if gpt.LogLevel == "debug" {
		log.Printf("DEBUG [gpt.AnalyzeCommentsWithPrompt:Response]: %s\n", resp.String())
	}

	// Zpracování odpovědi
	var result struct {
		Model   string                 `json:"model"`
		Usage   map[string]interface{} `json:"usage"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		if gpt.LogLevel == "debug" {
			log.Println("DEBUG [gpt.AnalyzeCommentsWithPrompt]: failed to parse OpenAI response:" + err.Error())
		}
		return response, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(result.Choices) == 0 {
		if gpt.LogLevel == "debug" {
			log.Printf("DEBUG [gpt.AnalyzeCommentsWithPrompt:Response]: no response from OpenAI API\n")
		}
		return response, fmt.Errorf("no response from OpenAI API")
	}

	// Převod odpovědi z GPT na JSON
	//var analysis []map[string]interface{}
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &response.Data); err != nil {
		if gpt.LogLevel == "debug" {
			log.Println("DEBUG [gpt.AnalyzeCommentsWithPrompt]: failed to parse GPT output as JSON: " + err.Error())
		}
		return response, fmt.Errorf("failed to parse GPT output as JSON: %w", err)
	}
	response.Usage = result.Usage
	response.FinishReason = result.Choices[0].FinishReason

	if gpt.LogLevel == "debug" {
		log.Printf("DEBUG [gpt.AnalyzeCommentsWithPrompt:ReturnedData]: %T\n", response)
	}

	return response, nil
}
