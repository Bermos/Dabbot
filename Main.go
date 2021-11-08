package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
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

	file := &TBot.Video{File: TBot.FromDisk("./dabs/" + filename + ".mp4")}
	bot.Send(recipient, file)
}

func main() {
	tokenEnv := os.Getenv("TOKEN")
	if tokenEnv == "" {
		log.Fatal("ERROR - No TOKEN environment var set")
	}
	bot := initialize(tokenEnv)

	// Handle dabs
	bot.Handle("/dab", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "dab")
	})

	bot.Handle("/rev_dab", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "rev_dab")
	})

	bot.Handle("/space_dab", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "space_dab")
	})

	bot.Handle("/rev_space_dab", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "rev_space_dab")
	})

	bot.Handle("/ht", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "ht")
	})

	bot.Handle("/ella", func(m *TBot.Message) {
		sendDab(bot, m.Chat, "ella")
	})

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
