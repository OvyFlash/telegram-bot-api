package tgbotapi

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRequestFileDataSources(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "upload-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}

	_, err = tmp.WriteString("temp-data")
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	err = tmp.Close()
	if err != nil {
		t.Fatalf("close temp file: %v", err)
	}

	cases := []struct {
		name         string
		data         RequestFileData
		expectUpload bool
		expectValue  string
	}{
		{
			name:         "bytes upload",
			data:         FileBytes{Name: "data.bin", Bytes: []byte("content")},
			expectUpload: true,
		},
		{
			name:         "reader upload",
			data:         FileReader{Name: "reader.dat", Reader: bytes.NewBufferString("stream")},
			expectUpload: true,
		},
		{
			name:         "path upload",
			data:         FilePath(tmp.Name()),
			expectUpload: true,
		},
		{
			name:        "remote url",
			data:        FileURL("https://example.com/demo"),
			expectValue: "https://example.com/demo",
		},
		{
			name:        "file id",
			data:        FileID("ABC123"),
			expectValue: "ABC123",
		},
		{
			name:        "attach value",
			data:        fileAttach("attach://demo"),
			expectValue: "attach://demo",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if tc.data.NeedsUpload() != tc.expectUpload {
				t.Fatalf("unexpected upload flag: %v", tc.data.NeedsUpload())
			}

			if tc.expectUpload {
				name, reader, err := tc.data.UploadData()
				if err != nil {
					t.Fatalf("upload error: %v", err)
				}

				if name == "" {
					t.Fatalf("expected upload name")
				}

				all, err := io.ReadAll(reader)
				if err != nil {
					t.Fatalf("read upload: %v", err)
				}

				if closer, ok := reader.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						t.Fatalf("close upload reader: %v", err)
					}
				}

				if len(all) == 0 {
					t.Fatalf("expected upload payload")
				}
			} else {
				if _, _, err := tc.data.UploadData(); err == nil {
					t.Fatalf("expected upload error")
				}

				value := tc.data.SendData()
				if value != tc.expectValue {
					t.Fatalf("unexpected value: %q", value)
				}
			}
		})
	}
}
