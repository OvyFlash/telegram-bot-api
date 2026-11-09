package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	api "github.com/OvyFlash/telegram-bot-api"
)

const exampleChatID int64 = 123456

func main() {
	examples := map[string]func(){
		"polling_echo":                 polling_echo,
		"polling_inline_query":         polling_inline_query,
		"polling_guest_bot":            polling_guest_bot,
		"webhook_echo":                 webhook_echo,
		"webhook_custom_handler":       webhook_custom_handler,
		"send_live_photo":              send_live_photo,
		"send_poll_with_media":         send_poll_with_media,
		"remove_message_reactions":     remove_message_reactions,
		"send_structured_rich_message": func() { send_structured_rich_message(newBot(), exampleChatID) },
		"configure_ephemeral_command":  func() { configure_ephemeral_command(newBot()) },
		"handle_bot_api_10_2_update":   pollBotAPI102Updates,
	}

	var name string
	if len(os.Args) == 2 {
		name = os.Args[1]
	}

	run, ok := examples[name]
	if !ok {
		names := strings.Join(slices.Sorted(maps.Keys(examples)), ", ")
		fmt.Fprintf(os.Stderr, "usage: go run ./examples <example>\nexamples: %s\n", names)
		os.Exit(2)
	}

	run()
}

func newBot() *api.BotAPI {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	return bot
}

func pollBotAPI102Updates() {
	bot := newBot()

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60
	updateConfig.AllowedUpdates = []string{api.UpdateTypeMessage, api.UpdateTypeSubscription}

	for update := range bot.GetUpdatesChan(updateConfig) {
		handle_bot_api_10_2_update(bot, update)
	}
}
