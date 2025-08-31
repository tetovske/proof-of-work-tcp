package hashcash

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ChallengeHappyPath(t *testing.T) {
	secret := "test_secret"

	ch, err := CreateChallenge(24, []byte(secret))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nonce, err := Solve(ctx, ch)
	assert.NoError(t, err)

	err = Validate([]byte(secret), ch, nonce, time.Minute)
	assert.NoError(t, err)
}

func Test_ChallengeInvalidHmac(t *testing.T) {
	secret := "test_secret"

	ch, err := CreateChallenge(24, []byte(secret))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nonce, err := Solve(ctx, ch)
	assert.NoError(t, err)

	err = Validate([]byte(secret+"1"), ch, nonce, time.Minute)
	assert.Error(t, err)
}

func Test_ChallengeInvalidHInvalidChallengeSize(t *testing.T) {
	secret := "test_secret"

	ch, err := CreateChallenge(24, []byte(secret))
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nonce, err := Solve(ctx, ch)
	assert.NoError(t, err)

	err = Validate([]byte(secret+"1"), []byte("as"), nonce, time.Minute)
	assert.Error(t, err)
}
