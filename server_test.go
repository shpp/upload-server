package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	halfGig    int = 1024 * 1024 * 500
	chunkSize      = 1024 * 1024
	serverPort     = "19502"
)

type ServerResp struct {
	code   int    `json:"-"`
	ID     string `json:"id,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Bytes  int64  `json:"bytes,omitempty"`
	Name   string `json:"name,omitempty"`
}

func TestUpload(t *testing.T) {
	go startServer(serverPort, "./content/files")

	data := make([]byte, halfGig)
	start, end, size := 0, chunkSize, len(data)
	sessID, sessOffset := "0", int64(0)
	randomBytes(data)

	log.Printf("Starting file upload\n")

	for end != size {
		if end = start + chunkSize; end > size {
			end = size
		}
		chunk := data[start:end]
		start += chunkSize

		url := fmt.Sprintf("/upload/%s/%d", sessID, sessOffset)
		if resp, err := doRequest("PUT", url, chunk); err != nil {
			t.Fatalf("Failed to upload chunk: %s\n", err)
		} else {
			log.Printf("Got response: %+v\n", resp)
			sessID, sessOffset = resp.ID, resp.Offset
		}
	}

	filename := "testfile.txt"
	url := fmt.Sprintf("/upload/%s/%s", sessID, filename)
	resp, err := doRequest("POST", url, nil)

	if err != nil {
		t.Fatalf("Failed to commit upload: %s\n", err)
	} else {
		log.Printf("Commit response: %+v\n", resp)
	}
	if err := os.Remove(filename); err != nil {
		t.Fatalf("Failed to remove test file %s (size %d bytes)\n", filename, halfGig)
	}
}

func doRequest(method, url string, body []byte) (ServerResp, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://0.0.0.0:"+serverPort+url, bytes.NewReader(body))

	if err != nil {
		return ServerResp{}, err
	}
	if resp, err := client.Do(req); err != nil {
		return ServerResp{}, err
	} else {
		defer resp.Body.Close()

		if rbody, err := ioutil.ReadAll(resp.Body); err != nil {
			return ServerResp{}, err
		} else {
			sresp := ServerResp{}

			if err := json.Unmarshal(rbody, &sresp); err != nil {
				return ServerResp{}, err
			}
			return sresp, nil
		}
	}
	return ServerResp{}, nil
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
