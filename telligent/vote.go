package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Vote struct {
}

func (f *Forum) Votes(ctx context.Context) <-chan forum.Vote {
	ch := make(chan forum.Vote)

	go f.votes(ctx, ch)

	return ch
}

func (f *Forum) votes(ctx context.Context, ch chan<- forum.Vote) {
	defer close(ch)
}
