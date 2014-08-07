package upload

type Uploader struct {
	contentpath string
	sessions    []*Session
}

func NewUploader(dir string) *Uploader {
	return &Uploader{contentpath: dir}
}

func (u *Uploader) Path() string {
	return u.contentpath
}

// Session returns a session associated with id.
func (u *Uploader) Session(id string) *Session {
	_, sess := u.findSession(id)
	return sess
}

// AddSession creates a new upload session, initializes it
// and adds it to sessions list if no error occurred.
func (u *Uploader) AddSession() (*Session, error) {
	sess := NewSession(u.contentpath)

	if err := sess.Init(); err != nil {
		return nil, err
	}
	u.sessions = append(u.sessions, sess)
	return sess, nil
}

func (u *Uploader) CleanupSession(id string) {
	i, _ := u.findSession(id)

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
