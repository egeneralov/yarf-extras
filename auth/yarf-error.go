package auth

type UnauthorizedError struct {}

// Implements the error interface returning the ErrorMsg value of each error.
func (e *UnauthorizedError) Error() string {
	return "Unauthorized"
}

// Code returns the error's HTTP code to be used in the response.
func (e *UnauthorizedError) Code() int {
	return 401
}

// ID returns the error's ID for further reference.
func (e *UnauthorizedError) ID() int {
	return 401
}

// Msg returns the error's message, used to implement the Error interface.
func (e *UnauthorizedError) Msg() string {
	return "Unauthorized"
}

// Body returns the error's content body, if needed, to be returned in the HTTP response.
func (e *UnauthorizedError) Body() string {
	return "Unauthorized"
}