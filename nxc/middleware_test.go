package nxc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/auth"
	"go.b8s.dev/nucleus/data"
	testUtil "go.b8s.dev/nucleus/internal/testing"
)

func TestMiddlewareEnsureAuthorizationWithoutAuth(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		c, res := setupRequest()

		middleware.EnsureAuthorization()(c)

		if res.Code != 400 {
			t.Errorf("EnsureAuthorization did not ensure authorization. Code: %v", res.Code)
		}
	})
}

func TestMiddlewareEnsureAuthorizationWithAuthInvalid(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		c, res := setupRequest()
		c.Request.Header.Add("Authorization", "Bearer idk")

		middleware.EnsureAuthorization()(c)

		if res.Code != 400 {
			t.Errorf("EnsureAuthorization did not ensure authorization. Code: %v", res.Code)
		}
	})
}

func TestMiddlewareEnsureAuthorizationWithAuthNotFound(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		c, res := setupRequest()
		c.Request.SetBasicAuth("admin", "password")

		middleware.EnsureAuthorization()(c)

		if res.Code != 401 {
			t.Errorf("EnsureAuthorization did not ensure authorization. Code: %v", res.Code)
		}
	})
}


func TestMiddlewareEnsureAuthorizationWithAuthValid(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		session, _ := CreateNextcloudAuthSession(ctx)
		CreateNextcloudAppPassword(ctx, session, user)

		c, res := setupRequest()
		c.Request.SetBasicAuth("admin", session.RawAppPassword)

		middleware.EnsureAuthorization()(c)

		if res.Code != 200 {
			t.Errorf("EnsureAuthorization failed for valid authorization. Code: %v", res.Code)
		}
	})
}

func TestMiddlewareGetCurrentUser(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		user := &auth.User{Name: "Tester", Username: "tester", EmailAddress: "tester@example.com"}
		user, _ = auth.CreateUser(ctx, user)
		session, _ := CreateNextcloudAuthSession(ctx)
		CreateNextcloudAppPassword(ctx, session, user)
		c, _ := setupRequest()
		c.Request.SetBasicAuth("admin", session.RawAppPassword)
		middleware.EnsureAuthorization()(c)

		requestUser := middleware.GetCurrentUser(c)

		if requestUser.ID != user.ID {
			t.Errorf("GetCurrentUser did not get the correct user! Got ID: %v", requestUser.ID)
		}
	})
}

func TestMiddlewareGetCurrentUserWithoutMiddleware(t *testing.T) {
	testUtil.WithData(func(ctx *data.Context) {
		middleware := &Middleware{DBContext: ctx}
		c, _ := setupRequest()
		defer func() {
			if recover() == nil {
				t.Errorf("GetCurrentUser did not panic when no user was set in context!")
			}
		}()

		middleware.GetCurrentUser(c)
	})
}

func setupRequest() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		URL: &url.URL{},
		Header: make(http.Header),
	}
	return c, w
}
