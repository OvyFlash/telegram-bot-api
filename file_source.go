package tgbotapi

import (
	"bytes"
	"errors"
	"io"
	"os"
)

var (
	errFileSourceNoUpload    = errors.New("file source does not support uploads")
	errFileSourceNoReference = errors.New("file source does not support reference values")
)

type fileSourceKind int

const (
	fileSourceUpload fileSourceKind = iota
	fileSourceFileID
	fileSourceURL
	fileSourceAttach
	fileSourceInline
)

type uploadDescriptor struct {
	name   string
	reader io.ReadCloser
}

type fileSource struct {
	kind        fileSourceKind
	uploadFn    func() (uploadDescriptor, error)
	referenceFn func() (string, error)
}

func (s fileSource) kindIsUpload() bool {
	return s.kind == fileSourceUpload
}

func (s fileSource) openUpload() (uploadDescriptor, error) {
	if s.uploadFn == nil {
		return uploadDescriptor{}, errFileSourceNoUpload
	}

	return s.uploadFn()
}

func (s fileSource) referenceValue() (string, error) {
	if s.referenceFn == nil {
		return "", errFileSourceNoReference
	}

	return s.referenceFn()
}

func newBytesSource(name string, data []byte) fileSource {
	return fileSource{
		kind: fileSourceUpload,
		uploadFn: func() (uploadDescriptor, error) {
			return uploadDescriptor{
				name:   name,
				reader: io.NopCloser(bytes.NewReader(data)),
			}, nil
		},
	}
}

func newReaderSource(name string, reader io.Reader) fileSource {
	return fileSource{
		kind: fileSourceUpload,
		uploadFn: func() (uploadDescriptor, error) {
			if rc, ok := reader.(io.ReadCloser); ok {
				return uploadDescriptor{name: name, reader: rc}, nil
			}

			return uploadDescriptor{
				name:   name,
				reader: io.NopCloser(reader),
			}, nil
		},
	}
}

func newPathSource(path string) fileSource {
	return fileSource{
		kind: fileSourceUpload,
		uploadFn: func() (uploadDescriptor, error) {
			handle, err := os.Open(path)
			if err != nil {
				return uploadDescriptor{}, err
			}

			return uploadDescriptor{
				name:   handle.Name(),
				reader: handle,
			}, nil
		},
	}
}

func newURLSource(raw string) fileSource {
	return fileSource{
		kind: fileSourceURL,
		referenceFn: func() (string, error) {
			return raw, nil
		},
	}
}

func newFileIDSource(id string) fileSource {
	return fileSource{
		kind: fileSourceFileID,
		referenceFn: func() (string, error) {
			return id, nil
		},
	}
}

func newAttachSource(value string) fileSource {
	return fileSource{
		kind: fileSourceAttach,
		referenceFn: func() (string, error) {
			return value, nil
		},
	}
}

type fileSourceProvider interface {
	descriptor() fileSource
}

func resolveRequestFileData(data RequestFileData) (fileSource, error) {
	if provider, ok := data.(fileSourceProvider); ok {
		return provider.descriptor(), nil
	}

	if data == nil {
		return fileSource{}, errors.New("file data is nil")
	}

	if data.NeedsUpload() {
		return fileSource{
			kind: fileSourceUpload,
			uploadFn: func() (uploadDescriptor, error) {
				name, reader, err := data.UploadData()
				if err != nil {
					return uploadDescriptor{}, err
				}

				if rc, ok := reader.(io.ReadCloser); ok {
					return uploadDescriptor{name: name, reader: rc}, nil
				}

				return uploadDescriptor{
					name:   name,
					reader: io.NopCloser(reader),
				}, nil
			},
		}, nil
	}

	value := data.SendData()

	return fileSource{
		kind: fileSourceInline,
		referenceFn: func() (string, error) {
			return value, nil
		},
	}, nil
}
