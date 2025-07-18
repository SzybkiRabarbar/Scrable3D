package dto

type Validatable interface {
	Validate() error
}
