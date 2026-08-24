package tgbotapi

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func TestBotAPI102InputRichBlockJSONContract(t *testing.T) {
	animation := NewInputMediaAnimation(FileID("animation"))
	audio := NewInputMediaAudio(FileID("audio"))
	photo := NewInputMediaPhoto(FileID("photo"))
	video := NewInputMediaVideo(FileID("video"))
	voiceNote := NewInputMediaVoiceNote(FileID("voice"))

	tests := []struct {
		name       string
		block      InputRichBlock
		typeName   string
		fieldMatch string
	}{
		{name: "paragraph", block: InputRichBlockParagraph{Type: "paragraph", Text: "text"}, typeName: "paragraph", fieldMatch: `"text":"text"`},
		{name: "section heading", block: InputRichBlockSectionHeading{Type: "heading", Text: "heading", Size: 2}, typeName: "heading", fieldMatch: `"size":2`},
		{name: "preformatted", block: InputRichBlockPreformatted{Type: "pre", Text: "code", Language: "go"}, typeName: "pre", fieldMatch: `"language":"go"`},
		{name: "footer", block: InputRichBlockFooter{Type: "footer", Text: "footer"}, typeName: "footer", fieldMatch: `"text":"footer"`},
		{name: "divider", block: InputRichBlockDivider{Type: "divider"}, typeName: "divider"},
		{name: "mathematical expression", block: InputRichBlockMathematicalExpression{Type: "mathematical_expression", Expression: "x^2"}, typeName: "mathematical_expression", fieldMatch: `"expression":"x^2"`},
		{name: "anchor", block: InputRichBlockAnchor{Type: "anchor", Name: "intro"}, typeName: "anchor", fieldMatch: `"name":"intro"`},
		{name: "list", block: InputRichBlockList{Type: "list", Items: []InputRichBlockListItem{{Blocks: []InputRichBlock{InputRichBlockParagraph{Type: "paragraph", Text: "item"}}, HasCheckbox: true}}}, typeName: "list", fieldMatch: `"has_checkbox":true`},
		{name: "block quotation", block: InputRichBlockBlockQuotation{Type: "blockquote", Blocks: []InputRichBlock{InputRichBlockParagraph{Type: "paragraph", Text: "quote"}}, Credit: "author"}, typeName: "blockquote", fieldMatch: `"credit":"author"`},
		{name: "expandable quotation", block: InputRichBlockExpandableBlockQuotation{Type: "expandable_blockquote", Text: "quote", Credit: "author"}, typeName: "expandable_blockquote", fieldMatch: `"credit":"author"`},
		{name: "pull quotation", block: InputRichBlockPullQuotation{Type: "pullquote", Text: "quote", Credit: "author"}, typeName: "pullquote", fieldMatch: `"credit":"author"`},
		{name: "collage", block: InputRichBlockCollage{Type: "collage", Blocks: []InputRichBlock{InputRichBlockPhoto{Type: "photo", Photo: photo}}}, typeName: "collage", fieldMatch: `"blocks":[`},
		{name: "slideshow", block: InputRichBlockSlideshow{Type: "slideshow", Blocks: []InputRichBlock{InputRichBlockPhoto{Type: "photo", Photo: photo}}}, typeName: "slideshow", fieldMatch: `"blocks":[`},
		{name: "table", block: InputRichBlockTable{Type: "table", Cells: [][]RichBlockTableCell{{{Text: "cell", Align: "left", Valign: "middle"}}}, IsBordered: true, IsCompact: true}, typeName: "table", fieldMatch: `"is_compact":true`},
		{name: "details", block: InputRichBlockDetails{Type: "details", Summary: "summary", Blocks: []InputRichBlock{InputRichBlockParagraph{Type: "paragraph", Text: "body"}}, IsOpen: true}, typeName: "details", fieldMatch: `"is_open":true`},
		{name: "map", block: InputRichBlockMap{Type: "map", Location: Location{Latitude: 10.5, Longitude: 20.25}, Zoom: 12, Width: 640, Height: 480}, typeName: "map", fieldMatch: `"zoom":12`},
		{name: "buttons", block: InputRichBlockButtons{Type: "buttons", Buttons: []RichMessageButton{{Text: "Go", CallbackData: "go"}}, Align: "center"}, typeName: "buttons", fieldMatch: `"align":"center"`},
		{name: "animation", block: InputRichBlockAnimation{Type: "animation", Animation: animation}, typeName: "animation", fieldMatch: `"animation":{"type":"animation","media":"animation"}`},
		{name: "audio", block: InputRichBlockAudio{Type: "audio", Audio: audio}, typeName: "audio", fieldMatch: `"audio":{"type":"audio","media":"audio"}`},
		{name: "document", block: InputRichBlockDocument{Type: "document", Document: NewInputMediaDocument(FileID("document"))}, typeName: "document", fieldMatch: `"document":{"type":"document","media":"document"}`},
		{name: "photo", block: InputRichBlockPhoto{Type: "photo", Photo: photo}, typeName: "photo", fieldMatch: `"photo":{"type":"photo","media":"photo"}`},
		{name: "video", block: InputRichBlockVideo{Type: "video", Video: video}, typeName: "video", fieldMatch: `"video":{"type":"video","media":"video"}`},
		{name: "voice note", block: InputRichBlockVoiceNote{Type: "voice_note", VoiceNote: voiceNote}, typeName: "voice_note", fieldMatch: `"voice_note":{"type":"voice_note","media":"voice"}`},
		{name: "thinking", block: InputRichBlockThinking{Type: "thinking", Text: "thinking"}, typeName: "thinking", fieldMatch: `"text":"thinking"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.block)
			if err != nil {
				t.Fatalf("marshal block: %v", err)
			}

			payload := string(data)
			if !strings.Contains(payload, `"type":"`+test.typeName+`"`) {
				t.Fatalf("missing type discriminator in %s", payload)
			}
			if test.fieldMatch != "" && !strings.Contains(payload, test.fieldMatch) {
				t.Fatalf("missing %s in %s", test.fieldMatch, payload)
			}
		})
	}
}

func TestBotAPI102InputRichMessageForms(t *testing.T) {
	photo := NewInputMediaPhoto(FileID("photo-id"))
	tests := []struct {
		name    string
		message InputRichMessage
		match   string
	}{
		{name: "html", message: NewInputRichMessageHTML("<b>Hello</b>"), match: `"html":"\u003cb\u003eHello\u003c/b\u003e"`},
		{name: "markdown", message: NewInputRichMessageMarkdown("**Hello**"), match: `"markdown":"**Hello**"`},
		{name: "blocks", message: NewInputRichMessageBlocks(InputRichBlockParagraph{Type: "paragraph", Text: "Hello"}), match: `"blocks":[{"type":"paragraph","text":"Hello"}]`},
		{name: "media", message: InputRichMessage{HTML: `<img src="tg://photo?id=hero">`, Media: []InputRichMessageMedia{{ID: "hero", Media: &photo}}}, match: `"media":[{"id":"hero","media":{"type":"photo","media":"photo-id"}}]`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.message)
			if err != nil {
				t.Fatalf("marshal rich message: %v", err)
			}
			if payload := string(data); !strings.Contains(payload, test.match) {
				t.Fatalf("missing %s in %s", test.match, payload)
			}
		})
	}
}

func TestBotAPI102IncomingTypes(t *testing.T) {
	var update Update
	if err := json.Unmarshal([]byte(`{"update_id":1,"subscription":{"user":{"id":9,"is_bot":false,"first_name":"Ada"},"invoice_payload":"plan","state":"active"}}`), &update); err != nil {
		t.Fatalf("unmarshal subscription update: %v", err)
	}
	if update.Subscription == nil || update.Subscription.InvoicePayload != "plan" || update.Subscription.State != "active" {
		t.Fatalf("subscription mismatch: %#v", update.Subscription)
	}
	if from := update.SentFrom(); from == nil || from.ID != 9 {
		t.Fatalf("subscription sender mismatch: %#v", from)
	}

	var message Message
	if err := json.Unmarshal([]byte(`{"message_id":0,"receiver_user":{"id":10,"is_bot":false,"first_name":"Grace"},"ephemeral_message_id":42,"community_chat_added":{"community":{"id":1000000000000,"name":"Go"}},"community_chat_removed":{}}`), &message); err != nil {
		t.Fatalf("unmarshal ephemeral message: %v", err)
	}
	if message.MessageID != 0 || message.EphemeralMessageID != 42 || message.ReceiverUser == nil || message.ReceiverUser.ID != 10 {
		t.Fatalf("ephemeral message mismatch: %#v", message)
	}
	if message.CommunityChatAdded == nil || message.CommunityChatAdded.Community.ID != 1000000000000 || message.CommunityChatRemoved == nil {
		t.Fatalf("community service fields mismatch: %#v", message)
	}

	var chat ChatFullInfo
	if err := json.Unmarshal([]byte(`{"id":-1001,"type":"supergroup","community":{"id":2000000000000,"name":"Backend"}}`), &chat); err != nil {
		t.Fatalf("unmarshal chat community: %v", err)
	}
	if chat.Community == nil || chat.Community.ID != 2000000000000 || chat.Community.Name != "Backend" {
		t.Fatalf("chat community mismatch: %#v", chat.Community)
	}

	data, err := json.Marshal(BotCommand{Command: "quick", Description: "Quick reply", IsEphemeral: true})
	if err != nil {
		t.Fatalf("marshal ephemeral command: %v", err)
	}
	if !strings.Contains(string(data), `"is_ephemeral":true`) {
		t.Fatalf("missing ephemeral command flag: %s", data)
	}
}

func TestBotAPI102ReplyParametersIdentifiers(t *testing.T) {
	tests := []struct {
		name       string
		parameters ReplyParameters
		match      string
		notMatch   string
	}{
		{name: "regular message", parameters: ReplyParameters{MessageID: 123}, match: `"message_id":123`, notMatch: `"ephemeral_message_id"`},
		{name: "ephemeral message", parameters: ReplyParameters{EphemeralMessageID: 456}, match: `"ephemeral_message_id":456`, notMatch: `"message_id"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.parameters)
			if err != nil {
				t.Fatalf("marshal reply parameters: %v", err)
			}
			payload := string(data)
			if !strings.Contains(payload, test.match) || strings.Contains(payload, test.notMatch) {
				t.Fatalf("identifier contract mismatch: %s", payload)
			}
		})
	}
}

