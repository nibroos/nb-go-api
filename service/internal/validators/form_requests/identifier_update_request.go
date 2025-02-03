package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/nb-go-api/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// IdentifierUpdateRequest handles the validation for the RegisterRequest.
type IdentifierUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of IdentifierUpdateRequest.
func NewIdentifierUpdateRequest() *IdentifierUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &IdentifierUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *IdentifierUpdateRequest) Validate(req *dtos.UpdateIdentifierRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"type_identifier_id": []string{"exists:mix_values,id"},
		"user_id":            []string{"required", "exists:users,id"},
		"ref_num":            []string{"required", fmt.Sprintf("unique_ig:identifiers,id,%d", req.ID)},
		"status":             []string{"required"},
	}

	opts := govalidator.Options{
		Data:  req,
		Rules: rules,
	}

	v := govalidator.New(opts)
	mappedErrors := v.ValidateStruct()

	if len(mappedErrors) == 0 {
		return nil
	}

	errors := make(map[string]string)
	for field, err := range mappedErrors {
		errors[field] = err[0]
	}
	return errors
}
