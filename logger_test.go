package maildoor_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

// stringLogger is a logger that writes to a string
// it serves as a way to demonstrate how to implement
// a custom logger if the application needs to log to
// other than stdout.
type stringLogger struct {
	content string
}

func (sl *stringLogger) Info(els ...interface{}) {
	sl.content += "level=info "
	sl.content += fmt.Sprint(els...)
	sl.content += "\n"
}

func (sl *stringLogger) Infof(format string, args ...interface{}) {
	sl.content += fmt.Sprintf("level=info %v \n"+format, args...)
}

func (sl *stringLogger) Error(els ...interface{}) {
	sl.content += "level=error "
	sl.content += fmt.Sprint(els...)
	sl.content += "\n"
}

func (sl *stringLogger) Errorf(format string, args ...interface{}) {
	sl.content += fmt.Sprintf("level=error %v \n"+format, args...)
}

func TestCustomLogger(t *testing.T) {
	lg := &stringLogger{}
	h, err := maildoor.NewWithOptions("secret", maildoor.UseLogger(lg))

	testhelpers.NoError(t, err)
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/auth/login/", nil)
	h.ServeHTTP(w, req)
	testhelpers.Equals(t, http.StatusOK, w.Code)

	testhelpers.Contains(t, lg.content, "level=info")
	testhelpers.Contains(t, lg.content, "/auth/login")
}
