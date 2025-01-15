package errors

import (
	"strings"

	"github.com/jellydator/ttlcache/v3"
	log "github.com/sirupsen/logrus"
)

type CoreError struct {
	Key        string               `json:"key"`
	Message    string               `json:"message"`
	Attributes []CoreAttributeError `json:"attributes"`
}

type CoreAttributeError struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type CoreErrorField struct {
	Field   string `json:"field"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

func New(keys ...string) *CoreError {
	var cacheMsg *ttlcache.Item[string, string]
	msgKey := keys[0]
	var message string

	if len(keys) > 1 {
		for i := 1; i < len(keys); i++ {
			message = ConcatenateStrings(message, keys[i])
		}
	}
	cacheMsg = C.Get(msgKey)

	if cacheMsg == nil {
		log.Errorf("%s", ConcatenateStrings("error not in errors.json please contact the dev team:", msgKey))
		return &CoreError{Key: msgKey}
	}

	if StringIsNotEmpty(message) {
		fullMessageError := ConcatenateStrings(cacheMsg.Value(), " ", message)
		return &CoreError{Key: cacheMsg.Key(), Message: fullMessageError}

	}
	return &CoreError{Key: cacheMsg.Key(), Message: cacheMsg.Value()}
}

func ConvertTo(err interface{}) *CoreError {
	errOut, ok := err.(*CoreError)
	if !ok {
		errDefault, _ := err.(error)
		errOut = New("error.unmapped", errDefault.Error())
	}

	return errOut
}

func (e *CoreError) Error() string {
	return ConcatenateStrings(e.Key, " | ", e.Message)
}

func ConcatenateStrings(values ...string) string {
	var sb strings.Builder

	for _, str := range values {
		sb.WriteString(str)
	}

	return sb.String()
}

func StringIsNotEmpty(value string) bool {
	return len(strings.TrimSpace(value)) > 0
}
