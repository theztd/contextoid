package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"theztd/contextoid/pkg/gpt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Načtení .env souboru
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables.")
	}

	// Ověření, zda je dostupný OpenAI API klíč
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set. Please set it in .env or as an environment variable.")
	}

	logLevel := os.Getenv("LOG_LEVEL")

	chatGPT := gpt.New("https://api.openai.com/", apiKey, logLevel)
	// Inicializace Gin routeru
	r := gin.Default()

	r.POST("/v1/gpt/analyze_comments", func(c *gin.Context) {
		var req gpt.AnalyzeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Zavolání analýzy komentářů pomocí GPT
		analysis, err := chatGPT.AnalyzeCommentsWithPrompt(req.TaskText, req.Comments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Odeslání odpovědi zpět klientovi
		c.JSON(http.StatusOK, analysis)
	})

	fmt.Println("API server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
