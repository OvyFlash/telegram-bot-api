package tgbotapi

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
)

// Telegram constants
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// Constant values for ChatActions
const (
	ChatTyping          = "typing"
	ChatUploadPhoto     = "upload_photo"
	ChatRecordVideo     = "record_video"
	ChatUploadVideo     = "upload_video"
	ChatRecordVoice     = "record_voice"
	ChatUploadVoice     = "upload_voice"
	ChatUploadDocument  = "upload_document"
	ChatChooseSticker   = "choose_sticker"
	ChatFindLocation    = "find_location"
	ChatRecordVideoNote = "record_video_note"
	ChatUploadVideoNote = "upload_video_note"
)

// API errors
const (
	// ErrAPIForbidden happens when a token is bad
	ErrAPIForbidden = "forbidden"
)

// Constant values for ParseMode in MessageConfig
const (
	ModeMarkdown   = "Markdown"
	ModeMarkdownV2 = "MarkdownV2"
	ModeHTML       = "HTML"
)

// Constant values for update types
const (
	// UpdateTypeMessage is new incoming message of any kind — text, photo, sticker, etc.
	UpdateTypeMessage = "message"

	// UpdateTypeEditedMessage is new version of a message that is known to the bot and was edited
	UpdateTypeEditedMessage = "edited_message"

	// UpdateTypeChannelPost is new incoming channel post of any kind — text, photo, sticker, etc.
	UpdateTypeChannelPost = "channel_post"

	// UpdateTypeEditedChannelPost is new version of a channel post that is known to the bot and was edited
	UpdateTypeEditedChannelPost = "edited_channel_post"

	// UpdateTypeBusinessConnection is the bot was connected to or disconnected from a business account,
	// or a user edited an existing connection with the bot
	UpdateTypeBusinessConnection = "business_connection"

	// UpdateTypeBusinessMessage is a new non-service message from a connected business account
	UpdateTypeBusinessMessage = "business_message"

	// UpdateTypeEditedBusinessMessage is a new version of a message from a connected business account
	UpdateTypeEditedBusinessMessage = "edited_business_message"

	// UpdateTypeDeletedBusinessMessages are the messages were deleted from a connected business account
	UpdateTypeDeletedBusinessMessages = "deleted_business_messages"

	// UpdateTypeMessageReactionis is a reaction to a message was changed by a user
	UpdateTypeMessageReaction = "message_reaction"

	// UpdateTypeMessageReactionCount are reactions to a message with anonymous reactions were changed
	UpdateTypeMessageReactionCount = "message_reaction_count"

	// UpdateTypeInlineQuery is new incoming inline query
	UpdateTypeInlineQuery = "inline_query"

	// UpdateTypeChosenInlineResult i the result of an inline query that was chosen by a user and sent to their
	// chat partner. Please see the documentation on the feedback collecting for
	// details on how to enable these updates for your bot.
	UpdateTypeChosenInlineResult = "chosen_inline_result"

	// UpdateTypeCallbackQuery is new incoming callback query
	UpdateTypeCallbackQuery = "callback_query"

	// UpdateTypeShippingQuery is new incoming shipping query. Only for invoices with flexible price
	UpdateTypeShippingQuery = "shipping_query"

	// UpdateTypePreCheckoutQuery is new incoming pre-checkout query. Contains full information about checkout
	UpdateTypePreCheckoutQuery = "pre_checkout_query"

	// UpdateTypePurchasedPaidMedia is a user purchased paid media with a non-empty payload
	// sent by the bot in a non-channel chat
	UpdateTypePurchasedPaidMedia = "purchased_paid_media"

	// UpdateTypePoll is new poll state. Bots receive only updates about stopped polls and polls
	// which are sent by the bot
	UpdateTypePoll = "poll"

	// UpdateTypePollAnswer is when user changed their answer in a non-anonymous poll. Bots receive new votes
	// only in polls that were sent by the bot itself.
	UpdateTypePollAnswer = "poll_answer"

	// UpdateTypeMyChatMember is when the bot's chat member status was updated in a chat. For private chats, this
	// update is received only when the bot is blocked or unblocked by the user.
	UpdateTypeMyChatMember = "my_chat_member"

	// UpdateTypeChatMember is when the bot must be an administrator in the chat and must explicitly specify
	// this update in the list of allowed_updates to receive these updates.
	UpdateTypeChatMember = "chat_member"

	// UpdateTypeChatJoinRequest is request to join the chat has been sent.
	// The bot must have the can_invite_users administrator right in the chat to receive these updates.
	UpdateTypeChatJoinRequest = "chat_join_request"

	// UpdateTypeChatBoost is chat boost was added or changed.
	// The bot must be an administrator in the chat to receive these updates.
	UpdateTypeChatBoost = "chat_boost"

	// UpdateTypeRemovedChatBoost is boost was removed from a chat.
	// The bot must be an administrator in the chat to receive these updates.
	UpdateTypeRemovedChatBoost = "removed_chat_boost"
)

// Library errors
const (
	ErrBadURL = "bad or empty url"
)

// Chattable is any config type that can be sent.
type Chattable interface {
	params() (Params, error)
	method() string
}

// Fileable is any config type that can be sent that includes a file.
type Fileable interface {
	Chattable
	files() []RequestFile
}

// RequestFile represents a file associated with a field name.
type RequestFile struct {
	// The file field name.
	Name string
	// The file data to include.
	Data RequestFileData
}

// RequestFileData represents the data to be used for a file.
type RequestFileData interface {
	// NeedsUpload shows if the file needs to be uploaded.
	NeedsUpload() bool

	// UploadData gets the file name and an `io.Reader` for the file to be uploaded. This
	// must only be called when the file needs to be uploaded.
	UploadData() (string, io.Reader, error)
	// SendData gets the file data to send when a file does not need to be uploaded. This
	// must only be called when the file does not need to be uploaded.
	SendData() string
}

// FileBytes contains information about a set of bytes to upload
// as a File.
type FileBytes struct {
	Name  string
	Bytes []byte
}

func (fb FileBytes) NeedsUpload() bool {
	return true
}

func (fb FileBytes) UploadData() (string, io.Reader, error) {
	return fb.Name, bytes.NewReader(fb.Bytes), nil
}

func (fb FileBytes) SendData() string {
	panic("FileBytes must be uploaded")
}

// FileReader contains information about a reader to upload as a File.
type FileReader struct {
	Name   string
	Reader io.Reader
}

func (fr FileReader) NeedsUpload() bool {
	return true
}

func (fr FileReader) UploadData() (string, io.Reader, error) {
	return fr.Name, fr.Reader, nil
}

func (fr FileReader) SendData() string {
	panic("FileReader must be uploaded")
}

// FilePath is a path to a local file.
type FilePath string

func (fp FilePath) NeedsUpload() bool {
	return true
}

func (fp FilePath) UploadData() (string, io.Reader, error) {
	fileHandle, err := os.Open(string(fp))
	if err != nil {
		return "", nil, err
	}

	name := fileHandle.Name()
	return name, fileHandle, err
}

func (fp FilePath) SendData() string {
	panic("FilePath must be uploaded")
}

// FileURL is a URL to use as a file for a request.
type FileURL string

func (fu FileURL) NeedsUpload() bool {
	return false
}

func (fu FileURL) UploadData() (string, io.Reader, error) {
	panic("FileURL cannot be uploaded")
}

func (fu FileURL) SendData() string {
	return string(fu)
}

// FileID is an ID of a file already uploaded to Telegram.
type FileID string

func (fi FileID) NeedsUpload() bool {
	return false
}

func (fi FileID) UploadData() (string, io.Reader, error) {
	panic("FileID cannot be uploaded")
}

func (fi FileID) SendData() string {
	return string(fi)
}

// fileAttach is an internal file type used for processed media groups.
type fileAttach string

func (fa fileAttach) NeedsUpload() bool {
	return false
}

func (fa fileAttach) UploadData() (string, io.Reader, error) {
	panic("fileAttach cannot be uploaded")
}

func (fa fileAttach) SendData() string {
	return string(fa)
}

// LogOutConfig is a request to log out of the cloud Bot API server.
//
// Note that you may not log back in for at least 10 minutes.
type LogOutConfig struct{}

func (LogOutConfig) method() string {
	return "logOut"
}

func (LogOutConfig) params() (Params, error) {
	return nil, nil
}

// CloseConfig is a request to close the bot instance on a local server.
//
// Note that you may not close an instance for the first 10 minutes after the
// bot has started.
type CloseConfig struct{}

func (CloseConfig) method() string {
	return "close"
}

func (CloseConfig) params() (Params, error) {
	return nil, nil
}

// MessageConfig contains information about a SendMessage request.
type MessageConfig struct {
	BaseChat
	Text               string
	ParseMode          string
	Entities           []MessageEntity
	LinkPreviewOptions LinkPreviewOptions
}

func (config MessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("text", config.Text)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("entities", config.Entities)
	if err != nil {
		return params, err
	}
	err = params.AddInterface("link_preview_options", config.LinkPreviewOptions)

	return params, err
}

func (config MessageConfig) method() string {
	return "sendMessage"
}

// ForwardConfig contains information about a ForwardMessage request.
type ForwardConfig struct {
	BaseChat
	FromChat            ChatConfig
	MessageID           int // required
	VideoStartTimestamp int64
}

func (config ForwardConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
	p1, err := config.FromChat.paramsWithKey("from_chat_id")
	if err != nil {
		return params, err
	}
	params.Merge(p1)
	params.AddNonZero("message_id", config.MessageID)
	params.AddNonZero64("video_start_timestamp", config.VideoStartTimestamp)

	return params, nil
}

func (config ForwardConfig) method() string {
	return "forwardMessage"
}

// ForwardMessagesConfig contains information about a ForwardMessages request.
type ForwardMessagesConfig struct {
	BaseChat
	FromChat   ChatConfig
	MessageIDs []int // required
}

func (config ForwardMessagesConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	p1, err := config.FromChat.paramsWithKey("from_chat_id")
	if err != nil {
		return params, err
	}
	params.Merge(p1)
	err = params.AddInterface("message_ids", config.MessageIDs)

	return params, err
}

func (config ForwardMessagesConfig) method() string {
	return "forwardMessages"
}

// CopyMessageConfig contains information about a copyMessage request.
// Service messages, paid media messages, giveaway messages, giveaway winners messages, and invoice messages can't be copied.
type CopyMessageConfig struct {
	BaseChat
	FromChat              ChatConfig
	MessageID             int
	VideoStartTimestamp   int64
	Caption               string
	ParseMode             string
	CaptionEntities       []MessageEntity
	ShowCaptionAboveMedia bool
}

