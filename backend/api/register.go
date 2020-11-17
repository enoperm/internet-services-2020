package api // TODO: Move each API implementation into its own package

// TODO: Split file
import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "unicode"

    "github.com/gorilla/mux"

    "github.com/enoperm/internet-services-2020/db"
    "github.com/enoperm/internet-services-2020/model"
)

type RegisterArguments struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	UserDB db.UserDatabase
	*mux.Router
}
var _ http.Handler = &Register{}

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "/register: ", log.LstdFlags | log.LUTC | log.Lmsgprefix)
}

func NewRegisterApi(router *mux.Router, udb db.UserDatabase) *Register {
	udb.InitializeUserSchema()
	r := Register{
		UserDB: udb,
		Router: router,
	}
	r.HandleFunc("", r.registerUser).Methods("POST")
	return &r
}

func (regApi *Register) registerUser(rw http.ResponseWriter, req *http.Request) {
	var args RegisterArguments

	fail := func(err error) {
		logger.Println(err)
		switch err {
		case model.ErrRegisterInvalidUsername: fallthrough
		case model.ErrRegisterInvalidPassword:
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))

		case model.ErrRegisterUsernameUnavailable:
			rw.WriteHeader(http.StatusConflict)
			rw.Write([]byte(err.Error()))

		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}

	validUsernameChar := func(r rune) rune {
		if unicode.IsPrint(r) && r <= unicode.MaxASCII {
			return r
		}
		return -r
	}

	readBuffer, err := ioutil.ReadAll(req.Body)
	if err != nil { fail(err); return }

	err = json.Unmarshal(readBuffer, &args)
	if err != nil { fail(err); return }

	filteredUsername := strings.Map(validUsernameChar, args.Username)

	switch {
	case len(args.Username) < 1: fallthrough
	case len(args.Username) != len(filteredUsername):
		fail(model.ErrRegisterInvalidUsername)
		return

	case len(args.Password) < 1:
		fail(model.ErrRegisterInvalidPassword)
		return
	}

	salt, err := model.GenerateRandomSalt()
	if err != nil { fail(err); return }

	password, err := model.NewPassword([]byte(args.Password), salt[:])
	if err != nil { fail(err); return }

	err = regApi.UserDB.InsertUser(model.User{
		Name: args.Username,
		Password: password,
	})
	if err != nil { fail(err); return }

	logger.Println("succ ->", args.Username)
	rw.WriteHeader(http.StatusCreated)
}
