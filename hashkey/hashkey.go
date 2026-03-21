// Package hashkey implements a simple random length key that can be
// hashed, encoded to string and decoded back
package hashkey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type HashKey []byte

// NewHashKey create a HashKey with the specified length
func NewHashKey(length int) (HashKey, error) {
	keyBytes := make([]byte, length)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}

	return HashKey(keyBytes), nil
}

// ParseHashKey parse a HashKey from a token (that is returned by calling String)
func ParseHashKey(token string) (HashKey, error) {
	keyBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, errors.New("invalid hash key: " + err.Error())
	}

	return HashKey(keyBytes), nil
}

// String encode the HashKey as a URL encoded string (this supposed to be returned to the user NOT STORED IN DATABASE)
func (t HashKey) String() string {
	return base64.URLEncoding.EncodeToString(t)
}

// Hash compute a sha256 hashsum on the key (for storing in databases)
func (t HashKey) Hash() []byte {
	hash := sha256.Sum256(t)
	return hash[:]
}

// B64Hash computes a sha256 checksum on the HashKey and then URL encodes it
func (t HashKey) B64Hash() string {
	return base64.URLEncoding.EncodeToString(t.Hash())
}
