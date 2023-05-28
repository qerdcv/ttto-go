package domain

import validation "github.com/go-ozzo/ozzo-validation"

type Step struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (s *Step) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Row, validation.Min(0), validation.Max(2)),
		validation.Field(&s.Col, validation.Min(0), validation.Max(2)),
	)
}
