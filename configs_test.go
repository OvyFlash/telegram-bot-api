package tgbotapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestAnswerGuestQueryConfigParams(t *testing.T) {
	result := NewInlineQueryResultArticle("guest-result", "Answer", "Hello")
	config := NewAnswerGuestQuery("guest-query", result)

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["guest_query_id"] != "guest-query" {
		t.Fatalf("guest_query_id mismatch: %#v", params)
	}
	if !strings.Contains(params["result"], `"type":"article"`) || !strings.Contains(params["result"], `"id":"guest-result"`) {
		t.Fatalf("result payload mismatch: %#v", params)
	}
}

func TestAnswerChatJoinRequestQueryConfigParams(t *testing.T) {
	config := NewAnswerChatJoinRequestQuery("join-query", "approve")

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["chat_join_request_query_id"] != "join-query" {
		t.Fatalf("chat_join_request_query_id mismatch: %#v", params)
	}
	if params["result"] != "approve" {
		t.Fatalf("result mismatch: %#v", params)
	}
}

func TestSendChatJoinRequestWebAppConfigParams(t *testing.T) {
	config := NewSendChatJoinRequestWebApp("join-query", "https://example.com/app")

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["chat_join_request_query_id"] != "join-query" {
		t.Fatalf("chat_join_request_query_id mismatch: %#v", params)
	}
	if params["web_app_url"] != "https://example.com/app" {
		t.Fatalf("web_app_url mismatch: %#v", params)
	}
}

func TestSendRichMessageConfigParams(t *testing.T) {
	config := NewSendRichMessage(123, NewInputRichMessageHTML("<p>Hello</p>"))
	config.MessageThreadID = 456
	config.DirectMessagesTopicID = 789
	config.AllowPaidBroadcast = true
	config.SuggestedPostParameters = &SuggestedPostParameters{Price: &SuggestedPostPrice{Currency: "XTR", Amount: 1}}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["chat_id"] != "123" {
		t.Fatalf("chat_id mismatch: %#v", params)
	}
	if params["message_thread_id"] != "456" || params["direct_messages_topic_id"] != "789" {
		t.Fatalf("thread params mismatch: %#v", params)
	}
	if params["allow_paid_broadcast"] != "true" {
		t.Fatalf("allow_paid_broadcast mismatch: %#v", params)
	}
	if !strings.Contains(params["rich_message"], `"html":"\u003cp\u003eHello\u003c/p\u003e"`) {
		t.Fatalf("rich_message mismatch: %#v", params)
	}
	if !strings.Contains(params["suggested_post_parameters"], `"price":{"currency":"XTR","amount":1}`) {
		t.Fatalf("suggested_post_parameters mismatch: %#v", params)
	}
}

func TestSendRichMessageDraftConfigParams(t *testing.T) {
	config := NewSendRichMessageDraft(123, 456, NewInputRichMessageMarkdown("**Hello**"))
	config.MessageThreadID = 789

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["chat_id"] != "123" {
		t.Fatalf("chat_id mismatch: %#v", params)
	}
	if params["draft_id"] != "456" || params["message_thread_id"] != "789" {
		t.Fatalf("draft params mismatch: %#v", params)
	}
	if !strings.Contains(params["rich_message"], `"markdown":"**Hello**"`) {
		t.Fatalf("rich_message mismatch: %#v", params)
	}
}

func TestInputMediaLinkSerialization(t *testing.T) {
	option := NewPollOptionWithMedia("docs", NewInputMediaLink("https://core.telegram.org/bots/api"))
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question: "q",
		Options:  []InputPollOption{option},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if !strings.Contains(params["options"], `"media":{"type":"link","url":"https://core.telegram.org/bots/api"}`) {
		t.Fatalf("options media mismatch: %#v", params["options"])
	}
	if files := config.files(); len(files) != 0 {
		t.Fatalf("unexpected files for link media: %+v", files)
	}
}

func TestSendLivePhotoConfigUploadSerialization(t *testing.T) {
	config := NewLivePhoto(1, FilePath("tests/video.mp4"), FilePath("tests/image.jpg"))
	config.Caption = "live"
	config.HasSpoiler = true

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}
	if params["chat_id"] != "1" || params["caption"] != "live" || params["has_spoiler"] != "true" {
		t.Fatalf("params mismatch: %#v", params)
	}

	files := config.files()
	if len(files) != 2 || files[0].Name != "live_photo" || files[1].Name != "photo" {
		t.Fatalf("unexpected files payload: %+v", files)
	}
}

func TestVideoAndAnimationConfigDimensionParams(t *testing.T) {
	tests := []struct {
		name           string
		params         func() (Params, error)
		expectedWidth  string
		expectedHeight string
	}{
		{
			name: "video",
			params: func() (Params, error) {
				config := NewVideo(1, FileID("video-file-id"))
				config.Width = 1920
				config.Height = 1080
				return config.params()
			},
			expectedWidth:  "1920",
			expectedHeight: "1080",
		},
		{
			name: "animation",
			params: func() (Params, error) {
				config := NewAnimation(1, FileID("animation-file-id"))
				config.Width = 640
				config.Height = 360
				return config.params()
			},
			expectedWidth:  "640",
			expectedHeight: "360",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := tt.params()
			if err != nil {
				t.Fatalf("params failed: %v", err)
			}
			if params["width"] != tt.expectedWidth || params["height"] != tt.expectedHeight {
				t.Fatalf("dimension params mismatch: %#v", params)
			}
		})
	}
}

func TestSendPollConfigBotAPI10MediaSerialization(t *testing.T) {
	optionMedia := &InputMediaSticker{
		Type:  "sticker",
		Media: FilePath("tests/image.jpg"),
		Emoji: ":)",
	}
	explanationMedia := &InputMediaLivePhoto{
		BaseInputMedia: BaseInputMedia{
			Type:  "live_photo",
			Media: FilePath("tests/video.mp4"),
		},
		Photo: FilePath("tests/image.jpg"),
	}
	pollMedia := &InputMediaLocation{
		Type:      "location",
		Latitude:  10.5,
		Longitude: 20.25,
	}
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question:         "q",
		Options:          []InputPollOption{NewPollOptionWithMedia("a", optionMedia), NewPollOption("b")},
		ExplanationMedia: explanationMedia,
		Media:            pollMedia,
		MembersOnly:      true,
		CountryCodes:     []string{"US", "PL"},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if !strings.Contains(params["options"], `"media":{"type":"sticker","media":"attach://poll-option-0-0","emoji":":)"}`) {
		t.Fatalf("options media mismatch: %#v", params["options"])
	}
	if !strings.Contains(params["explanation_media"], `"media":"attach://explanation-media-0"`) ||
		!strings.Contains(params["explanation_media"], `"photo":"attach://explanation-media-0-photo"`) {
		t.Fatalf("explanation_media mismatch: %#v", params["explanation_media"])
	}
	if !strings.Contains(params["media"], `"type":"location"`) {
		t.Fatalf("poll media mismatch: %#v", params["media"])
	}
	if params["members_only"] != "true" || params["country_codes"] != `["US","PL"]` {
		t.Fatalf("poll limit params mismatch: %#v", params)
	}

	files := config.files()
	if len(files) != 3 {
		t.Fatalf("unexpected files count: %+v", files)
	}
	if optionMedia.Media != FilePath("tests/image.jpg") || explanationMedia.Media != FilePath("tests/video.mp4") {
		t.Fatalf("original media was mutated")
	}
}

