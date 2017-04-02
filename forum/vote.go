package forum // import "github.com/BenLubar/wtdwtf-science/forum"

type Vote interface {
	UserID() int64
	PostID() int64
	Up() bool
}
