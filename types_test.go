package tgbotapi

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUserStringWith(t *testing.T) {
	user := User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "",
		LanguageCode: "en",
		IsBot:        false,
	}

	if user.String() != "Test Test" {
		t.Fail()
	}
}

func TestUserStringWithUserName(t *testing.T) {
	user := User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "@test",
		LanguageCode: "en",
	}

	if user.String() != "@test" {
		t.Fail()
	}
}

func TestMessageTime(t *testing.T) {
	message := Message{Date: 0}

	date := time.Unix(0, 0)
	if message.Time() != date {
		t.Fail()
	}
}

func TestMessageTimeAfter2038(t *testing.T) {
	const messageUnixDate int64 = 2208988800
	message := Message{Date: messageUnixDate}

	date := time.Unix(messageUnixDate, 0)
	if message.Time() != date {
		t.Fail()
	}
}

func TestMessageUnmarshalDateAfter2038(t *testing.T) {
	const body = `{"message_id":1,"date":2208988800,"chat":{"id":1,"type":"private"}}`

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		t.Fatal(err)
	}

	if message.Date != 2208988800 {
		t.Fatalf("incorrect date: %d", message.Date)
	}
}

func TestBotAPI101FieldsUnmarshal(t *testing.T) {
	const userBody = `{"id":1,"is_bot":true,"first_name":"Bot","supports_join_request_queries":true}`

	var user User
	if err := json.Unmarshal([]byte(userBody), &user); err != nil {
		t.Fatal(err)
	}
	if !user.SupportsJoinRequestQueries {
		t.Fatalf("expected supports_join_request_queries: %#v", user)
	}

	const chatJoinRequestBody = `{
		"chat":{"id":-1,"type":"supergroup","title":"Chat"},
		"from":{"id":2,"is_bot":false,"first_name":"Alice"},
		"user_chat_id":2,
		"date":123,
		"query_id":"join-query"
	}`

	var chatJoinRequest ChatJoinRequest
	if err := json.Unmarshal([]byte(chatJoinRequestBody), &chatJoinRequest); err != nil {
		t.Fatal(err)
	}
	if chatJoinRequest.QueryID != "join-query" {
		t.Fatalf("unexpected query_id: %#v", chatJoinRequest)
	}

	const pollMediaBody = `{"link":{"url":"https://core.telegram.org/bots/api"}}`

	var pollMedia PollMedia
	if err := json.Unmarshal([]byte(pollMediaBody), &pollMedia); err != nil {
		t.Fatal(err)
	}
	if pollMedia.Link == nil || pollMedia.Link.URL != "https://core.telegram.org/bots/api" {
		t.Fatalf("unexpected poll media link: %#v", pollMedia)
	}

	const chatFullInfoBody = `{"id":-1,"type":"supergroup","title":"Chat","guard_bot":{"id":3,"is_bot":true,"first_name":"Guard"}}`

	var chatFullInfo ChatFullInfo
	if err := json.Unmarshal([]byte(chatFullInfoBody), &chatFullInfo); err != nil {
		t.Fatal(err)
	}
	if chatFullInfo.GuardBot == nil || chatFullInfo.GuardBot.FirstName != "Guard" {
		t.Fatalf("unexpected guard_bot: %#v", chatFullInfo)
	}
}

func TestRichMessageJSON(t *testing.T) {
	richMessage := RichMessage{
		Blocks: []RichBlock{
			RichBlockParagraph{
				Type: "paragraph",
				Text: "Hello",
			},
			RichBlockSectionHeading{
				Type: "heading",
				Text: []RichText{
					RichTextBold{
						Type: "bold",
						Text: "News",
					},
				},
				Size: 2,
			},
		},
		IsRTL: true,
	}

	data, err := json.Marshal(richMessage)
	if err != nil {
		t.Fatal(err)
	}
	payload := string(data)

	for _, expected := range []string{
		`"blocks":[`,
		`"type":"paragraph"`,
		`"text":"Hello"`,
		`"type":"heading"`,
		`"size":2`,
		`"is_rtl":true`,
	} {
		if !strings.Contains(payload, expected) {
			t.Fatalf("expected %q in %s", expected, payload)
		}
	}
}