func (config CopyMessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	p1, err := config.FromChat.paramsWithKey("from_chat_id")
	if err != nil {
		return params, err
	}
	params.Merge(p1)
	params.AddNonZero("message_id", config.MessageID)
	params.AddNonZero64("video_start_timestamp", config.VideoStartTimestamp)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)
	err = params.AddInterface("caption_entities", config.CaptionEntities)

	return params, err
}

func (config CopyMessageConfig) method() string {
	return "copyMessage"
}

// CopyMessagesConfig contains information about a copyMessages request.
// Service messages, paid media messages, giveaway messages, giveaway winners messages, and invoice messages can't be copied.
type CopyMessagesConfig struct {
	BaseChat
	FromChat      ChatConfig
	MessageIDs    []int
	RemoveCaption bool
}

func (config CopyMessagesConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	p1, err := config.FromChat.paramsWithKey("from_chat_id")
	if err != nil {
		return params, err
	}
	params.Merge(p1)
	params.AddBool("remove_caption", config.RemoveCaption)
	err = params.AddInterface("message_ids", config.MessageIDs)

	return params, err
}

func (config CopyMessagesConfig) method() string {
	return "copyMessages"
}

// PhotoConfig contains information about a SendPhoto request.
type PhotoConfig struct {
	BaseFile
	BaseSpoiler
	Thumb                 RequestFileData
	Caption               string
	ParseMode             string
	CaptionEntities       []MessageEntity
	ShowCaptionAboveMedia bool
}

func (config PhotoConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)
	err = params.AddInterface("caption_entities", config.CaptionEntities)
	if err != nil {
		return params, err
	}

	p1, err := config.BaseSpoiler.params()
	if err != nil {
		return params, err
	}
	params.Merge(p1)

	return params, err
}

func (config PhotoConfig) method() string {
	return "sendPhoto"
}

func (config PhotoConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "photo",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// AudioConfig contains information about a SendAudio request.
type AudioConfig struct {
	BaseFile
	Thumb           RequestFileData
	Caption         string
	ParseMode       string
	CaptionEntities []MessageEntity
	Duration        int
	Performer       string
	Title           string
}

func (config AudioConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero("duration", config.Duration)
	params.AddNonEmpty("performer", config.Performer)
	params.AddNonEmpty("title", config.Title)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("caption_entities", config.CaptionEntities)

	return params, err
}

func (config AudioConfig) method() string {
	return "sendAudio"
}

func (config AudioConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "audio",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// DocumentConfig contains information about a SendDocument request.
type DocumentConfig struct {
	BaseFile
	Thumb                       RequestFileData
	Caption                     string
	ParseMode                   string
	CaptionEntities             []MessageEntity
	DisableContentTypeDetection bool
}

func (config DocumentConfig) params() (Params, error) {
	params, err := config.BaseFile.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("disable_content_type_detection", config.DisableContentTypeDetection)
	err = params.AddInterface("caption_entities", config.CaptionEntities)
	if err != nil {
		return params, err
	}

	return params, err
}

func (config DocumentConfig) method() string {
	return "sendDocument"
}

func (config DocumentConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "document",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// StickerConfig contains information about a SendSticker request.
type StickerConfig struct {
	// Emoji associated with the sticker; only for just uploaded stickers
	Emoji string
	BaseFile
}

func (config StickerConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("emoji", config.Emoji)
	return params, err
}

func (config StickerConfig) method() string {
	return "sendSticker"
}

func (config StickerConfig) files() []RequestFile {
	return []RequestFile{{
		Name: "sticker",
		Data: config.File,
	}}
}

// VideoConfig contains information about a SendVideo request.
type VideoConfig struct {
	BaseFile
	BaseSpoiler
	Thumb                 RequestFileData
	Duration              int
	Cover                 RequestFileData
	StartTimestamp        int64
	Caption               string
	ParseMode             string
	CaptionEntities       []MessageEntity
	ShowCaptionAboveMedia bool
	SupportsStreaming     bool
}

func (config VideoConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero("duration", config.Duration)
	params.AddNonZero64("start_timestamp", config.StartTimestamp)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("supports_streaming", config.SupportsStreaming)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)
	err = params.AddInterface("caption_entities", config.CaptionEntities)
	if err != nil {
		return params, err
	}

	p1, err := config.BaseSpoiler.params()
	if err != nil {
		return params, err
	}
	params.Merge(p1)

	return params, err
}

func (config VideoConfig) method() string {
	return "sendVideo"
}

func (config VideoConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "video",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	if config.Cover != nil {
		files = append(files, RequestFile{
			Name: "cover",
			Data: config.Cover,
		})
	}
	return files
}

// AnimationConfig contains information about a SendAnimation request.
type AnimationConfig struct {
	BaseFile
	BaseSpoiler
	Duration              int
	Thumb                 RequestFileData
	Caption               string
	ParseMode             string
	CaptionEntities       []MessageEntity
	ShowCaptionAboveMedia bool
}

func (config AnimationConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero("duration", config.Duration)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)
	err = params.AddInterface("caption_entities", config.CaptionEntities)
	if err != nil {
		return params, err
	}

	p1, err := config.BaseSpoiler.params()
	if err != nil {
		return params, err
	}
	params.Merge(p1)

	return params, err
}

func (config AnimationConfig) method() string {
	return "sendAnimation"
}

func (config AnimationConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "animation",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// VideoNoteConfig contains information about a SendVideoNote request.
type VideoNoteConfig struct {
	BaseFile
	Thumb    RequestFileData
	Duration int
	Length   int
}

func (config VideoNoteConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params.AddNonZero("duration", config.Duration)
	params.AddNonZero("length", config.Length)

	return params, err
}

func (config VideoNoteConfig) method() string {
	return "sendVideoNote"
}

func (config VideoNoteConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "video_note",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// Use this method to send paid media to channel chats. On success, the sent Message is returned.
type PaidMediaConfig struct {
	BaseChat
	StarCount             int64
	Media                 *InputPaidMedia
	Caption               string          // optional
	ParseMode             string          // optional
	CaptionEntities       []MessageEntity // optional
	ShowCaptionAboveMedia bool            // optional
}

func (config PaidMediaConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero64("star_count", config.StarCount)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)

	media := []InputMedia{config.Media}
	newMedia := prepareInputMediaForParams(media)
	err = params.AddInterface("media", newMedia[0])
	if err != nil {
		return params, err
	}
	err = params.AddInterface("caption_entities", config.CaptionEntities)
	return params, err
}

func (config PaidMediaConfig) files() []RequestFile {
	files := []RequestFile{}

	if config.Media.getMedia().NeedsUpload() {
		files = append(files, RequestFile{
			Name: "file-0",
			Data: config.Media.getMedia(),
		})
	}

	if thumb := config.Media.getThumb(); thumb != nil && thumb.NeedsUpload() {
		files = append(files, RequestFile{
			Name: "file-0-thumb",
			Data: thumb,
		})
	}

	return files
}

func (config PaidMediaConfig) method() string {
	return "sendPaidMedia"
}

// VoiceConfig contains information about a SendVoice request.
type VoiceConfig struct {
	BaseFile
	Thumb           RequestFileData
	Caption         string
	ParseMode       string
	CaptionEntities []MessageEntity
	Duration        int
}

func (config VoiceConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero("duration", config.Duration)
	params.AddNonEmpty("caption", config.Caption)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("caption_entities", config.CaptionEntities)

	return params, err
}

func (config VoiceConfig) method() string {
	return "sendVoice"
}

func (config VoiceConfig) files() []RequestFile {
	files := []RequestFile{{
		Name: "voice",
		Data: config.File,
	}}

	if config.Thumb != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumb,
		})
	}

	return files
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	BaseChat
	Latitude             float64 // required
	Longitude            float64 // required
	HorizontalAccuracy   float64 // optional
	LivePeriod           int     // optional
	Heading              int     // optional
	ProximityAlertRadius int     // optional
}

func (config LocationConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params.AddNonZeroFloat("horizontal_accuracy", config.HorizontalAccuracy)
	params.AddNonZero("live_period", config.LivePeriod)
	params.AddNonZero("heading", config.Heading)
	params.AddNonZero("proximity_alert_radius", config.ProximityAlertRadius)

	return params, err
}

func (config LocationConfig) method() string {
	return "sendLocation"
}

// EditMessageLiveLocationConfig allows you to update a live location.
type EditMessageLiveLocationConfig struct {
	BaseEdit
	Latitude             float64 // required
	Longitude            float64 // required
	LivePeriod           int     // optional
	HorizontalAccuracy   float64 // optional
	Heading              int     // optional
	ProximityAlertRadius int     // optional
}

func (config EditMessageLiveLocationConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params.AddNonZeroFloat("horizontal_accuracy", config.HorizontalAccuracy)
	params.AddNonZero("heading", config.Heading)
	params.AddNonZero("live_period", config.LivePeriod)
	params.AddNonZero("proximity_alert_radius", config.ProximityAlertRadius)

	return params, err
}

func (config EditMessageLiveLocationConfig) method() string {
	return "editMessageLiveLocation"
}

// StopMessageLiveLocationConfig stops updating a live location.
type StopMessageLiveLocationConfig struct {
	BaseEdit
}

