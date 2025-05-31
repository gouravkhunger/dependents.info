package models

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type IngestRequest struct {
	Dependents []Dependent `json:"dependents" validate:"required,min=1,dive"`
}

type Dependent struct {
	Name  string `json:"name" validate:"required,contains=/"`
	Image string `json:"image" validate:"required,startswith=data:image/,contains=base64"`
}

var errorMessages = map[string]string{
	"default":                    "Invalid request format",
	"Dependents.min":             "At least one dependent is required",
	"Dependents.required":        "Dependents list is required",
	"Dependent.Image.required":   "Each dependent must have an image",
	"Dependent.Image.startswith": "Image must start with 'data:image/,'",
	"Dependent.Name.required":    "Each dependent must have a repository name",
	"Dependent.Name.contains":    "Repository name must be in the format 'owner/repo'",
}

func (r *IngestRequest) Validate() error {
	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			key := err.Field() + "." + err.Tag()
			if msg, ok := errorMessages[key]; ok {
				return errors.New(msg)
			} else {
				return errors.New(err.Error())
			}
		}
	}
	return nil
}
