package urls

type NotFoundError struct{}

func NewNotFoundError() NotFoundError {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "failed to get document, a record was not found"
}
