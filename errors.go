package maildoor

// ecodes holds the error messages for the supported error codes.
// these get rendered in the login page error box.
var ecodes = map[string]string{
	"E1": "😥  something happened while trying to find a user account with the given email. Please try again.",
	"E2": "We're sorry, the specified token has already expired. Please enter your email again to receive a new one.",
	"E3": "The token you have entered is invalid. Please enter your email again to receive a new one.",
	"E4": "🤔 something was out of order with your previous login attempt. Please try again.",
	"E5": "😥 something happened while attempting to send the login email. Please try again.",
	"E6": "😥 an error ocurred while generating authentication token. Please try again.",
	"E7": "😥 an error ocurred login in specified user. Please try again.",
}
