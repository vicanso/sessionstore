package sessionstore

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vicanso/keygrip"

	"github.com/vicanso/cod"

	"github.com/vicanso/cookies"
)

var (
	sessionKey = "X-Session-Id"
	signKeys   = []string{
		"secret1",
		"secret2",
	}
)

func TestGetID(t *testing.T) {
	sessionID := "abcd"
	t.Run("get from cookie not signed", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  sessionKey,
			Value: sessionID,
		})
		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL:           30 * time.Second,
				Key:           sessionKey,
				CookieOptions: &cookies.Options{},
			},
			ctx: ctx,
		}
		if s.GetID() != sessionID {
			t.Fatalf("get id fail")
		}
	})

	t.Run("get from cookie signed", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)

		kg := keygrip.New(signKeys)

		req.AddCookie(&http.Cookie{
			Name:  sessionKey,
			Value: sessionID,
		})

		// 添加sign cookie
		signValue := string(kg.Sign([]byte(sessionKey + "=" + sessionID)))
		req.AddCookie(&http.Cookie{
			Name:  sessionKey + cookies.SigSuffix,
			Value: signValue,
		})
		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL:           30 * time.Second,
				Key:           sessionKey,
				SignKeys:      signKeys,
				CookieOptions: &cookies.Options{},
			},
			ctx: ctx,
		}
		if s.GetID() != sessionID {
			t.Fatalf("get id fail")
		}
	})

	t.Run("get from request header not sign", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(sessionKey, sessionID)
		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL: 30 * time.Second,
				Key: sessionKey,
			},
			ctx: ctx,
		}
		if s.GetID() != sessionID {
			t.Fatalf("get id fail")
		}
	})

	t.Run("get from request header sign", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)
		req.Header.Set(sessionKey, sessionID)
		kg := keygrip.New(signKeys)
		req.Header.Set(sessionKey+cookies.SigSuffix, string(kg.Sign([]byte(sessionID))))
		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL:      30 * time.Second,
				Key:      sessionKey,
				SignKeys: signKeys,
			},
			ctx: ctx,
		}
		if s.GetID() != sessionID {
			t.Fatalf("get id fail")
		}
	})
}

func TestCreateID(t *testing.T) {
	t.Run("create id for cookie", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)

		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL:           30 * time.Second,
				Key:           sessionKey,
				SignKeys:      signKeys,
				CookieOptions: &cookies.Options{},
				IDGenerator: func() string {
					return "abcd"
				},
			},
			ctx: ctx,
		}
		id, err := s.CreateID()
		if err != nil {
			t.Fatalf("create id fail, %v", err)
		}
		v := ctx.Header()[cod.HeaderSetCookie][1]
		if id != "abcd" ||
			v != "X-Session-Id.sig=5Fak8l1ZZhsf2opVDLzxGoHJdnI" {
			t.Fatalf("create id fail")
		}
	})

	t.Run("create id for response header", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/me", nil)

		ctx := cod.NewContext(resp, req)
		s := Store{
			opts: &Options{
				TTL:      30 * time.Second,
				Key:      sessionKey,
				SignKeys: signKeys,
				IDGenerator: func() string {
					return "abcd"
				},
			},
			ctx: ctx,
		}
		id, err := s.CreateID()
		if err != nil {
			t.Fatalf("create id fail, %v", err)
		}
		v := ctx.GetHeader("X-Session-Id.sig")
		if id != "abcd" ||
			v != "Lv7WcJA3wCPst85_53_8B0jUiVY" {
			t.Fatalf("create id fail")
		}
	})
}