func TestPaidMediaLivePhotoSerialization(t *testing.T) {
	livePhoto := NewInputMediaLivePhoto(FilePath("tests/video.mp4"), FilePath("tests/image.jpg"))
	paid := NewInputPaidMediaLivePhoto(&livePhoto)
	config := NewPaidMedia(1, 10, &paid)

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}
	if !strings.Contains(params["media"], `"type":"live_photo"`) ||
		!strings.Contains(params["media"], `"media":"attach://file-0"`) ||
		!strings.Contains(params["media"], `"photo":"attach://file-0-photo"`) {
		t.Fatalf("paid media payload mismatch: %#v", params["media"])
	}

	files := config.files()
	if len(files) != 2 || files[0].Name != "file-0" || files[1].Name != "file-0-photo" {
		t.Fatalf("unexpected files payload: %+v", files)
	}
}

func TestPaidMediaConfigSerializesMediaArray(t *testing.T) {
	photo := NewInputMediaPhoto(FileID("paid-photo-id"))
	paid := NewInputPaidMediaPhoto(&photo)
	config := NewPaidMedia(1, 10, &paid)

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	var media []map[string]any
	if err := json.Unmarshal([]byte(params["media"]), &media); err != nil {
		t.Fatalf("media must be a JSON array, got %q: %v", params["media"], err)
	}
	if len(media) != 1 || media[0]["type"] != "photo" || media[0]["media"] != "paid-photo-id" {
		t.Fatalf("unexpected paid media payload: %#v", media)
	}
}

func TestPaidMediaGroupConfigSerializesMultipleMedia(t *testing.T) {
	photo := NewInputMediaPhoto(FileID("paid-photo-id"))
	paidPhoto := NewInputPaidMediaPhoto(&photo)
	video := NewInputMediaVideo(FileID("paid-video-id"))
	paidVideo := NewInputPaidMediaVideo(&video)
	config := NewPaidMediaGroup(1, 10, paidPhoto, paidVideo)

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	var media []map[string]any
	if err := json.Unmarshal([]byte(params["media"]), &media); err != nil {
		t.Fatalf("media must be a JSON array, got %q: %v", params["media"], err)
	}
	if len(media) != 2 || media[0]["type"] != "photo" || media[1]["type"] != "video" {
		t.Fatalf("unexpected paid media payload: %#v", media)
	}
}

func TestPaidMediaVideoMetadataSurvivesUploadPreparation(t *testing.T) {
	video := NewInputMediaVideo(FilePath("tests/video.mp4"))
	paid := NewInputPaidMediaVideo(&video)
	paid.Cover = "cover-file-id"
	paid.StartTimestamp = 12

	prepared := prepareInputMediaForParams([]InputMedia{&paid})
	if len(prepared) != 1 {
		t.Fatalf("unexpected prepared media: %+v", prepared)
	}

	params := Params{}
	if err := params.AddInterface("media", prepared[0]); err != nil {
		t.Fatalf("marshal prepared media: %v", err)
	}
	if !strings.Contains(params["media"], `"cover":"cover-file-id"`) || !strings.Contains(params["media"], `"start_timestamp":12`) {
		t.Fatalf("paid media metadata was not preserved: %s", params["media"])
	}
}

func TestManagedBotAccessSettingsConfigFalseParam(t *testing.T) {
	config := NewSetManagedBotAccessSettings(42, false, 1, 2)

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}
	if params["user_id"] != "42" || params["is_access_restricted"] != "false" || params["added_user_ids"] != "[1,2]" {
		t.Fatalf("managed access params mismatch: %#v", params)
	}
}

func TestDeleteReactionConfigsAndReturnBotsParams(t *testing.T) {
	deleteReaction := NewDeleteMessageReaction(1, 2)
	deleteReaction.UserID = 3
	params, err := deleteReaction.params()
	if err != nil {
		t.Fatalf("delete reaction params failed: %v", err)
	}
	if params["chat_id"] != "1" || params["message_id"] != "2" || params["user_id"] != "3" {
		t.Fatalf("delete reaction params mismatch: %#v", params)
	}

	deleteAll := NewDeleteAllMessageReactions(1)
	deleteAll.ActorChatID = 4
	params, err = deleteAll.params()
	if err != nil {
		t.Fatalf("delete all reactions params failed: %v", err)
	}
	if params["chat_id"] != "1" || params["actor_chat_id"] != "4" {
		t.Fatalf("delete all reactions params mismatch: %#v", params)
	}

	admins := NewChatAdministrators(1)
	admins.ReturnBots = true
	params, err = admins.params()
	if err != nil {
		t.Fatalf("chat administrators params failed: %v", err)
	}
	if params["return_bots"] != "true" {
		t.Fatalf("return_bots mismatch: %#v", params)
	}
}

func TestSendPollConfigCloseDate64BitParam(t *testing.T) {
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question:  "q",
		Options:   []InputPollOption{{Text: "a"}},
		CloseDate: 2208988800,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["close_date"] != "2208988800" {
		t.Fatalf("close_date mismatch: %s", params["close_date"])
	}
}

func TestSendPollConfigBotAPI96Params(t *testing.T) {
	allowsRevoting := false
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question:               "q",
		Options:                []InputPollOption{{Text: "a"}, {Text: "b"}, {Text: "c"}},
		Type:                   "quiz",
		AllowsRevoting:         &allowsRevoting,
		ShuffleOptions:         true,
		AllowAddingOptions:     true,
		HideResultsUntilCloses: true,
		CorrectOptionIDs:       []int{0, 2},
		Description:            "desc",
		DescriptionEntities:    []MessageEntity{{Type: "bold", Offset: 0, Length: 4}},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["allows_revoting"] != "false" {
		t.Fatalf("allows_revoting mismatch: %#v", params)
	}
	if params["shuffle_options"] != "true" {
		t.Fatalf("shuffle_options mismatch: %#v", params)
	}
	if params["allow_adding_options"] != "true" {
		t.Fatalf("allow_adding_options mismatch: %#v", params)
	}
	if params["hide_results_until_closes"] != "true" {
		t.Fatalf("hide_results_until_closes mismatch: %#v", params)
	}
	if params["correct_option_ids"] != "[0,2]" {
		t.Fatalf("correct_option_ids mismatch: %#v", params)
	}
	if _, ok := params["correct_option_id"]; ok {
		t.Fatalf("unexpected legacy correct_option_id key: %#v", params)
	}
	if params["description"] != "desc" {
		t.Fatalf("description mismatch: %#v", params)
	}
	if params["description_entities"] != `[{"type":"bold","offset":0,"length":4}]` {
		t.Fatalf("description_entities mismatch: %#v", params)
	}
}

