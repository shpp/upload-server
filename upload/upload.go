package upload

import ()

type Uploader struct {
	contentpath string
	sessions    []*Session
}

func NewUploader(dir string) *Uploader {
	return &Uploader{contentpath: dir}
}

// Session returns a session associated with id.
func (u *Uploader) Session(id string) *Session {
	for _, sess := range u.sessions {
		if sess.ID() == id {
			return sess
		}
	}
	return nil
}

func (u *Uploader) AddSession() *Session {
	sess := NewSession(u.contentpath)
	u.sessions = append(u.sessions, sess)
	return sess
}
