package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	TBot "gopkg.in/tucnak/telebot.v2"
)

func initialize(token string) *TBot.Bot {
	log.Println("INFO - Initializing bot...")

	bot, err := TBot.NewBot(TBot.Settings{
		Token:  token,
		Poller: &TBot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("INFO - Initialize successful")

	return bot
}

func sendVideo(bot *TBot.Bot, recipient *TBot.Chat, filename string) {
	log.Printf("[Video] %-13s requested by %s", filename, recipient.Username)

	file := &TBot.Video{File: TBot.FromDisk("./dabs/" + filename)}
	_, _ = bot.Send(recipient, file)
}

func sendPicture(bot *TBot.Bot, recipient *TBot.Chat, filename string) {
	log.Printf("[Picture] %-13s requested by %s", filename, recipient.Username)

	file := &TBot.Video{File: TBot.FromDisk("./dabs/" + filename)}
	_, _ = bot.Send(recipient, file)
}

func loadFiles(bot *TBot.Bot) {
	files, err := os.ReadDir("./dabs/")

	if err != nil {
		log.Fatal("ERROR - Could not read dab dir.")
	}

	for i, file := range files {
		log.Printf("DEBUG - Loading dab %d: %s", i, file.Name())
		fileExt := filepath.Ext(file.Name())
		filename := strings.TrimSuffix(file.Name(), fileExt)

		switch fileExt {
		case ".mp4":
			bot.Handle(fmt.Sprintf("/%s", filename), func(m *TBot.Message) {
				sendVideo(bot, m.Chat, fmt.Sprintf("%s%s", filename, fileExt))
			})
			log.Printf("DEBUG - Video '/%s' loaded and registered", filename)

		case ".jpg":
			bot.Handle(fmt.Sprintf("/%s", filename), func(m *TBot.Message) {
				sendPicture(bot, m.Chat, fmt.Sprintf("%s%s", filename, fileExt))
			})
			log.Printf("DEBUG - Picture '/%s' loaded and registered", filename)

		default:
			log.Printf("ERROR - file not loaded, file extension is not recognised but '%s'", fileExt)
		}

	}
}

func getToken() (string, error) {
	tokenEnv := os.Getenv("TOKEN")
	tokenFileEnv := os.Getenv("TOKEN_FILE")

	// TOKEN set, use that
	if tokenEnv != "" {
		if tokenFileEnv != "" {
			log.Print("WARNING - TOKEN and TOKEN_FILE env set, TOKEN will take precedence.")
		}
		return tokenEnv, nil
	}

	// TOKEN and TOKEN_FILE not set, no token -> crash
	if tokenFileEnv == "" {
		return "", errors.New("no TOKEN or TOKEN_FILE environment var set")
	}

	token, err := os.ReadFile(tokenFileEnv)
	if err != nil {
		return "", errors.New(fmt.Sprintf("TOKEN_FILE env - %v", err))
	}

	if string(token) == "" {
		return "", errors.New("token read from TOKEN_FILE file is empty")
	}

	return string(token), nil
}

func main() {
	token, err := getToken()
	if err != nil {
		log.Fatalf("ERROR - %v", err)
	}
	bot := initialize(token)

	// Load and register files to handle
	loadFiles(bot)

	// Handle poster requests
	bot.Handle("/poster", func(m *TBot.Message) {
		args := strings.Split(m.Payload, ".")

		if len(args) != 2 {
			_, _ = bot.Send(m.Chat, fmt.Sprintf("Wrong number of arguments. Required: 2. Found: %v", len(args)))
			return
		}
		log.Printf("[Poster] request received by %s with args \"%s\"", m.Sender.Username, args)

		args[0] = url.PathEscape(strings.TrimSpace(args[0]))
		args[1] = url.PathEscape(strings.TrimSpace(args[1]))
		escapedUrl := fmt.Sprintf("https://punkt.felunka.de/generate.php?text=%s&text2=%s&color=c", url.QueryEscape(args[0]), url.QueryEscape(args[1]))

		file := &TBot.Photo{File: TBot.FromURL(escapedUrl)}
		_, err := bot.Send(m.Chat, file)
		if err != nil {
			log.Println("ERROR - Poster could not be sent. See error below.")
			log.Println(err)
		}
	})

	// Register listener for term signal and gracefully shut down
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("INFO - Term signal received. Shutting down...")
		bot.Stop()
		os.Exit(0)
	}()

	log.Print("INFO - Starting bot")
	bot.Start()
}
