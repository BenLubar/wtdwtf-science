package forum // import "github.com/BenLubar/wtdwtf-science/forum"

type User interface {
	ID() int64
	Name() string
}
