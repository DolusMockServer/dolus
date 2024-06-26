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
	name:       string    @go(Name)
	value:      string    @go(Value)
	path:       string    @go(Path)
	domain:     string    @go(Domain)
	expires:    time.Time @go(Expires)
	rawExpires: string    @go(RawExpires)

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	maxAge:   int       @go(MaxAge)
	secure:   bool      @go(Secure)
	httpOnly: bool      @go(HttpOnly)
	sameSite: #SameSite @go(SameSite)
	raw:      string    @go(Raw)
	unparsed: [...string] @go(Unparsed,[]string)
}