func TestBotAPI102NewEphemeralMessageDirectAdminParams(t *testing.T) {
	config := NewEphemeralMessage(-1001, 42, "private")
	params, err := config.params()
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	if params["chat_id"] != "-1001" || params["text"] != "private" {
		t.Fatalf("direct admin ephemeral params mismatch: %#v", params)
	}
	raw, ok := params["ephemeral_message_parameters"]
	if !ok {
		t.Fatalf("missing ephemeral_message_parameters: %#v", params)
	}
	if !strings.Contains(raw, `"receiver_user_id":42`) {
		t.Fatalf("unexpected ephemeral_message_parameters: %s", raw)
	}
	if strings.Contains(raw, "callback_query_id") {
		t.Fatalf("direct admin request unexpectedly contains callback query: %s", raw)
	}
	if _, ok := params["reply_parameters"]; ok {
		t.Fatalf("direct admin request unexpectedly contains reply parameters: %#v", params)
	}
}

func TestBotAPI102EphemeralSendParams(t *testing.T) {
	ephemeral := &EphemeralMessageParameters{ReceiverUserID: 42, CallbackQueryID: "callback"}
	message := NewMessage(1, "text")
	message.EphemeralMessageParameters = ephemeral
	animation := NewAnimation(1, FileID("animation"))
	animation.EphemeralMessageParameters = ephemeral
	audio := NewAudio(1, FileID("audio"))
	audio.EphemeralMessageParameters = ephemeral
	document := NewDocument(1, FileID("document"))
	document.EphemeralMessageParameters = ephemeral
	photo := NewPhoto(1, FileID("photo"))
	photo.EphemeralMessageParameters = ephemeral
	livePhoto := NewLivePhoto(1, FileID("live-photo"), FileID("photo"))
	livePhoto.EphemeralMessageParameters = ephemeral
	sticker := NewSticker(1, FileID("sticker"))
	sticker.EphemeralMessageParameters = ephemeral
	video := NewVideo(1, FileID("video"))
	video.EphemeralMessageParameters = ephemeral
	videoNote := NewVideoNote(1, 10, FileID("video-note"))
	videoNote.EphemeralMessageParameters = ephemeral
	voice := NewVoice(1, FileID("voice"))
	voice.EphemeralMessageParameters = ephemeral
	contact := NewContact(1, "+12025550123", "Ada")
	contact.EphemeralMessageParameters = ephemeral
	location := NewLocation(1, 10.5, 20.25)
	location.EphemeralMessageParameters = ephemeral
	venue := NewVenue(1, "Office", "Main Street", 10.5, 20.25)
	venue.EphemeralMessageParameters = ephemeral

	configs := []struct {
		method string
		config Chattable
	}{
		{method: "sendMessage", config: message},
		{method: "sendAnimation", config: animation},
		{method: "sendAudio", config: audio},
		{method: "sendDocument", config: document},
		{method: "sendPhoto", config: photo},
		{method: "sendLivePhoto", config: livePhoto},
		{method: "sendSticker", config: sticker},
		{method: "sendVideo", config: video},
		{method: "sendVideoNote", config: videoNote},
		{method: "sendVoice", config: voice},
		{method: "sendContact", config: contact},
		{method: "sendLocation", config: location},
		{method: "sendVenue", config: venue},
	}
	for _, test := range configs {
		t.Run(test.method, func(t *testing.T) {
			if method := test.config.method(); method != test.method {
				t.Fatalf("method mismatch: got %q, want %q", method, test.method)
			}
			params, err := test.config.params()
			if err != nil {
				t.Fatalf("params: %v", err)
			}
			raw, ok := params["ephemeral_message_parameters"]
			if !ok {
				t.Fatalf("missing ephemeral_message_parameters: %#v", params)
			}
			if !strings.Contains(raw, `"receiver_user_id":42`) || !strings.Contains(raw, `"callback_query_id":"callback"`) {
				t.Fatalf("ephemeral send params mismatch: %s", raw)
			}
		})
	}

	chatAction := NewChatAction(1, ChatTyping)
	params, err := chatAction.params()
	if err != nil {
		t.Fatalf("chat action params: %v", err)
	}
	if _, ok := params["ephemeral_message_parameters"]; ok {
		t.Fatalf("ephemeral_message_parameters leaked into unsupported method: %#v", params)
	}
}

