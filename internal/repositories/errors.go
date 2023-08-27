package repositories

import "fmt"

type ErrNotFound struct {
	document string
	id       any
}

func (e ErrNotFound) Error() string {
	return fmt.Errorf("%s with %s not found", e.document, e.id).Error()
}
