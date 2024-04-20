package expectation

import "time"

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

// struct taking from net/http package
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path"`       // optional
	Domain     string    `json:"domain"`     // optional
	Expires    time.Time `json:"expires"`    // optional
	RawExpires string    `json:"rawExpires"` // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int      `json:"maxAge"`
	Secure   bool     `json:"secure"`
	HttpOnly bool     `json:"httpOnly"`
	SameSite SameSite `json:"sameSite"`
	Raw      string   `json:"raw"`
	Unparsed []string `json:"unparsed"` // Raw text of unparsed attribute-value pairs
}