func TestBotAPI102EphemeralLifecycleParams(t *testing.T) {
	markup := NewInlineKeyboardMarkup(NewInlineKeyboardRow(NewInlineKeyboardButtonData("Done", "done")))
	text := NewEditEphemeralMessageText(1, 2, 3, "updated")
	text.ParseMode = ModeMarkdownV2
	mediaValue := NewInputMediaPhoto(FileID("photo"))
	media := NewEditEphemeralMessageMedia(1, 2, 3, &mediaValue)
	caption := NewEditEphemeralMessageCaption(1, 2, 3, "")
	replyMarkup := NewEditEphemeralMessageReplyMarkup(1, 2, 3, markup)
	deleteMessage := NewDeleteEphemeralMessage(1, 2, 3)

	tests := []struct {
		config Chattable
		method string
		key    string
		value  string
	}{
		{config: text, method: "editEphemeralMessageText", key: "text", value: "updated"},
		{config: media, method: "editEphemeralMessageMedia", key: "media", value: `{"type":"photo","media":"photo"}`},
		{config: caption, method: "editEphemeralMessageCaption", key: "caption", value: ""},
		{config: replyMarkup, method: "editEphemeralMessageReplyMarkup", key: "reply_markup", value: `{"inline_keyboard":[[{"text":"Done","callback_data":"done"}]]}`},
		{config: deleteMessage, method: "deleteEphemeralMessage"},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			if method := test.config.method(); method != test.method {
				t.Fatalf("method mismatch: got %q, want %q", method, test.method)
			}
			params, err := test.config.params()
			if err != nil {
				t.Fatalf("params: %v", err)
			}
			if params["chat_id"] != "1" || params["receiver_user_id"] != "2" || params["ephemeral_message_id"] != "3" {
				t.Fatalf("ephemeral identifiers mismatch: %#v", params)
			}
			if test.key != "" {
				value, ok := params[test.key]
				if !ok || value != test.value {
					t.Fatalf("%s mismatch: got %q want %q in %#v", test.key, value, test.value, params)
				}
			}
		})
	}

	uploadMedia := NewInputMediaPhoto(FileBytes{Name: "photo.jpg", Bytes: []byte("photo")})
	uploadConfig := NewEditEphemeralMessageMedia(1, 2, 3, &uploadMedia)
	fileable, ok := any(uploadConfig).(Fileable)
	if !ok {
		t.Fatal("ephemeral media edit must support direct file uploads")
	}
	files := fileable.files()
	if len(files) != 1 || files[0].Name != "file-0" {
		t.Fatalf("unexpected ephemeral media upload files: %#v", files)
	}
}

