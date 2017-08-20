package mx

// LoginError описывает ошибку авторизации пользователя.
type LoginError struct {
	Code       uint8  `xml:"Code,attr" json:"code,omitempty"`
	SN         string `xml:"sn,attr" json:"sn,omitempty"`
	APIVersion uint8  `xml:"apiversion,attr" json:"apiVersion,omitempty"`
	Message    string `xml:",chardata" json:"message,omitempty"`
}

// Error возвращает строку с описанием причины ошибки авторизации.
func (e *LoginError) Error() string { return e.Message }

// CSTAError описывает ошибку CSTA.
type CSTAError struct {
	Message string `xml:",any"`
}

// Error возвращает текстовое описание ошибки.
func (e *CSTAError) Error() string { return e.Message }

// ErrTimeout возвращается когда ответ от сервера на команду не получен за время
// ReadTimeout.
var ErrTimeout error = timeoutError{}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

type emptyError struct{}

func (emptyError) Error() string { return "" }
