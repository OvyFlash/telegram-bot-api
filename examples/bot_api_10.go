package main

import (
	"log"

	api "github.com/OvyFlash/telegram-bot-api"
)

func polling_guest_bot() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60
	updateConfig.AllowedUpdates = []string{api.UpdateTypeGuestMessage}

	updatesChannel := bot.GetUpdatesChan(updateConfig)
	for update := range updatesChannel {
		if update.GuestMessage == nil || update.GuestMessage.GuestQueryID == "" {
			continue
		}

		result := api.NewInlineQueryResultArticle("guest-reply", "Reply", "Hello from guest mode")
		if _, err := bot.AnswerGuestQuery(api.NewAnswerGuestQuery(update.GuestMessage.GuestQueryID, result)); err != nil {
			log.Println(err)
		}
	}
}

func send_live_photo() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	msg := api.NewLivePhoto(123456, api.FilePath("live-photo.mp4"), api.FilePath("live-photo.jpg"))
	msg.Caption = "Live photo"

	if _, err := bot.SendLivePhoto(msg); err != nil {
		log.Println(err)
	}
}

func send_poll_with_media() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	optionMedia := api.NewInputMediaSticker(api.FileID("sticker_file_id"))
	poll := api.NewPoll(
		123456,
		"Where should we meet?",
		api.NewPollOptionWithMedia("At the cafe", &optionMedia),
		api.NewPollOption("At the office"),
	)
	location := api.NewInputMediaLocation(40.7128, -74.0060)
	poll.Media = &location
	poll.MembersOnly = true
	poll.CountryCodes = []string{"US"}

	if _, err := bot.Send(poll); err != nil {
		log.Println(err)
	}
}

func remove_message_reactions() {
	bot, err := api.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		panic(err)
	}

	reaction := api.NewDeleteMessageReaction(123456, 10)
	reaction.UserID = 777000
	if _, err := bot.DeleteMessageReaction(reaction); err != nil {
		log.Println(err)
	}

	clearRecent := api.NewDeleteAllMessageReactions(123456)
	clearRecent.UserID = 777000
	if _, err := bot.DeleteAllMessageReactions(clearRecent); err != nil {
		log.Println(err)
	}
}

func send_structured_rich_message(bot *api.BotAPI, chatID int64) {
	photo := api.NewInputMediaPhoto(api.FilePath("diagram.jpg"))
	richMessage := api.NewInputRichMessageBlocks(
		api.InputRichBlockSectionHeading{
			Type: "heading",
			Text: "Release summary",
			Size: 2,
		},
		api.InputRichBlockParagraph{
			Type: "paragraph",
			Text: "The migration completed successfully.",
		},
		api.InputRichBlockPhoto{
			Type:  "photo",
			Photo: photo,
		},
	)

	if _, err := bot.SendRichMessage(api.NewSendRichMessage(chatID, richMessage)); err != nil {
		log.Println(err)
	}
}

func configure_ephemeral_command(bot *api.BotAPI) {
	command := api.BotCommand{
		Command:     "private_status",
		Description: "Show status privately",
		IsEphemeral: true,
	}
	if _, err := bot.Request(api.NewSetMyCommands(command)); err != nil {
		log.Println(err)
	}
}

func handle_bot_api_10_2_update(bot *api.BotAPI, update api.Update) {
	if update.Subscription != nil {
		log.Printf("subscription %s for user %d", update.Subscription.State, update.Subscription.User.ID)
	}
	if update.Message == nil {
		return
	}
	if update.Message.CommunityChatAdded != nil {
		community := update.Message.CommunityChatAdded.Community
		log.Printf("joined community %s (%d)", community.Name, community.ID)
	}
	if update.Message.CommunityChatRemoved != nil {
		log.Print("removed from community")
	}
	if update.Message.EphemeralMessageID == 0 || update.Message.From == nil {
		return
	}

	reply := api.NewMessage(update.Message.Chat.ID, "Checking...")
	reply.ReceiverUserID = update.Message.From.ID
	reply.ReplyParameters.EphemeralMessageID = update.Message.EphemeralMessageID
	sent, err := bot.Send(reply)
	if err != nil {
		log.Println(err)
		return
	}

	edit := api.NewEditEphemeralMessageText(
		update.Message.Chat.ID,
		update.Message.From.ID,
		sent.EphemeralMessageID,
		"Ready",
	)
	if _, err := bot.EditEphemeralMessageText(edit); err != nil {
		log.Println(err)
	}
}
