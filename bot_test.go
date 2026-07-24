package tgbotapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	ChatID                  = 76918703
	Channel                 = "@tgbotapitest"
	SupergroupChatID        = -1001120141283
	ReplyToMessageID        = 35
	ExistingPhotoFileID     = "AgACAgIAAxkDAAEBFUZhIALQ9pZN4BUe8ZSzUU_2foSo1AACnrMxG0BucEhezsBWOgcikQEAAwIAA20AAyAE"
	ExistingDocumentFileID  = "BQADAgADOQADjMcoCcioX1GrDvp3Ag"
	ExistingAudioFileID     = "BQADAgADRgADjMcoCdXg3lSIN49lAg"
	ExistingVoiceFileID     = "AwADAgADWQADjMcoCeul6r_q52IyAg"
	ExistingVideoFileID     = "BAADAgADZgADjMcoCav432kYe0FRAg"
	ExistingVideoNoteFileID = "DQADAgADdQAD70cQSUK41dLsRMqfAg"
	ExistingStickerFileID   = "BQADAgADcwADjMcoCbdl-6eB--YPAg"
)

type testLogger struct {
	t *testing.T
}

type fakeHTTPClient struct {
	do func(*http.Request) (*http.Response, error)
}

type recordingBotLogger struct {
	entries []string
}

type panicBotLogger struct{}

type recordingSlogHandler struct {
	records *[]slog.Record
}

func (c fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func (l *recordingBotLogger) Println(v ...any) {
	l.entries = append(l.entries, fmt.Sprintln(v...))
}

func (l *recordingBotLogger) Printf(format string, v ...any) {
	l.entries = append(l.entries, fmt.Sprintf(format, v...))
}

func (panicBotLogger) Println(v ...any) {
	panic(fmt.Sprintln(v...))
}

func (panicBotLogger) Printf(format string, v ...any) {
	panic(fmt.Sprintf(format, v...))
}

func (h recordingSlogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h recordingSlogHandler) Handle(_ context.Context, record slog.Record) error {
	*h.records = append(*h.records, record.Clone())
	return nil
}

func (h recordingSlogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h recordingSlogHandler) WithGroup(string) slog.Handler {
	return h
}

func newFakeBot(client HTTPClient) *BotAPI {
	return &BotAPI{
		Token:       "token",
		Client:      client,
		apiEndpoint: "https://example.com/bot%s/%s",
	}
}

func okGetMeResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"id":123,"is_bot":true,"first_name":"Test","username":"test_bot"}}`)),
	}
}

func okAPIResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`)),
	}
}

func (t testLogger) Println(v ...any) {
	t.t.Log(v...)
}

func (t testLogger) Printf(format string, v ...any) {
	t.t.Logf(format, v...)
}

func TestNewBotAPIWithOptionsConfiguresBot(t *testing.T) {
	var requests []string
	client := fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			requests = append(requests, req.URL.String())
			return okGetMeResponse(), nil
		},
	}

	bot, err := NewBotAPIWithOptions(
		"token",
		WithAPIEndpoint("https://api.example/bot%s/%s"),
		WithFileEndpoint("https://files.example/file/bot%s/%s"),
		WithHTTPClient(client),
		WithDebug(true),
		WithUpdatesBuffer(42),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}

	if len(requests) != 1 || requests[0] != "https://api.example/bottoken/getMe" {
		t.Fatalf("unexpected constructor requests: %#v", requests)
	}
	if bot.Self.UserName != "test_bot" {
		t.Fatalf("unexpected self user: %+v", bot.Self)
	}
	if !bot.Debug {
		t.Fatalf("expected debug to be enabled")
	}
	if bot.Buffer != 42 {
		t.Fatalf("expected buffer 42, got %d", bot.Buffer)
	}
	if got := bot.FileURL(File{FilePath: "docs/file.txt"}); got != "https://files.example/file/bottoken/docs/file.txt" {
		t.Fatalf("unexpected file URL: %q", got)
	}
}