func TestBotAPI102RichMessageMultipartPreparation(t *testing.T) {
	topPhoto := NewInputMediaPhoto(FileBytes{Name: "top.jpg", Bytes: []byte("top")})
	topRichMessage := InputRichMessage{
		HTML:  `<img src="tg://photo?id=hero">`,
		Media: []InputRichMessageMedia{{ID: "hero", Media: &topPhoto}},
	}
	topConfig := NewSendRichMessage(1, topRichMessage)
	assertRichMessageUpload(t, topConfig, []string{"rich-message-media-0"})
	if _, ok := topPhoto.Media.(FileBytes); !ok {
		t.Fatalf("top-level user media was mutated: %#v", topPhoto.Media)
	}

	video := NewInputMediaVideo(FileBytes{Name: "video.mp4", Bytes: []byte("video")})
	video.Thumb = FileBytes{Name: "thumb.jpg", Bytes: []byte("thumb")}
	voiceNote := NewInputMediaVoiceNote(FileBytes{Name: "voice.ogg", Bytes: []byte("voice")})
	photo := NewInputMediaPhoto(FileBytes{Name: "photo.jpg", Bytes: []byte("photo")})
	videoBlock := &InputRichBlockVideo{Type: "video", Video: video}
	voiceNoteBlock := &InputRichBlockVoiceNote{Type: "voice_note", VoiceNote: voiceNote}
	photoBlock := &InputRichBlockPhoto{Type: "photo", Photo: photo}
	details := &InputRichBlockDetails{
		Type:    "details",
		Summary: "Media",
		Blocks: []InputRichBlock{
			videoBlock,
		},
	}
	list := &InputRichBlockList{
		Type: "list",
		Items: []InputRichBlockListItem{{
			Blocks: []InputRichBlock{details},
		}},
	}
	collage := &InputRichBlockCollage{
		Type: "collage",
		Blocks: []InputRichBlock{
			list,
			voiceNoteBlock,
			&InputRichBlockBlockQuotation{
				Type: "blockquote",
				Blocks: []InputRichBlock{
					&InputRichBlockSlideshow{
						Type: "slideshow",
						Blocks: []InputRichBlock{
							photoBlock,
						},
					},
				},
			},
		},
	}
	blocksConfig := NewSendRichMessage(1, NewInputRichMessageBlocks(collage))
	expectedFiles := []string{
		"rich-message-block-0-block-0-item-0-block-0-block-0",
		"rich-message-block-0-block-0-item-0-block-0-block-0-thumb",
		"rich-message-block-0-block-1",
		"rich-message-block-0-block-2-block-0-block-0",
	}
	assertRichMessageUpload(t, blocksConfig, expectedFiles)
	if _, ok := videoBlock.Video.Media.(FileBytes); !ok {
		t.Fatalf("nested user media was mutated: %#v", videoBlock.Video.Media)
	}
	if _, ok := videoBlock.Video.Thumb.(FileBytes); !ok {
		t.Fatalf("nested user thumbnail was mutated: %#v", videoBlock.Video.Thumb)
	}
	if _, ok := voiceNoteBlock.VoiceNote.Media.(FileBytes); !ok {
		t.Fatalf("nested user voice note was mutated: %#v", voiceNoteBlock.VoiceNote.Media)
	}
	if _, ok := photoBlock.Photo.Media.(FileBytes); !ok {
		t.Fatalf("nested user photo was mutated: %#v", photoBlock.Photo.Media)
	}

	edit := NewEditMessageText(1, 2, "")
	edit.RichMessage = topRichMessage
	assertRichMessageUpload(t, edit, []string{"rich-message-media-0"})
	if _, ok := topPhoto.Media.(FileBytes); !ok {
		t.Fatalf("regular edit mutated user media: %#v", topPhoto.Media)
	}

	inlinePhoto := NewInputMediaPhoto(FileBytes{Name: "inline.jpg", Bytes: []byte("inline")})
	inline := EditMessageTextConfig{
		BaseEdit: BaseEdit{InlineMessageID: "inline"},
		RichMessage: InputRichMessage{
			HTML:  `<img src="tg://photo?id=hero">`,
			Media: []InputRichMessageMedia{{ID: "hero", Media: &inlinePhoto}},
		},
	}
	params, err := inline.params()
	if err != nil {
		t.Fatalf("inline edit params: %v", err)
	}
	if files := inline.files(); len(files) != 0 {
		t.Fatalf("inline edit unexpectedly supports uploads: %+v", files)
	}
	if strings.Contains(params["rich_message"], "attach://") {
		t.Fatalf("inline rich message was rewritten: %s", params["rich_message"])
	}
	if _, ok := inlinePhoto.Media.(FileBytes); !ok {
		t.Fatalf("inline user media was mutated: %#v", inlinePhoto.Media)
	}

	draft := NewSendRichMessageDraft(1, 2, topRichMessage)
	if _, ok := any(draft).(Fileable); ok {
		t.Fatal("rich message drafts must not support direct file uploads")
	}
}

