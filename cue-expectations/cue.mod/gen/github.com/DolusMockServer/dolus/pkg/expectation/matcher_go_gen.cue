// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go github.com/DolusMockServer/dolus/pkg/expectation --exclude=ExpectationError,ExpectationFieldError,Route

package expectation

#Matcher: {
	match:  string   @go(Match)
	value?: null | _ @go(Value,*T)
}