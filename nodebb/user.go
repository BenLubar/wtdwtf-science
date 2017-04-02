package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type User struct {
}

func (f *Forum) Users(ctx context.Context) <-chan forum.User {
	ch := make(chan forum.User)

	go f.users(ctx, ch)

	return ch
}

func (f *Forum) users(ctx context.Context, ch chan<- forum.User) {
	defer close(ch)
}