func TestBotAPI102RichMediaBlockUploads(t *testing.T) {
	tests := []struct {
		name  string
		block InputRichBlock
	}{
		{name: "animation", block: InputRichBlockAnimation{Type: "animation", Animation: NewInputMediaAnimation(FileBytes{Name: "animation.mp4", Bytes: []byte("animation")})}},
		{name: "audio", block: InputRichBlockAudio{Type: "audio", Audio: NewInputMediaAudio(FileBytes{Name: "audio.mp3", Bytes: []byte("audio")})}},
		{name: "photo", block: InputRichBlockPhoto{Type: "photo", Photo: NewInputMediaPhoto(FileBytes{Name: "photo.jpg", Bytes: []byte("photo")})}},
		{name: "video", block: InputRichBlockVideo{Type: "video", Video: NewInputMediaVideo(FileBytes{Name: "video.mp4", Bytes: []byte("video")})}},
		{name: "voice note", block: InputRichBlockVoiceNote{Type: "voice_note", VoiceNote: NewInputMediaVoiceNote(FileBytes{Name: "voice.ogg", Bytes: []byte("voice")})}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := NewSendRichMessage(1, NewInputRichMessageBlocks(test.block))
			assertRichMessageUpload(t, config, []string{"rich-message-block-0"})

			data, err := json.Marshal(test.block)
			if err != nil {
				t.Fatalf("marshal original block: %v", err)
			}
			if strings.Contains(string(data), "attach://") {
				t.Fatalf("user block was mutated: %s", data)
			}
		})
	}
}

