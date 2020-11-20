package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/enoperm/internet-services-2020/api"
	"github.com/enoperm/internet-services-2020/db"
	"github.com/enoperm/internet-services-2020/model"
)

var (
	ErrMacMismatch = errors.New("hmac mismatch")
)

type Session struct {
	SessionDB     db.SessionDatabase
	sessionSecret []byte
}

func NewSession(sdb db.SessionDatabase, secret []byte) *Session {
	return &Session{
		SessionDB:     sdb,
		sessionSecret: secret,
	}
}

func (sm *Session) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var nextReq = req
		sessionCookie, err := req.Cookie(api.COOKIE_SESSION)
		log.Println(err)
		if err == nil && sessionCookie.Expires.Before(time.Now()) {
			nextReq, err = sm.trySetSession(sessionCookie.Value, req)
			log.Printf("middleware/session: %s", err)
		}
		log.Printf("middleware/session: context-v: %#v", nextReq.Context().Value(api.CONTEXT_SESSION))
		next.ServeHTTP(rw, nextReq)
	})
}

func (sm *Session) trySetSession(b64Cookie string, req *http.Request) (*http.Request, error) {
	jsonCookie, err := base64.RawStdEncoding.DecodeString(b64Cookie)
	if err != nil {
		return req, err
	}

	var reqSess model.Session
	err = json.Unmarshal(jsonCookie, &reqSess)
	if err != nil {
		return req, err
	}

	sess, err := sm.SessionDB.FetchSession(reqSess.SessionID)
	if err != nil {
		return req, err
	}

	if sess.ValidateMac(reqSess.Mac) {
		ctx := context.WithValue(req.Context(), api.CONTEXT_SESSION, sess)
		return req.Clone(ctx), nil
	}

	return req, ErrMacMismatch
}
