// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go github.com/DolusMockServer/dolus/pkg/expectation --exclude=ExpectationError,ExpectationFieldError,Route

package expectation

import "time"

#SameSite: int // #enumSameSite

#enumSameSite:
	#SameSiteDefaultMode |
	#SameSiteLaxMode |
	#SameSiteStrictMode |
	#SameSiteNoneMode

#values_SameSite: {
	SameSiteDefaultMode: #SameSiteDefaultMode
	SameSiteLaxMode:     #SameSiteLaxMode
	SameSiteStrictMode:  #SameSiteStrictMode
	SameSiteNoneMode:    #SameSiteNoneMode
}

#SameSiteDefaultMode: #SameSite & 1
#SameSiteLaxMode:     #SameSite & 2
#SameSiteStrictMode:  #SameSite & 3
#SameSiteNoneMode:    #SameSite & 4

// struct taking from net/http package
#Cookie: {
	Name:       string
	Value:      string
	Path:       string
	Domain:     string
	Expires:    time.Time
	RawExpires: string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge:   int
	Secure:   bool
	HttpOnly: bool
	SameSite: #SameSite
	Raw:      string
	Unparsed: [...string] @go(,[]string)
}