func (config StopMessageLiveLocationConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (config StopMessageLiveLocationConfig) method() string {
	return "stopMessageLiveLocation"
}

// VenueConfig contains information about a SendVenue request.
type VenueConfig struct {
	BaseChat
	Latitude        float64 // required
	Longitude       float64 // required
	Title           string  // required
	Address         string  // required
	FoursquareID    string
	FoursquareType  string
	GooglePlaceID   string
	GooglePlaceType string
}

func (config VenueConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params["title"] = config.Title
	params["address"] = config.Address
	params.AddNonEmpty("foursquare_id", config.FoursquareID)
	params.AddNonEmpty("foursquare_type", config.FoursquareType)
	params.AddNonEmpty("google_place_id", config.GooglePlaceID)
	params.AddNonEmpty("google_place_type", config.GooglePlaceType)

	return params, err
}

func (config VenueConfig) method() string {
	return "sendVenue"
}

// ContactConfig allows you to send a contact.
type ContactConfig struct {
	BaseChat
	PhoneNumber string
	FirstName   string
	LastName    string
	VCard       string
}

func (config ContactConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params["phone_number"] = config.PhoneNumber
	params["first_name"] = config.FirstName

	params.AddNonEmpty("last_name", config.LastName)
	params.AddNonEmpty("vcard", config.VCard)

	return params, err
}

func (config ContactConfig) method() string {
	return "sendContact"
}

// SendPollConfig allows you to send a poll.
type SendPollConfig struct {
	BaseChat
	Question              string
	QuestionParseMode     string          // optional
	QuestionEntities      []MessageEntity // optional
	Options               []InputPollOption
	IsAnonymous           bool
	Type                  string
	AllowsMultipleAnswers bool
	CorrectOptionID       int64
	Explanation           string
	ExplanationParseMode  string
	ExplanationEntities   []MessageEntity
	OpenPeriod            int
	CloseDate             int
	IsClosed              bool
}

func (config SendPollConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params["question"] = config.Question
	params.AddNonEmpty("question_parse_mode", config.QuestionParseMode)
	if err = params.AddInterface("question_entities", config.QuestionEntities); err != nil {
		return params, err
	}
	if err = params.AddInterface("options", config.Options); err != nil {
		return params, err
	}
	params["is_anonymous"] = strconv.FormatBool(config.IsAnonymous)
	params.AddNonEmpty("type", config.Type)
	params["allows_multiple_answers"] = strconv.FormatBool(config.AllowsMultipleAnswers)
	params["correct_option_id"] = strconv.FormatInt(config.CorrectOptionID, 10)
	params.AddBool("is_closed", config.IsClosed)
	params.AddNonEmpty("explanation", config.Explanation)
	params.AddNonEmpty("explanation_parse_mode", config.ExplanationParseMode)
	params.AddNonZero("open_period", config.OpenPeriod)
	params.AddNonZero("close_date", config.CloseDate)
	err = params.AddInterface("explanation_entities", config.ExplanationEntities)

	return params, err
}

func (SendPollConfig) method() string {
	return "sendPoll"
}

// GameConfig allows you to send a game.
type GameConfig struct {
	BaseChat
	GameShortName string
}

func (config GameConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params["game_short_name"] = config.GameShortName

	return params, err
}

func (config GameConfig) method() string {
	return "sendGame"
}

// SetGameScoreConfig allows you to update the game score in a chat.
type SetGameScoreConfig struct {
	BaseChatMessage

	UserID             int64
	Score              int
	Force              bool
	DisableEditMessage bool
	InlineMessageID    string
}

func (config SetGameScoreConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params.AddNonZero("scrore", config.Score)
	params.AddBool("disable_edit_message", config.DisableEditMessage)

	if config.InlineMessageID != "" {
		params["inline_message_id"] = config.InlineMessageID
	} else {
		p1, err := config.BaseChatMessage.params()
		if err != nil {
			return params, err
		}
		params.Merge(p1)
	}

	return params, nil
}

func (config SetGameScoreConfig) method() string {
	return "setGameScore"
}

// GetGameHighScoresConfig allows you to fetch the high scores for a game.
type GetGameHighScoresConfig struct {
	BaseChatMessage

	UserID          int64
	InlineMessageID string
}

func (config GetGameHighScoresConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)

	if config.InlineMessageID != "" {
		params["inline_message_id"] = config.InlineMessageID
	} else {
		p1, err := config.BaseChatMessage.params()
		if err != nil {
			return params, err
		}
		params.Merge(p1)
	}

	return params, nil
}

func (config GetGameHighScoresConfig) method() string {
	return "getGameHighScores"
}

// ChatActionConfig contains information about a SendChatAction request.
type ChatActionConfig struct {
	BaseChat
	MessageThreadID int
	Action          string // required
}

func (config ChatActionConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params["action"] = config.Action
	params.AddNonZero("message_thread_id", config.MessageThreadID)

	return params, err
}

func (config ChatActionConfig) method() string {
	return "sendChatAction"
}

// EditMessageTextConfig allows you to modify the text in a message.
type EditMessageTextConfig struct {
	BaseEdit
	Text               string
	ParseMode          string
	Entities           []MessageEntity
	LinkPreviewOptions LinkPreviewOptions
}

func (config EditMessageTextConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()
	if err != nil {
		return params, err
	}

	params["text"] = config.Text
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("entities", config.Entities)
	if err != nil {
		return params, err
	}
	err = params.AddInterface("link_preview_options", config.LinkPreviewOptions)

	return params, err
}

func (config EditMessageTextConfig) method() string {
	return "editMessageText"
}

// EditMessageCaptionConfig allows you to modify the caption of a message.
type EditMessageCaptionConfig struct {
	BaseEdit
	Caption               string
	ParseMode             string
	CaptionEntities       []MessageEntity
	ShowCaptionAboveMedia bool
}

func (config EditMessageCaptionConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()
	if err != nil {
		return params, err
	}

	params["caption"] = config.Caption
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("show_caption_above_media", config.ShowCaptionAboveMedia)
	err = params.AddInterface("caption_entities", config.CaptionEntities)

	return params, err
}

func (config EditMessageCaptionConfig) method() string {
	return "editMessageCaption"
}

// EditMessageMediaConfig allows you to make an editMessageMedia request.
type EditMessageMediaConfig struct {
	BaseEdit

	Media InputMedia
}

func (EditMessageMediaConfig) method() string {
	return "editMessageMedia"
}

func (config EditMessageMediaConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()
	if err != nil {
		return params, err
	}

	preparedMedia := prepareInputMediaForParams([]InputMedia{config.Media})

	err = params.AddInterface("media", preparedMedia[0])

	return params, err
}

func (config EditMessageMediaConfig) files() []RequestFile {
	return prepareInputMediaForFiles([]InputMedia{config.Media})
}

// EditMessageReplyMarkupConfig allows you to modify the reply markup
// of a message.
type EditMessageReplyMarkupConfig struct {
	BaseEdit
}

func (config EditMessageReplyMarkupConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (config EditMessageReplyMarkupConfig) method() string {
	return "editMessageReplyMarkup"
}

// StopPollConfig allows you to stop a poll sent by the bot.
type StopPollConfig struct {
	BaseEdit
}

func (config StopPollConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (StopPollConfig) method() string {
	return "stopPoll"
}

// SetMessageReactionConfig changes reactions on a message. Returns true on success.
type SetMessageReactionConfig struct {
	BaseChatMessage
	Reaction []ReactionType
	IsBig    bool
}

func (config SetMessageReactionConfig) params() (Params, error) {
	params, err := config.BaseChatMessage.params()
	if err != nil {
		return params, err
	}
	params.AddBool("is_big", config.IsBig)
	err = params.AddInterface("reaction", config.Reaction)

	return params, err
}

func (SetMessageReactionConfig) method() string {
	return "setMessageReaction"
}

// UserProfilePhotosConfig contains information about a
// GetUserProfilePhotos request.
type UserProfilePhotosConfig struct {
	UserID int64
	Offset int
	Limit  int
}

func (UserProfilePhotosConfig) method() string {
	return "getUserProfilePhotos"
}

func (config UserProfilePhotosConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params.AddNonZero("offset", config.Offset)
	params.AddNonZero("limit", config.Limit)

	return params, nil
}

// SetUserEmojiStatusConfig changes the emoji status for a given user that
// previously allowed the bot to manage their emoji status via
// the Mini App method requestEmojiStatusAccess.
// Returns True on success.
type SetUserEmojiStatusConfig struct {
	UserID                    int64 // required
	EmojiStatusCustomEmojiID  string
	EmojiStatusExpirationDate int64
}

func (SetUserEmojiStatusConfig) method() string {
	return "setUserEmojiStatus"
}

func (config SetUserEmojiStatusConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params.AddNonEmpty("emoji_status_custom_emoji_id", config.EmojiStatusCustomEmojiID)
	params.AddNonZero64("emoji_status_expiration_date	", config.EmojiStatusExpirationDate)

	return params, nil
}

// FileConfig has information about a file hosted on Telegram.
type FileConfig struct {
	FileID string
}

func (FileConfig) method() string {
	return "getFile"
}

func (config FileConfig) params() (Params, error) {
	params := make(Params)

	params["file_id"] = config.FileID

	return params, nil
}

// UpdateConfig contains information about a GetUpdates request.
type UpdateConfig struct {
	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates []string
}

func (UpdateConfig) method() string {
	return "getUpdates"
}

func (config UpdateConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero("offset", config.Offset)
	params.AddNonZero("limit", config.Limit)
	params.AddNonZero("timeout", config.Timeout)
	params.AddInterface("allowed_updates", config.AllowedUpdates)

	return params, nil
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	URL                *url.URL
	Certificate        RequestFileData
	IPAddress          string
	MaxConnections     int
	AllowedUpdates     []string
	DropPendingUpdates bool
	SecretToken        string
}

func (config WebhookConfig) method() string {
	return "setWebhook"
}

func (config WebhookConfig) params() (Params, error) {
	params := make(Params)

	if config.URL != nil {
		params["url"] = config.URL.String()
	}

	params.AddNonEmpty("ip_address", config.IPAddress)
	params.AddNonZero("max_connections", config.MaxConnections)
	err := params.AddInterface("allowed_updates", config.AllowedUpdates)
	params.AddBool("drop_pending_updates", config.DropPendingUpdates)
	params.AddNonEmpty("secret_token", config.SecretToken)

	return params, err
}

func (config WebhookConfig) files() []RequestFile {
	if config.Certificate != nil {
		return []RequestFile{{
			Name: "certificate",
			Data: config.Certificate,
		}}
	}

	return nil
}

// DeleteWebhookConfig is a helper to delete a webhook.
type DeleteWebhookConfig struct {
	DropPendingUpdates bool
}

func (config DeleteWebhookConfig) method() string {
	return "deleteWebhook"
}

func (config DeleteWebhookConfig) params() (Params, error) {
	params := make(Params)

	params.AddBool("drop_pending_updates", config.DropPendingUpdates)

	return params, nil
}

// InlineQueryResultsButton represents a button to be shown above inline query results. You must use exactly one of the optional fields.
type InlineQueryResultsButton struct {
	// Label text on the button
	Text string `json:"text"`
	//Description of the Web App that will be launched when the user presses the button. The Web App will be able to switch back to the inline mode using the method switchInlineQuery inside the Web App.
	//
	//Optional
	WebApp *WebAppInfo `json:"web_app,omitempty"`
	// Deep-linking parameter for the /start message sent to the bot when a user presses the button. 1-64 characters, only A-Z, a-z, 0-9, _ and - are allowed.
	//
	//Optional
	StartParam string `json:"start_parameter,omitempty"`
}

// InlineConfig contains information on making an InlineQuery response.
type InlineConfig struct {
	InlineQueryID string                    `json:"inline_query_id"`
	Results       []interface{}             `json:"results"`
	CacheTime     int                       `json:"cache_time"`
	IsPersonal    bool                      `json:"is_personal"`
	NextOffset    string                    `json:"next_offset"`
	Button        *InlineQueryResultsButton `json:"button,omitempty"`
}

func (config InlineConfig) method() string {
	return "answerInlineQuery"
}

func (config InlineConfig) params() (Params, error) {
	params := make(Params)

	params["inline_query_id"] = config.InlineQueryID
	params.AddNonZero("cache_time", config.CacheTime)
	params.AddBool("is_personal", config.IsPersonal)
	params.AddNonEmpty("next_offset", config.NextOffset)
	err := params.AddInterface("button", config.Button)
	if err != nil {
		return params, err
	}
	err = params.AddInterface("results", config.Results)

	return params, err
}

// AnswerWebAppQueryConfig is used to set the result of an interaction with a
// Web App and send a corresponding message on behalf of the user to the chat
// from which the query originated.
type AnswerWebAppQueryConfig struct {
	// WebAppQueryID is the unique identifier for the query to be answered.
	WebAppQueryID string `json:"web_app_query_id"`
	// Result is an InlineQueryResult object describing the message to be sent.
	Result interface{} `json:"result"`
}

func (config AnswerWebAppQueryConfig) method() string {
	return "answerWebAppQuery"
}

func (config AnswerWebAppQueryConfig) params() (Params, error) {
	params := make(Params)

	params["web_app_query_id"] = config.WebAppQueryID
	err := params.AddInterface("result", config.Result)

	return params, err
}

// SavePreparedInlineMessageConfig stores a message that can be sent by a user of a Mini App.
// Returns a PreparedInlineMessage object.
type SavePreparedInlineMessageConfig[T InlineQueryResults] struct {
	UserID            int64 // required
	Result            T     // required
	AllowUserChats    bool
	AllowBotChats     bool
	AllowGroupChats   bool
	AllowChannelChats bool
}

func (config SavePreparedInlineMessageConfig[T]) method() string {
	return "savePreparedInlineMessage"
}

func (config SavePreparedInlineMessageConfig[T]) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	err := params.AddInterface("result", config.Result)
	if err != nil {
		return params, err
	}

	params.AddBool("allow_user_chats", config.AllowUserChats)
	params.AddBool("allow_bot_chats", config.AllowBotChats)
	params.AddBool("allow_group_chats", config.AllowGroupChats)
	params.AddBool("allow_channel_chats", config.AllowChannelChats)

	return params, err
}

// CallbackConfig contains information on making a CallbackQuery response.
type CallbackConfig struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
	URL             string `json:"url"`
	CacheTime       int    `json:"cache_time"`
}

