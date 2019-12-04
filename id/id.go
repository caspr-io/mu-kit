package id

import (
	"strings"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
	"github.com/martinlindhe/base36"
)

type base36Encoder struct{}

func (enc base36Encoder) Encode(u uuid.UUID) string {
	bytes, _ := u.MarshalBinary()
	return strings.ToLower(base36.EncodeBytes(bytes))
}

func (enc base36Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base36.DecodeToBytes(strings.ToUpper(s)))
}

func New(prefix string) string {
	enc := base36Encoder{}
	id := shortuuid.NewWithEncoder(enc)

	for len(id) < 25 {
		id = "0" + id
	}

	return prefix + "-" + id
}
