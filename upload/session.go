package upload

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type Session struct {
	id string

	// Start and end time of a session. Session is
	// considered expired if now > end.
	start, end time.Time

	// chunks is a slice of chunk paths to be combined
	// into a single file
	chunks []string

	// Offset is the current offset in bytes that
	// represents the amount of data written
	Offset int64

	// path is a directory where where chunks
	// for an upload session are stored
	path string
	// rootpath is a path relative to which session
	// directories are created
	rootpath string
}

func NewSession(storagepath string) *Session {
	return &Session{rootpath: storagepath}
}

// ID returns the id of current upload session.
func (s *Session) ID() string {
	return s.id
}

// Init assings id, start/end time of a session and
// creates directory where chunks will be stored.
func (s *Session) Init() error {
	now := time.Now()
	s.id = fmt.Sprintf("%d", now.UnixNano())
	s.start = now
	s.end = now.AddDate(0, 0, 1)
	s.path = path.Join(s.rootpath, s.id)

	if err := os.Mkdir(s.path, 0775); err != nil {
		return err
	}
	return nil
}

func (s *Session) Expired() bool {
	return time.Now().After(s.end)
}

// Put writes a file chunk to disk in a separate file.
func (s *Session) Put(chunk []byte) error {
	tmppath := path.Join(s.path, fmt.Sprintf("%d.tmp", s.Offset))
	chunkpath := path.Join(s.path, fmt.Sprintf("%d.chunk", s.Offset))

	if err := s.write(tmppath, chunk); err != nil {
		return err
	}
	if err := os.Rename(tmppath, chunkpath); err != nil {
		return err
	}

	s.chunks = append(s.chunks, chunkpath)
	return nil
}

// Commit finishes an upload session by combining all its chunk into
// final destination file.
func (s *Session) Commit(filepath string) error {
	if file, err := os.Create(filepath); err != nil {
		return err
	} else {
		file.Close()
	}
	dst, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		return err
	}
	defer dst.Close()

	for _, chunk := range s.chunks {
		if file, err := os.Open(chunk); err != nil {
			return err
		} else {
			io.Copy(dst, file)
			if err := file.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Cleanup removes all files and directories associated with a session
func (s *Session) Cleanup() error {
	return os.RemoveAll(s.path)
}

func (s *Session) write(fpath string, data []byte) error {
	if file, err := os.Create(fpath); err != nil {
		return err
	} else {
		defer file.Close()

		if n, err := file.Write(data); err != nil {
			return err
		} else {
			s.Offset += int64(n)
			log.Printf("[Session %s] %d bytes written\n", s.id, n)
		}
	}
	return nil
}
