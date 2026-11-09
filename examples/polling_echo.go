package main

import (
	"log"
	"time"

	api "github.com/OvyFlash/telegram-bot-api"
)

func polling_echo() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60

	// Empty list allows all updates excluding
	// UpdateTypeChatMember, UpdateTypeMessageReaction and UpdateTypeMessageReactionCount
	updateConfig.AllowedUpdates = []string{
		api.UpdateTypeMessage,
		api.UpdateTypeEditedMessage,
		api.UpdateTypeChannelPost,
		api.UpdateTypeEditedChannelPost,
		api.UpdateTypeBusinessConnection,
		api.UpdateTypeBusinessMessage,
		api.UpdateTypeEditedBusinessMessage,
		api.UpdateTypeDeletedBusinessMessages,
		api.UpdateTypeMessageReaction,
		api.UpdateTypeMessageReactionCount,
		api.UpdateTypeInlineQuery,
		api.UpdateTypeChosenInlineResult,
		api.UpdateTypeCallbackQuery,
		api.UpdateTypeShippingQuery,
		api.UpdateTypePreCheckoutQuery,
		api.UpdateTypePurchasedPaidMedia,
		api.UpdateTypePoll,
		api.UpdateTypePollAnswer,
		api.UpdateTypeMyChatMember,
		api.UpdateTypeChatMember,
		api.UpdateTypeChatJoinRequest,
		api.UpdateTypeChatBoost,
		api.UpdateTypeRemovedChatBoost,
	}

	updatesChannel := bot.GetUpdatesChan(updateConfig)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updatesChannel.Clear()

	for update := range updatesChannel {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Send an echo reply to the message
		msg := api.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyParameters.MessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
