package api

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e AuthenticationError) Error() string {
	if e.Message == "" {
		return "Authentication error"
	}
	return e.Message
}

type AuthorizationError struct {
	Message string
}

func (e AuthorizationError) Error() string {
	if e.Message == "" {
		return "Authorization error"
	}
	return e.Message
}