func TestGetFileDirectURLUsesBotFileEndpoint(t *testing.T) {
	client := fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			switch {
			case strings.HasSuffix(req.URL.Path, "/getMe"):
				return okGetMeResponse(), nil
			case strings.HasSuffix(req.URL.Path, "/getFile"):
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"file_id":"file-id","file_path":"photos/file.jpg"}}`)),
				}, nil
			default:
				t.Fatalf("unexpected request URL: %s", req.URL.String())
				return nil, nil
			}
		},
	}
	bot, err := NewBotAPIWithOptions(
		"token",
		WithAPIEndpoint("https://api.example/bot%s/%s"),
		WithFileEndpoint("https://files.example/file/bot%s/%s"),
		WithHTTPClient(client),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}

	link, err := bot.GetFileDirectURL("file-id")
	if err != nil {
		t.Fatalf("get file URL: %v", err)
	}
	if link != "https://files.example/file/bottoken/photos/file.jpg" {
		t.Fatalf("unexpected file URL: %q", link)
	}

	if got := (&File{FilePath: "photos/file.jpg"}).Link("token"); got != "https://api.telegram.org/file/bottoken/photos/file.jpg" {
		t.Fatalf("legacy File.Link changed: %q", got)
	}
}

func TestBotUsesInstanceSlogLoggerForDebug(t *testing.T) {
	previousLogger := log
	log = panicBotLogger{}
	t.Cleanup(func() {
		log = previousLogger
	})

	var records []slog.Record
	_, err := NewBotAPIWithOptions(
		"token",
		WithAPIEndpoint("https://api.example/bot%s/%s"),
		WithHTTPClient(fakeHTTPClient{do: func(*http.Request) (*http.Response, error) {
			return okGetMeResponse(), nil
		}}),
		WithDebug(true),
		WithLogger(slog.New(recordingSlogHandler{records: &records})),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("expected slog records")
	}
}

func TestBotCanDisableLogging(t *testing.T) {
	previousLogger := log
	log = panicBotLogger{}
	t.Cleanup(func() {
		log = previousLogger
	})

	bot, err := NewBotAPIWithOptions(
		"token",
		WithAPIEndpoint("https://api.example/bot%s/%s"),
		WithHTTPClient(fakeHTTPClient{do: func(*http.Request) (*http.Response, error) {
			return okGetMeResponse(), nil
		}}),
		WithDebug(true),
		WithLoggingDisabled(),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}

	bot.StopReceivingUpdates()
}

func TestBotUsesLegacyGlobalLoggerFallback(t *testing.T) {
	previousLogger := log
	logger := &recordingBotLogger{}
	log = logger
	t.Cleanup(func() {
		log = previousLogger
	})

	bot := newFakeBot(fakeHTTPClient{do: func(*http.Request) (*http.Response, error) {
		return okAPIResponse(), nil
	}})
	bot.Debug = true

	if _, err := bot.MakeRequest("getMe", nil); err != nil {
		t.Fatalf("make request: %v", err)
	}
	if len(logger.entries) == 0 {
		t.Fatalf("expected legacy logger entries")
	}
}

func TestRequestWithContextCarriesContextForUploads(t *testing.T) {
	type contextKey struct{}

	ctx := context.WithValue(context.Background(), contextKey{}, "upload-request")
	bot := newFakeBot(fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			if got := req.Context().Value(contextKey{}); got != "upload-request" {
				t.Fatalf("upload request lost context value: %#v", got)
			}
			return okAPIResponse(), nil
		},
	})

	config := NewPhoto(123, FileBytes{Name: "photo.jpg", Bytes: []byte("image")})

	if _, err := bot.RequestWithContext(ctx, config); err != nil {
		t.Fatalf("request with context failed: %v", err)
	}
}

func TestUploadFilesPreservesAPIErrorCode(t *testing.T) {
	bot := newFakeBot(fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"ok":false,"error_code":400,"description":"bad request"}`)),
			}, nil
		},
	})

	_, err := bot.UploadFiles("sendPhoto", Params{"chat_id": "1"}, []RequestFile{{
		Name: "photo",
		Data: FileBytes{Name: "photo.jpg", Bytes: []byte("image")},
	}})
	if err == nil {
		t.Fatalf("expected API error")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.Code != http.StatusBadRequest {
		t.Fatalf("expected error code 400, got %d", apiErr.Code)
	}
}