func TestMessageUnmarshalRichMessage(t *testing.T) {
	const body = `{
		"message_id":1,
		"date":123,
		"chat":{"id":1,"type":"private"},
		"rich_message":{
			"blocks":[{"type":"paragraph","text":"Hello"}],
			"is_rtl":true
		}
	}`

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		t.Fatal(err)
	}

	if message.RichMessage == nil || len(message.RichMessage.Blocks) != 1 || !message.RichMessage.IsRTL {
		t.Fatalf("unexpected rich_message: %#v", message.RichMessage)
	}
}

func TestPollUnmarshalBotAPI96(t *testing.T) {
	const body = `{
		"id":"poll-1",
		"question":"quiz",
		"options":[{"persistent_id":"opt-1","text":"A","voter_count":1,"addition_date":123,"added_by_user":{"id":2,"is_bot":false,"first_name":"Alice"}}],
		"total_voter_count":1,
		"is_closed":false,
		"is_anonymous":false,
		"type":"quiz",
		"allows_multiple_answers":true,
		"allows_revoting":true,
		"correct_option_ids":[0],
		"description":"desc",
		"description_entities":[{"type":"bold","offset":0,"length":4}]
	}`

	var poll Poll
	if err := json.Unmarshal([]byte(body), &poll); err != nil {
		t.Fatal(err)
	}

	if !poll.AllowsRevoting {
		t.Fatalf("expected allows_revoting to be true: %#v", poll)
	}
	if len(poll.CorrectOptionIDs) != 1 || poll.CorrectOptionIDs[0] != 0 {
		t.Fatalf("unexpected correct_option_ids: %#v", poll.CorrectOptionIDs)
	}
	if poll.CorrectOptionID != 0 {
		t.Fatalf("unexpected compatibility correct_option_id: %d", poll.CorrectOptionID)
	}
	if poll.Description != "desc" {
		t.Fatalf("unexpected description: %#v", poll)
	}
	if len(poll.DescriptionEntities) != 1 || poll.DescriptionEntities[0].Type != "bold" {
		t.Fatalf("unexpected description_entities: %#v", poll.DescriptionEntities)
	}
	if len(poll.Options) != 1 || poll.Options[0].PersistentID != "opt-1" {
		t.Fatalf("unexpected poll options: %#v", poll.Options)
	}
	if poll.Options[0].AddedByUser == nil || poll.Options[0].AddedByUser.FirstName != "Alice" {
		t.Fatalf("unexpected added_by_user: %#v", poll.Options[0].AddedByUser)
	}
}

func TestPollUnmarshalLegacyCorrectOptionID(t *testing.T) {
	const body = `{
		"id":"poll-legacy",
		"question":"quiz",
		"options":[],
		"total_voter_count":0,
		"is_closed":false,
		"is_anonymous":true,
		"type":"quiz",
		"allows_multiple_answers":false,
		"correct_option_id":0
	}`

	var poll Poll
	if err := json.Unmarshal([]byte(body), &poll); err != nil {
		t.Fatal(err)
	}

	if len(poll.CorrectOptionIDs) != 1 || poll.CorrectOptionIDs[0] != 0 {
		t.Fatalf("unexpected correct_option_ids: %#v", poll.CorrectOptionIDs)
	}
}

func TestGuestMessageUnmarshalBotAPI10(t *testing.T) {
	const body = `{
		"update_id":1,
		"guest_message":{
			"message_id":10,
			"guest_query_id":"guest-query",
			"guest_bot_caller_user":{"id":2,"is_bot":false,"first_name":"Alice"},
			"guest_bot_caller_chat":{"id":3,"type":"group","title":"Team"},
			"from":{"id":4,"is_bot":false,"first_name":"Bob"},
			"date":1,
			"chat":{"id":5,"type":"group","title":"Work"},
			"text":"hello"
		}
	}`

	var update Update
	if err := json.Unmarshal([]byte(body), &update); err != nil {
		t.Fatal(err)
	}

	if update.GuestMessage == nil || update.GuestMessage.GuestQueryID != "guest-query" {
		t.Fatalf("unexpected guest_message payload: %#v", update.GuestMessage)
	}
	if update.GuestMessage.GuestBotCallerUser == nil || update.GuestMessage.GuestBotCallerUser.ID != 2 {
		t.Fatalf("unexpected guest caller user: %#v", update.GuestMessage.GuestBotCallerUser)
	}
	if update.SentFrom() == nil || update.SentFrom().ID != 4 {
		t.Fatalf("unexpected SentFrom result: %#v", update.SentFrom())
	}
	if update.FromChat() == nil || update.FromChat().ID != 5 {
		t.Fatalf("unexpected FromChat result: %#v", update.FromChat())
	}
}

