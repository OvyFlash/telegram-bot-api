package main

import (
	"log"
	"net/http"
	"sync"

	api "github.com/OvyFlash/telegram-bot-api"
)

func webhook_custom_handler() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Listen for updates from the webhook and handle them yourself
	http.HandleFunc("/"+bot.Token, func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			log.Printf("%+v\n", err.Error())
		} else {
			log.Printf("%+v\n", *update)
		}

		if update.Message == nil {
			return
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Send an echo reply to the message
		msg := api.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyParameters.MessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	})

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil); err != nil {
			log.Printf("Error listening on TLS: %v", err)
		}
	}()

	// Set the webhook
	webHook, err := api.NewWebhookWithCert("https://your-bot-host-url/"+bot.Token, api.FilePath("cert.pem"))
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	apiResponse, err := bot.Request(webHook)
	if err != nil {
		panic(err)
	}

	if apiResponse.Ok {
		log.Printf("Webhook set successfully")
	} else {
		log.Printf("Failed to set webhook: %s", apiResponse.Description)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		panic(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Failed to get webhook info: %s", info.LastErrorMessage)
	}

	// Wait for the server to terminate
	waitGroup.Wait()
}
