# Bot API 10.x

## Rich Messages

Bot API 10.2 supports structured outgoing blocks and media uploads in rich messages.

```go
photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath("diagram.jpg"))
richMessage := tgbotapi.NewInputRichMessageBlocks(
	tgbotapi.InputRichBlockSectionHeading{
		Type: "heading",
		Text: "Release summary",
		Size: 2,
	},
	tgbotapi.InputRichBlockParagraph{
		Type: "paragraph",
		Text: "The migration completed successfully.",
	},
	tgbotapi.InputRichBlockPhoto{
		Type:  "photo",
		Photo: photo,
	},
)

_, err := bot.SendRichMessage(tgbotapi.NewSendRichMessage(chatID, richMessage))
```

`SendRichMessageConfig` and regular, non-inline `EditMessageTextConfig` recursively upload media from top-level `media` entries and nested collage, slideshow, list, quotation, details, and media blocks. The library replaces uploaded values with stable `attach://...` references without changing the caller's values.

For HTML or Markdown, associate an explicit media identifier with the corresponding `tg://photo?id=`, `tg://video?id=`, or `tg://audio?id=` link:

```go
photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath("diagram.jpg"))
richMessage := tgbotapi.NewInputRichMessageHTML(`<img src="tg://photo?id=diagram">`)
richMessage.Media = []tgbotapi.InputRichMessageMedia{
	{ID: "diagram", Media: &photo},
}

_, err := bot.SendRichMessage(tgbotapi.NewSendRichMessage(chatID, richMessage))
```

Exactly one of `HTML`, `Markdown`, or `Blocks` must be used. Direct upload of new files is not available for `SendRichMessageDraftConfig`, inline `EditMessageTextConfig`, or `EditEphemeralMessageMediaConfig`; use a Telegram `file_id` or an HTTP URL there.

## Ephemeral Commands and Messages

Register an ephemeral command by setting `BotCommand.IsEphemeral`:

```go
_, err := bot.Request(tgbotapi.NewSetMyCommands(tgbotapi.BotCommand{
	Command:     "private_status",
	Description: "Show status privately",
	IsEphemeral: true,
}))
```

An incoming ephemeral command has `Message.MessageID == 0` and a separate `Message.EphemeralMessageID`. Use that identifier in `ReplyParameters` and in the ephemeral edit/delete methods:

```go
incoming := update.Message
reply := tgbotapi.NewMessage(incoming.Chat.ID, "Checking...")
reply.ReceiverUserID = incoming.From.ID
reply.ReplyParameters.EphemeralMessageID = incoming.EphemeralMessageID

sent, err := bot.Send(reply)
if err != nil {
	log.Println(err)
} else {
	_, err = bot.EditEphemeralMessageText(tgbotapi.NewEditEphemeralMessageText(
		incoming.Chat.ID,
		incoming.From.ID,
		sent.EphemeralMessageID,
		"Ready",
	))
}
```

For a callback-query-triggered response, set `CallbackQueryID` on the selected send config instead of replying to an ephemeral message. Non-administrator bots must respond within 15 seconds of the eligible action. Telegram does not guarantee delivery, especially when the receiver is offline.

The edit and delete methods require all three identifiers: chat, receiver user, and ephemeral message. `EditEphemeralMessageCaption` accepts an empty caption to remove it. New uploads are not supported by `EditEphemeralMessageMedia`.

## Subscription and Community Updates

Request `UpdateTypeSubscription` to receive payment subscription changes. Community membership changes arrive as service messages:

```go
updateConfig := tgbotapi.NewUpdate(0)
updateConfig.AllowedUpdates = []string{
	tgbotapi.UpdateTypeMessage,
	tgbotapi.UpdateTypeSubscription,
}

for update := range bot.GetUpdatesChan(updateConfig) {
	if update.Subscription != nil {
		log.Printf("subscription %s for user %d", update.Subscription.State, update.Subscription.User.ID)
	}
	if update.Message != nil && update.Message.CommunityChatAdded != nil {
		community := update.Message.CommunityChatAdded.Community
		log.Printf("joined community %s (%d)", community.Name, community.ID)
	}
	if update.Message != nil && update.Message.CommunityChatRemoved != nil {
		log.Print("removed from community")
	}
}
```

`Update.SentFrom()` returns the subscriber for subscription updates. `ChatFullInfo.Community` identifies the community associated with a chat when Telegram returns it.

## Mini App Origin Protection

Telegram enables origin protection for all Mini Apps on July 20, 2026. Mini App methods are then accepted only from the original Mini App domain. This is configured through BotFather and requires no library code; opting out makes the bot responsible for avoiding links to untrusted sites.

## Guest Bots

Guest Mode lets a user invoke a bot in any chat without adding it as a member.
Request the `guest_message` update and answer with `answerGuestQuery`.

```go
updateConfig := tgbotapi.NewUpdate(0)
updateConfig.AllowedUpdates = []string{tgbotapi.UpdateTypeGuestMessage}

for update := range bot.GetUpdatesChan(updateConfig) {
	if update.GuestMessage == nil || update.GuestMessage.GuestQueryID == "" {
		continue
	}

	result := tgbotapi.NewInlineQueryResultArticle("guest-reply", "Reply", "Hello")
	_, err := bot.AnswerGuestQuery(tgbotapi.NewAnswerGuestQuery(update.GuestMessage.GuestQueryID, result))
	if err != nil {
		log.Println(err)
	}
}
```

`Message.GuestBotCallerUser` and `Message.GuestBotCallerChat` contain caller context when Telegram provides it.

## Live Photos

Live photos use two file values: the video part and the static photo.

```go
msg := tgbotapi.NewLivePhoto(chatID, tgbotapi.FilePath("live.mp4"), tgbotapi.FilePath("live.jpg"))
msg.Caption = "Live photo"
_, err := bot.SendLivePhoto(msg)
```

For paid media, wrap `InputMediaLivePhoto` with `NewInputPaidMediaLivePhoto`.

## Poll Media

Polls can include media on the poll itself, quiz explanations, and options.

```go
sticker := tgbotapi.NewInputMediaSticker(tgbotapi.FileID("sticker_file_id"))
poll := tgbotapi.NewPoll(chatID, "Choose one", tgbotapi.NewPollOptionWithMedia("A", &sticker), tgbotapi.NewPollOption("B"))

location := tgbotapi.NewInputMediaLocation(40.7128, -74.0060)
poll.Media = &location
poll.MembersOnly = true
poll.CountryCodes = []string{"US"}

_, err := bot.Send(poll)
```

## Reaction Moderation

Admin bots can remove a reaction from one message, or clear recent reactions from a user or actor chat.

```go
reaction := tgbotapi.NewDeleteMessageReaction(chatID, messageID)
reaction.UserID = userID
_, err := bot.DeleteMessageReaction(reaction)

clearRecent := tgbotapi.NewDeleteAllMessageReactions(chatID)
clearRecent.UserID = userID
_, err = bot.DeleteAllMessageReactions(clearRecent)
```