func TestBotAPI10MediaUnmarshal(t *testing.T) {
	const body = `{
		"message_id":1,
		"date":1,
		"chat":{"id":1,"type":"private"},
		"live_photo":{
			"file_id":"live",
			"file_unique_id":"unique-live",
			"width":640,
			"height":480,
			"duration":2,
			"photo":[{"file_id":"photo","file_unique_id":"unique-photo","width":640,"height":480}]
		},
		"paid_media":{
			"star_count":5,
			"paid_media":[{
				"type":"live_photo",
				"live_photo":{
					"file_id":"paid-live",
					"file_unique_id":"paid-unique",
					"width":640,
					"height":480,
					"duration":2
				}
			}]
		}
	}`

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		t.Fatal(err)
	}

	if message.LivePhoto == nil || message.LivePhoto.FileID != "live" || len(message.LivePhoto.Photo) != 1 {
		t.Fatalf("unexpected live_photo: %#v", message.LivePhoto)
	}
	if message.PaidMedia == nil || len(message.PaidMedia.PaidMedia) != 1 ||
		message.PaidMedia.PaidMedia[0].LivePhoto == nil ||
		message.PaidMedia.PaidMedia[0].LivePhoto.FileID != "paid-live" {
		t.Fatalf("unexpected paid media live photo: %#v", message.PaidMedia)
	}
}

func TestPollUnmarshalBotAPI10MediaAndLimits(t *testing.T) {
	const body = `{
		"id":"poll-10",
		"question":"q",
		"options":[{
			"persistent_id":"opt-1",
			"text":"A",
			"voter_count":1,
			"media":{"sticker":{"file_id":"sticker","file_unique_id":"sticker-unique","type":"regular","width":512,"height":512,"is_animated":false,"is_video":false}}
		}],
		"total_voter_count":1,
		"is_closed":false,
		"is_anonymous":false,
		"type":"regular",
		"allows_multiple_answers":false,
		"members_only":true,
		"country_codes":["US","PL"],
		"media":{"location":{"latitude":10.5,"longitude":20.25}},
		"explanation_media":{"live_photo":{"file_id":"live","file_unique_id":"unique","width":1,"height":1,"duration":1}}
	}`

	var poll Poll
	if err := json.Unmarshal([]byte(body), &poll); err != nil {
		t.Fatal(err)
	}

	if !poll.MembersOnly || len(poll.CountryCodes) != 2 || poll.Media == nil || poll.Media.Location == nil {
		t.Fatalf("unexpected poll media/limits: %#v", poll)
	}
	if poll.ExplanationMedia == nil || poll.ExplanationMedia.LivePhoto == nil {
		t.Fatalf("unexpected explanation media: %#v", poll.ExplanationMedia)
	}
	if len(poll.Options) != 1 || poll.Options[0].Media == nil || poll.Options[0].Media.Sticker == nil {
		t.Fatalf("unexpected option media: %#v", poll.Options)
	}
}

func TestBotAPI10PermissionsUnmarshal(t *testing.T) {
	const permissionsBody = `{"can_send_messages":true,"can_react_to_messages":true}`
	var permissions ChatPermissions
	if err := json.Unmarshal([]byte(permissionsBody), &permissions); err != nil {
		t.Fatal(err)
	}
	if !permissions.CanReactToMessages {
		t.Fatalf("expected can_react_to_messages: %#v", permissions)
	}

	const memberBody = `{"user":{"id":1,"is_bot":false,"first_name":"A"},"status":"restricted","is_member":true,"can_react_to_messages":true}`
	var member ChatMember
	if err := json.Unmarshal([]byte(memberBody), &member); err != nil {
		t.Fatal(err)
	}
	if !member.CanReactToMessages {
		t.Fatalf("expected member can_react_to_messages: %#v", member)
	}
}

