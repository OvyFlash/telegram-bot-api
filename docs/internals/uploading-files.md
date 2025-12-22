# Uploading Files

To make files work as expected, there's a lot going on behind the scenes. Make
sure to read through the [Files](../getting-started/files.md) section in
Getting Started first as we'll be building on that information.

This section only talks about file uploading. For non-uploaded files such as
URLs and file IDs, you just need to pass a string.

## Fields

Let's start by talking about how the library represents files as part of a
Config.

### Static Fields

Most endpoints use static file fields. For example, `sendPhoto` expects a single
file named `photo`. All we have to do is set that single field with the correct
value (either a string or multipart file). Methods like `sendDocument` take two
file uploads, a `document` and a `thumb`. These are pretty straightforward.

Remembering that the `Fileable` interface only requires one method, we expose a
`filePayload` helper that declares the intent for each field. The helper uses a
builder that decides whether an entry becomes an inline parameter or a streamed
upload.

```go
func (config DocumentConfig) filePayload() uploadPayload {
	payload := newUploadPayload()
	payload.Add("document", config.File)
	payload.Add("thumbnail", config.Thumb)
	return payload
}

func (config DocumentConfig) files() []RequestFile {
	return config.filePayload().filesSlice()
}
```

Calling `payload.Add` automatically promotes remote references (for example
`FileID`) into inline params while keeping uploadable variants in the returned
slice. This keeps the configuration declarative and avoids the panic guards
that used to exist on individual file types.

Telegram also supports the `attach://` syntax (discussed more later) for
thumbnails, but there's no reason to make things more complicated.

### Dynamic Fields

Of course, not everything can be so simple. Methods like `sendMediaGroup`
can accept many files, and each file can have custom markup. Using a static
field isn't possible because we need to specify which field is attached to each
item. Telegram introduced the `attach://` syntax for this.

Let's follow through creating a new media group with string and file uploads.

First, we start by creating some `InputMediaPhoto`.

```go
photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath("tests/image.jpg"))
url := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg"))
```

This created a new `InputMediaPhoto` struct, with a type of `photo` and the
media interface that we specified.

We'll now create our media group with the photo and URL.

```go
mediaGroup := NewMediaGroup(ChatID, []interface{}{
    photo,
    url,
})
```

A `MediaGroupConfig` stores all the media in an array of interfaces. We now
have all the data we need to upload, but how do we figure out field names for
uploads? We didn't specify `attach://unique-file` anywhere.

When the library goes to upload the files, it materializes the payload builder
and inspects which entries require streaming. The prepared media is rewritten
to reference `attach://file-%d` slots, and the corresponding uploads are added
to the multipart request under matching names. Remote references stay in the
params map, so the calling code does not have to manually synchronize field
names or worry about transitioning between upload and reference modes.
