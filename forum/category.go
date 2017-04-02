package forum // import "github.com/BenLubar/wtdwtf-science/forum"

type Category interface {
	ID() int64
	Parent() int64
	Name() string
	Slug() string
	Description() string
	Order() int
	Imported() map[Forum]int64
}
