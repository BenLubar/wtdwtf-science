package forum // import "github.com/BenLubar/wtdwtf-science/forum"

import (
	"time"

	"github.com/lib/pq"
)

type User interface {
	ID() int64
	Login() string
	DisplayName() string
	Email() string
	Slug() string
	CreatedAt() time.Time
	LastSeen() pq.NullTime
	Signature() string
	Location() string
	Bio() string
	WebAddress() string
	DateOfBirth() pq.NullTime
	Imported() map[Forum]int64
}
