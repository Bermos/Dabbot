package main

import (
	"fmt"
	"io/ioutil"
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

func sendDab(bot *TBot.Bot, recipient *TBot.Chat, filename string) {
	log.Printf("[Dab] %-13s requested by %s", filename, recipient.Username)

	file := &TBot.Video{File: TBot.FromDisk("./dabs/" + filename)}
	bot.Send(recipient, file)
}

func loadDabs(bot *TBot.Bot) {
	files, err := ioutil.ReadDir("./dabs/")

	if err != nil {
		log.Fatal("ERROR - Could not read dab dir.")
	}

	for i, file := range files {
		log.Printf("DEBUG - Loading dab %d: %s", i, file.Name())
		fileExt := filepath.Ext(file.Name())
		filename := strings.TrimSuffix(file.Name(), fileExt)

		if fileExt != ".mp4" {
			log.Printf("ERROR - dab not loaded, file extension is not '.mp4' but '%s'", fileExt)
			continue
		}

		bot.Handle(fmt.Sprintf("/%s", filename), func(m *TBot.Message) {
			sendDab(bot, m.Chat, file.Name())
		})
		log.Printf("DEBUG - Dab '%s' loaded and registered", filename)
	}
}

func main() {
	tokenEnv := os.Getenv("TOKEN")
	if tokenEnv == "" {
		log.Fatal("ERROR - No TOKEN environment var set")
	}
	bot := initialize(tokenEnv)

	// Load and register dabs to handle
	loadDabs(bot)

	// Handle poster requests
	bot.Handle("/poster", func(m *TBot.Message) {
		args := strings.Split(m.Payload, ".")

		if len(args) != 2 {
			bot.Send(m.Chat, fmt.Sprintf("Wrong number of arguments. Required: 2. Found: %v", len(args)))
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
