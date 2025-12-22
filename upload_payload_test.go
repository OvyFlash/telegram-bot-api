package tgbotapi

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"testing"
)

func TestUploadPayloadBuilder(t *testing.T) {
	payload := newUploadPayload()

	payload.Add("photo", FileBytes{Name: "pic.jpg", Bytes: []byte("data")})
	payload.Add("thumb", FileID("file-id"))
	payload.AddUploadOnly("skip", FileID("unused"))

	if !payload.needsUpload() {
		t.Fatalf("expected upload payload to require upload")
	}

	files := payload.filesSlice()
	if len(files) != 1 {
		t.Fatalf("expected single upload file, got %d", len(files))
	}

	if files[0].Name != "photo" {
		t.Fatalf("unexpected upload field %q", files[0].Name)
	}

	params := payload.applyInline(nil)
	if params["thumb"] != "file-id" {
		t.Fatalf("expected inline thumb value, got %q", params["thumb"])
	}

	if _, ok := params["skip"]; ok {
		t.Fatalf("did not expect upload-only field in inline params")
	}
}

func TestPrepareInputMedia(t *testing.T) {
	photo := NewInputMediaPhoto(FileBytes{Name: "image.png", Bytes: []byte("media")})
	video := NewInputMediaVideo(FileBytes{Name: "video.mp4", Bytes: []byte("clip")})
	video.Thumb = FileBytes{Name: "thumb.jpg", Bytes: []byte("thumb")}

	prepared, payload := prepareInputMedia([]InputMedia{&photo, &video})

	if prepared[0].getMedia().SendData() != "attach://file-0" {
		t.Fatalf("unexpected media ref: %q", prepared[0].getMedia().SendData())
	}

	if prepared[1].getMedia().SendData() != "attach://file-1" {
		t.Fatalf("unexpected video ref: %q", prepared[1].getMedia().SendData())
	}

	if prepared[1].getThumb().SendData() != "attach://file-1-thumb" {
		t.Fatalf("unexpected thumb ref: %q", prepared[1].getThumb().SendData())
	}

	files := payload.filesSlice()
	if len(files) != 3 {
		t.Fatalf("expected 3 upload parts, got %d", len(files))
	}

	expectedNames := map[string]struct{}{
		"file-0":       {},
		"file-1":       {},
		"file-1-thumb": {},
	}

	for _, f := range files {
		if _, ok := expectedNames[f.Name]; !ok {
			t.Fatalf("unexpected upload name %q", f.Name)
		}
		delete(expectedNames, f.Name)
	}

	if len(expectedNames) != 0 {
		t.Fatalf("missing upload fields: %v", expectedNames)
	}
}

func TestBuildMultipartPayload(t *testing.T) {
	params := Params{
		"text": "hello",
	}

	files := []RequestFile{
		{
			Name: "photo",
			Data: FileBytes{Name: "img.jpg", Bytes: []byte("jpeg-data")},
		},
		{
			Name: "thumbnail",
			Data: FileID("remote-thumb"),
		},
	}

	payload, err := buildMultipartPayload(params, files)
	if err != nil {
		t.Fatalf("build multipart payload: %v", err)
	}

	mediaType, attrs, err := mime.ParseMediaType(payload.contentType)
	if err != nil {
		t.Fatalf("parse media type: %v", err)
	}

	if mediaType != "multipart/form-data" {
		t.Fatalf("unexpected media type %q", mediaType)
	}

	boundary := attrs["boundary"]
	if boundary == "" {
		t.Fatalf("missing boundary")
	}

	rawBody, err := io.ReadAll(payload.body)
	if err != nil {
		t.Fatalf("read payload body: %v", err)
	}

	if len(rawBody) == 0 {
		t.Fatalf("empty multipart payload")
	}

	reader := multipart.NewReader(bytes.NewReader(rawBody), boundary)

	var (
		foundText      bool
		foundThumb     bool
		foundUpload    bool
		uploadContents []byte
	)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read multipart: %v", err)
		}

		data, err := io.ReadAll(part)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("read part: %v", err)
		}

		switch part.FormName() {
		case "text":
			foundText = bytes.Equal(data, []byte("hello"))
		case "thumbnail":
			foundThumb = bytes.Equal(data, []byte("remote-thumb"))
		case "photo":
			foundUpload = part.FileName() == "img.jpg"
			uploadContents = append([]byte(nil), data...)
		}
	}

	if !foundText {
		t.Fatalf("missing form field value")
	}

	if !foundThumb {
		t.Fatalf("missing inline thumb value")
	}

	if !foundUpload {
		t.Fatalf("missing upload part")
	}

	if !bytes.Equal(uploadContents, []byte("jpeg-data")) {
		t.Fatalf("unexpected upload payload %q", string(uploadContents))
	}
}
