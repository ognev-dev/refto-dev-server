package error

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/refto/server/config"
)

func New400(message string) error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func New401(message string) error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func New422(message string) error {
	return Error{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func New404(message string) error {
	return Error{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

type Error struct {
	Code    int
	Message string
	Err     error
}

type List []string
type Input map[string]string

func (l List) Error() string {
	s := strings.Builder{}
	for i, v := range l {
		s.WriteString(fmt.Sprintf("%d. %s\n", i+1, v))
	}
	return s.String()
}

func (i Input) Add(k, v string) {
	i[k] = v
}
func (i Input) Has() bool {
	return len(i) > 0
}

func (i Input) Error() string {
	s := strings.Builder{}
	for k, v := range i {
		s.WriteString(k + ": " + v + "\n")
	}
	return s.String()
}

func (e Error) Error() string {
	message := e.Message
	if !config.IsReleaseEnv() && e.Err != nil {
		message = e.Err.Error()
	}

	if message == "" {
		message = http.StatusText(e.Code)
	}

	return message
}
