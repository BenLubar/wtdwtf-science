package forum // import "github.com/BenLubar/wtdwtf-science/forum"

import (
	"context"
	"sync"
)

type DialFunc func(context.Context) (Forum, error)

type Shared struct {
	ID      string
	err     error
	errLock sync.Mutex
}

func (s *Shared) Name() string { return s.ID }
func (s *Shared) Err() error   { return s.err }
func (s *Shared) Check(err error) bool {
	if err != nil {
		s.errLock.Lock()
		if s.err == nil {
			s.err = err
		}
		s.errLock.Unlock()
		return true
	}
	return false
}

type Forum interface {
	Close() error
	Err() error
	Name() string
	SetPreviousForums([]Forum)
	Users(context.Context) <-chan User
	Groups(context.Context) <-chan Group
	Categories(context.Context) <-chan Category
	Topics(context.Context) <-chan Topic
	Posts(context.Context) <-chan Post
	Votes(context.Context) <-chan Vote
}
