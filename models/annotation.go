package models

import "github.com/go-swagger/go-swagger/strfmt"

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

/*Annotation Annotation

swagger:model Annotation
*/
type Annotation struct {

	/* Endpoint endpoint
	 */
	Endpoint *Endpoint `json:"endpoint,omitempty"`

	/* Timestamp timestamp
	 */
	Timestamp *int64 `json:"timestamp,omitempty"`

	/* Value value
	 */
	Value *string `json:"value,omitempty"`
}

// Validate validates this annotation
func (m *Annotation) Validate(formats strfmt.Registry) error {
	return nil
}