func TestRequestSerializesInlineFileReferencesAsFormFields(t *testing.T) {
	bot := newFakeBot(fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			if got := req.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
				t.Fatalf("expected form request, got %q", got)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			values, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatalf("parse request body: %v", err)
			}
			if values.Get("chat_id") != "123" || values.Get("photo") != "existing-photo" || values.Get("caption") != "hello" {
				t.Fatalf("unexpected form values: %v", values)
			}

			return okAPIResponse(), nil
		},
	})

	config := NewPhoto(123, FileID("existing-photo"))
	config.Caption = "hello"

	if _, err := bot.Request(config); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestWriteToHTTPResponseIncludesInlineFileReferences(t *testing.T) {
	recorder := httptest.NewRecorder()
	config := NewPhoto(123, FileID("existing-photo"))
	config.Caption = "hello"

	if err := WriteToHTTPResponse(recorder, config); err != nil {
		t.Fatalf("write response failed: %v", err)
	}

	values, err := url.ParseQuery(recorder.Body.String())
	if err != nil {
		t.Fatalf("parse response body: %v", err)
	}
	if values.Get("method") != "sendPhoto" || values.Get("chat_id") != "123" || values.Get("photo") != "existing-photo" || values.Get("caption") != "hello" {
		t.Fatalf("unexpected response values: %v", values)
	}
}

func TestUploadFilesSerializesMixedMultipartPayload(t *testing.T) {
	bot := newFakeBot(fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			mediaType, attrs, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil {
				t.Fatalf("parse content type: %v", err)
			}
			if mediaType != "multipart/form-data" {
				t.Fatalf("expected multipart request, got %q", mediaType)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}

			reader := multipart.NewReader(bytes.NewReader(body), attrs["boundary"])
			parts := map[string]string{}
			filenames := map[string]string{}

			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("read multipart part: %v", err)
				}
				data, err := io.ReadAll(part)
				if err != nil {
					t.Fatalf("read multipart data: %v", err)
				}
				parts[part.FormName()] = string(data)
				filenames[part.FormName()] = part.FileName()
			}

			if parts["caption"] != "hello" || parts["thumbnail"] != "existing-thumb" || parts["photo"] != "image-bytes" {
				t.Fatalf("unexpected multipart fields: %#v", parts)
			}
			if filenames["photo"] != "photo.jpg" {
				t.Fatalf("unexpected upload filename: %#v", filenames)
			}

			return okAPIResponse(), nil
		},
	})

	_, err := bot.UploadFiles("sendPhoto", Params{"caption": "hello"}, []RequestFile{
		{Name: "photo", Data: FileBytes{Name: "photo.jpg", Bytes: []byte("image-bytes")}},
		{Name: "thumbnail", Data: FileID("existing-thumb")},
	})
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
}

func TestRequestMultipartIncludesVideoCover(t *testing.T) {
	bot := newFakeBot(fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			_, attrs, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil {
				t.Fatalf("parse content type: %v", err)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}

			reader := multipart.NewReader(bytes.NewReader(body), attrs["boundary"])
			names := map[string]bool{}

			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("read multipart part: %v", err)
				}
				names[part.FormName()] = true
			}

			for _, name := range []string{"video", "thumbnail", "cover"} {
				if !names[name] {
					t.Fatalf("missing multipart field %q in %#v", name, names)
				}
			}

			return okAPIResponse(), nil
		},
	})

	config := NewVideo(123, FileBytes{Name: "video.mp4", Bytes: []byte("video")})
	config.Thumb = FileBytes{Name: "thumb.jpg", Bytes: []byte("thumb")}
	config.Cover = FileBytes{Name: "cover.jpg", Bytes: []byte("cover")}

	if _, err := bot.Request(config); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestFileDataWrongDirectionPanics(t *testing.T) {
	uploadable := []RequestFileData{
		FileBytes{Name: "data.bin", Bytes: []byte("data")},
		FileReader{Name: "reader.bin", Reader: strings.NewReader("data")},
		FilePath("tests/image.jpg"),
	}
	for _, data := range uploadable {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%T SendData did not panic", data)
				}
			}()
			_ = data.SendData()
		}()
	}

	references := []RequestFileData{
		FileID("file-id"),
		FileURL("https://example.com/file.jpg"),
		fileAttach("attach://file-0"),
	}
	for _, data := range references {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%T UploadData did not panic", data)
				}
			}()
			_, _, _ = data.UploadData()
		}()
	}
}

