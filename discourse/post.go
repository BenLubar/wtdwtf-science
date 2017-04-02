package discourse // import "github.com/BenLubar/wtdwtf-science/discourse"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Post struct {
}

func (f *Forum) Posts(ctx context.Context) <-chan forum.Post {
	ch := make(chan forum.Post)

	go f.posts(ctx, ch)

	return ch
}

func (f *Forum) posts(ctx context.Context, ch chan<- forum.Post) {
	defer close(ch)
}
