package tgbotapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBotAPI103EphemeralMessageParameters(t *testing.T) {
	config := NewSendRichMessage(1, NewInputRichMessageHTML("<p>hi</p>"))
	config.EphemeralMessageParameters = EphemeralMessageParameters{
		ReceiverUserID:              42,
		CallbackQueryID:             "cbq",
		ReplaceCallbackQueryMessage: true,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	raw := params["ephemeral_message_parameters"]
	if !strings.Contains(raw, `"receiver_user_id":42`) ||
		!strings.Contains(raw, `"callback_query_id":"cbq"`) ||
		!strings.Contains(raw, `"replace_callback_query_message":true`) {
		t.Fatalf("unexpected ephemeral_message_parameters: %s", raw)
	}
}

func TestBotAPI103DraftStopParams(t *testing.T) {
	draft := SendMessageDraftConfig{
		ChatConfig: ChatConfig{ChatID: 1},
		DraftID:    7,
		Text:       "partial",
		CanStop:    true,
		KeepOnStop: true,
	}
	params, err := draft.params()
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if params["can_stop"] != "true" || params["keep_on_stop"] != "true" {
		t.Fatalf("draft stop params mismatch: %#v", params)
	}

	richDraft := NewSendRichMessageDraft(1, 7, NewInputRichMessageMarkdown("**hi**"))
	richDraft.CanStop = true
	richDraft.KeepOnStop = true
	params, err = richDraft.params()
	if err != nil {
		t.Fatalf("rich draft params: %v", err)
	}
	if params["can_stop"] != "true" || params["keep_on_stop"] != "true" {
		t.Fatalf("rich draft stop params mismatch: %#v", params)
	}
}

func TestBotAPI103EditEphemeralExtensions(t *testing.T) {
	rich := NewInputRichMessageHTML("<p>hi</p>")
	text := NewEditEphemeralMessageText(1, 2, 3, "")
	text.RichMessage = &rich
	params, err := text.params()
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if _, ok := params["text"]; ok {
		t.Fatalf("empty text should be omitted when rich_message is set: %#v", params)
	}
	if !strings.Contains(params["rich_message"], `"html"`) || !strings.Contains(params["rich_message"], "hi") {
		t.Fatalf("unexpected rich_message: %s", params["rich_message"])
	}

	caption := NewEditEphemeralMessageCaption(1, 2, 3, "cap")
	caption.ShowCaptionAboveMedia = true
	params, err = caption.params()
	if err != nil {
		t.Fatalf("caption params: %v", err)
	}
	if params["show_caption_above_media"] != "true" {
		t.Fatalf("show_caption_above_media missing: %#v", params)
	}
}

func TestBotAPI103PromoteWelcomeMessages(t *testing.T) {
	config := PromoteChatMemberConfig{
		ChatMemberConfig: ChatMemberConfig{
			ChatConfig: ChatConfig{ChatID: 1},
			UserID:     2,
		},
		CanSendWelcomeMessages: true,
	}
	params, err := config.params()
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if params["can_send_welcome_messages"] != "true" {
		t.Fatalf("can_send_welcome_messages missing: %#v", params)
	}
}

func TestBotAPI103NewTypesJSON(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name: "rich message button",
			value: RichMessageButton{
				Text:         "Go",
				Style:        "primary",
				CallbackData: "go",
				Disabled:     &DisabledButton{},
			},
			want: `"callback_data":"go"`,
		},
		{
			name:  "rich text button",
			value: RichTextButton{Type: "button", Button: RichMessageButton{Text: "Go", URL: "https://t.me"}},
			want:  `"type":"button"`,
		},
		{
			name:  "buttons block",
			value: RichBlockButtons{Type: "buttons", Buttons: []RichMessageButton{{Text: "A", CallbackData: "a"}}, Align: RichBlockButtonsAlignCenter},
			want:  `"align":"center"`,
		},
		{
			name:  "expandable quotation",
			value: RichBlockExpandableBlockQuotation{Type: "expandable_blockquote", Text: "quote"},
			want:  `"type":"expandable_blockquote"`,
		},
		{
			name:  "document block",
			value: RichBlockDocument{Type: "document", Document: Document{FileID: "doc"}},
			want:  `"type":"document"`,
		},
		{
			name:  "compact table",
			value: RichBlockTable{Type: "table", Cells: [][]RichBlockTableCell{{{Align: "left", Valign: "top"}}}, IsCompact: true},
			want:  `"is_compact":true`,
		},
		{
			name: "input document block",
			value: InputRichBlockDocument{
				Type:     "document",
				Document: NewInputMediaDocument(FileID("doc")),
			},
			want: `"document":{"type":"document","media":"doc"}`,
		},
		{
			name: "unique gift info",
			value: UniqueGiftInfo{
				Gift:      UniqueGift{BaseName: "Gift"},
				Origin:    "transfer",
				Text:      "hello",
				IsPrivate: true,
			},
			want: `"is_private":true`,
		},
		{
			name: "message generation stopped",
			value: MessageGenerationStopped{
				Chat:    Chat{ID: 1},
				DraftID: 9,
			},
			want: `"draft_id":9`,
		},
		{
			name:  "community chat joined",
			value: CommunityChatJoined{Community: Community{ID: 5, Name: "Ops"}},
			want:  `"name":"Ops"`,
		},
		{
			name: "keyboard force reply",
			value: InlineKeyboardMarkup{
				InlineKeyboard: [][]InlineKeyboardButton{{{Text: "A", CallbackData: stringPtr("a")}}},
				ForceReply:     true,
			},
			want: `"force_reply":true`,
		},
		{
			name: "disabled inline button",
			value: InlineKeyboardButton{
				Text:     "Nope",
				Disabled: &DisabledButton{},
			},
			want: `"disabled":{}`,
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

func TestBotAPI103RichBlockButtonsAlign(t *testing.T) {
	block := RichBlockButtons{Align: RichBlockButtonsAlignCenter}
	if !block.Center() || block.Left() || block.Right() {
		t.Fatalf("unexpected align helpers: %#v", block)
	}

	input := InputRichBlockButtons{Align: RichBlockButtonsAlignRight}
	if !input.Right() || input.Left() || input.Center() {
		t.Fatalf("unexpected input align helpers: %#v", input)
	}
}

func TestBotAPI103RichMessageButtonStyle(t *testing.T) {
	button := RichMessageButton{Style: RichMessageButtonStylePrimary}
	if !button.Primary() || button.Danger() || button.LinkStyle() {
		t.Fatalf("unexpected style helpers: %#v", button)
	}
}

func TestBotAPI103UpdateStoppedMessageGeneration(t *testing.T) {
	var update Update
	if err := json.Unmarshal([]byte(`{
		"update_id":1,
		"stopped_message_generation":{"chat":{"id":10,"type":"private"},"draft_id":3}
	}`), &update); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if update.StoppedMessageGeneration == nil || update.StoppedMessageGeneration.DraftID != 3 {
		t.Fatalf("unexpected update: %#v", update.StoppedMessageGeneration)
	}
	if chat := update.FromChat(); chat == nil || chat.ID != 10 {
		t.Fatalf("FromChat mismatch: %#v", chat)
	}
}

func TestBotAPI103RichDocumentUpload(t *testing.T) {
	document := NewInputMediaDocument(FileBytes{Name: "notes.pdf", Bytes: []byte("%PDF")})
	block := InputRichBlockDocument{
		Type:     "document",
		Document: document,
	}
	config := NewSendRichMessage(1, NewInputRichMessageBlocks(block))
	assertRichMessageUpload(t, config, []string{"rich-message-block-0"})

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "attach://") {
		t.Fatalf("user block was mutated: %s", data)
	}
}

func stringPtr(v string) *string {
	return &v
}
