package model

import (
	"crypto/hmac"
	"crypto/sha512" // TODO: Make hash function injectable
	"encoding/binary"
	"hash"
	"time"
)

type Session struct {
	UserID    int64  `json:"-"`
	SessionID int64  `json:"id"`
	Mac       []byte `json:"mac"`

	CreatedAt time.Time `json:"-"`
	LastSeen  time.Time `json:"-"`
}

func macAlgorithm(key []byte) hash.Hash {
	return hmac.New(sha512.New, key)
}

func NewSession(uid int64, key []byte) (*Session, error) {
	mac := macAlgorithm(key)

	sess := Session{
		UserID:    uid,
		CreatedAt: time.Now(),
	}

	binary.Write(mac, binary.BigEndian, sess.UserID)
	binary.Write(mac, binary.BigEndian, sess.CreatedAt.UnixNano())

	sess.Mac = mac.Sum([]byte{})

	return &sess, nil
}

func (s Session) ValidateMac(mac []byte) bool {
	return hmac.Equal(s.Mac, mac)
}
