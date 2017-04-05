package forum // import "github.com/BenLubar/wtdwtf-science/forum"

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type DialFunc func(context.Context) (Forum, error)

type Shared struct {
	ID      string
	err     error
	errLock sync.Mutex
}

func (s *Shared) Name() string { return s.ID }
func (s *Shared) Err() error   { return s.err }
func (s *Shared) Check(err error, message string, args ...interface{}) bool {
	if err != nil {
		s.errLock.Lock()
		if s.err == nil {
			s.err = errors.Wrapf(err, "%s %s", s.ID, fmt.Sprintf(message, args...))
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