func TestPollAnswerUnmarshalBotAPI96(t *testing.T) {
	const body = `{"poll_id":"poll-1","option_ids":[0],"option_persistent_ids":["opt-1"]}`

	var answer PollAnswer
	if err := json.Unmarshal([]byte(body), &answer); err != nil {
		t.Fatal(err)
	}

	if len(answer.OptionPersistentIDs) != 1 || answer.OptionPersistentIDs[0] != "opt-1" {
		t.Fatalf("unexpected option_persistent_ids: %#v", answer.OptionPersistentIDs)
	}
}

func TestManagedBotUpdatedUnmarshal(t *testing.T) {
	const body = `{
		"update_id":1,
		"managed_bot":{
			"user":{"id":1,"is_bot":false,"first_name":"Owner"},
			"bot":{"id":2,"is_bot":true,"first_name":"Managed"}
		}
	}`

	var update Update
	if err := json.Unmarshal([]byte(body), &update); err != nil {
		t.Fatal(err)
	}

	if update.ManagedBot == nil || update.ManagedBot.Bot.ID != 2 || update.ManagedBot.User.ID != 1 {
		t.Fatalf("unexpected managed_bot payload: %#v", update.ManagedBot)
	}
}

func TestManagedBotCreatedUnmarshal(t *testing.T) {
	const body = `{
		"message_id":1,
		"date":1,
		"chat":{"id":1,"type":"private"},
		"managed_bot_created":{"bot":{"id":2,"is_bot":true,"first_name":"Managed"}}
	}`

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		t.Fatal(err)
	}

	if message.ManagedBotCreated == nil || message.ManagedBotCreated.Bot.ID != 2 {
		t.Fatalf("unexpected managed_bot_created payload: %#v", message.ManagedBotCreated)
	}
}

func TestPreparedKeyboardButtonUnmarshal(t *testing.T) {
	const body = `{"id":"prepared-button"}`

	var button PreparedKeyboardButton
	if err := json.Unmarshal([]byte(body), &button); err != nil {
		t.Fatal(err)
	}

	if button.ID != "prepared-button" {
		t.Fatalf("unexpected prepared keyboard button: %#v", button)
	}
}

func TestPollOptionServiceMessagesUnmarshal(t *testing.T) {
	const addedBody = `{
		"poll_message":{"message_id":10,"date":1,"chat":{"id":1,"type":"private"}},
		"option_persistent_id":"opt-1",
		"option_text":"A",
		"option_text_entities":[{"type":"bold","offset":0,"length":1}]
	}`
	const deletedBody = `{
		"poll_message":{"message_id":11,"date":1,"chat":{"id":1,"type":"private"}},
		"option_persistent_id":"opt-2",
		"option_text":"B"
	}`

	var added PollOptionAdded
	if err := json.Unmarshal([]byte(addedBody), &added); err != nil {
		t.Fatal(err)
	}
	if added.PollMessage == nil || added.OptionPersistentID != "opt-1" || added.OptionText != "A" {
		t.Fatalf("unexpected poll_option_added payload: %#v", added)
	}
	if len(added.OptionTextEntities) != 1 || added.OptionTextEntities[0].Type != "bold" {
		t.Fatalf("unexpected option_text_entities: %#v", added.OptionTextEntities)
	}

	var deleted PollOptionDeleted
	if err := json.Unmarshal([]byte(deletedBody), &deleted); err != nil {
		t.Fatal(err)
	}
	if deleted.PollMessage == nil || deleted.OptionPersistentID != "opt-2" || deleted.OptionText != "B" {
		t.Fatalf("unexpected poll_option_deleted payload: %#v", deleted)
	}
}

func TestVideoChatScheduledTimeAfter2038(t *testing.T) {
	const startDate int64 = 2208988800
	scheduled := VideoChatScheduled{StartDate: startDate}

	date := time.Unix(startDate, 0)
	if scheduled.Time() != date {
		t.Fail()
	}
}

