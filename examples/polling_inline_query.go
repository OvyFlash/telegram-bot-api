package main

import (
	"log"

	api "github.com/OvyFlash/telegram-bot-api"
)

// func main() { polling_inline_query() }

func polling_inline_query() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken") // create new bot
	if err != nil {
		panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60

	updatesChannel := bot.GetUpdatesChan(updateConfig)

	for update := range updatesChannel {
		if update.InlineQuery == nil { // if no inline query, skip update
			continue
		}

		article := api.NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
		article.Description = update.InlineQuery.Query

		inlineConfig := api.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       []any{article},
		}

		if _, err := bot.Request(inlineConfig); err != nil {
			log.Println(err)
		}
	}
}