func (config CallbackConfig) method() string {
	return "answerCallbackQuery"
}

func (config CallbackConfig) params() (Params, error) {
	params := make(Params)

	params["callback_query_id"] = config.CallbackQueryID
	params.AddNonEmpty("text", config.Text)
	params.AddBool("show_alert", config.ShowAlert)
	params.AddNonEmpty("url", config.URL)
	params.AddNonZero("cache_time", config.CacheTime)

	return params, nil
}

// ChatMemberConfig contains information about a user in a chat for use
// with administrative functions such as kicking or unbanning a user.
type ChatMemberConfig struct {
	ChatConfig
	UserID int64
}

func (config ChatMemberConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero64("user_id", config.UserID)
	return params, nil
}

// UnbanChatMemberConfig allows you to unban a user.
type UnbanChatMemberConfig struct {
	ChatMemberConfig
	OnlyIfBanned bool
}

func (config UnbanChatMemberConfig) method() string {
	return "unbanChatMember"
}

func (config UnbanChatMemberConfig) params() (Params, error) {
	params, err := config.ChatMemberConfig.params()
	if err != nil {
		return params, err
	}

	params.AddBool("only_if_banned", config.OnlyIfBanned)

	return params, nil
}

// BanChatMemberConfig contains extra fields to kick user.
type BanChatMemberConfig struct {
	ChatMemberConfig
	UntilDate      int64
	RevokeMessages bool
}

func (config BanChatMemberConfig) method() string {
	return "banChatMember"
}

func (config BanChatMemberConfig) params() (Params, error) {
	params, err := config.ChatMemberConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero64("until_date", config.UntilDate)
	params.AddBool("revoke_messages", config.RevokeMessages)

	return params, nil
}

// KickChatMemberConfig contains extra fields to ban user.
//
// This was renamed to BanChatMember in later versions of the Telegram Bot API.
type KickChatMemberConfig = BanChatMemberConfig

// RestrictChatMemberConfig contains fields to restrict members of chat
type RestrictChatMemberConfig struct {
	ChatMemberConfig
	UntilDate                     int64
	UseIndependentChatPermissions bool
	Permissions                   *ChatPermissions
}

func (config RestrictChatMemberConfig) method() string {
	return "restrictChatMember"
}

func (config RestrictChatMemberConfig) params() (Params, error) {
	params, err := config.ChatMemberConfig.params()
	if err != nil {
		return params, err
	}

	params.AddBool("use_independent_chat_permissions", config.UseIndependentChatPermissions)
	params.AddNonZero64("until_date", config.UntilDate)
	err = params.AddInterface("permissions", config.Permissions)

	return params, err
}

// PromoteChatMemberConfig contains fields to promote members of chat
type PromoteChatMemberConfig struct {
	ChatMemberConfig
	IsAnonymous         bool
	CanManageChat       bool
	CanChangeInfo       bool
	CanPostMessages     bool
	CanEditMessages     bool
	CanDeleteMessages   bool
	CanManageVideoChats bool
	CanInviteUsers      bool
	CanRestrictMembers  bool
	CanPinMessages      bool
	CanPromoteMembers   bool
	CanPostStories      bool
	CanEditStories      bool
	CanDeleteStories    bool
	CanManageTopics     bool
}

func (config PromoteChatMemberConfig) method() string {
	return "promoteChatMember"
}

func (config PromoteChatMemberConfig) params() (Params, error) {
	params, err := config.ChatMemberConfig.params()
	if err != nil {
		return params, err
	}

	params.AddBool("is_anonymous", config.IsAnonymous)
	params.AddBool("can_manage_chat", config.CanManageChat)
	params.AddBool("can_change_info", config.CanChangeInfo)
	params.AddBool("can_post_messages", config.CanPostMessages)
	params.AddBool("can_edit_messages", config.CanEditMessages)
	params.AddBool("can_delete_messages", config.CanDeleteMessages)
	params.AddBool("can_manage_video_chats", config.CanManageVideoChats)
	params.AddBool("can_invite_users", config.CanInviteUsers)
	params.AddBool("can_restrict_members", config.CanRestrictMembers)
	params.AddBool("can_pin_messages", config.CanPinMessages)
	params.AddBool("can_promote_members", config.CanPromoteMembers)
	params.AddBool("can_post_stories", config.CanPostStories)
	params.AddBool("can_edit_stories", config.CanEditStories)
	params.AddBool("can_delete_stories", config.CanDeleteStories)
	params.AddBool("can_manage_topics", config.CanManageTopics)

	return params, nil
}

// SetChatAdministratorCustomTitle sets the title of an administrative user
// promoted by the bot for a chat.
type SetChatAdministratorCustomTitle struct {
	ChatMemberConfig
	CustomTitle string
}

func (SetChatAdministratorCustomTitle) method() string {
	return "setChatAdministratorCustomTitle"
}

func (config SetChatAdministratorCustomTitle) params() (Params, error) {
	params, err := config.ChatMemberConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("custom_title", config.CustomTitle)

	return params, nil
}

// BanChatSenderChatConfig bans a channel chat in a supergroup or a channel. The
// owner of the chat will not be able to send messages and join live streams on
// behalf of the chat, unless it is unbanned first. The bot must be an
// administrator in the supergroup or channel for this to work and must have the
// appropriate administrator rights.
type BanChatSenderChatConfig struct {
	ChatConfig
	SenderChatID int64
	UntilDate    int
}

func (config BanChatSenderChatConfig) method() string {
	return "banChatSenderChat"
}

func (config BanChatSenderChatConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero64("sender_chat_id", config.SenderChatID)
	params.AddNonZero("until_date", config.UntilDate)

	return params, nil
}

// UnbanChatSenderChatConfig unbans a previously banned channel chat in a
// supergroup or channel. The bot must be an administrator for this to work and
// must have the appropriate administrator rights.
type UnbanChatSenderChatConfig struct {
	ChatConfig
	SenderChatID int64
}

func (config UnbanChatSenderChatConfig) method() string {
	return "unbanChatSenderChat"
}

func (config UnbanChatSenderChatConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero64("sender_chat_id", config.SenderChatID)

	return params, nil
}

// ChatInfoConfig contains information about getting chat information.
type ChatInfoConfig struct {
	ChatConfig
}

func (ChatInfoConfig) method() string {
	return "getChat"
}

// ChatMemberCountConfig contains information about getting the number of users in a chat.
type ChatMemberCountConfig struct {
	ChatConfig
}

func (ChatMemberCountConfig) method() string {
	return "getChatMembersCount"
}

// ChatAdministratorsConfig contains information about getting chat administrators.
type ChatAdministratorsConfig struct {
	ChatConfig
}

func (ChatAdministratorsConfig) method() string {
	return "getChatAdministrators"
}

// SetChatPermissionsConfig allows you to set default permissions for the
// members in a group. The bot must be an administrator and have rights to
// restrict members.
type SetChatPermissionsConfig struct {
	ChatConfig
	UseIndependentChatPermissions bool
	Permissions                   *ChatPermissions
}

func (SetChatPermissionsConfig) method() string {
	return "setChatPermissions"
}

func (config SetChatPermissionsConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddBool("use_independent_chat_permissions", config.UseIndependentChatPermissions)
	err = params.AddInterface("permissions", config.Permissions)

	return params, err
}

// ChatInviteLinkConfig contains information about getting a chat link.
//
// Note that generating a new link will revoke any previous links.
type ChatInviteLinkConfig struct {
	ChatConfig
}

func (ChatInviteLinkConfig) method() string {
	return "exportChatInviteLink"
}