func TestUploadPlanAppliesInlineFieldsAndFiles(t *testing.T) {
	plan := newUploadPlan()
	plan.AddField("photo", FileBytes{Name: "photo.jpg", Bytes: []byte("image")})
	plan.AddField("thumbnail", FileID("existing-thumb"))
	plan.AddUploadOnly("cover", FileID("existing-cover"))

	params := plan.Apply(Params{"caption": "hello"})
	if params["caption"] != "hello" || params["thumbnail"] != "existing-thumb" {
		t.Fatalf("unexpected params: %#v", params)
	}
	if _, ok := params["cover"]; ok {
		t.Fatalf("upload-only field leaked into params: %#v", params)
	}

	files := plan.Files()
	if len(files) != 1 || files[0].Name != "photo" {
		t.Fatalf("unexpected files: %+v", files)
	}
	if !plan.NeedsUpload() {
		t.Fatalf("expected upload plan to need upload")
	}
}

func TestUploadPlanPreparesNestedMediaOnce(t *testing.T) {
	livePhoto := NewInputMediaLivePhoto(FileBytes{Name: "live.mp4", Bytes: []byte("live")}, FileBytes{Name: "still.jpg", Bytes: []byte("still")})
	media := []InputMedia{&livePhoto}

	prepared, plan := prepareInputMediaUploadPlan(media, "file")
	if len(prepared) != 1 {
		t.Fatalf("unexpected prepared media: %+v", prepared)
	}

	data, err := json.Marshal(prepared[0])
	if err != nil {
		t.Fatalf("marshal prepared media: %v", err)
	}
	payload := string(data)
	if !strings.Contains(payload, `"media":"attach://file-0"`) || !strings.Contains(payload, `"photo":"attach://file-0-photo"`) {
		t.Fatalf("unexpected prepared payload: %s", payload)
	}

	files := plan.Files()
	if len(files) != 2 || files[0].Name != "file-0" || files[1].Name != "file-0-photo" {
		t.Fatalf("unexpected files: %+v", files)
	}
	if _, ok := livePhoto.Media.(FileBytes); !ok {
		t.Fatalf("original media was mutated")
	}
}

func getBot(t *testing.T) (*BotAPI, error) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Skip("TEST_TOKEN is not set")
	}
	bot, err := NewBotAPI(token)
	bot.Debug = true

	logger := testLogger{t}
	SetLogger(logger)

	if err != nil {
		t.Error(err)
	}

	return bot, err
}

func TestBotWithCustomBuffer(t *testing.T) {
	bot, _ := getBot(t)

	customValue := 200
	bot.SetUpdatesBuffer(customValue)

	assertEq(t, bot.Buffer, customValue)
}

func TestNewBotAPI_notoken(t *testing.T) {
	_, err := NewBotAPI("")

	if err == nil {
		t.Error(err)
	}
}

