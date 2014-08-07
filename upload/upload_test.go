package upload

import (
	"bytes"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	Gig       int = 1024 * 1024 * 1024
	chunkSize     = 1024 * 1024
)

var (
	uploader = NewUploader(".")
)

func TestUpload(t *testing.T) {
	sess, err := uploader.AddSession()

	if err != nil {
		t.Fatal("Failed to add session", err)
	}
	data := make([]byte, Gig)
	start, end := 0, chunkSize

	randomBytes(data)

	for end != len(data) {
		if end = start + chunkSize; end > len(data) {
			end = len(data)
		}
		chunk := data[start:end]
		start += chunkSize

		if err := sess.Put(bytes.NewReader(chunk)); err != nil {
			t.Fatal("Failed to upload chunk", err)
		}
	}
	fpath := "./randomfile.txt"

	if err := sess.Commit(fpath); err != nil {
		t.Fatal("Failed to commit upload", err)
	}
	os.Remove(fpath)
}

func randomBytes(p []byte) {
	r := rand.NewSource(time.Now().UnixNano())
	todo := len(p)
	offset := 0

	for {
		val := int64(r.Int63())

		for i := 0; i < 8; i++ {
			p[offset] = byte(val)
			todo--

			if todo == 0 {
				return
			}
			offset++
			val >>= 8
		}
	}
}
