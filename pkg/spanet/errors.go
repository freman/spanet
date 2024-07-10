package spanet

import "fmt"

type ErrUnexpectedResponse struct {
	Expected string
	Got      string
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response, wanted: %s, got %s (w: %x, g: %x)", e.Expected, e.Got, e.Expected, e.Got)
}

type ErrValueOutOfRange struct {
	Min   int
	Max   int
	Value int
	Name  string
}

func (e ErrValueOutOfRange) Error() string {
	return fmt.Sprintf("%s outside of permitted range %d<%d>%d", e.Name, e.Min, e.Value, e.Max)
}
