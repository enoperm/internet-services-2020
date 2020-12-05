package api

// TODO: Split file
import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/enoperm/internet-services-2020/db"
	"github.com/enoperm/internet-services-2020/model"
)

type Profile struct {
	UserDB        db.UserDatabase
	ProfileDB     db.ProfileDatabase
	*mux.Router
}

var _ http.Handler = &Profile{}

func NewProfileApi(router *mux.Router, udb db.UserDatabase, pdb db.ProfileDatabase) *Profile {
	udb.InitializeUserSchema()
	pdb.InitializeProfileSchema()

	r := Profile{
		UserDB:        udb,
		ProfileDB:     pdb,

		Router: router,
	}
	r.HandleFunc("", r.FetchProfile).Methods("GET")
	r.HandleFunc("", r.UpdateProfile).Methods("POST")
	return &r
}

func (profApi *Profile) FetchProfile(rw http.ResponseWriter, req *http.Request) {
	fail := func(err error) {
		logger.Println("profiles/fetch:", err)
		switch err {
		default:
			rw.WriteHeader(http.StatusUnauthorized)
		}
	}

	sessIf := req.Context().Value(CONTEXT_SESSION)
	if sessIf == nil {
		fail(fmt.Errorf("no session"))
		return
	}

	sess := sessIf.(*model.Session)
	uid := sess.UserID

	prof, err := profApi.ProfileDB.FetchProfile(uid)
	if err != nil {
		fail(err)
		return
	}

	rankTier, _ := profApi.ProfileDB.FetchRank(uid)

	pws := prof.WithStats()
	pws.RankBelow = uint32(rankTier)

	ser, err := json.Marshal(pws)
	if err != nil {
		fail(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(ser)
}

func (profApi *Profile) UpdateProfile(rw http.ResponseWriter, req *http.Request) {
	var args model.Profile

	fail := func(err error) {
		switch err {
		default:
			logger.Println("profiles/update:", err)
			rw.WriteHeader(http.StatusUnauthorized)
		}
	}

	sessIf := req.Context().Value(CONTEXT_SESSION)
	if sessIf == nil {
		fail(fmt.Errorf("no session"))
		return
	}

	sess := sessIf.(*model.Session)
	uid := sess.UserID

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(err)
		return
	}

	json.Unmarshal(bodyBytes, &args)
	if err != nil {
		fail(err)
		return
	}

	errs := args.Validate()
	if len(errs) > 0 { fail(errs[0]); return } // TODO

	err = profApi.ProfileDB.UpdateProfile(uid, args)
	if err != nil {
		fail(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