func TestMessageIsCommandWithCommand(t *testing.T) {
	message := Message{Text: "/command"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	if !message.IsCommand() {
		t.Fail()
	}
}

func TestIsCommandWithText(t *testing.T) {
	message := Message{Text: "some text"}

	if message.IsCommand() {
		t.Fail()
	}
}

func TestIsCommandWithEmptyText(t *testing.T) {
	message := Message{Text: ""}

	if message.IsCommand() {
		t.Fail()
	}
}

func TestCommandWithCommand(t *testing.T) {
	message := Message{Text: "/command"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	if message.Command() != "command" {
		t.Fail()
	}
}

func TestCommandWithEmptyText(t *testing.T) {
	message := Message{Text: ""}

	if message.Command() != "" {
		t.Fail()
	}
}

func TestCommandWithNonCommand(t *testing.T) {
	message := Message{Text: "test text"}

	if message.Command() != "" {
		t.Fail()
	}
}

func TestCommandWithBotName(t *testing.T) {
	message := Message{Text: "/command@testbot"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 16}}

	if message.Command() != "command" {
		t.Fail()
	}
}

func TestCommandWithAtWithBotName(t *testing.T) {
	message := Message{Text: "/command@testbot"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 16}}

	if message.CommandWithAt() != "command@testbot" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithArguments(t *testing.T) {
	message := Message{Text: "/command with arguments"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}
	if message.CommandArguments() != "with arguments" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithMalformedArguments(t *testing.T) {
	message := Message{Text: "/command-without argument space"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}
	if message.CommandArguments() != "without argument space" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsWithoutArguments(t *testing.T) {
	message := Message{Text: "/command"}
	if message.CommandArguments() != "" {
		t.Fail()
	}
}

func TestMessageCommandArgumentsForNonCommand(t *testing.T) {
	message := Message{Text: "test text"}
	if message.CommandArguments() != "" {
		t.Fail()
	}
}

func TestMessageEntityParseURLGood(t *testing.T) {
	entity := MessageEntity{URL: "https://www.google.com"}

	if _, err := entity.ParseURL(); err != nil {
		t.Fail()
	}
}

func TestMessageEntityParseURLBad(t *testing.T) {
	entity := MessageEntity{URL: ""}

	if _, err := entity.ParseURL(); err == nil {
		t.Fail()
	}
}

func TestChatIsPrivate(t *testing.T) {
	chat := Chat{ID: 10, Type: "private"}

	if !chat.IsPrivate() {
		t.Fail()
	}
}

func TestChatIsGroup(t *testing.T) {
	chat := Chat{ID: 10, Type: "group"}

	if !chat.IsGroup() {
		t.Fail()
	}
}

func TestChatIsChannel(t *testing.T) {
	chat := Chat{ID: 10, Type: "channel"}

	if !chat.IsChannel() {
		t.Fail()
	}
}

func TestChatIsSuperGroup(t *testing.T) {
	chat := Chat{ID: 10, Type: "supergroup"}

	if !chat.IsSuperGroup() {
		t.Fail()
	}
}

func TestMessageEntityIsMention(t *testing.T) {
	entity := MessageEntity{Type: "mention"}

	if !entity.IsMention() {
		t.Fail()
	}
}

func TestMessageEntityIsHashtag(t *testing.T) {
	entity := MessageEntity{Type: "hashtag"}

	if !entity.IsHashtag() {
		t.Fail()
	}
}

func TestMessageEntityIsBotCommand(t *testing.T) {
	entity := MessageEntity{Type: "bot_command"}

	if !entity.IsCommand() {
		t.Fail()
	}
}

func TestMessageEntityIsUrl(t *testing.T) {
	entity := MessageEntity{Type: "url"}

	if !entity.IsURL() {
		t.Fail()
	}
}

func TestMessageEntityIsEmail(t *testing.T) {
	entity := MessageEntity{Type: "email"}

	if !entity.IsEmail() {
		t.Fail()
	}
}

func TestMessageEntityIsBold(t *testing.T) {
	entity := MessageEntity{Type: "bold"}

	if !entity.IsBold() {
		t.Fail()
	}
}

func TestMessageEntityIsItalic(t *testing.T) {
	entity := MessageEntity{Type: "italic"}

	if !entity.IsItalic() {
		t.Fail()
	}
}

func TestMessageEntityIsCode(t *testing.T) {
	entity := MessageEntity{Type: "code"}

	if !entity.IsCode() {
		t.Fail()
	}
}

func TestMessageEntityIsPre(t *testing.T) {
	entity := MessageEntity{Type: "pre"}

	if !entity.IsPre() {
		t.Fail()
	}
}

func TestMessageEntityIsTextLink(t *testing.T) {
	entity := MessageEntity{Type: "text_link"}

	if !entity.IsTextLink() {
		t.Fail()
	}
}

func TestFileLink(t *testing.T) {
	file := File{FilePath: "test/test.txt"}

	if file.Link("token") != "https://api.telegram.org/file/bottoken/test/test.txt" {
		t.Fail()
	}
}

// Ensure all configs are sendable
var (
	_ Chattable = AddStickerConfig{}
	_ Chattable = AnimationConfig{}
	_ Chattable = AnswerWebAppQueryConfig{}
	_ Chattable = AudioConfig{}
	_ Chattable = BanChatMemberConfig{}
	_ Chattable = BanChatSenderChatConfig{}
	_ Chattable = CallbackConfig{}
	_ Chattable = ChatActionConfig{}
	_ Chattable = ChatAdministratorsConfig{}
	_ Chattable = ChatInfoConfig{}
	_ Chattable = ChatInviteLinkConfig{}
	_ Chattable = CloseConfig{}
	_ Chattable = CloseForumTopicConfig{}
	_ Chattable = CloseGeneralForumTopicConfig{}
	_ Chattable = ContactConfig{}
	_ Chattable = CopyMessageConfig{}
	_ Chattable = CreateChatInviteLinkConfig{}
	_ Chattable = CreateChatSubscriptionLinkConfig{}
	_ Chattable = CreateForumTopicConfig{}
	_ Chattable = DeleteChatPhotoConfig{}
	_ Chattable = DeleteChatStickerSetConfig{}
	_ Chattable = DeleteForumTopicConfig{}
	_ Chattable = DeleteMessageConfig{}
	_ Chattable = DeleteMyCommandsConfig{}
	_ Chattable = DeleteStickerSetConfig{}
	_ Chattable = DeleteWebhookConfig{}
	_ Chattable = DocumentConfig{}
	_ Chattable = EditChatInviteLinkConfig{}
	_ Chattable = EditChatSubscriptionLinkConfig{}
	_ Chattable = EditForumTopicConfig{}
	_ Chattable = EditGeneralForumTopicConfig{}
	_ Chattable = EditMessageCaptionConfig{}
	_ Chattable = EditMessageLiveLocationConfig{}
	_ Chattable = EditMessageMediaConfig{}
	_ Chattable = EditMessageReplyMarkupConfig{}
	_ Chattable = EditMessageTextConfig{}
	_ Chattable = FileConfig{}
	_ Chattable = ForwardConfig{}
	_ Chattable = GameConfig{}
	_ Chattable = GetAvailableGiftsConfig{}
	_ Chattable = GetBusinessConnectionConfig{}
	_ Chattable = GetChatMemberConfig{}
	_ Chattable = GetChatMenuButtonConfig{}
	_ Chattable = GetForumTopicIconStickersConfig{}
	_ Chattable = GetGameHighScoresConfig{}
	_ Chattable = GetMyDefaultAdministratorRightsConfig{}
	_ Chattable = GetMyDescriptionConfig{}
	_ Chattable = GetMyNameConfig{}
	_ Chattable = GetMyShortDescriptionConfig{}
	_ Chattable = GetStarTransactionsConfig{}
	_ Chattable = HideGeneralForumTopicConfig{}
	_ Chattable = InlineConfig{}
	_ Chattable = InvoiceConfig{}
	_ Chattable = KickChatMemberConfig{}
	_ Chattable = LeaveChatConfig{}
	_ Chattable = LocationConfig{}
	_ Chattable = LogOutConfig{}
	_ Chattable = MediaGroupConfig{}
	_ Chattable = MessageConfig{}
	_ Chattable = PaidMediaConfig{}
	_ Chattable = PhotoConfig{}
	_ Chattable = PinChatMessageConfig{}
	_ Chattable = PreCheckoutConfig{}
	_ Chattable = PromoteChatMemberConfig{}
	_ Chattable = RefundStarPaymentConfig{}
	_ Chattable = RemoveChatVerificationConfig{}
	_ Chattable = RemoveUserVerificationConfig{}
	_ Chattable = ReopenForumTopicConfig{}
	_ Chattable = ReopenGeneralForumTopicConfig{}
	_ Chattable = ReplaceStickerInSetConfig{}
	_ Chattable = RestrictChatMemberConfig{}
	_ Chattable = RevokeChatInviteLinkConfig{}
	_ Chattable = SendPollConfig{}
	_ Chattable = SetChatDescriptionConfig{}
	_ Chattable = SetChatMemberTagConfig{}
	_ Chattable = SetChatMenuButtonConfig{}
	_ Chattable = SetChatPhotoConfig{}
	_ Chattable = SetChatTitleConfig{}
	_ Chattable = SetCustomEmojiStickerSetThumbnailConfig{}
	_ Chattable = SetGameScoreConfig{}
	_ Chattable = SetMessageReactionConfig{}
	_ Chattable = SetMyDefaultAdministratorRightsConfig{}
	_ Chattable = SetMyDescriptionConfig{}
	_ Chattable = SetMyNameConfig{}
	_ Chattable = SetMyShortDescriptionConfig{}
	_ Chattable = SetStickerEmojiListConfig{}
	_ Chattable = SetStickerKeywordsConfig{}
	_ Chattable = SetStickerMaskPositionConfig{}
	_ Chattable = SetStickerSetTitleConfig{}
	_ Chattable = SetUserEmojiStatusConfig{}
	_ Chattable = ShippingConfig{}
	_ Chattable = StickerConfig{}
	_ Chattable = StopMessageLiveLocationConfig{}
	_ Chattable = StopPollConfig{}
	_ Chattable = UnbanChatMemberConfig{}
	_ Chattable = UnbanChatSenderChatConfig{}
	_ Chattable = UnhideGeneralForumTopicConfig{}
	_ Chattable = UnpinAllForumTopicMessagesConfig{}
	_ Chattable = UnpinAllGeneralForumTopicMessagesConfig{}
	_ Chattable = UnpinChatMessageConfig{}
	_ Chattable = UpdateConfig{}
	_ Chattable = UserProfilePhotosConfig{}
	_ Chattable = VenueConfig{}
	_ Chattable = VerifyChatConfig{}
	_ Chattable = VerifyUserConfig{}
	_ Chattable = VideoConfig{}
	_ Chattable = VideoNoteConfig{}
	_ Chattable = VoiceConfig{}
	_ Chattable = WebhookConfig{}
)

// Ensure all Fileable types are correct.
var (
	_ Fileable = (*PhotoConfig)(nil)
	_ Fileable = (*AudioConfig)(nil)
	_ Fileable = (*DocumentConfig)(nil)
	_ Fileable = (*StickerConfig)(nil)
	_ Fileable = (*VideoConfig)(nil)
	_ Fileable = (*AnimationConfig)(nil)
	_ Fileable = (*VideoNoteConfig)(nil)
	_ Fileable = (*VoiceConfig)(nil)
	_ Fileable = (*SetChatPhotoConfig)(nil)
	_ Fileable = (*EditMessageMediaConfig)(nil)
	_ Fileable = (*SetChatPhotoConfig)(nil)
	_ Fileable = (*UploadStickerConfig)(nil)
	_ Fileable = (*NewStickerSetConfig)(nil)
	_ Fileable = (*AddStickerConfig)(nil)
	_ Fileable = (*MediaGroupConfig)(nil)
	_ Fileable = (*WebhookConfig)(nil)
	_ Fileable = (*SetStickerSetThumbConfig)(nil)
	_ Fileable = (*PaidMediaConfig)(nil)
)

// Ensure all RequestFileData types are correct.
var (
	_ RequestFileData = (*FilePath)(nil)
	_ RequestFileData = (*FileBytes)(nil)
	_ RequestFileData = (*FileReader)(nil)
	_ RequestFileData = (*FileURL)(nil)
	_ RequestFileData = (*FileID)(nil)
	_ RequestFileData = (*fileAttach)(nil)
)