func TestGetUpdates(t *testing.T) {
	bot, _ := getBot(t)

	u := NewUpdate(0)

	_, err := bot.GetUpdates(u)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = ModeMarkdown
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessageReply(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ReplyParameters.MessageID = ReplyToMessageID
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessageForward(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewForward(ChatID, ChatID, ReplyToMessageID)
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestCopyMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	message, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}

	copyMessageConfig := NewCopyMessage(SupergroupChatID, message.Chat.ID, message.MessageID)
	messageID, err := bot.CopyMessage(copyMessageConfig)
	if err != nil {
		t.Error(err)
	}

	if messageID.MessageID == message.MessageID {
		t.Error("copied message ID was the same as original message")
	}
}

func TestSendWithNewPhoto(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhoto(ChatID, FilePath("tests/image.jpg"))
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoWithFileBytes(t *testing.T) {
	bot, _ := getBot(t)

	data, _ := os.ReadFile("tests/image.jpg")
	b := FileBytes{Name: "image.jpg", Bytes: data}

	msg := NewPhoto(ChatID, b)
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoWithFileReader(t *testing.T) {
	bot, _ := getBot(t)

	f, _ := os.Open("tests/image.jpg")
	reader := FileReader{Name: "image.jpg", Reader: f}

	msg := NewPhoto(ChatID, reader)
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewPhotoReply(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhoto(ChatID, FilePath("tests/image.jpg"))
	msg.ReplyParameters.MessageID = ReplyToMessageID

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendNewPhotoToChannel(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhotoToChannel(Channel, FilePath("tests/image.jpg"))
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendNewPhotoToChannelFileBytes(t *testing.T) {
	bot, _ := getBot(t)

	data, _ := os.ReadFile("tests/image.jpg")
	b := FileBytes{Name: "image.jpg", Bytes: data}

	msg := NewPhotoToChannel(Channel, b)
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendNewPhotoToChannelFileReader(t *testing.T) {
	bot, _ := getBot(t)

	f, _ := os.Open("tests/image.jpg")
	reader := FileReader{Name: "image.jpg", Reader: f}

	msg := NewPhotoToChannel(Channel, reader)
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendWithExistingPhoto(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewPhoto(ChatID, FileID(ExistingPhotoFileID))
	msg.Caption = "Test"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewDocument(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDocument(ChatID, FilePath("tests/image.jpg"))
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewDocumentAndThumb(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDocument(ChatID, FilePath("tests/voice.ogg"))
	msg.Thumb = FilePath("tests/image.jpg")
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingDocument(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDocument(ChatID, FileID(ExistingDocumentFileID))
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewAudio(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewAudio(ChatID, FilePath("tests/audio.mp3"))
	msg.Title = "TEST"
	msg.Duration = 10
	msg.Performer = "TEST"
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingAudio(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewAudio(ChatID, FileID(ExistingAudioFileID))
	msg.Title = "TEST"
	msg.Duration = 10
	msg.Performer = "TEST"

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVoice(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVoice(ChatID, FilePath("tests/voice.ogg"))
	msg.Duration = 10
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVoice(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVoice(ChatID, FileID(ExistingVoiceFileID))
	msg.Duration = 10
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithContact(t *testing.T) {
	bot, _ := getBot(t)

	contact := NewContact(ChatID, "5551234567", "Test")

	if _, err := bot.Send(contact); err != nil {
		t.Error(err)
	}
}

func TestSendWithLocation(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.Send(NewLocation(ChatID, 40, 40))
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithVenue(t *testing.T) {
	bot, _ := getBot(t)

	venue := NewVenue(ChatID, "A Test Location", "123 Test Street", 40, 40)

	if _, err := bot.Send(venue); err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVideo(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideo(ChatID, FilePath("tests/video.mp4"))
	msg.Duration = 10
	msg.Caption = "TEST"

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVideo(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideo(ChatID, FileID(ExistingVideoFileID))
	msg.Duration = 10
	msg.Caption = "TEST"

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewVideoNote(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoNote(ChatID, 240, FilePath("tests/videonote.mp4"))
	msg.Duration = 10

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingVideoNote(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewVideoNote(ChatID, 240, FileID(ExistingVideoNoteFileID))
	msg.Duration = 10

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewSticker(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewSticker(ChatID, FilePath("tests/image.jpg"))

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingSticker(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewSticker(ChatID, FileID(ExistingStickerFileID))

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithNewStickerAndKeyboardHide(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewSticker(ChatID, FilePath("tests/image.jpg"))
	msg.ReplyMarkup = ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithExistingStickerAndKeyboardHide(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewSticker(ChatID, FileID(ExistingStickerFileID))
	msg.ReplyMarkup = ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}

	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithDice(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDice(ChatID)
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSendWithDiceWithEmoji(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewDiceWithEmoji(ChatID, "🏀")
	_, err := bot.Send(msg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetFile(t *testing.T) {
	bot, _ := getBot(t)

	file := FileConfig{
		FileID: ExistingPhotoFileID,
	}

	_, err := bot.GetFile(file)
	if err != nil {
		t.Error(err)
	}
}

func TestSendChatConfig(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.Request(NewChatAction(ChatID, ChatTyping))
	if err != nil {
		t.Error(err)
	}
}

// TODO: identify why this isn't working
// func TestSendEditMessage(t *testing.T) {
// 	bot, _ := getBot(t)

// 	msg, err := bot.Send(NewMessage(ChatID, "Testing editing."))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	edit := EditMessageTextConfig{
// 		BaseEdit: BaseEdit{
// 			ChatID:    ChatID,
// 			MessageID: msg.MessageID,
// 		},
// 		Text: "Updated text.",
// 	}

// 	_, err = bot.Send(edit)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestGetUserProfilePhotos(t *testing.T) {
	bot, _ := getBot(t)

	_, err := bot.GetUserProfilePhotos(NewUserProfilePhotos(ChatID))
	if err != nil {
		t.Error(err)
	}
}

func TestSetWebhookWithCert(t *testing.T) {
	bot, _ := getBot(t)

	time.Sleep(time.Second * 2)

	bot.Request(DeleteWebhookConfig{})

	wh, err := NewWebhookWithCert("https://example.com/tgbotapi-test/"+bot.Token, FilePath("tests/cert.pem"))
	if err != nil {
		t.Error(err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		t.Error(err)
	}

	_, err = bot.GetWebhookInfo()
	if err != nil {
		t.Error(err)
	}

	bot.Request(DeleteWebhookConfig{})
}

func TestSetWebhookWithoutCert(t *testing.T) {
	bot, _ := getBot(t)

	time.Sleep(time.Second * 2)

	bot.Request(DeleteWebhookConfig{})

	wh, err := NewWebhook("https://example.com/tgbotapi-test/" + bot.Token)
	if err != nil {
		t.Error(err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		t.Error(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		t.Error(err)
	}
	if info.MaxConnections == 0 {
		t.Errorf("Expected maximum connections to be greater than 0")
	}
	if info.LastErrorDate != 0 {
		t.Errorf("failed to set webhook: %s", info.LastErrorMessage)
	}

	bot.Request(DeleteWebhookConfig{})
}

func TestSendWithMediaGroupPhotoVideo(t *testing.T) {
	bot, _ := getBot(t)

	cfg := NewMediaGroup(ChatID, []InputMedia{
		ptr(NewInputMediaPhoto(FileURL("https://github.com/go-telegram-bot-api/telegram-bot-api/raw/0a3a1c8716c4cd8d26a262af9f12dcbab7f3f28c/tests/image.jpg"))),
		ptr(NewInputMediaPhoto(FilePath("tests/image.jpg"))),
		ptr(NewInputMediaVideo(FilePath("tests/video.mp4"))),
	})

	messages, err := bot.SendMediaGroup(cfg)
	if err != nil {
		t.Error(err)
	}

	if messages == nil {
		t.Error("No received messages")
	}

	if len(messages) != len(cfg.Media) {
		t.Errorf("Different number of messages: %d", len(messages))
	}
}

func TestSendWithMediaGroupDocument(t *testing.T) {
	bot, _ := getBot(t)

	cfg := NewMediaGroup(ChatID, []InputMedia{
		ptr(NewInputMediaDocument(FileURL("https://i.imgur.com/unQLJIb.jpg"))),
		ptr(NewInputMediaDocument(FilePath("tests/image.jpg"))),
	})

	messages, err := bot.SendMediaGroup(cfg)
	if err != nil {
		t.Error(err)
	}

	if messages == nil {
		t.Error("No received messages")
	}

	if len(messages) != len(cfg.Media) {
		t.Errorf("Different number of messages: %d", len(messages))
	}
}

func TestSendWithMediaGroupAudio(t *testing.T) {
	bot, _ := getBot(t)

	cfg := NewMediaGroup(ChatID, []InputMedia{
		ptr(NewInputMediaAudio(FilePath("tests/audio.mp3"))),
		ptr(NewInputMediaAudio(FilePath("tests/audio.mp3"))),
	})

	messages, err := bot.SendMediaGroup(cfg)
	if err != nil {
		t.Error(err)
	}

	if messages == nil {
		t.Error("No received messages")
	}

	if len(messages) != len(cfg.Media) {
		t.Errorf("Different number of messages: %d", len(messages))
	}
}

func TestDeleteMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(ChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = ModeMarkdown
	message, _ := bot.Send(msg)

	deleteMessageConfig := DeleteMessageConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: message.Chat.ID,
			},
			MessageID: message.MessageID,
		},
	}
	_, err := bot.Request(deleteMessageConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestPinChatMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(SupergroupChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = ModeMarkdown
	message, _ := bot.Send(msg)

	pinChatMessageConfig := PinChatMessageConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: ChatID,
			},
			MessageID: message.MessageID,
		},
		DisableNotification: false,
	}
	_, err := bot.Request(pinChatMessageConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestUnpinChatMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(SupergroupChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = ModeMarkdown
	message, _ := bot.Send(msg)

	// We need pin message to unpin something
	pinChatMessageConfig := PinChatMessageConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: message.Chat.ID,
			},
			MessageID: message.MessageID,
		},
		DisableNotification: false,
	}

	if _, err := bot.Request(pinChatMessageConfig); err != nil {
		t.Error(err)
	}

	unpinChatMessageConfig := UnpinChatMessageConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: message.Chat.ID,
			},
			MessageID: message.MessageID,
		},
	}

	if _, err := bot.Request(unpinChatMessageConfig); err != nil {
		t.Error(err)
	}
}

func TestUnpinAllChatMessages(t *testing.T) {
	bot, _ := getBot(t)

	msg := NewMessage(SupergroupChatID, "A test message from the test library in telegram-bot-api")
	msg.ParseMode = ModeMarkdown
	message, _ := bot.Send(msg)

	pinChatMessageConfig := PinChatMessageConfig{
		BaseChatMessage: BaseChatMessage{
			ChatConfig: ChatConfig{
				ChatID: message.Chat.ID,
			},
			MessageID: message.MessageID,
		},
		DisableNotification: true,
	}

	if _, err := bot.Request(pinChatMessageConfig); err != nil {
		t.Error(err)
	}

	unpinAllChatMessagesConfig := UnpinAllChatMessagesConfig{
		ChatConfig: ChatConfig{ChatID: message.Chat.ID},
	}

	if _, err := bot.Request(unpinAllChatMessagesConfig); err != nil {
		t.Error(err)
	}
}

func TestPolls(t *testing.T) {
	bot, _ := getBot(t)

	poll := NewPoll(SupergroupChatID, "Are polls working?", NewPollOption("Yes"), NewPollOption("No"))

	msg, err := bot.Send(poll)
	if err != nil {
		t.Error(err)
	}

	result, err := bot.StopPoll(NewStopPoll(SupergroupChatID, msg.MessageID))
	if err != nil {
		t.Error(err)
	}

	if result.Question != "Are polls working?" {
		t.Error("Poll question did not match")
	}

	if !result.IsClosed {
		t.Error("Poll did not end")
	}

	if result.Options[0].Text != "Yes" || result.Options[0].VoterCount != 0 || result.Options[1].Text != "No" || result.Options[1].VoterCount != 0 {
		t.Error("Poll options were incorrect")
	}
}

func TestSendDice(t *testing.T) {
	bot, _ := getBot(t)

	dice := NewDice(ChatID)

	msg, err := bot.Send(dice)
	if err != nil {
		t.Error("Unable to send dice roll")
	}

	if msg.Dice == nil {
		t.Error("Dice roll was not received")
	}
}

func TestCommands(t *testing.T) {
	bot, _ := getBot(t)

	setCommands := NewSetMyCommands(BotCommand{
		Command:     "test",
		Description: "a test command",
	})

	if _, err := bot.Request(setCommands); err != nil {
		t.Error("Unable to set commands")
	}

	commands, err := bot.GetMyCommands()
	if err != nil {
		t.Error("Unable to get commands")
	}

	if len(commands) != 1 {
		t.Error("Incorrect number of commands returned")
	}

	if commands[0].Command != "test" || commands[0].Description != "a test command" {
		t.Error("Commands were incorrectly set")
	}

	setCommands = NewSetMyCommandsWithScope(NewBotCommandScopeAllPrivateChats(), BotCommand{
		Command:     "private",
		Description: "a private command",
	})

	if _, err := bot.Request(setCommands); err != nil {
		t.Error("Unable to set commands")
	}

	commands, err = bot.GetMyCommandsWithConfig(NewGetMyCommandsWithScope(NewBotCommandScopeAllPrivateChats()))
	if err != nil {
		t.Error("Unable to get commands")
	}

	if len(commands) != 1 {
		t.Error("Incorrect number of commands returned")
	}

	if commands[0].Command != "private" || commands[0].Description != "a private command" {
		t.Error("Commands were incorrectly set")
	}
}
