package forum // import "github.com/BenLubar/wtdwtf-science/forum"

type Topic interface {
	ID() int64
	Name() string
}
