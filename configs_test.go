package tgbotapi

import (
	"fmt"
	"strings"
	"testing"
)

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
	for i := 0; i < len(inputMedia); i++ {
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