func TestSendPollConfigAllowsRevotingOmittedWhenNil(t *testing.T) {
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question: "q",
		Options:  []InputPollOption{{Text: "a"}, {Text: "b"}},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if _, ok := params["allows_revoting"]; ok {
		t.Fatalf("unexpected allows_revoting param: %#v", params)
	}
	if _, ok := params["correct_option_ids"]; ok {
		t.Fatalf("unexpected correct_option_ids param: %#v", params)
	}
}

func TestSendPollConfigLegacyCorrectOptionIDCompat(t *testing.T) {
	config := SendPollConfig{
		BaseChat: BaseChat{
			ChatConfig: ChatConfig{ChatID: 1},
		},
		Question:        "q",
		Options:         []InputPollOption{{Text: "a"}, {Text: "b"}},
		Type:            "quiz",
		CorrectOptionID: 0,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["correct_option_ids"] != "[0]" {
		t.Fatalf("correct_option_ids mismatch: %#v", params)
	}
	if _, ok := params["correct_option_id"]; ok {
		t.Fatalf("unexpected legacy correct_option_id key: %#v", params)
	}
}

func TestSavePreparedKeyboardButtonConfigParams(t *testing.T) {
	config := SavePreparedKeyboardButtonConfig{
		UserID: 42,
		Button: KeyboardButton{
			Text: "Create bot",
			RequestManagedBot: &KeyboardButtonRequestManagedBot{
				RequestID:         7,
				SuggestedName:     "Demo",
				SuggestedUsername: "demo_helper_bot",
			},
		},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["user_id"] != "42" {
		t.Fatalf("user_id mismatch: %#v", params)
	}
	if !strings.Contains(params["button"], `"request_managed_bot":{"request_id":7`) {
		t.Fatalf("button payload mismatch: %#v", params)
	}
}

func TestSetChatMemberTagConfigTagParam(t *testing.T) {
	config := SetChatMemberTagConfig{
		ChatMemberConfig: ChatMemberConfig{
			ChatConfig: ChatConfig{ChatID: 1},
			UserID:     777,
		},
		Tag: "text tag param",
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}
	if params["tag"] != "text tag param" {
		t.Fatalf("tag mismatch: %s", params["tag"])
	}
}

func TestBanChatSenderChatConfigUntilDate64BitParam(t *testing.T) {
	config := BanChatSenderChatConfig{
		ChatConfig:   ChatConfig{ChatID: 1},
		SenderChatID: 2,
		UntilDate:    2208988800,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["until_date"] != "2208988800" {
		t.Fatalf("until_date mismatch: %s", params["until_date"])
	}
}

func TestCreateChatInviteLinkConfigExpireDate64BitParam(t *testing.T) {
	config := CreateChatInviteLinkConfig{
		ChatConfig: ChatConfig{ChatID: 1},
		Name:       "name",
		ExpireDate: 2208988800,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["expire_date"] != "2208988800" {
		t.Fatalf("expire_date mismatch: %s", params["expire_date"])
	}
}

func TestEditChatInviteLinkConfigExpireDate64BitParam(t *testing.T) {
	config := EditChatInviteLinkConfig{
		ChatConfig: ChatConfig{ChatID: 1},
		InviteLink: "https://t.me/+abc",
		ExpireDate: 2208988800,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params failed: %v", err)
	}

	if params["expire_date"] != "2208988800" {
		t.Fatalf("expire_date mismatch: %s", params["expire_date"])
	}
}

func TestPrepareInputMediaForParams(t *testing.T) {
	tests := []struct {
		name               string
		inputMedia         []InputMedia
		expectedMediaPaths []string
		expectedThumbPaths []string
	}{
		{
			name: "photo that needs upload",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/image.jpg"),
					},
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{""},
		},
		{
			name: "video with thumbnail that need upload",
			inputMedia: []InputMedia{
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FilePath("tests/video.mp4"),
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name: "multiple media items",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/image.jpg"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"), // URL doesn't need upload
					},
				},
				&InputMediaDocument{
					BaseInputMedia: BaseInputMedia{
						Type:  "document",
						Media: FilePath("tests/audio.mp3"),
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0", "https://example.com/video.mp4", "attach://file-2"},
			expectedThumbPaths: []string{"", "", "attach://file-2-thumb"},
		},
		{
			name: "items that don't need upload",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FileID("photo123"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"),
					},
					Thumb: FileID("thumb123"),
				},
			},
			expectedMediaPaths: []string{"photo123", "https://example.com/video.mp4"},
			expectedThumbPaths: []string{"", "thumb123"},
		},
		{
			name: "paid media",
			inputMedia: []InputMedia{
				&InputPaidMedia{
					Type: "video",
					Media: &InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:  "video",
							Media: FilePath("tests/video.mp4"),
						},
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name: "complex nested media with metadata",
			inputMedia: []InputMedia{
				&InputPaidMedia{
					Type: "video",
					Media: &InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:    "video",
							Media:   FilePath("tests/complex_video.mp4"),
							Caption: "Complex video with metadata",
						},
						Width:             1920,
						Height:            1080,
						Duration:          60,
						SupportsStreaming: true,
						Thumb:             FilePath("tests/inner_thumb.jpg"),
					},
					Thumb:             FilePath("tests/outer_thumb.jpg"),
					Width:             1920,
					Height:            1080,
					Duration:          60,
					SupportsStreaming: true,
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name: "animation with caption entities and spoiler",
			inputMedia: []InputMedia{
				&InputMediaAnimation{
					BaseInputMedia: BaseInputMedia{
						Type:    "animation",
						Media:   FilePath("tests/animation.gif"),
						Caption: "Animation with *formatting*",
						CaptionEntities: []MessageEntity{
							{
								Type:   "bold",
								Offset: 14,
								Length: 11,
							},
						},
						ParseMode:  "Markdown",
						HasSpoiler: true,
					},
					Thumb: FilePath("tests/anim_thumb.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name: "audio with all media types",
			inputMedia: []InputMedia{
				&InputMediaAudio{
					BaseInputMedia: BaseInputMedia{
						Type:  "audio",
						Media: FilePath("tests/audio.mp3"),
					},
					Thumb:     FilePath("tests/audio_thumb.jpg"),
					Duration:  180,
					Performer: "Test Artist",
					Title:     "Test Track",
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name:       "empty media array",
			inputMedia: []InputMedia{
				// Empty slice
			},
			expectedMediaPaths: []string{},
			expectedThumbPaths: []string{},
		},
		{
			name: "all media types with mixed sources",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/photo.jpg"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileID("existing_video"),
					},
					Thumb: FilePath("tests/video_thumb.jpg"),
				},
				&InputMediaAnimation{
					BaseInputMedia: BaseInputMedia{
						Type:  "animation",
						Media: FileURL("https://example.com/animation.gif"),
					},
				},
				&InputMediaAudio{
					BaseInputMedia: BaseInputMedia{
						Type:  "audio",
						Media: FilePath("tests/audio.mp3"),
					},
				},
				&InputMediaDocument{
					BaseInputMedia: BaseInputMedia{
						Type:  "document",
						Media: FileID("existing_document"),
					},
				},
				&InputPaidMedia{
					Type: "photo",
					Media: &InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:  "photo",
							Media: FilePath("tests/paid_photo.jpg"),
						},
					},
				},
			},
			expectedMediaPaths: []string{
				"attach://file-0",
				"existing_video",
				"https://example.com/animation.gif",
				"attach://file-3",
				"existing_document",
				"attach://file-5",
			},
			expectedThumbPaths: []string{
				"",
				"attach://file-1-thumb",
				"",
				"",
				"",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prepareInputMediaForParams(tt.inputMedia)

			if len(result) != len(tt.inputMedia) {
				t.Errorf("Expected result length %d, got %d", len(tt.inputMedia), len(result))
			}

			for i, media := range result {
				if i >= len(tt.expectedMediaPaths) {
					break // Safety check for empty media array test
				}

				mediaPath := media.getMedia().SendData()
				if mediaPath != tt.expectedMediaPaths[i] {
					t.Errorf("Media path at index %d: expected %s, got %s", i, tt.expectedMediaPaths[i], mediaPath)
				}

				thumb := media.getThumb()
				var thumbPath string
				if thumb != nil {
					thumbPath = thumb.SendData()
				}

				expectedThumb := tt.expectedThumbPaths[i]
				if thumbPath != expectedThumb {
					t.Errorf("Thumb path at index %d: expected %s, got %s", i, expectedThumb, thumbPath)
				}
			}

			if len(tt.inputMedia) > 0 && &result[0] == &tt.inputMedia[0] {
				t.Error("Result should be a deep copy, not a reference to the original slice")
			}
		})
	}
}

func TestPrepareInputMediaForFiles(t *testing.T) {
	tests := []struct {
		name          string
		inputMedia    []InputMedia
		expectedFiles []string // Just the file names we expect
	}{
		{
			name: "basic media mix",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/image.jpg"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"), // This doesn't need upload
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FilePath("tests/video.mp4"),
					},
					Thumb: FilePath("tests/image.jpg"),
				},
				&InputPaidMedia{
					Type: "video",
					Media: &InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:  "video",
							Media: FilePath("tests/video.mp4"),
						},
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedFiles: []string{
				"file-0", "file-2", "file-2-thumb", "file-3", "file-3-thumb",
			},
		},
		{
			name:       "empty media array",
			inputMedia: []InputMedia{
				// Empty slice
			},
			expectedFiles: []string{},
		},
		{
			name: "only remote media (nothing to upload)",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FileID("existing_photo"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"),
					},
					Thumb: FileID("existing_thumb"),
				},
			},
			expectedFiles: []string{},
		},
		{
			name: "complex nested media with both local and remote files",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FileID("existing_photo"),
					},
				},
				&InputPaidMedia{
					Type: "video",
					Media: &InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:  "video",
							Media: FilePath("tests/video.mp4"),
						},
						Thumb: FileID("existing_thumb"), // Remote thumb
					},
					Thumb: FilePath("tests/outer_thumb.jpg"), // Local thumb
				},
				&InputPaidMedia{
					Type: "photo",
					Media: &InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:  "photo",
							Media: FileURL("https://example.com/photo.jpg"), // Remote media
						},
					},
					Thumb: FilePath("tests/thumb.jpg"), // Local thumb
				},
			},
			expectedFiles: []string{
				"file-1", "file-1-thumb", "file-2-thumb",
			},
		},
		{
			name: "all media types needing upload",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/photo.jpg"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FilePath("tests/video.mp4"),
					},
					Thumb: FilePath("tests/video_thumb.jpg"),
				},
				&InputMediaAnimation{
					BaseInputMedia: BaseInputMedia{
						Type:  "animation",
						Media: FilePath("tests/animation.gif"),
					},
					Thumb: FilePath("tests/animation_thumb.jpg"),
				},
				&InputMediaAudio{
					BaseInputMedia: BaseInputMedia{
						Type:  "audio",
						Media: FilePath("tests/audio.mp3"),
					},
					Thumb: FilePath("tests/audio_thumb.jpg"),
				},
				&InputMediaDocument{
					BaseInputMedia: BaseInputMedia{
						Type:  "document",
						Media: FilePath("tests/document.pdf"),
					},
					Thumb: FilePath("tests/document_thumb.jpg"),
				},
			},
			expectedFiles: []string{
				"file-0",
				"file-1", "file-1-thumb",
				"file-2", "file-2-thumb",
				"file-3", "file-3-thumb",
				"file-4", "file-4-thumb",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := prepareInputMediaForFiles(tt.inputMedia)

			if len(files) != len(tt.expectedFiles) {
				t.Errorf("Expected %d files, got %d", len(tt.expectedFiles), len(files))
				t.Logf("Expected: %v", tt.expectedFiles)
				t.Logf("Got: %v", files)
			}

			// Check that all expected filenames are in the result
			fileMap := make(map[string]bool)
			for _, file := range files {
				fileMap[file.Name] = true
			}

			for _, expectedName := range tt.expectedFiles {
				if !fileMap[expectedName] {
					t.Errorf("Expected file %s not found in result", expectedName)
				}
			}
		})
	}
}

