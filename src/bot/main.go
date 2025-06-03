package main

import (
	"flag"
	"log"
	"os"

	"tgbot-numerologist/ai"
	"tgbot-numerologist/communicate"
	"tgbot-numerologist/database"
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logFilePath := flag.String("logfile", "/logs/logs.log", "Path to the log file")
	flag.Parse()
	err := utils.InitLogger(*logFilePath)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer utils.CloseLogger()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST environment variable is not set")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("REDIS_PORT environment variable is not set")
	}
	proxyURL := os.Getenv("PROXY_URL")
	if proxyURL == "" {
		log.Fatal("PROXY_URL environment variable is not set")
	}
	gptKey := os.Getenv("CHATGPT_KEY")
	if gptKey == "" {
		log.Fatal("CHATGPT_KEY environment variable is not set")
	}

	database.InitRDB(redisHost, redisPort)
	err = ai.Init(gptKey, proxyURL)
	if err != nil {
		log.Fatal("Couldn't init ai client: %w", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = os.Getenv("DEBUG") == "true"
	bot.Self.CanJoinGroups = false

	utils.Log("Authorized on account %s", bot.Self.UserName)

	communicate.StartReceivingUpdates(bot)
}
