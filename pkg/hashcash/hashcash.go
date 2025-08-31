package hashcash

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	timestampSize = 8
	saltSize      = 8
	targetSize    = 8
	hmacSize      = 32
	prefixSize    = 24
)

const (
	ChallengeSize = prefixSize + hmacSize
	NonceSize     = 8
)

var (
	ErrBadHMAC              = errors.New("bad hmac")
	ErrExpired              = errors.New("expired")
	ErrInvalidChallengeSize = errors.New("invalid challenge size")
	ErrInvalidNonceSize     = errors.New("invalid nonce size")
	ErrVerification         = errors.New("invalid target size")
)

// CreateChallenge creates challenge signed with hmac.
// challenge structure: timestamp (8 bytes) || salt (8 bytes) || target (8 bytes) || hmac (32 bytes).
func CreateChallenge(complexity int, secret []byte) ([]byte, error) {
	if complexity < 1 || complexity > 64 {
		return nil, fmt.Errorf("bad difficulty: %d", complexity)
	}

	k := uint(complexity)
	var target uint64 = 1 << (64 - k)

	prefix := make([]byte, prefixSize)

	binary.BigEndian.PutUint64(prefix[:timestampSize], uint64(time.Now().UTC().Unix()))

	if _, err := rand.Read(prefix[timestampSize : timestampSize+saltSize]); err != nil {
		return nil, fmt.Errorf("salt: %w", err)
	}

	binary.BigEndian.PutUint64(prefix[timestampSize+saltSize:prefixSize], target)

	mac := hmac.New(sha256.New, secret)
	mac.Write(prefix)
	hmacSig := mac.Sum(nil)

	ch := make([]byte, 0, len(prefix)+mac.Size()+1)
	ch = append(ch, prefix...)
	ch = append(ch, hmacSig...)

	return ch, nil
}

// Solve solves challenge.
func Solve(ctx context.Context, challenge []byte) ([]byte, error) {
	if len(challenge) < prefixSize+hmacSize {
		return nil, ErrInvalidChallengeSize
	}

	prefix := challenge[:prefixSize]

	targetStart := timestampSize + saltSize
	target := binary.BigEndian.Uint64(challenge[targetStart : targetStart+targetSize])

	nonce := make([]byte, targetSize)

	h := sha256.New()

	for n := uint64(0); ; n++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		binary.BigEndian.PutUint64(nonce, n)

		h.Reset()
		h.Write(prefix)
		h.Write(nonce)
		sum := h.Sum(nil)

		if binary.BigEndian.Uint64(sum[:targetSize]) < target {
			return nonce, nil
		}
	}
}

// Validate validates challenge with nonce and hmac.
func Validate(secret []byte, challenge, nonce []byte, ttl time.Duration) error {
	if len(challenge) < prefixSize+hmacSize {
		return ErrInvalidChallengeSize
	}

	if len(nonce) < NonceSize {
		return ErrInvalidNonceSize
	}

	prefix := challenge[:prefixSize]
	userHmac := challenge[len(challenge)-hmacSize:]

	m := hmac.New(sha256.New, secret)
	m.Write(prefix)

	if !hmac.Equal(userHmac, m.Sum(nil)) {
		return ErrBadHMAC
	}

	ts := int64(binary.BigEndian.Uint64(challenge[:timestampSize]))
	now := time.Now().UTC()
	dt := now.Sub(time.Unix(ts, 0))
	if dt > ttl {
		return ErrExpired
	}

	targetStart := timestampSize + saltSize
	target := binary.BigEndian.Uint64(challenge[targetStart : targetStart+targetSize])

	h := sha256.New()
	h.Write(prefix)
	h.Write(nonce)
	sum := h.Sum(nil)

	if binary.BigEndian.Uint64(sum[:targetSize]) >= target {
		return ErrVerification
	}

	return nil
}
