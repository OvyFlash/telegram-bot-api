package tgbotapi

type uploadPayload struct {
	files  []RequestFile
	inline map[string]string
}

func newUploadPayload() uploadPayload {
	return uploadPayload{
		inline: map[string]string{},
	}
}

func (p *uploadPayload) Add(field string, data RequestFileData) {
	if data == nil {
		return
	}

	source, err := resolveRequestFileData(data)
	if err != nil {
		return
	}

	if source.kindIsUpload() {
		p.files = append(p.files, RequestFile{
			Name: field,
			Data: data,
		})
		return
	}

	value, err := source.referenceValue()
	if err != nil {
		return
	}

	if p.inline == nil {
		p.inline = map[string]string{}
	}

	p.inline[field] = value
}

func (p *uploadPayload) AddUploadOnly(field string, data RequestFileData) {
	if data == nil {
		return
	}

	source, err := resolveRequestFileData(data)
	if err != nil {
		return
	}

	if source.kindIsUpload() {
		p.files = append(p.files, RequestFile{
			Name: field,
			Data: data,
		})
	}
}

func (p uploadPayload) needsUpload() bool {
	return len(p.files) > 0
}

func (p uploadPayload) filesSlice() []RequestFile {
	return p.files
}

func (p uploadPayload) applyInline(params Params) Params {
	if len(p.inline) == 0 {
		return params
	}

	if params == nil {
		params = Params{}
	}

	for key, value := range p.inline {
		params[key] = value
	}

	return params
}

func payloadFromFileable(f Fileable) uploadPayload {
	if provider, ok := f.(interface{ filePayload() uploadPayload }); ok {
		return provider.filePayload()
	}

	payload := newUploadPayload()

	for _, file := range f.files() {
		payload.Add(file.Name, file.Data)
	}

	return payload
}
