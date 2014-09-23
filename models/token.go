package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"strconv"

	"time"
)

const (
	tokenExpiry = time.Hour * 24 * 7
)

type Token struct {
	Id        int64
	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `sql:"-" json:"expired_at"`
	UserId    int64     `json:"user_id"`
	Ip        IP        `json:"ip"`
	IsCurrent bool      `sql:"-" json:"is_current"`
}

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"key":        t.Key,
		"updated_at": ISOTime(t.UpdatedAt),
		"expired_at": ISOTime(t.GetExpiredTime()),
		"user_id":    t.UserId,
	})
}

func (t *Token) BeforeSave() error {
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Token) BeforeCreate() error {
	t.CreatedAt = time.Now().UTC()

	h := sha256.New()
	raw := strconv.FormatInt(t.UserId, 10) + "/" + strconv.FormatInt(t.CreatedAt.UnixNano(), 10)
	io.WriteString(h, raw)

	t.Key = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (t *Token) GetExpiredTime() time.Time {
	return t.UpdatedAt.Add(tokenExpiry)
}

func (t *Token) IsExpired() bool {
	return t.GetExpiredTime().Before(time.Now())
}

func (t *Token) SetIP(addr string) {
	t.Ip = IP{net.ParseIP(addr)}
}