func assertRichMessageUpload(t *testing.T, config Fileable, expectedFiles []string) {
	t.Helper()

	params, err := config.params()
	if err != nil {
		t.Fatalf("rich message params: %v", err)
	}
	for _, name := range expectedFiles {
		if !strings.Contains(params["rich_message"], `"attach://`+name+`"`) {
			t.Fatalf("missing attachment %q in %s", name, params["rich_message"])
		}
	}

	files := config.files()
	names := make([]string, len(files))
	for idx, file := range files {
		names[idx] = file.Name
	}
	for _, name := range expectedFiles {
		if !slices.Contains(names, name) {
			t.Fatalf("missing file %q in %v", name, names)
		}
	}
	if len(files) != len(expectedFiles) {
		t.Fatalf("unexpected files: %v", names)
	}
}

var (
	_ func(*BotAPI, EditEphemeralMessageTextConfig) (bool, error)        = (*BotAPI).EditEphemeralMessageText
	_ func(*BotAPI, EditEphemeralMessageMediaConfig) (bool, error)       = (*BotAPI).EditEphemeralMessageMedia
	_ func(*BotAPI, EditEphemeralMessageCaptionConfig) (bool, error)     = (*BotAPI).EditEphemeralMessageCaption
	_ func(*BotAPI, EditEphemeralMessageReplyMarkupConfig) (bool, error) = (*BotAPI).EditEphemeralMessageReplyMarkup
	_ func(*BotAPI, DeleteEphemeralMessageConfig) (bool, error)          = (*BotAPI).DeleteEphemeralMessage
)
