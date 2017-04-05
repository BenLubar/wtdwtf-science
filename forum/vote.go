package forum // import "github.com/BenLubar/wtdwtf-science/forum"

type Vote interface {
	User() int64
	Post() int64
	Up() bool
}
