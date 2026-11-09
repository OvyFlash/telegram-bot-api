package tgbotapi

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestBotAPI104ConfigParams(t *testing.T) {
	widget := NewSendWidget(1, "<b>hi</b>")
	widget.MessageThreadID = 4
	widget.Width = 320
	widget.Height = 200
	widget.PublicData = "pub"
	widget.ChatData = "chat"
	widget.Caption = "cap"
	widget.ParseMode = ModeHTML
	widget.CaptionEntities = []MessageEntity{{Type: "bold", Offset: 0, Length: 3}}
	widget.ShowCaptionAboveMedia = true
	widget.AllowQueries = true
	widget.IsShared = true
	widget.ProtectContent = true

	sharedWidget := NewSendSharedWidget(1, "w1")
	sharedWidget.ChatData = "chat"
	sharedWidget.Caption = "cap"
	sharedWidget.ShowCaptionAboveMedia = true

	editWidget := NewEditMessageWidgetData(1, 2)
	editWidget.PublicData = "pub"

	editInlineWidget := EditMessageWidgetDataConfig{
		ChatConfig:       ChatConfig{ChatID: 1},
		MessageID:        2,
		InlineMessageID:  "inline",
		RemovePublicData: true,
		ChatData:         "ignored",
		RemoveChatData:   true,
	}

	editSharedWidget := NewEditSharedWidgetData("w1", "")
	editSharedWidget.CacheTime = 60

	verifyUser := NewVerifyUser(5, "Partner")
	verifyUser.CustomDescriptionParseMode = ModeHTML
	verifyUser.CustomDescriptionEntities = []MessageEntity{{Type: "text_link", Offset: 0, Length: 7, URL: "https://t.me"}}

	verifyChat := NewVerifyChat(ChatConfig{ChannelUsername: "@channel"}, "Partner")
	verifyChat.CustomDescriptionParseMode = ModeMarkdownV2

	tests := []struct {
		name   string
		config Chattable
		method string
		want   Params
	}{
		{
			name:   "send widget",
			config: widget,
			method: "sendWidget",
			want: Params{
				"chat_id":                  "1",
				"message_thread_id":        "4",
				"html":                     "<b>hi</b>",
				"width":                    "320",
				"height":                   "200",
				"public_data":              "pub",
				"chat_data":                "chat",
				"caption":                  "cap",
				"parse_mode":               ModeHTML,
				"caption_entities":         `[{"type":"bold","offset":0,"length":3}]`,
				"show_caption_above_media": "true",
				"allow_queries":            "true",
				"is_shared":                "true",
				"protect_content":          "true",
			},
		},
		{
			name:   "send shared widget",
			config: sharedWidget,
			method: "sendSharedWidget",
			want: Params{
				"chat_id":                  "1",
				"widget_id":                "w1",
				"chat_data":                "chat",
				"caption":                  "cap",
				"show_caption_above_media": "true",
			},
		},
		{
			name:   "answer widget query",
			config: NewAnswerWidgetQuery("q1", "reply"),
			method: "answerWidgetQuery",
			want:   Params{"widget_query_id": "q1", "data": "reply"},
		},
		{
			name:   "edit message widget data keeps omitted chat data",
			config: editWidget,
			method: "editMessageWidgetData",
			want:   Params{"chat_id": "1", "message_id": "2", "public_data": "pub"},
		},
		{
			name:   "edit inline widget data removes data",
			config: editInlineWidget,
			method: "editMessageWidgetData",
			want:   Params{"inline_message_id": "inline", "public_data": "", "chat_data": ""},
		},
		{
			name:   "edit shared widget data removes public data",
			config: editSharedWidget,
			method: "editSharedWidgetData",
			want:   Params{"widget_id": "w1", "public_data": "", "cache_time": "60"},
		},
		{
			name:   "delete business account gift",
			config: DeleteBusinessAccountGiftConfig{BusinessConnectionID: "bc", OwnedGiftID: "g1"},
			method: "deleteBusinessAccountGift",
			want:   Params{"business_connection_id": "bc", "owned_gift_id": "g1"},
		},
		{
			name:   "get unique gift",
			config: GetUniqueGiftConfig{Name: "PlushPepe-1"},
			method: "getUniqueGift",
			want:   Params{"name": "PlushPepe-1"},
		},
		{
			name:   "get chat share theme backdrops",
			config: GetChatShareThemeBackdropsConfig{},
			method: "getChatShareThemeBackdrops",
			want:   Params{},
		},
		{
			name:   "verify user with entities",
			config: verifyUser,
			method: "verifyUser",
			want: Params{
				"user_id":                       "5",
				"custom_description":            "Partner",
				"custom_description_parse_mode": ModeHTML,
				"custom_description_entities":   `[{"type":"text_link","offset":0,"length":7,"url":"https://t.me"}]`,
			},
		},
		{
			name:   "verify chat with parse mode",
			config: verifyChat,
			method: "verifyChat",
			want: Params{
				"chat_id":                       "@channel",
				"custom_description":            "Partner",
				"custom_description_parse_mode": ModeMarkdownV2,
			},
		},
		{
			name:   "ephemeral callback reply omits chat and receiver",
			config: NewEphemeralCallbackReply("cbq", "only you"),
			method: "sendMessage",
			want: Params{
				"text":                         "only you",
				"ephemeral_message_parameters": `{"callback_query_id":"cbq"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.config.method(); got != test.method {
				t.Fatalf("method = %q, want %q", got, test.method)
			}
			params, err := test.config.params()
			if err != nil {
				t.Fatalf("params: %v", err)
			}
			if params == nil {
				params = Params{}
			}
			if !reflect.DeepEqual(params, test.want) {
				t.Fatalf("params = %#v, want %#v", params, test.want)
			}
		})
	}
}

func TestBotAPI104WidgetQueryUpdates(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		check func(t *testing.T, update Update)
	}{
		{
			name: "widget query with signed data",
			raw: `{
				"update_id": 1,
				"widget_query": {
					"id": "q1",
					"from": {"id": 7, "first_name": "Ann"},
					"chat": {"id": 10, "type": "private"},
					"message_id": 3,
					"chat_instance": "ci",
					"type": {
						"type": "sign_data",
						"wallet": {"address": "EQ1", "amount": "100", "public_key": "ab", "state_init": "te6"},
						"sign_data": "{\"payload\":1}"
					},
					"data": "sign",
					"widget_id": "w1"
				}
			}`,
			check: func(t *testing.T, update Update) {
				query := update.WidgetQuery
				if query == nil || query.ID != "q1" || query.MessageID != 3 || query.WidgetID != "w1" {
					t.Fatalf("unexpected widget query: %#v", query)
				}
				if !query.Type.IsSignData() || query.Type.IsData() || query.Type.IsSentTransaction() {
					t.Fatalf("unexpected query type: %#v", query.Type)
				}
				if query.Type.Wallet == nil || query.Type.Wallet.StateInit != "te6" || query.Type.SignData != `{"payload":1}` {
					t.Fatalf("unexpected signed data: %#v", query.Type)
				}
				if from := update.SentFrom(); from == nil || from.ID != 7 {
					t.Fatalf("SentFrom mismatch: %#v", from)
				}
				if chat := update.FromChat(); chat == nil || chat.ID != 10 {
					t.Fatalf("FromChat mismatch: %#v", chat)
				}
			},
		},
		{
			name: "inline widget query with wallet card",
			raw: `{
				"update_id": 2,
				"widget_query": {
					"id": "q2",
					"from": {"id": 8, "first_name": "Bob"},
					"inline_message_id": "inline",
					"chat_instance": "ci",
					"chat_type": "sender",
					"type": {
						"type": "data",
						"wallet_card": {
							"address": "EQ2",
							"amount": "5",
							"stickers": [{
								"position": {"x_percentage": 10, "y_percentage": 20, "width_percentage": 30, "height_percentage": 40, "rotation_angle": 90, "corner_radius_percentage": 5},
								"nft": {"address": "NFT", "index": "1", "name": "Pepe", "description": "Frog", "attributes": [{"trait_type": "Model", "value": "Gold"}]}
							}]
						}
					},
					"data": "ping"
				}
			}`,
			check: func(t *testing.T, update Update) {
				query := update.WidgetQuery
				if query == nil || query.InlineMessageID != "inline" || query.ChatType != "sender" {
					t.Fatalf("unexpected widget query: %#v", query)
				}
				if !query.Type.IsData() || query.Type.WalletCard == nil {
					t.Fatalf("unexpected query type: %#v", query.Type)
				}
				sticker := query.Type.WalletCard.Stickers[0]
				if sticker.Position.RotationAngle != 90 || sticker.Nft.Attributes[0].Value != "Gold" {
					t.Fatalf("unexpected wallet card sticker: %#v", sticker)
				}
				if chat := update.FromChat(); chat != nil {
					t.Fatalf("FromChat for inline widget query = %#v, want nil", chat)
				}
			},
		},
		{
			name: "shared widget query with sent transaction",
			raw: `{
				"update_id": 3,
				"shared_widget_query": {
					"id": "q3",
					"from": {"id": 9, "first_name": "Eve"},
					"type": {"type": "sent_transaction", "wallet": {"address": "EQ3", "amount": "1"}, "msg_hash": "hash"},
					"data": "paid",
					"widget_id": "w2"
				}
			}`,
			check: func(t *testing.T, update Update) {
				query := update.SharedWidgetQuery
				if query == nil || query.WidgetID != "w2" || !query.Type.IsSentTransaction() || query.Type.MsgHash != "hash" {
					t.Fatalf("unexpected shared widget query: %#v", query)
				}
				if from := update.SentFrom(); from == nil || from.ID != 9 {
					t.Fatalf("SentFrom mismatch: %#v", from)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var update Update
			if err := json.Unmarshal([]byte(test.raw), &update); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			test.check(t, update)
		})
	}
}

func TestBotAPI104MessageFields(t *testing.T) {
	var message Message
	if err := json.Unmarshal([]byte(`{
		"message_id": 1,
		"date": 1,
		"chat": {"id": 10, "type": "private"},
		"widget": {"id": "w1", "width": 300, "html": "<p>hi</p>", "public_data": "pub", "chat_data": "chat"},
		"external_reply": {
			"origin": {"type": "user", "date": 1, "sender_user": {"id": 2, "first_name": "A"}},
			"widget": {"html": "<p>quoted</p>"}
		},
		"quote": {"text": "", "rich_texts": ["plain", {"type": "bold", "text": "bold"}], "position": 4},
		"ton_wallet_transfer": {"transaction_id": "", "peer_address": "EQ", "amount": "-15", "comment": "thanks"},
		"paid_messages_waived": {}
	}`), &message); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if message.Widget == nil || message.Widget.ID != "w1" || message.Widget.Width != 300 || message.Widget.ChatData != "chat" {
		t.Fatalf("unexpected widget: %#v", message.Widget)
	}
	if message.ExternalReply == nil || message.ExternalReply.Widget == nil || message.ExternalReply.Widget.HTML != "<p>quoted</p>" {
		t.Fatalf("unexpected external reply widget: %#v", message.ExternalReply)
	}
	if message.Quote == nil || len(message.Quote.RichTexts) != 2 || message.Quote.Position != 4 {
		t.Fatalf("unexpected quote: %#v", message.Quote)
	}
	if message.TonWalletTransfer == nil || message.TonWalletTransfer.Amount != "-15" || message.TonWalletTransfer.Comment != "thanks" {
		t.Fatalf("unexpected TON wallet transfer: %#v", message.TonWalletTransfer)
	}
	if message.PaidMessagesWaived == nil {
		t.Fatal("paid_messages_waived service message was not decoded")
	}
}

func TestBotAPI104ChatFullInfoShareTheme(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		check func(t *testing.T, chat ChatFullInfo)
	}{
		{
			name: "emoji theme with wallet card",
			raw: `{
				"id": 1,
				"type": "private",
				"ton_wallet_card": {"address": "EQ", "amount": "0", "stickers": [], "background_nft_address": "NFT"},
				"share_theme": {
					"type": "emoji",
					"emojis": [{"type": "emoji", "emoji": "🔥"}, {"type": "custom_emoji", "custom_emoji_id": "123"}],
					"background_backdrop_id": 7
				}
			}`,
			check: func(t *testing.T, chat ChatFullInfo) {
				if chat.TonWalletCard == nil || chat.TonWalletCard.BackgroundNftAddress != "NFT" {
					t.Fatalf("unexpected wallet card: %#v", chat.TonWalletCard)
				}
				theme := chat.ShareTheme
				if theme == nil || !theme.IsEmoji() || theme.IsUniqueGift() || theme.BackgroundBackdropID != 7 {
					t.Fatalf("unexpected share theme: %#v", theme)
				}
				if !theme.Emojis[0].IsEmoji() || !theme.Emojis[1].IsCustomEmoji() || theme.Emojis[1].CustomEmojiID != "123" {
					t.Fatalf("unexpected theme emoji: %#v", theme.Emojis)
				}
			},
		},
		{
			name: "unique gift theme",
			raw: `{
				"id": 2,
				"type": "channel",
				"share_theme": {"type": "unique_gift", "gift": {"gift_id": "g", "base_name": "Pepe", "name": "Pepe-1", "number": 1}}
			}`,
			check: func(t *testing.T, chat ChatFullInfo) {
				theme := chat.ShareTheme
				if theme == nil || !theme.IsUniqueGift() || theme.Gift == nil || theme.Gift.Name != "Pepe-1" {
					t.Fatalf("unexpected share theme: %#v", theme)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var chat ChatFullInfo
			if err := json.Unmarshal([]byte(test.raw), &chat); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			test.check(t, chat)
		})
	}
}

func TestBotAPI104OutgoingJSON(t *testing.T) {
	widgetContent := NewInputWidgetMessageContent("Hello")
	widgetContent.AllowQueries = true
	widgetContent.IsShared = true

	sharedWidgetContent := NewInputSharedWidgetMessageContent("w1")
	sharedWidgetContent.Caption = "cap"

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "input widget message content",
			value: widgetContent,
			want:  `{"html":"Hello","allow_queries":true,"is_shared":true}`,
		},
		{
			name:  "input shared widget message content",
			value: sharedWidgetContent,
			want:  `{"widget_id":"w1","caption":"cap"}`,
		},
		{
			name: "inline article with widget content",
			value: InlineQueryResultArticle{
				Type:                "article",
				ID:                  "1",
				Title:               "Widget",
				InputMessageContent: widgetContent,
			},
			want: `"input_message_content":{"html":"Hello","allow_queries":true,"is_shared":true}`,
		},
		{
			name:  "reply parameters with rich quote",
			value: ReplyParameters{MessageID: 5, QuoteRichTexts: []RichText{"plain", RichTextBold{Type: "bold", Text: "bold"}}},
			want:  `{"message_id":5,"quote_rich_texts":["plain",{"type":"bold","text":"bold"}]}`,
		},
		{
			name:  "rich text TON address",
			value: RichTextTonAddress{Type: "ton_address", Text: "wallet", Address: "EQ"},
			want:  `{"type":"ton_address","text":"wallet","address":"EQ"}`,
		},
		{
			name:  "ephemeral parameters with receiver",
			value: EphemeralMessageParameters{ReceiverUserID: 42, CallbackQueryID: "cbq"},
			want:  `{"receiver_user_id":42,"callback_query_id":"cbq"}`,
		},
		{
			name:  "story area uses area position",
			value: StoryArea{Position: StoryAreaPosition{XPercentage: 1}, Type: StoryAreaType{Type: "link"}},
			want:  `"position":{"x_percentage":1,`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if !strings.Contains(string(raw), test.want) {
				t.Fatalf("got %s, want substring %q", raw, test.want)
			}
		})
	}
}