func TestCloneInputMedia(t *testing.T) {
	tests := []struct {
		name  string
		media InputMedia
	}{
		{
			name: "photo with caption and entities",
			media: &InputMediaPhoto{
				BaseInputMedia: BaseInputMedia{
					Type:    "photo",
					Media:   FilePath("tests/photo.jpg"),
					Caption: "Photo with *bold* text",
					CaptionEntities: []MessageEntity{
						{
							Type:   "bold",
							Offset: 11,
							Length: 4,
						},
					},
					ParseMode: "Markdown",
				},
			},
		},
		{
			name: "video with all fields",
			media: &InputMediaVideo{
				BaseInputMedia: BaseInputMedia{
					Type:    "video",
					Media:   FilePath("tests/video.mp4"),
					Caption: "Video caption",
				},
				Thumb:             FilePath("tests/thumb.jpg"),
				Width:             1280,
				Height:            720,
				Duration:          300,
				SupportsStreaming: true,
				HasSpoiler:        true,
			},
		},
		{
			name: "nested paid media with video",
			media: &InputPaidMedia{
				Type: "video",
				Media: &InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:    "video",
						Media:   FilePath("tests/video.mp4"),
						Caption: "Nested video caption",
					},
					Thumb:    FilePath("tests/inner_thumb.jpg"),
					Duration: 120,
				},
				Thumb:    FilePath("tests/outer_thumb.jpg"),
				Duration: 120,
				Width:    1920,
				Height:   1080,
			},
		},
		{
			name: "animation with spoiler",
			media: &InputMediaAnimation{
				BaseInputMedia: BaseInputMedia{
					Type:       "animation",
					Media:      FilePath("tests/animation.gif"),
					HasSpoiler: true,
				},
				Thumb:    FilePath("tests/anim_thumb.jpg"),
				Width:    480,
				Height:   320,
				Duration: 15,
			},
		},
		{
			name: "audio with metadata",
			media: &InputMediaAudio{
				BaseInputMedia: BaseInputMedia{
					Type:    "audio",
					Media:   FilePath("tests/audio.mp3"),
					Caption: "Audio track",
				},
				Thumb:     FilePath("tests/audio_thumb.jpg"),
				Duration:  240,
				Performer: "Test Artist",
				Title:     "Test Song",
			},
		},
		{
			name: "document with caption",
			media: &InputMediaDocument{
				BaseInputMedia: BaseInputMedia{
					Type:    "document",
					Media:   FilePath("tests/document.pdf"),
					Caption: "Document with caption",
				},
				Thumb:                       FilePath("tests/doc_thumb.jpg"),
				DisableContentTypeDetection: true,
			},
		},
		{
			name:  "nil media",
			media: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip this test case
			if tt.media == nil {
				clone := cloneInputMedia(tt.media)
				if clone != nil {
					t.Errorf("Expected nil clone for nil input")
				}
				return
			}

			clone := cloneInputMedia(tt.media)

			if clone == tt.media {
				t.Errorf("Clone is the same instance as original")
			}

			if fmt.Sprintf("%T", clone) != fmt.Sprintf("%T", tt.media) {
				t.Errorf("Clone type %T doesn't match original type %T", clone, tt.media)
			}

			// Test deep copy by modifying the clone
			switch m := clone.(type) {
			case *InputMediaPhoto:
				originalCaption := tt.media.(*InputMediaPhoto).Caption
				m.Caption = "Modified caption"
				if tt.media.(*InputMediaPhoto).Caption != originalCaption {
					t.Errorf("Original caption was modified, not a deep copy")
				}

			case *InputMediaVideo:
				originalCaption := tt.media.(*InputMediaVideo).Caption
				m.Caption = "Modified caption"
				if tt.media.(*InputMediaVideo).Caption != originalCaption {
					t.Errorf("Original caption was modified, not a deep copy")
				}

			case *InputPaidMedia:
				// For paid media, test the nested media is also cloned
				origMedia := tt.media.(*InputPaidMedia)
				if m.Media == origMedia.Media {
					t.Errorf("Nested media is not cloned")
				}

				// Modify the nested media caption if it's a video
				if v, ok := m.Media.(*InputMediaVideo); ok {
					original := origMedia.Media.(*InputMediaVideo)
					oldCaption := original.Caption
					v.Caption = "Modified nested caption"
					if original.Caption != oldCaption {
						t.Errorf("Original nested caption was modified, not a deep copy")
					}
				}
			}
		})
	}
}

