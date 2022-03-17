package binding

import "net/http"


const (
	MIMEJSON              = "application/json"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type Binding interface {
	Bind(*http.Request, interface{}) error
}

var (
	JSON     = jsonBinding{}
	FormPost = formPostBinding{}
)


func Defaultt(contentType string) Binding {
	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEMultipartPOSTForm:
		return FormPost
	default: // case MIMEPOSTForm:
		return FormPost
	}
}

