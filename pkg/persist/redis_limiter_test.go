package persist

import (
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

var (
	db *Redis
)

var (
	key            = "key"
	val      int64 = 1
	duration int64 = 10
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db = NewRedis("redis://" + mr.Addr())

	code := m.Run()
	os.Exit(code)
}

func TestSetNX(t *testing.T) {

	err := db.SetNX(key, val, duration)
	assert.NoError(t, err)

	err = db.SetNX(key, val, 10000000000000)
	assert.NoError(t, err)

	err = db.SetNX("", val, duration)
	assert.Error(t, err)

	err = db.SetNX(key, "", duration)
	assert.Error(t, err)

	err = db.SetNX(key, "", 0)
	assert.Error(t, err)
}

func TestReset(t *testing.T) {

	err := db.Reset(key, duration)
	assert.NoError(t, err)

	err = db.Reset("", duration)
	assert.Error(t, err)

	err = db.Reset(key, 0)
	assert.Error(t, err)
}

func TestIncr(t *testing.T) {

	res, err := db.Incr(key)
	assert.NoError(t, err)
	assert.Equal(t, val, res)

	res, err = db.Incr(key)
	assert.NoError(t, err)
	assert.Equal(t, val+1, res)

	_, err = db.Incr("")
	assert.Error(t, err)
}