func TestCloneMediaSlice(t *testing.T) {
	// Test that the cloning function creates proper deep copies
	inputMedia := []InputMedia{
		&InputMediaPhoto{
			BaseInputMedia: BaseInputMedia{
				Type:    "photo",
				Media:   FilePath("photo.jpg"),
				Caption: "Original photo caption",
			},
		},
		&InputMediaVideo{
			BaseInputMedia: BaseInputMedia{
				Type:    "video",
				Media:   FilePath("video.mp4"),
				Caption: "Original video caption",
			},
			Thumb:    FilePath("thumb.jpg"),
			Duration: 30,
		},
		&InputMediaAnimation{
			BaseInputMedia: BaseInputMedia{
				Type:    "animation",
				Media:   FilePath("anim.gif"),
				Caption: "Original animation caption",
			},
			Thumb: FilePath("anim-thumb.jpg"),
		},
		&InputMediaAudio{
			BaseInputMedia: BaseInputMedia{
				Type:    "audio",
				Media:   FilePath("audio.mp3"),
				Caption: "Original audio caption",
			},
			Thumb: FilePath("audio-thumb.jpg"),
		},
		&InputMediaDocument{
			BaseInputMedia: BaseInputMedia{
				Type:    "document",
				Media:   FilePath("doc.pdf"),
				Caption: "Original document caption",
			},
			Thumb: FilePath("doc-thumb.jpg"),
		},
		&InputPaidMedia{
			Type: "video",
			Media: &InputMediaVideo{
				BaseInputMedia: BaseInputMedia{
					Type:  "video",
					Media: FilePath("tests/video.mp4"),
				},
			},
			Thumb: FilePath("paid-thumb.jpg"),
		},
	}

	cloned := cloneMediaSlice(inputMedia)

	if len(cloned) != len(inputMedia) {
		t.Fatalf("Expected cloned slice length %d, got %d", len(inputMedia), len(cloned))
	}

	// Test that we have a deep copy (modifying one doesn't affect the other)
	for i := range inputMedia {
		if fmt.Sprintf("%T", cloned[i]) != fmt.Sprintf("%T", inputMedia[i]) {
			t.Errorf("Type mismatch at index %d: expected %T, got %T",
				i, inputMedia[i], cloned[i])
		}

		// Test deep copy by modifying caption in cloned and checking original
		switch media := cloned[i].(type) {
		case *InputMediaPhoto:
			originalMedia := inputMedia[i].(*InputMediaPhoto)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaVideo:
			originalMedia := inputMedia[i].(*InputMediaVideo)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaAnimation:
			originalMedia := inputMedia[i].(*InputMediaAnimation)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaAudio:
			originalMedia := inputMedia[i].(*InputMediaAudio)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaDocument:
			originalMedia := inputMedia[i].(*InputMediaDocument)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		}
	}
}

func TestPrepareInputProfilePhotoForParams(t *testing.T) {
	photo := &InputProfilePhotoStatic{
		Type:  "static",
		Photo: FilePath("tests/profile.jpg"),
	}

	prepared := prepareInputProfilePhotoForParams(photo)
	if prepared == nil {
		t.Fatalf("expected prepared profile photo")
	}
	if prepared == photo {
		t.Fatalf("expected prepared profile photo to be cloned")
	}
	if prepared.getMedia().SendData() != "attach://file-0" {
		t.Fatalf("expected attach reference, got %q", prepared.getMedia().SendData())
	}
	if _, ok := photo.Photo.(FilePath); !ok {
		t.Fatalf("expected original profile photo to remain unchanged")
	}
}

