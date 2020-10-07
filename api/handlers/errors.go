package handlers

//Error ...
type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

// NewError creates a new error instance
func NewError(text string) error {
	return &Error{Message: text}
}
