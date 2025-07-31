package method

type HTTPContentType string

const (
	ContentTypeJSON  HTTPContentType = "application/json"
	ContentTypeForm  HTTPContentType = "application/x-www-form-urlencoded"
	ContentTypeMulti HTTPContentType = "multipart/form-data"
	ContentTypeText  HTTPContentType = "text/plain"
)

func (c HTTPContentType) IsValid() bool {
	switch c {
	case ContentTypeJSON, ContentTypeForm, ContentTypeMulti, ContentTypeText:
		return true
	}
	return false
}