func TestPrepareInputProfilePhotoForFiles(t *testing.T) {
	files := prepareInputProfilePhotoForFiles(&InputProfilePhotoStatic{
		Type:  "static",
		Photo: FilePath("tests/profile.jpg"),
	})
	if len(files) != 1 {
		t.Fatalf("expected one file, got %d", len(files))
	}
	if files[0].Name != "file-0" {
		t.Fatalf("expected file name file-0, got %q", files[0].Name)
	}

	files = prepareInputProfilePhotoForFiles(&InputProfilePhotoStatic{
		Type:  "static",
		Photo: FileID("existing-profile"),
	})
	if len(files) != 0 {
		t.Fatalf("expected no files for existing file id, got %d", len(files))
	}
}

func TestPrepareInputStoryContentForParams(t *testing.T) {
	content := &InputStoryContentVideo{
		Type:  "video",
		Video: FilePath("tests/story.mp4"),
	}

	prepared := prepareInputStoryContentForParams(content)
	if prepared == nil {
		t.Fatalf("expected prepared story content")
	}
	if prepared == content {
		t.Fatalf("expected prepared story content to be cloned")
	}
	if prepared.getMedia().SendData() != "attach://file-0" {
		t.Fatalf("expected attach reference, got %q", prepared.getMedia().SendData())
	}
	if _, ok := content.Video.(FilePath); !ok {
		t.Fatalf("expected original story content to remain unchanged")
	}
}

func TestPrepareInputStoryContentForFiles(t *testing.T) {
	files := prepareInputStoryContentForFiles(&InputStoryContentVideo{
		Type:  "video",
		Video: FilePath("tests/story.mp4"),
	})
	if len(files) != 1 {
		t.Fatalf("expected one file, got %d", len(files))
	}
	if files[0].Name != "file-0" {
		t.Fatalf("expected file name file-0, got %q", files[0].Name)
	}

	files = prepareInputStoryContentForFiles(&InputStoryContentVideo{
		Type:  "video",
		Video: FileURL("https://example.com/story.mp4"),
	})
	if len(files) != 0 {
		t.Fatalf("expected no files for url source, got %d", len(files))
	}
}

func TestSetMyProfilePhotoConfigUploadSerialization(t *testing.T) {
	photo := &InputProfilePhotoStatic{
		Type:  "static",
		Photo: FilePath("tests/profile.jpg"),
	}
	config := SetMyProfilePhotoConfig{
		Photo: photo,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params error: %v", err)
	}
	if params["photo"] == "" {
		t.Fatalf("expected photo param to be present")
	}
	if !strings.Contains(params["photo"], "attach://file-0") {
		t.Fatalf("expected photo param to contain attach reference, got %q", params["photo"])
	}

	files := config.files()
	if len(files) != 1 || files[0].Name != "file-0" {
		t.Fatalf("unexpected files payload: %+v", files)
	}
	if _, ok := photo.Photo.(FilePath); !ok {
		t.Fatalf("expected original photo to remain unchanged")
	}
}

func TestPostStoryConfigUploadSerialization(t *testing.T) {
	content := &InputStoryContentPhoto{
		Type:  "photo",
		Photo: FilePath("tests/story.jpg"),
	}
	config := PostStoryConfig{
		BusinessConnectionID: "business",
		Content:              content,
		ActivePeriod:         86400,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("params error: %v", err)
	}
	if params["content"] == "" {
		t.Fatalf("expected content param to be present")
	}
	if !strings.Contains(params["content"], "attach://file-0") {
		t.Fatalf("expected content param to contain attach reference, got %q", params["content"])
	}

	files := config.files()
	if len(files) != 1 || files[0].Name != "file-0" {
		t.Fatalf("unexpected files payload: %+v", files)
	}
	if _, ok := content.Photo.(FilePath); !ok {
		t.Fatalf("expected original story content to remain unchanged")
	}
}

