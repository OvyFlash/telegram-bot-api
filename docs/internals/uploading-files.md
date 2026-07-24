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
file fields, a `document` and a `thumbnail`. These are pretty straightforward.

Remembering that the `Fileable` interface only requires one method, the config
only declares which API fields may contain file data.

```go
func (config DocumentConfig) files() []RequestFile {
	return requestFiles(
		requestFile("document", config.File),
		requestFile("thumbnail", config.Thumb),
	)
}
```

The request path turns those declarations into an upload plan. `FileID`,
`FileURL`, and existing `attach://` references become regular form fields.
`FilePath`, `FileBytes`, and `FileReader` become multipart file parts.

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

When the library prepares dynamic objects, it clones the caller-provided media
first. Uploadable file fields in that clone are rewritten to `attach://file-%d`
references, while the upload plan records multipart parts under matching field
names. File IDs and URLs stay inside the JSON payload as normal string values.
This keeps the JSON params and multipart files synchronized without mutating
the original media objects.

## Nested Upload Objects

Upload-capable nested objects like `InputProfilePhoto`, `InputStoryContent`,
poll media, and paid media use the same upload plan:

- Uploaded payload reference in params: `attach://file-0`
- Multipart field name in files list: `file-0`
- Original caller object is cloned before rewrite, so the input struct is not mutated

### Example: setMyProfilePhoto

```go
photo := tgbotapi.InputProfilePhotoStatic{
	Type:  "static",
	Photo: tgbotapi.FilePath("tests/profile.jpg"),
}

cfg := tgbotapi.NewSetMyProfilePhoto(&photo)
params, _ := cfg.params() // params["photo"] contains attach://file-0
files := cfg.files()      // files[0].Name == "file-0"
```

### Example: postStory

```go
content := tgbotapi.InputStoryContentPhoto{
	Type:  "photo",
	Photo: tgbotapi.FilePath("tests/story.jpg"),
}

cfg := tgbotapi.NewPostStory("business-connection-id", &content, 86400)
params, _ := cfg.params() // params["content"] contains attach://file-0
files := cfg.files()      // files[0].Name == "file-0"
```
