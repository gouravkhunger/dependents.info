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
	Total      int         `json:"total" validate:"min=0"`
	Dependents []Dependent `json:"dependents" validate:"required,dive"`
}

type Dependent struct {
	Image string `json:"image" validate:"required,startswith=data:image/,contains=base64"`
}

type RepoPage struct {
	Id    string
	Total int
	Owner string
	Repo  string
}

var errorMessages = map[string]string{
	"default":                    "Invalid request format",
	"Total.required":             "Total number of dependents is required",
	"Dependents.min":             "At least one dependent is required",
	"Dependents.required":        "Dependents list is required",
	"Dependent.Image.required":   "Each dependent must have an image",
	"Dependent.Image.startswith": "Image must start with 'data:image/,'",
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