func (config ChatInviteLinkConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

// CreateChatInviteLinkConfig allows you to create an additional invite link for
// a chat. The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights. The link can be revoked using the
// RevokeChatInviteLinkConfig.
type CreateChatInviteLinkConfig struct {
	ChatConfig
	Name               string
	ExpireDate         int
	MemberLimit        int
	CreatesJoinRequest bool
}

func (CreateChatInviteLinkConfig) method() string {
	return "createChatInviteLink"
}

func (config CreateChatInviteLinkConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("name", config.Name)
	params.AddNonZero("expire_date", config.ExpireDate)
	params.AddNonZero("member_limit", config.MemberLimit)
	params.AddBool("creates_join_request", config.CreatesJoinRequest)

	return params, nil
}

// EditChatInviteLinkConfig allows you to edit a non-primary invite link created
// by the bot. The bot must be an administrator in the chat for this to work and
// must have the appropriate admin rights.
type EditChatInviteLinkConfig struct {
	ChatConfig
	InviteLink         string
	Name               string
	ExpireDate         int
	MemberLimit        int
	CreatesJoinRequest bool
}

func (EditChatInviteLinkConfig) method() string {
	return "editChatInviteLink"
}

func (config EditChatInviteLinkConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("name", config.Name)
	params["invite_link"] = config.InviteLink
	params.AddNonZero("expire_date", config.ExpireDate)
	params.AddNonZero("member_limit", config.MemberLimit)
	params.AddBool("creates_join_request", config.CreatesJoinRequest)

	return params, nil
}

// CreateChatSubscriptionLinkConfig creates a subscription invite link for a channel chat.
// The bot must have the can_invite_users administrator rights.
// The link can be edited using the method editChatSubscriptionInviteLink or
// revoked using the method revokeChatInviteLink.
// Returns the new invite link as a ChatInviteLink object.
type CreateChatSubscriptionLinkConfig struct {
	ChatConfig
	Name               string
	SubscriptionPeriod int
	SubscriptionPrice  int
}

func (CreateChatSubscriptionLinkConfig) method() string {
	return "createChatSubscriptionInviteLink"
}

func (config CreateChatSubscriptionLinkConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("name", config.Name)
	params.AddNonZero("subscription_period", config.SubscriptionPeriod)
	params.AddNonZero("subscription_price", config.SubscriptionPrice)

	return params, nil
}

// EditChatSubscriptionLinkConfig edits a subscription invite link created by the bot.
// The bot must have the can_invite_users administrator rights.
// Returns the edited invite link as a ChatInviteLink object.
type EditChatSubscriptionLinkConfig struct {
	ChatConfig
	InviteLink string
	Name       string
}

func (EditChatSubscriptionLinkConfig) method() string {
	return "editChatSubscriptionInviteLink"
}

func (config EditChatSubscriptionLinkConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params["invite_link"] = config.InviteLink
	params.AddNonEmpty("name", config.Name)

	return params, nil
}

// RevokeChatInviteLinkConfig allows you to revoke an invite link created by the
// bot. If the primary link is revoked, a new link is automatically generated.
// The bot must be an administrator in the chat for this to work and must have
// the appropriate admin rights.
type RevokeChatInviteLinkConfig struct {
	ChatConfig
	InviteLink string
}

func (RevokeChatInviteLinkConfig) method() string {
	return "revokeChatInviteLink"
}

func (config RevokeChatInviteLinkConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params["invite_link"] = config.InviteLink

	return params, nil
}

// ApproveChatJoinRequestConfig allows you to approve a chat join request.
type ApproveChatJoinRequestConfig struct {
	ChatConfig
	UserID int64
}

func (ApproveChatJoinRequestConfig) method() string {
	return "approveChatJoinRequest"
}

func (config ApproveChatJoinRequestConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

// DeclineChatJoinRequest allows you to decline a chat join request.
type DeclineChatJoinRequest struct {
	ChatConfig
	UserID int64
}

func (DeclineChatJoinRequest) method() string {
	return "declineChatJoinRequest"
}

func (config DeclineChatJoinRequest) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

// LeaveChatConfig allows you to leave a chat.
type LeaveChatConfig struct {
	ChatConfig
}

func (config LeaveChatConfig) method() string {
	return "leaveChat"
}

func (config LeaveChatConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

// ChatConfigWithUser contains information about a chat and a user.
type ChatConfigWithUser struct {
	ChatConfig
	UserID int64
}

func (config ChatConfigWithUser) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

// GetChatMemberConfig is information about getting a specific member in a chat.
type GetChatMemberConfig struct {
	ChatConfigWithUser
}

func (GetChatMemberConfig) method() string {
	return "getChatMember"
}

// InvoiceConfig contains information for sendInvoice request.
type InvoiceConfig struct {
	BaseChat
	Title                     string         // required
	Description               string         // required
	Payload                   string         // required
	ProviderToken             string         // required
	Currency                  string         // required
	Prices                    []LabeledPrice // required
	MaxTipAmount              int
	SuggestedTipAmounts       []int
	StartParameter            string
	ProviderData              string
	PhotoURL                  string
	PhotoSize                 int
	PhotoWidth                int
	PhotoHeight               int
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
}

func (config InvoiceConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params["title"] = config.Title
	params["description"] = config.Description
	params["payload"] = config.Payload
	params["currency"] = config.Currency
	if err = params.AddInterface("prices", config.Prices); err != nil {
		return params, err
	}

	params.AddNonEmpty("provider_token", config.ProviderToken)
	params.AddNonZero("max_tip_amount", config.MaxTipAmount)
	if len(config.SuggestedTipAmounts) > 0 {
		err = params.AddInterface("suggested_tip_amounts", config.SuggestedTipAmounts)
		if err != nil {
			return params, err
		}
	}
	params.AddNonEmpty("start_parameter", config.StartParameter)
	params.AddNonEmpty("provider_data", config.ProviderData)
	params.AddNonEmpty("photo_url", config.PhotoURL)
	params.AddNonZero("photo_size", config.PhotoSize)
	params.AddNonZero("photo_width", config.PhotoWidth)
	params.AddNonZero("photo_height", config.PhotoHeight)
	params.AddBool("need_name", config.NeedName)
	params.AddBool("need_phone_number", config.NeedPhoneNumber)
	params.AddBool("need_email", config.NeedEmail)
	params.AddBool("need_shipping_address", config.NeedShippingAddress)
	params.AddBool("is_flexible", config.IsFlexible)
	params.AddBool("send_phone_number_to_provider", config.SendPhoneNumberToProvider)
	params.AddBool("send_email_to_provider", config.SendEmailToProvider)

	return params, err
}

func (config InvoiceConfig) method() string {
	return "sendInvoice"
}

// InvoiceLinkConfig contains information for createInvoiceLink method
type InvoiceLinkConfig struct {
	BusinessConnectionID      BusinessConnectionID
	Title                     string         // Required
	Description               string         // Required
	Payload                   string         // Required
	ProviderToken             string         // Required
	Currency                  string         // Required
	Prices                    []LabeledPrice // Required
	SubscriptionPeriod        int
	MaxTipAmount              int
	SuggestedTipAmounts       []int
	ProviderData              string
	PhotoURL                  string
	PhotoSize                 int
	PhotoWidth                int
	PhotoHeight               int
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
}

func (config InvoiceLinkConfig) params() (Params, error) {
	params, err := config.BusinessConnectionID.params()
	if err != nil {
		return params, err
	}

	params["title"] = config.Title
	params["description"] = config.Description
	params["payload"] = config.Payload
	params["currency"] = config.Currency
	if err := params.AddInterface("prices", config.Prices); err != nil {
		return params, err
	}

	params.AddNonZero("subscription_period", config.SubscriptionPeriod)
	params.AddNonEmpty("provider_token", config.ProviderToken)
	params.AddNonZero("max_tip_amount", config.MaxTipAmount)
	if len(config.SuggestedTipAmounts) > 0 {
		err := params.AddInterface("suggested_tip_amounts", config.SuggestedTipAmounts)
		if err != nil {
			return params, err
		}
	}
	params.AddNonEmpty("provider_data", config.ProviderData)
	params.AddNonEmpty("photo_url", config.PhotoURL)
	params.AddNonZero("photo_size", config.PhotoSize)
	params.AddNonZero("photo_width", config.PhotoWidth)
	params.AddNonZero("photo_height", config.PhotoHeight)
	params.AddBool("need_name", config.NeedName)
	params.AddBool("need_phone_number", config.NeedPhoneNumber)
	params.AddBool("need_email", config.NeedEmail)
	params.AddBool("need_shipping_address", config.NeedShippingAddress)
	params.AddBool("send_phone_number_to_provider", config.SendPhoneNumberToProvider)
	params.AddBool("send_email_to_provider", config.SendEmailToProvider)
	params.AddBool("is_flexible", config.IsFlexible)

	return params, nil
}

func (config InvoiceLinkConfig) method() string {
	return "createInvoiceLink"
}

// ShippingConfig contains information for answerShippingQuery request.
type ShippingConfig struct {
	ShippingQueryID string // required
	OK              bool   // required
	ShippingOptions []ShippingOption
	ErrorMessage    string
}

func (config ShippingConfig) method() string {
	return "answerShippingQuery"
}

func (config ShippingConfig) params() (Params, error) {
	params := make(Params)

	params["shipping_query_id"] = config.ShippingQueryID
	params.AddBool("ok", config.OK)
	err := params.AddInterface("shipping_options", config.ShippingOptions)
	params.AddNonEmpty("error_message", config.ErrorMessage)

	return params, err
}

// PreCheckoutConfig contains information for answerPreCheckoutQuery request.
type PreCheckoutConfig struct {
	PreCheckoutQueryID string // required
	OK                 bool   // required
	ErrorMessage       string
}

func (config PreCheckoutConfig) method() string {
	return "answerPreCheckoutQuery"
}

func (config PreCheckoutConfig) params() (Params, error) {
	params := make(Params)

	params["pre_checkout_query_id"] = config.PreCheckoutQueryID
	params.AddBool("ok", config.OK)
	params.AddNonEmpty("error_message", config.ErrorMessage)

	return params, nil
}

// Returns the bot's Telegram Star transactions in chronological order. On success, returns a StarTransactions object.
type GetStarTransactionsConfig struct {
	// Number of transactions to skip in the response
	Offset int64
	// The maximum number of transactions to be retrieved. Values between 1-100 are accepted. Defaults to 100.
	Limit int64
}

func (config GetStarTransactionsConfig) method() string {
	return "getStarTransactions"
}

func (config GetStarTransactionsConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("offset", config.Offset)
	params.AddNonZero64("limit", config.Limit)

	return params, nil
}

// RefundStarPaymentConfig refunds a successful payment in Telegram Stars.
// Returns True on success.
type RefundStarPaymentConfig struct {
	UserID                  int64  // required
	TelegramPaymentChargeID string // required
}

func (config RefundStarPaymentConfig) method() string {
	return "refundStarPayment"
}

func (config RefundStarPaymentConfig) params() (Params, error) {
	params := make(Params)

	params["telegram_payment_charge_id"] = config.TelegramPaymentChargeID
	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

// EditUserStarSubscriptionConfig allows the bot to cancel or re-enable extension
// of a subscription paid in Telegram Stars. Returns True on success.
type EditUserStarSubscriptionConfig struct {
	UserID                  int64  // required
	TelegramPaymentChargeID string // required
	IsCanceled              bool   // required
}

func (config EditUserStarSubscriptionConfig) method() string {
	return "editUserStarSubscription"
}

func (config EditUserStarSubscriptionConfig) params() (Params, error) {
	params := make(Params)

	params["telegram_payment_charge_id"] = config.TelegramPaymentChargeID
	params.AddNonZero64("user_id", config.UserID)
	params.AddBool("is_canceled", config.IsCanceled)

	return params, nil
}

// DeleteMessageConfig contains information of a message in a chat to delete.
type DeleteMessageConfig struct {
	BaseChatMessage
}

func (config DeleteMessageConfig) method() string {
	return "deleteMessage"
}

func (config DeleteMessageConfig) params() (Params, error) {
	return config.BaseChatMessage.params()
}

// DeleteMessageConfig contains information of a messages in a chat to delete.
type DeleteMessagesConfig struct {
	BaseChatMessages
}

func (config DeleteMessagesConfig) method() string {
	return "deleteMessages"
}

func (config DeleteMessagesConfig) params() (Params, error) {
	return config.BaseChatMessages.params()
}

// GetAvailableGiftsConfig returns the list of gifts that can be sent by the bot
// to users and channel chats. Requires no parameters. Returns a Gifts object.
type GetAvailableGiftsConfig struct{}

func (config GetAvailableGiftsConfig) method() string {
	return "getAvailableGifts"
}

func (config GetAvailableGiftsConfig) params() (Params, error) {
	return nil, nil
}

// SendGiftConfig sends a gift to the given user or channel chat.
// The gift can't be converted to Telegram Stars by the receiver.
// Returns True on success.
type SendGiftConfig struct {
	UserID        int64
	Chat          ChatConfig
	GiftID        string // required
	PayForUpgrade bool
	Text          string
	TextParseMode string
	TextEntities  []MessageEntity
}

func (config SendGiftConfig) method() string {
	return "sendGift"
}

func (config SendGiftConfig) params() (Params, error) {
	params := make(Params)
	params.AddNonZero64("user_id", config.UserID)

	p1, err := config.Chat.params()
	if err != nil {
		return params, err
	}
	params.Merge(p1)

	params.AddNonEmpty("gift_id", config.GiftID)
	params.AddBool("pay_for_upgrade", config.PayForUpgrade)
	params.AddNonEmpty("text", config.Text)
	params.AddNonEmpty("text_parse_mode", config.Text)
	params.AddInterface("text_entities", config.TextEntities)

	return params, nil
}

// VerifyUserConfig verifies a user on behalf of the organization
// which is represented by the bot.
// Returns True on success.
type VerifyUserConfig struct {
	UserID            int64 // required
	CustomDescription string
}

func (config VerifyUserConfig) method() string {
	return "verifyUser"
}

func (config VerifyUserConfig) params() (Params, error) {
	params := make(Params)
	params.AddNonZero64("user_id", config.UserID)
	params.AddNonEmpty("custom_description", config.CustomDescription)

	return params, nil
}

// VerifyChatConfig verifies a chat on behalf of the organization
// which is represented by the bot.
// Returns True on success.
type VerifyChatConfig struct {
	Chat              ChatConfig
	CustomDescription string
}

func (config VerifyChatConfig) method() string {
	return "verifyChat"
}

func (config VerifyChatConfig) params() (Params, error) {
	params, err := config.Chat.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("custom_description", config.CustomDescription)

	return params, nil
}

// RemoveUserVerificationConfig removes verification from a user who is currently
// verified on behalf of the organization represented by the bot.
// Returns True on success.
type RemoveUserVerificationConfig struct {
	UserID int64 // required
}

func (config RemoveUserVerificationConfig) method() string {
	return "removeUserVerification"
}

func (config RemoveUserVerificationConfig) params() (Params, error) {
	params := make(Params)
	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

// RemoveChatVerificationConfig removes verification from a chat who is currently
// verified on behalf of the organization represented by the bot.
// Returns True on success.
type RemoveChatVerificationConfig struct {
	Chat ChatConfig
}

func (config RemoveChatVerificationConfig) method() string {
	return "removeChatVerification"
}

func (config RemoveChatVerificationConfig) params() (Params, error) {
	return config.Chat.params()
}

// PinChatMessageConfig contains information of a message in a chat to pin.
type PinChatMessageConfig struct {
	BaseChatMessage
	DisableNotification bool
}

func (config PinChatMessageConfig) method() string {
	return "pinChatMessage"
}

func (config PinChatMessageConfig) params() (Params, error) {
	params, err := config.BaseChatMessage.params()
	if err != nil {
		return params, err
	}

	params.AddBool("disable_notification", config.DisableNotification)

	return params, nil
}

// UnpinChatMessageConfig contains information of a chat message to unpin.
//
// If MessageID is not specified, it will unpin the most recent pin.
type UnpinChatMessageConfig struct {
	BaseChatMessage
}

func (config UnpinChatMessageConfig) method() string {
	return "unpinChatMessage"
}

func (config UnpinChatMessageConfig) params() (Params, error) {
	return config.BaseChatMessage.params()
}

// UnpinAllChatMessagesConfig contains information of all messages to unpin in
// a chat.
type UnpinAllChatMessagesConfig struct {
	ChatConfig
}

func (config UnpinAllChatMessagesConfig) method() string {
	return "unpinAllChatMessages"
}

func (config UnpinAllChatMessagesConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

// SetChatPhotoConfig allows you to set a group, supergroup, or channel's photo.
type SetChatPhotoConfig struct {
	BaseFile
}

func (config SetChatPhotoConfig) method() string {
	return "setChatPhoto"
}

func (config SetChatPhotoConfig) files() []RequestFile {
	return []RequestFile{{
		Name: "photo",
		Data: config.File,
	}}
}

// DeleteChatPhotoConfig allows you to delete a group, supergroup, or channel's photo.
type DeleteChatPhotoConfig struct {
	ChatConfig
}

func (config DeleteChatPhotoConfig) method() string {
	return "deleteChatPhoto"
}

func (config DeleteChatPhotoConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

// SetChatTitleConfig allows you to set the title of something other than a private chat.
type SetChatTitleConfig struct {
	ChatConfig
	Title string
}

func (config SetChatTitleConfig) method() string {
	return "setChatTitle"
}

func (config SetChatTitleConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params["title"] = config.Title

	return params, nil
}

// SetChatDescriptionConfig allows you to set the description of a supergroup or channel.
type SetChatDescriptionConfig struct {
	ChatConfig
	Description string
}

func (config SetChatDescriptionConfig) method() string {
	return "setChatDescription"
}

func (config SetChatDescriptionConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params["description"] = config.Description

	return params, nil
}

// GetStickerSetConfig allows you to get the stickers in a set.
type GetStickerSetConfig struct {
	Name string
}

func (config GetStickerSetConfig) method() string {
	return "getStickerSet"
}

func (config GetStickerSetConfig) params() (Params, error) {
	params := make(Params)

	params["name"] = config.Name

	return params, nil
}

// GetCustomEmojiStickersConfig get information about
// custom emoji stickers by their identifiers.
type GetCustomEmojiStickersConfig struct {
	CustomEmojiIDs []string
}

func (config GetCustomEmojiStickersConfig) params() (Params, error) {
	params := make(Params)

	params.AddInterface("custom_emoji_ids", config.CustomEmojiIDs)

	return params, nil
}

func (config GetCustomEmojiStickersConfig) method() string {
	return "getCustomEmojiStickers"
}

// UploadStickerConfig allows you to upload a sticker for use in a set later.
type UploadStickerConfig struct {
	UserID        int64
	Sticker       RequestFile
	StickerFormat string
}

func (config UploadStickerConfig) method() string {
	return "uploadStickerFile"
}

func (config UploadStickerConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params["sticker_format"] = config.StickerFormat

	return params, nil
}

func (config UploadStickerConfig) files() []RequestFile {
	return []RequestFile{config.Sticker}
}

// NewStickerSetConfig allows creating a new sticker set.
type NewStickerSetConfig struct {
	UserID          int64
	Name            string
	Title           string
	Stickers        []InputSticker
	StickerType     string
	NeedsRepainting bool // optional; Pass True if stickers in the sticker set must be repainted to the color of text when used in messages, the accent color if used as emoji status, white on chat photos, or another appropriate color based on context; for custom emoji sticker sets only
}

func (config NewStickerSetConfig) method() string {
	return "createNewStickerSet"
}

func (config NewStickerSetConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params["name"] = config.Name
	params["title"] = config.Title

	params.AddBool("needs_repainting", config.NeedsRepainting)
	params.AddNonEmpty("sticker_type", string(config.StickerType))
	err := params.AddInterface("stickers", config.Stickers)

	return params, err
}

func (config NewStickerSetConfig) files() []RequestFile {
	requestFiles := []RequestFile{}
	for _, v := range config.Stickers {
		requestFiles = append(requestFiles, v.Sticker)
	}
	return requestFiles
}

// AddStickerConfig allows you to add a sticker to a set.
type AddStickerConfig struct {
	UserID  int64
	Name    string
	Sticker InputSticker
}

func (config AddStickerConfig) method() string {
	return "addStickerToSet"
}

func (config AddStickerConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params["name"] = config.Name
	err := params.AddInterface("sticker", config.Sticker)
	return params, err
}

func (config AddStickerConfig) files() []RequestFile {
	return []RequestFile{config.Sticker.Sticker}
}

// SetStickerPositionConfig allows you to change the position of a sticker in a set.
type SetStickerPositionConfig struct {
	Sticker  string
	Position int
}

func (config SetStickerPositionConfig) method() string {
	return "setStickerPositionInSet"
}

func (config SetStickerPositionConfig) params() (Params, error) {
	params := make(Params)

	params["sticker"] = config.Sticker
	params.AddNonZero("position", config.Position)

	return params, nil
}

// SetCustomEmojiStickerSetThumbnailConfig allows you to set the thumbnail of a custom emoji sticker set
type SetCustomEmojiStickerSetThumbnailConfig struct {
	Name          string
	CustomEmojiID string
}

func (config SetCustomEmojiStickerSetThumbnailConfig) method() string {
	return "setCustomEmojiStickerSetThumbnail"
}

func (config SetCustomEmojiStickerSetThumbnailConfig) params() (Params, error) {
	params := make(Params)

	params["name"] = config.Name
	params.AddNonEmpty("position", config.CustomEmojiID)

	return params, nil
}

// SetStickerSetTitle allows you to set the title of a created sticker set
type SetStickerSetTitleConfig struct {
	Name  string
	Title string
}

func (config SetStickerSetTitleConfig) method() string {
	return "setStickerSetTitle"
}

func (config SetStickerSetTitleConfig) params() (Params, error) {
	params := make(Params)

	params["name"] = config.Name
	params["title"] = config.Title

	return params, nil
}

// DeleteStickerSetConfig allows you to delete a sticker set that was created by the bot.
type DeleteStickerSetConfig struct {
	Name string
}

func (config DeleteStickerSetConfig) method() string {
	return "deleteStickerSet"
}

func (config DeleteStickerSetConfig) params() (Params, error) {
	params := make(Params)

	params["name"] = config.Name

	return params, nil
}

// DeleteStickerConfig allows you to delete a sticker from a set.
type DeleteStickerConfig struct {
	Sticker string
}

func (config DeleteStickerConfig) method() string {
	return "deleteStickerFromSet"
}

func (config DeleteStickerConfig) params() (Params, error) {
	params := make(Params)

	params["sticker"] = config.Sticker

	return params, nil
}

// ReplaceStickerInSetConfig allows you to replace an existing sticker in a sticker set
// with a new one. The method is equivalent to calling deleteStickerFromSet,
// then addStickerToSet, then setStickerPositionInSet.
// Returns True on success.
type ReplaceStickerInSetConfig struct {
	UserID     int64
	Name       string
	OldSticker string
	Sticker    InputSticker
}

func (config ReplaceStickerInSetConfig) method() string {
	return "replaceStickerInSet"
}

func (config ReplaceStickerInSetConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero64("user_id", config.UserID)
	params["name"] = config.Name
	params["old_sticker"] = config.OldSticker

	err := params.AddInterface("sticker", config.Sticker)

	return params, err
}

// SetStickerEmojiListConfig allows you to change the list of emoji assigned to a regular or custom emoji sticker. The sticker must belong to a sticker set created by the bot
type SetStickerEmojiListConfig struct {
	Sticker   string
	EmojiList []string
}

func (config SetStickerEmojiListConfig) method() string {
	return "setStickerEmojiList"
}

func (config SetStickerEmojiListConfig) params() (Params, error) {
	params := make(Params)

	params["sticker"] = config.Sticker
	err := params.AddInterface("emoji_list", config.EmojiList)

	return params, err
}

// SetStickerKeywordsConfig allows you to change search keywords assigned to a regular or custom emoji sticker. The sticker must belong to a sticker set created by the bot.
type SetStickerKeywordsConfig struct {
	Sticker  string
	Keywords []string
}

func (config SetStickerKeywordsConfig) method() string {
	return "setStickerKeywords"
}

func (config SetStickerKeywordsConfig) params() (Params, error) {
	params := make(Params)

	params["sticker"] = config.Sticker
	err := params.AddInterface("keywords", config.Keywords)

	return params, err
}

// SetStickerMaskPositionConfig allows you to  change the mask position of a mask sticker. The sticker must belong to a sticker set that was created by the bot
type SetStickerMaskPositionConfig struct {
	Sticker      string
	MaskPosition *MaskPosition
}

func (config SetStickerMaskPositionConfig) method() string {
	return "setStickerMaskPosition"
}

func (config SetStickerMaskPositionConfig) params() (Params, error) {
	params := make(Params)

	params["sticker"] = config.Sticker
	err := params.AddInterface("keywords", config.MaskPosition)

	return params, err
}

// SetStickerSetThumbConfig allows you to set the thumbnail for a sticker set.
type SetStickerSetThumbConfig struct {
	Name   string
	UserID int64
	Thumb  RequestFileData
	Format string
}

func (config SetStickerSetThumbConfig) method() string {
	return "setStickerSetThumbnail"
}

func (config SetStickerSetThumbConfig) params() (Params, error) {
	params := make(Params)

	params["name"] = config.Name
	params["format"] = config.Format

	params.AddNonZero64("user_id", config.UserID)

	return params, nil
}

func (config SetStickerSetThumbConfig) files() []RequestFile {
	return []RequestFile{{
		Name: "thumbnail",
		Data: config.Thumb,
	}}
}

// SetChatStickerSetConfig allows you to set the sticker set for a supergroup.
type SetChatStickerSetConfig struct {
	ChatConfig

	StickerSetName string
}

func (config SetChatStickerSetConfig) method() string {
	return "setChatStickerSet"
}

func (config SetChatStickerSetConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params["sticker_set_name"] = config.StickerSetName

	return params, nil
}

// DeleteChatStickerSetConfig allows you to remove a supergroup's sticker set.
type DeleteChatStickerSetConfig struct {
	ChatConfig
}

func (config DeleteChatStickerSetConfig) method() string {
	return "deleteChatStickerSet"
}

func (config DeleteChatStickerSetConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

// GetForumTopicIconStickersConfig allows you to get custom emoji stickers,
// which can be used as a forum topic icon by any user.
type GetForumTopicIconStickersConfig struct{}

func (config GetForumTopicIconStickersConfig) method() string {
	return "getForumTopicIconStickers"
}

func (config GetForumTopicIconStickersConfig) params() (Params, error) {
	return nil, nil
}

// CreateForumTopicConfig allows you to create a topic
// in a forum supergroup chat.
type CreateForumTopicConfig struct {
	ChatConfig
	Name              string
	IconColor         int
	IconCustomEmojiID string
}

func (config CreateForumTopicConfig) method() string {
	return "createForumTopic"
}

func (config CreateForumTopicConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("name", config.Name)
	params.AddNonZero("icon_color", config.IconColor)
	params.AddNonEmpty("icon_custom_emoji_id", config.IconCustomEmojiID)

	return params, nil
}

type BaseForum struct {
	ChatConfig
	MessageThreadID int
}

func (base BaseForum) params() (Params, error) {
	params, err := base.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero("message_thread_id", base.MessageThreadID)

	return params, nil
}

// EditForumTopicConfig allows you to edit
// name and icon of a topic in a forum supergroup chat.
type EditForumTopicConfig struct {
	BaseForum
	Name              string
	IconCustomEmojiID string
}

func (config EditForumTopicConfig) method() string {
	return "editForumTopic"
}

func (config EditForumTopicConfig) params() (Params, error) {
	params, err := config.BaseForum.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("name", config.Name)
	params.AddNonEmpty("icon_custom_emoji_id", config.IconCustomEmojiID)

	return params, nil
}

// CloseForumTopicConfig allows you to close
// an open topic in a forum supergroup chat.
type CloseForumTopicConfig struct{ BaseForum }

func (config CloseForumTopicConfig) method() string {
	return "closeForumTopic"
}

// ReopenForumTopicConfig allows you to reopen
// an closed topic in a forum supergroup chat.
type ReopenForumTopicConfig struct{ BaseForum }

func (config ReopenForumTopicConfig) method() string {
	return "reopenForumTopic"
}

// DeleteForumTopicConfig allows you to delete a forum topic
// along with all its messages in a forum supergroup chat.
type DeleteForumTopicConfig struct{ BaseForum }

func (config DeleteForumTopicConfig) method() string {
	return "deleteForumTopic"
}

// UnpinAllForumTopicMessagesConfig allows you to clear the list
// of pinned messages in a forum topic.
type UnpinAllForumTopicMessagesConfig struct{ BaseForum }

func (config UnpinAllForumTopicMessagesConfig) method() string {
	return "unpinAllForumTopicMessages"
}

// UnpinAllForumTopicMessagesConfig allows you to edit the name of
// the 'General' topic in a forum supergroup chat.
// The bot must be an administrator in the chat for this to work
// and must have can_manage_topics administrator rights. Returns True on success.
type EditGeneralForumTopicConfig struct {
	BaseForum
	Name string
}

func (config EditGeneralForumTopicConfig) method() string {
	return "editGeneralForumTopic"
}

func (config EditGeneralForumTopicConfig) params() (Params, error) {
	params, err := config.BaseForum.params()
	if err != nil {
		return params, err
	}
	params.AddNonEmpty("name", config.Name)

	return params, nil
}

// CloseGeneralForumTopicConfig allows you to to close an open 'General' topic
// in a forum supergroup chat. The bot must be an administrator in the chat
// for this to work and must have the can_manage_topics administrator rights.
// Returns True on success.
type CloseGeneralForumTopicConfig struct{ BaseForum }

func (config CloseGeneralForumTopicConfig) method() string {
	return "closeGeneralForumTopic"
}

// CloseGeneralForumTopicConfig allows you to reopen a closed 'General' topic
// in a forum supergroup chat. The bot must be an administrator in the chat
// for this to work and must have the can_manage_topics administrator rights.
// The topic will be automatically unhidden if it was hidden.
// Returns True on success.
type ReopenGeneralForumTopicConfig struct{ BaseForum }

func (config ReopenGeneralForumTopicConfig) method() string {
	return "reopenGeneralForumTopic"
}

// HideGeneralForumTopicConfig allows you to hide the 'General' topic
// in a forum supergroup chat. The bot must be an administrator in the chat
// for this to work and must have the can_manage_topics administrator rights.
// The topic will be automatically closed if it was open.
// Returns True on success.
type HideGeneralForumTopicConfig struct{ BaseForum }

func (config HideGeneralForumTopicConfig) method() string {
	return "hideGeneralForumTopic"
}

// UnhideGeneralForumTopicConfig allows you to unhide the 'General' topic
// in a forum supergroup chat. The bot must be an administrator in the chat
// for this to work and must have the can_manage_topics administrator rights.
// Returns True on success.
type UnhideGeneralForumTopicConfig struct{ BaseForum }

func (config UnhideGeneralForumTopicConfig) method() string {
	return "unhideGeneralForumTopic"
}

// UnpinAllGeneralForumTopicMessagesConfig allows you to to clear
// the list of pinned messages in a General forum topic.
// The bot must be an administrator in the chat for this to work
// and must have the can_pin_messages administrator right in the supergroup.
// Returns True on success.
type UnpinAllGeneralForumTopicMessagesConfig struct{ BaseForum }

func (config UnpinAllGeneralForumTopicMessagesConfig) method() string {
	return "unpinAllGeneralForumTopicMessages"
}

// MediaGroupConfig allows you to send a group of media.
//
// Media consist of InputMedia items (InputMediaPhoto, InputMediaVideo).
type MediaGroupConfig struct {
	BaseChat
	Media []InputMedia
}

func (config MediaGroupConfig) method() string {
	return "sendMediaGroup"
}

func (config MediaGroupConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return nil, err
	}

	err = params.AddInterface("media", prepareInputMediaForParams(config.Media))

	return params, err
}

func (config MediaGroupConfig) files() []RequestFile {
	return prepareInputMediaForFiles(config.Media)
}

// DiceConfig contains information about a sendDice request.
type DiceConfig struct {
	BaseChat
	// Emoji on which the dice throw animation is based.
	// Currently, must be one of 🎲, 🎯, 🏀, ⚽, 🎳, or 🎰.
	// Dice can have values 1-6 for 🎲, 🎯, and 🎳, values 1-5 for 🏀 and ⚽,
	// and values 1-64 for 🎰.
	// Defaults to "🎲"
	Emoji string
}

func (config DiceConfig) method() string {
	return "sendDice"
}

func (config DiceConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("emoji", config.Emoji)

	return params, err
}

type GetUserChatBoostsConfig struct {
	ChatConfig
	UserID int64
}

func (config GetUserChatBoostsConfig) method() string {
	return "getUserChatBoosts"
}

func (config GetUserChatBoostsConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero64("user_id", config.UserID)

	return params, err
}

type (
	GetBusinessConnectionConfig struct {
		BusinessConnectionID BusinessConnectionID
	}
	BusinessConnectionID string
)

func (GetBusinessConnectionConfig) method() string {
	return "getBusinessConnection"
}

func (config GetBusinessConnectionConfig) params() (Params, error) {
	return config.BusinessConnectionID.params()
}

func (config BusinessConnectionID) params() (Params, error) {
	params := make(Params)

	params["business_connection_id"] = string(config)

	return params, nil
}

// GetMyCommandsConfig gets a list of the currently registered commands.
type GetMyCommandsConfig struct {
	Scope        *BotCommandScope
	LanguageCode string
}

func (config GetMyCommandsConfig) method() string {
	return "getMyCommands"
}

func (config GetMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

// SetMyCommandsConfig sets a list of commands the bot understands.
type SetMyCommandsConfig struct {
	Commands     []BotCommand
	Scope        *BotCommandScope
	LanguageCode string
}

func (config SetMyCommandsConfig) method() string {
	return "setMyCommands"
}

func (config SetMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	if err := params.AddInterface("commands", config.Commands); err != nil {
		return params, err
	}
	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

type DeleteMyCommandsConfig struct {
	Scope        *BotCommandScope
	LanguageCode string
}

func (config DeleteMyCommandsConfig) method() string {
	return "deleteMyCommands"
}

func (config DeleteMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

// SetMyNameConfig change the bot's name
type SetMyNameConfig struct {
	Name         string
	LanguageCode string
}

func (config SetMyNameConfig) method() string {
	return "setMyName"
}

func (config SetMyNameConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("name", config.Name)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

type GetMyNameConfig struct {
	LanguageCode string
}

func (config GetMyNameConfig) method() string {
	return "getMyName"
}

func (config GetMyNameConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

// GetMyDescriptionConfig get the current bot description for the given user language
type GetMyDescriptionConfig struct {
	LanguageCode string
}

func (config GetMyDescriptionConfig) method() string {
	return "getMyDescription"
}

func (config GetMyDescriptionConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

// SetMyDescroptionConfig sets the bot's description, which is shown in the chat with the bot if the chat is empty
type SetMyDescriptionConfig struct {
	// Pass an empty string to remove the dedicated description for the given language.
	Description string
	// If empty, the description will be applied to all users for whose language there is no dedicated description.
	LanguageCode string
}

func (config SetMyDescriptionConfig) method() string {
	return "setMyDescription"
}

func (config SetMyDescriptionConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("description", config.Description)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

// GetMyShortDescriptionConfig get the current bot short description for the given user language
type GetMyShortDescriptionConfig struct {
	LanguageCode string
}

func (config GetMyShortDescriptionConfig) method() string {
	return "getMyShortDescription"
}

func (config GetMyShortDescriptionConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

// SetMyDescroptionConfig sets the bot's short description, which is shown on the bot's profile page and is sent together with the link when users share the bot.
type SetMyShortDescriptionConfig struct {
	// New short description for the bot; 0-120 characters.
	//
	//Pass an empty string to remove the dedicated short description for the given language.
	ShortDescription string
	//A two-letter ISO 639-1 language code.
	//
	//If empty, the short description will be applied to all users for whose language there is no dedicated short description.
	LanguageCode string
}

func (config SetMyShortDescriptionConfig) method() string {
	return "setMyShortDescription"
}

func (config SetMyShortDescriptionConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("short_description", config.ShortDescription)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, nil
}

// SetChatMenuButtonConfig changes the bot's menu button in a private chat,
// or the default menu button.
type SetChatMenuButtonConfig struct {
	ChatConfig

	MenuButton *MenuButton
}

func (config SetChatMenuButtonConfig) method() string {
	return "setChatMenuButton"
}

func (config SetChatMenuButtonConfig) params() (Params, error) {
	params, err := config.ChatConfig.params()
	if err != nil {
		return params, err
	}

	err = params.AddInterface("menu_button", config.MenuButton)

	return params, err
}

type GetChatMenuButtonConfig struct {
	ChatConfig
}

func (config GetChatMenuButtonConfig) method() string {
	return "getChatMenuButton"
}

func (config GetChatMenuButtonConfig) params() (Params, error) {
	return config.ChatConfig.params()
}

type SetMyDefaultAdministratorRightsConfig struct {
	Rights      ChatAdministratorRights
	ForChannels bool
}

func (config SetMyDefaultAdministratorRightsConfig) method() string {
	return "setMyDefaultAdministratorRights"
}

func (config SetMyDefaultAdministratorRightsConfig) params() (Params, error) {
	params := make(Params)

	err := params.AddInterface("rights", config.Rights)
	params.AddBool("for_channels", config.ForChannels)

	return params, err
}

type GetMyDefaultAdministratorRightsConfig struct {
	ForChannels bool
}

func (config GetMyDefaultAdministratorRightsConfig) method() string {
	return "getMyDefaultAdministratorRights"
}

func (config GetMyDefaultAdministratorRightsConfig) params() (Params, error) {
	params := make(Params)

	params.AddBool("for_channels", config.ForChannels)

	return params, nil
}

// prepareInputMediaForParams processes media items for API parameters.
// It creates a copy of the media array with files prepared for upload.
func prepareInputMediaForParams(inputMedia []InputMedia) []InputMedia {
	newMedias := cloneMediaSlice(inputMedia)
	for idx, media := range newMedias {
		if media.getMedia().NeedsUpload() {
			media.setUploadMedia(fmt.Sprintf("attach://file-%d", idx))
		}

		if thumb := media.getThumb(); thumb != nil && thumb.NeedsUpload() {
			media.setUploadThumb(fmt.Sprintf("attach://file-%d-thumb", idx))
		}

		newMedias[idx] = media
	}

	return newMedias
}

// prepareInputMediaForFiles generates RequestFile objects for media items
// that need to be uploaded.
func prepareInputMediaForFiles(inputMedia []InputMedia) []RequestFile {
	files := []RequestFile{}

	for idx, media := range inputMedia {
		if media.getMedia() != nil && media.getMedia().NeedsUpload() {
			files = append(files, RequestFile{
				Name: fmt.Sprintf("file-%d", idx),
				Data: media.getMedia(),
			})
		}

		if thumb := media.getThumb(); thumb != nil && thumb.NeedsUpload() {
			files = append(files, RequestFile{
				Name: fmt.Sprintf("file-%d-thumb", idx),
				Data: thumb,
			})
		}
	}

	return files
}

func ptr[T any](v T) *T {
	return &v
}

func cloneMediaSlice(media []InputMedia) []InputMedia {
	cloned := make([]InputMedia, len(media))
	for i, m := range media {
		cloned[i] = cloneInputMedia(m)
	}
	return cloned
}

func cloneInputMedia(media InputMedia) InputMedia {
	if media == nil {
		return nil
	}

	switch m := media.(type) {
	case *InputMediaPhoto:
		return ptr(*m)
	case *InputMediaVideo:
		return ptr(*m)
	case *InputMediaAnimation:
		return ptr(*m)
	case *InputMediaAudio:
		return ptr(*m)
	case *InputMediaDocument:
		return ptr(*m)
	case *InputPaidMedia:
		clone := &InputPaidMedia{
			Type:              m.Type,
			Thumb:             m.Thumb,
			Width:             m.Width,
			Height:            m.Height,
			Duration:          m.Duration,
			SupportsStreaming: m.SupportsStreaming,
		}
		if m.Media != nil {
			clone.Media = cloneInputMedia(m.Media)
		}
		return clone
	case *PaidMediaConfig:
		clone := &PaidMediaConfig{
			BaseChat:              m.BaseChat,
			StarCount:             m.StarCount,
			Caption:               m.Caption,
			ParseMode:             m.ParseMode,
			CaptionEntities:       m.CaptionEntities,
			ShowCaptionAboveMedia: m.ShowCaptionAboveMedia,
		}
		if m.Media != nil {
			clone.Media = &InputPaidMedia{
				Type:              m.Media.Type,
				Thumb:             m.Media.Thumb,
				Width:             m.Media.Width,
				Height:            m.Media.Height,
				Duration:          m.Media.Duration,
				SupportsStreaming: m.Media.SupportsStreaming,
			}
			if m.Media.Media != nil {
				clone.Media.Media = cloneInputMedia(m.Media.Media)
			}
		}
		return clone
	}
	return nil
}
