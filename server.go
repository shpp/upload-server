package main

import (
	"encoding/json"
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/shpp/upload-server/upload"
	"log"
	"net/http"
	"path"
)

var (
	Uploader *upload.Uploader
)

type HttpResp struct {
	code   int    `json:"-"`
	ID     string `json:"id,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Bytes  int64  `json:"bytes,omitempty"`
	Name   string `json:"name,omitempty"`
}

func main() {
	router := httprouter.New()
	content_path := flag.String("content-path", ".", "Path to storage directory")
	port := flag.String("port", "19502", "Server port")
	Uploader = upload.NewUploader(*content_path)

	flag.Parse()

	router.PUT("/upload/:id/:offset", handleFilePut)
	router.POST("/upload/:id/:name", handleCommit)

	log.Println("Listening on port", *port)
	log.Fatalln(http.ListenAndServe(":"+*port, router))
}

func encodeResponse(response HttpResp) ([]byte, error) {
	data, err := json.Marshal(response)

	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func sendHttpResponse(rw http.ResponseWriter, response HttpResp) {
	resp, err := encodeResponse(response)

	rw.Header().Set("Server", "Shpp video server")
	rw.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(rw, "{}", http.StatusInternalServerError)
		return
	}
	if response.code == http.StatusOK {
		if _, err := rw.Write(resp); err != nil {
			log.Println("HTTP write error", err)
		}
	} else {
		http.Error(rw, string(resp), response.code)
	}
}

func handleFilePut(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id, offset := p.ByName("id"), p.ByName("offset")
	var (
		sess *upload.Session
		err  error
	)
	defer req.Body.Close()

	if id == "0" {
		sess, err = Uploader.AddSession()

		if err != nil {
			log.Println("Session add:", err)
			sendHttpResponse(rw, HttpResp{code: http.StatusInternalServerError})
			return
		}
		if err = sess.Put(req.Body); err != nil {
			log.Println("Put:", err)
			sendHttpResponse(rw, HttpResp{
				code:   http.StatusInternalServerError,
				ID:     sess.ID(),
				Offset: sess.Offset(),
			})
			return
		}
	} else {
		sess = Uploader.Session(id)

		if sess == nil || sess.Expired() {
			sendHttpResponse(rw, HttpResp{code: http.StatusNotFound})
			return
		}
		if sess.OffsetStr() != offset {
			sendHttpResponse(rw, HttpResp{
				code:   http.StatusBadRequest,
				ID:     sess.ID(),
				Offset: sess.Offset(),
			})
			return
		}
		if err = sess.Put(req.Body); err != nil {
			log.Println("Put:", err)
			sendHttpResponse(rw, HttpResp{code: http.StatusInternalServerError})
			return
		}
	}

	sendHttpResponse(rw, HttpResp{
		code:   http.StatusOK,
		ID:     sess.ID(),
		Offset: sess.Offset(),
	})
}

func handleCommit(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	sess := Uploader.Session(p.ByName("id"))
	name := p.ByName("name")

	if sess == nil {
		sendHttpResponse(rw, HttpResp{code: http.StatusNotFound})
		return
	}
	fpath := path.Join(Uploader.Path(), name)

	if err := sess.Commit(fpath); err != nil {
		log.Println("Commit:", err)
		sendHttpResponse(rw, HttpResp{code: http.StatusInternalServerError})
	} else {
		Uploader.CleanupSession(sess.ID())

		sendHttpResponse(rw, HttpResp{
			code:  http.StatusOK,
			ID:    sess.ID(),
			Bytes: sess.Offset(),
			Name:  name,
		})
	}
}
