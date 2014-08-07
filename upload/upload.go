package upload

type Uploader struct {
	contentpath string
	sessions    []*Session
}

func NewUploader(dir string) *Uploader {
	return &Uploader{contentpath: dir}
}

// Session returns a session associated with id.
func (u *Uploader) Session(id string) *Session {
	_, sess := u.findSession(id)
	return sess
}

// AddSession creates a new upload session and adds
// it to sessions list.
func (u *Uploader) AddSession() *Session {
	sess := NewSession(u.contentpath)
	u.sessions = append(u.sessions, sess)
	return sess
}

func (u *Uploader) CleanupSession(id string) {
	i, sess := u.findSession(id)

	if i >= 0 {
		u.sessions[i] = u.sessions[len(u.sessions)-1]
		u.sessions[len(u.sessions)-1] = nil
		u.sessions = u.sessions[:len(u.sessions)-1]
	}
}

// findSession searches session list and returns found session
// and its index in the list.
func (u *Uploader) findSession(id string) (int, *Session) {
	for i, sess := range u.sessions {
		if sess.id == id {
			return i, sess
		}
	}
	return -1, nil
}
