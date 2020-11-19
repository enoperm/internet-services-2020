package api

// TODO: Split file
import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/enoperm/internet-services-2020/db"
	"github.com/enoperm/internet-services-2020/model"
)

type SessionArguments struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	UserDB        db.UserDatabase
	SessionDB     db.SessionDatabase
	sessionSecret []byte

	*mux.Router
}

var _ http.Handler = &Session{}

const COOKIE_SESSION = "session"

func NewSessionApi(router *mux.Router, udb db.UserDatabase, sdb db.SessionDatabase, secret []byte) *Session {
	udb.InitializeUserSchema()
	sdb.InitializeSessionSchema()

	r := Session{
		UserDB:        udb,
		SessionDB:     sdb,
		sessionSecret: secret,

		Router: router,
	}
	r.HandleFunc("", r.OpenSession).Methods("POST")
	r.HandleFunc("", r.OpenSession).Methods("DELETE")
	return &r
}

func (sessApi *Session) OpenSession(rw http.ResponseWriter, req *http.Request) {
	var args SessionArguments

	fail := func(err error) {
		switch err {
		default:
			rw.WriteHeader(http.StatusUnauthorized)
		}
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(err)
		return
	}

	err = json.Unmarshal(bodyBytes, &args)
	if err != nil {
		fail(err)
		return
	}

	user, err := sessApi.UserDB.FetchUser(args.Username)
	if err != nil {
		fail(err)
		return
	}

	err = user.Password.Check([]byte(args.Password))

	switch err {
	case model.ErrAuthenticationOk:
		break
	default:
		fail(err)
		return
	}

	sess, err := model.NewSession(user.Id, sessApi.sessionSecret)
	if err != nil {
		fail(err)
		return
	}

	sess, err = sessApi.SessionDB.InsertSession(*sess)
	if err != nil {
		fail(err)
		return
	}

	serSess, err := json.Marshal(sess)
	if err != nil {
		fail(err)
		return
	}

	cv := base64.RawStdEncoding.EncodeToString(serSess)

	// TODO: Set cookie attributes in middleware.
	http.SetCookie(rw, &http.Cookie{
		Name:  COOKIE_SESSION,
		Value: cv,
	})
	rw.WriteHeader(http.StatusCreated)
}

func (sessApi *Session) TerminateSession(rw http.ResponseWriter, req *http.Request) {
	// TODO: Set "session-id" object through middleware, remove session from DB.
	http.SetCookie(rw, &http.Cookie{
		Name:  COOKIE_SESSION,
		MaxAge: -1,
	})
	rw.WriteHeader(http.StatusNoContent)
}
