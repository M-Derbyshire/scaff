package customerrors

// Used to represent an error that occured when validating a model's data/structure
type ValidationError struct {
	Message string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}
