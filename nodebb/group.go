package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Group struct {
}

func (f *Forum) Groups(ctx context.Context) <-chan forum.Group {
	ch := make(chan forum.Group)

	go f.groups(ctx, ch)

	return ch
}

func (f *Forum) groups(ctx context.Context, ch chan<- forum.Group) {
	defer close(ch)
}