func TestAPIParityRegressionFixes(t *testing.T) {
	setGameScore := SetGameScoreConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: 123,
			},
			MessageID: 456,
		},
		UserID: 111,
		Score:  42,
		Force:  true,
	}
	params, err := setGameScore.params()
	if err != nil {
		t.Fatalf("setGameScore params error: %v", err)
	}
	if params["score"] != "42" {
		t.Fatalf("expected score param, got %#v", params)
	}
	if _, ok := params["scrore"]; ok {
		t.Fatalf("unexpected legacy typo key in params")
	}
	if params["force"] != "true" {
		t.Fatalf("expected force param, got %#v", params)
	}

	emojiStatus := SetUserEmojiStatusConfig{
		UserID:                    999,
		EmojiStatusExpirationDate: 123456,
	}
	params, err = emojiStatus.params()
	if err != nil {
		t.Fatalf("setUserEmojiStatus params error: %v", err)
	}
	if params["emoji_status_expiration_date"] != "123456" {
		t.Fatalf("expected emoji_status_expiration_date param, got %#v", params)
	}
	if _, ok := params["emoji_status_expiration_date\t"]; ok {
		t.Fatalf("unexpected malformed emoji status key in params")
	}

	sendGift := SendGiftConfig{
		UserID:        42,
		Chat:          ChatConfig{ChatID: 4242},
		GiftID:        "gift-id",
		Text:          "hello",
		TextParseMode: "MarkdownV2",
	}
	params, err = sendGift.params()
	if err != nil {
		t.Fatalf("sendGift params error: %v", err)
	}
	if params["text_parse_mode"] != "MarkdownV2" {
		t.Fatalf("expected text_parse_mode from TextParseMode field, got %#v", params)
	}

	if got := (ChatMemberCountConfig{}).method(); got != "getChatMemberCount" {
		t.Fatalf("expected getChatMemberCount method, got %q", got)
	}

	editText := EditMessageTextConfig{
		BaseEdit: BaseEdit{
			BaseChatMessage: BaseChatMessage{
				ChatConfig: ChatConfig{ChatID: 1},
				MessageID:  2,
			},
		},
		RichMessage: NewInputRichMessageMarkdown("**updated**"),
	}
	params, err = editText.params()
	if err != nil {
		t.Fatalf("editMessageText params error: %v", err)
	}
	if !strings.Contains(params["rich_message"], `"markdown":"**updated**"`) {
		t.Fatalf("expected rich_message param, got %#v", params)
	}
	if _, ok := params["text"]; ok {
		t.Fatalf("unexpected empty text param with rich_message: %#v", params)
	}

	paidPhotoMedia := NewInputMediaPhoto(FileID("paid-photo-id"))
	paid := NewInputPaidMediaPhoto(&paidPhotoMedia)
	paidMedia := NewPaidMedia(1, 10, &paid)
	paidMedia.Payload = "paid-payload"
	params, err = paidMedia.params()
	if err != nil {
		t.Fatalf("sendPaidMedia params error: %v", err)
	}
	if params["payload"] != "paid-payload" {
		t.Fatalf("expected payload param, got %#v", params)
	}

	customEmojiThumbnail := NewCustomEmojiStickerSetThumbnal("emoji_set", "custom-emoji-id")
	params, err = customEmojiThumbnail.params()
	if err != nil {
		t.Fatalf("setCustomEmojiStickerSetThumbnail params error: %v", err)
	}
	if params["custom_emoji_id"] != "custom-emoji-id" {
		t.Fatalf("expected custom_emoji_id param, got %#v", params)
	}
	if _, ok := params["position"]; ok {
		t.Fatalf("unexpected position param in custom emoji thumbnail params")
	}

	maskPosition := SetStickerMaskPositionConfig{
		Sticker:      "sticker-file-id",
		MaskPosition: &MaskPosition{Point: "forehead", XShift: 0.1, YShift: 0.2, Scale: 1.3},
	}
	params, err = maskPosition.params()
	if err != nil {
		t.Fatalf("setStickerMaskPosition params error: %v", err)
	}
	if !strings.Contains(params["mask_position"], `"point":"forehead"`) {
		t.Fatalf("expected mask_position param, got %#v", params)
	}
	if _, ok := params["keywords"]; ok {
		t.Fatalf("unexpected keywords param in sticker mask position params")
	}

	addSticker := AddStickerConfig{
		UserID: 1,
		Name:   "stickers",
		Sticker: InputSticker{
			Sticker:   RequestFile{Name: "sticker", Data: FileID("sticker-file-id")},
			Format:    "static",
			EmojiList: []string{"\u263a"},
		},
	}
	params, err = addSticker.params()
	if err != nil {
		t.Fatalf("addStickerToSet params error: %v", err)
	}
	if !strings.Contains(params["sticker"], `"sticker":"sticker-file-id"`) || strings.Contains(params["sticker"], `"Name"`) {
		t.Fatalf("expected InputSticker.sticker to serialize as a file string, got %#v", params["sticker"])
	}
	if files := addSticker.files(); len(files) != 0 {
		t.Fatalf("expected file_id sticker not to be treated as upload file, got %+v", files)
	}

	newStickerSet := NewStickerSetConfig{
		UserID: 1,
		Name:   "stickers",
		Title:  "Stickers",
		Stickers: []InputSticker{
			{
				Sticker:   RequestFile{Data: FilePath("tests/sticker.webp")},
				Format:    "static",
				EmojiList: []string{"\u263a"},
			},
		},
	}
	params, err = newStickerSet.params()
	if err != nil {
		t.Fatalf("createNewStickerSet params error: %v", err)
	}
	if !strings.Contains(params["stickers"], `"sticker":"attach://sticker-0"`) {
		t.Fatalf("expected uploaded InputSticker to use attach reference, got %#v", params["stickers"])
	}
	stickerFiles := newStickerSet.files()
	if len(stickerFiles) != 1 || stickerFiles[0].Name != "sticker-0" {
		t.Fatalf("expected sticker upload file field, got %+v", stickerFiles)
	}

	shipping := ShippingConfig{
		ShippingQueryID: "shipping-query",
		OK:              false,
		ErrorMessage:    "no shipping",
	}
	params, err = shipping.params()
	if err != nil {
		t.Fatalf("answerShippingQuery params error: %v", err)
	}
	if params["ok"] != "false" {
		t.Fatalf("expected required ok=false param, got %#v", params)
	}

	preCheckout := PreCheckoutConfig{
		PreCheckoutQueryID: "pre-checkout-query",
		OK:                 false,
		ErrorMessage:       "no payment",
	}
	params, err = preCheckout.params()
	if err != nil {
		t.Fatalf("answerPreCheckoutQuery params error: %v", err)
	}
	if params["ok"] != "false" {
		t.Fatalf("expected required ok=false param, got %#v", params)
	}

	starSubscription := EditUserStarSubscriptionConfig{
		UserID:                  1,
		TelegramPaymentChargeID: "charge-id",
		IsCanceled:              false,
	}
	params, err = starSubscription.params()
	if err != nil {
		t.Fatalf("editUserStarSubscription params error: %v", err)
	}
	if params["is_canceled"] != "false" {
		t.Fatalf("expected required is_canceled=false param, got %#v", params)
	}

	message := NewMessage(1, "hello")
	params, err = message.params()
	if err != nil {
		t.Fatalf("sendMessage params error: %v", err)
	}
	for _, key := range []string{"business_connection_id", "reply_parameters", "link_preview_options", "entities"} {
		if _, ok := params[key]; ok {
			t.Fatalf("unexpected empty %s param in sendMessage: %#v", key, params)
		}
	}

	message.BusinessConnectionID = "business-connection"
	message.ReplyParameters = ReplyParameters{MessageID: 123}
	message.LinkPreviewOptions = LinkPreviewOptions{IsDisabled: true}
	params, err = message.params()
	if err != nil {
		t.Fatalf("sendMessage params error: %v", err)
	}
	if params["business_connection_id"] != "business-connection" ||
		!strings.Contains(params["reply_parameters"], `"message_id":123`) ||
		!strings.Contains(params["link_preview_options"], `"is_disabled":true`) {
		t.Fatalf("expected explicit optional params in sendMessage: %#v", params)
	}

	draft := SendMessageDraftConfig{
		ChatConfig:          ChatConfig{ChatID: 1},
		DraftID:             2,
		ThinkingPlaceholder: true,
	}
	params, err = draft.params()
	if err != nil {
		t.Fatalf("sendMessageDraft params error: %v", err)
	}
	if text, ok := params["text"]; !ok || text != "" {
		t.Fatalf("expected explicit empty text for thinking placeholder, got %#v", params)
	}

	removeEmojiStatus := SetUserEmojiStatusConfig{UserID: 1, RemoveStatus: true}
	params, err = removeEmojiStatus.params()
	if err != nil {
		t.Fatalf("setUserEmojiStatus params error: %v", err)
	}
	if status, ok := params["emoji_status_custom_emoji_id"]; !ok || status != "" {
		t.Fatalf("expected explicit empty emoji status, got %#v", params)
	}

	dropThumbnail := SetCustomEmojiStickerSetThumbnailConfig{Name: "emoji_set", DropThumbnail: true}
	params, err = dropThumbnail.params()
	if err != nil {
		t.Fatalf("setCustomEmojiStickerSetThumbnail params error: %v", err)
	}
	if thumbnail, ok := params["custom_emoji_id"]; !ok || thumbnail != "" {
		t.Fatalf("expected explicit empty custom_emoji_id, got %#v", params)
	}

	removeForumIcon := EditForumTopicConfig{
		BaseForum: BaseForum{
			ChatConfig:      ChatConfig{ChatID: 1},
			MessageThreadID: 2,
		},
		RemoveIcon: true,
	}
	params, err = removeForumIcon.params()
	if err != nil {
		t.Fatalf("editForumTopic params error: %v", err)
	}
	if icon, ok := params["icon_custom_emoji_id"]; !ok || icon != "" {
		t.Fatalf("expected explicit empty icon_custom_emoji_id, got %#v", params)
	}

	removeName := SetMyNameConfig{LanguageCode: "en", RemoveName: true}
	params, err = removeName.params()
	if err != nil {
		t.Fatalf("setMyName params error: %v", err)
	}
	if name, ok := params["name"]; !ok || name != "" {
		t.Fatalf("expected explicit empty name, got %#v", params)
	}

	removeDescription := SetMyDescriptionConfig{LanguageCode: "en", RemoveDescription: true}
	params, err = removeDescription.params()
	if err != nil {
		t.Fatalf("setMyDescription params error: %v", err)
	}
	if description, ok := params["description"]; !ok || description != "" {
		t.Fatalf("expected explicit empty description, got %#v", params)
	}

	removeShortDescription := SetMyShortDescriptionConfig{LanguageCode: "en", RemoveShortDescription: true}
	params, err = removeShortDescription.params()
	if err != nil {
		t.Fatalf("setMyShortDescription params error: %v", err)
	}
	if description, ok := params["short_description"]; !ok || description != "" {
		t.Fatalf("expected explicit empty short_description, got %#v", params)
	}

	uploadSticker := UploadStickerConfig{
		UserID:        42,
		Sticker:       RequestFile{Name: "custom-name", Data: FileBytes{Name: "sticker.webp", Bytes: []byte("sticker")}},
		StickerFormat: "static",
	}
	files := uploadSticker.files()
	if len(files) != 1 || files[0].Name != "sticker" {
		t.Fatalf("expected uploadStickerFile sticker file field, got %+v", files)
	}
}

func TestMediaGroupConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         MediaGroupConfig
		expectedMethod string
		expectedParams map[string]string
		expectedFiles  int
	}{
		{
			name: "basic media group with photos",
			config: MediaGroupConfig{
				BaseChat: BaseChat{
					ChatConfig: ChatConfig{
						ChatID: 12345,
					},
				},
				Media: []InputMedia{
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:    "photo",
							Media:   FilePath("tests/image1.jpg"),
							Caption: "First photo",
						},
					},
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:    "photo",
							Media:   FilePath("tests/image2.jpg"),
							Caption: "Second photo",
						},
					},
				},
			},
			expectedMethod: "sendMediaGroup",
			expectedParams: map[string]string{
				"chat_id": "12345",
			},
			expectedFiles: 2,
		},
		{
			name: "media group with different media types",
			config: MediaGroupConfig{
				BaseChat: BaseChat{
					ChatConfig: ChatConfig{
						ChatID: 12345,
					},
				},
				Media: []InputMedia{
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:    "photo",
							Media:   FilePath("tests/image.jpg"),
							Caption: "A photo",
						},
					},
					&InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:    "video",
							Media:   FilePath("tests/video.mp4"),
							Caption: "A video",
						},
						Thumb: FilePath("tests/thumb.jpg"),
					},
					&InputMediaDocument{
						BaseInputMedia: BaseInputMedia{
							Type:    "document",
							Media:   FilePath("tests/document.pdf"),
							Caption: "A document",
						},
					},
				},
			},
			expectedMethod: "sendMediaGroup",
			expectedParams: map[string]string{
				"chat_id": "12345",
			},
			expectedFiles: 4, // 3 media files + 1 thumb
		},
		{
			name: "media group with mixed file sources",
			config: MediaGroupConfig{
				BaseChat: BaseChat{
					ChatConfig: ChatConfig{
						ChatID: 12345,
					},
				},
				Media: []InputMedia{
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:  "photo",
							Media: FileID("photo123"),
						},
					},
					&InputMediaVideo{
						BaseInputMedia: BaseInputMedia{
							Type:  "video",
							Media: FileURL("https://example.com/video.mp4"),
						},
					},
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:  "photo",
							Media: FilePath("tests/image.jpg"),
						},
					},
				},
			},
			expectedMethod: "sendMediaGroup",
			expectedParams: map[string]string{
				"chat_id": "12345",
			},
			expectedFiles: 1, // Only one local file to upload
		},
		{
			name: "media group with paid media",
			config: MediaGroupConfig{
				BaseChat: BaseChat{
					ChatConfig: ChatConfig{
						ChatID: 12345,
					},
				},
				Media: []InputMedia{
					&InputMediaPhoto{
						BaseInputMedia: BaseInputMedia{
							Type:  "photo",
							Media: FilePath("tests/image.jpg"),
						},
					},
					&InputPaidMedia{
						Type: "video",
						Media: &InputMediaVideo{
							BaseInputMedia: BaseInputMedia{
								Type:  "video",
								Media: FilePath("tests/paid_video.mp4"),
							},
						},
						Thumb: FilePath("tests/thumb.jpg"),
					},
				},
			},
			expectedMethod: "sendMediaGroup",
			expectedParams: map[string]string{
				"chat_id": "12345",
			},
			expectedFiles: 3, // 2 media files + 1 thumb
		},
		{
			name: "media group with multiple paid media",
			config: MediaGroupConfig{
				BaseChat: BaseChat{
					ChatConfig: ChatConfig{
						ChatID: 12345,
					},
				},
				Media: []InputMedia{
					&InputPaidMedia{
						Type: "photo",
						Media: &InputMediaPhoto{
							BaseInputMedia: BaseInputMedia{
								Type:  "photo",
								Media: FilePath("tests/paid_photo.jpg"),
							},
						},
					},
					&InputPaidMedia{
						Type: "video",
						Media: &InputMediaVideo{
							BaseInputMedia: BaseInputMedia{
								Type:  "video",
								Media: FilePath("tests/paid_video.mp4"),
							},
							Duration: 30,
							Width:    1280,
							Height:   720,
						},
						Thumb:    FilePath("tests/thumb.jpg"),
						Width:    1280,
						Height:   720,
						Duration: 30,
					},
				},
			},
			expectedMethod: "sendMediaGroup",
			expectedParams: map[string]string{
				"chat_id": "12345",
			},
			expectedFiles: 3, // 2 media files + 1 thumb
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := tt.config.method()
			if method != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, method)
			}

			params, err := tt.config.params()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for key, value := range tt.expectedParams {
				if params[key] != value {
					t.Errorf("Expected param %s to be %s, got %s", key, value, params[key])
				}
			}

			// Check that media field exists in params
			if _, ok := params["media"]; !ok {
				t.Error("Expected 'media' field in params")
			}

			files := tt.config.files()
			if len(files) != tt.expectedFiles {
				t.Errorf("Expected %d files, got %d", tt.expectedFiles, len(files))
			}

			// Verify that each file has a name
			for i, file := range files {
				if file.Name == "" {
					t.Errorf("File at index %d has empty name", i)
				}
			}
		})
	}
}
