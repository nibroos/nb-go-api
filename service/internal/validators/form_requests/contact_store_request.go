package form_requests

import (
	"context"

	"github.com/nibroos/nb-go-api/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ContactStoreRequest handles the validation for the RegisterRequest.
type ContactStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of ContactStoreRequest.
func NewContactStoreRequest() *ContactStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &ContactStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ContactStoreRequest) Validate(req *dtos.CreateContactRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"type_contact_id": []string{"required", "exists:mix_values,id"},
		"user_id":         []string{"required", "exists:users,id"},
		"ref_num":         []string{"required"},
		"status":          []string{"required"},
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
