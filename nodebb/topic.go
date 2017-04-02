package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Topic struct {
}

func (f *Forum) Topics(ctx context.Context) <-chan forum.Topic {
	ch := make(chan forum.Topic)

	go f.topics(ctx, ch)

	return ch
}

func (f *Forum) topics(ctx context.Context, ch chan<- forum.Topic) {
	defer close(ch)
}
