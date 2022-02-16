package csrf

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func TestGenerate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(""))
	conn, _ := buildTestConn(req)
	Generate()(conn)
	s := sessions.Default(conn)
	if s.Get("CSRFToken") == nil {
		t.Errorf("CSRF token was not present in session!")
	}
}

func TestValidateWhenTokensMatch(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("_csrf=supersecure"))
	req.Header.Add("content-type", gin.MIMEPOSTForm)
	conn, resp := buildTestConn(req)
	s := sessions.Default(conn)
	s.Set("CSRFToken", "supersecure")
	Validate()(conn)
	if resp.Code != 200 {
		t.Errorf(
			"CSRF tokens did not match: form=%v  session=%v",
			conn.PostForm("_csrf"),
			s.Get("CSRFToken"),
		)
	}
}

func TestValidateWhenSessionMissing(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("_csrf=supersecure"))
	req.Header.Add("content-type", gin.MIMEPOSTForm)
	conn, resp := buildTestConn(req)
	Validate()(conn)
	if resp.Code != 400 {
		t.Errorf("Wrong response code for missing session: %v", resp.Code)
	}
}

func TestValidateWhenPostParamMissing(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(""))
	req.Header.Add("content-type", gin.MIMEPOSTForm)
	conn, resp := buildTestConn(req)
	s := sessions.Default(conn)
	s.Set("CSRFToken", "supersecure")
	Validate()(conn)
	if resp.Code != 401 {
		t.Errorf("Wrong response code for missing session: %v", resp.Code)
	}
}

func TestValidateWhenPostParamAndSessionMismatch(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString("_csrf=notsecure"))
	req.Header.Add("content-type", gin.MIMEPOSTForm)
	conn, resp := buildTestConn(req)
	s := sessions.Default(conn)
	s.Set("CSRFToken", "supersecure")
	Validate()(conn)
	if resp.Code != 401 {
		t.Errorf("Wrong response code for missing session: %v", resp.Code)
	}
}

func buildTestConn(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	conn, _ := gin.CreateTestContext(w)
	conn.Request = req
	store := memstore.NewStore([]byte("testsecret"))
	sessions.Sessions("testsession", store)(conn)
	return conn, w
}
