package model

import (
	"fmt"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string `json:"name" validate:"nonzero"`
	RG   string `json:"rg" validate:"len:9, regex:^[0-9]*$"`
	CPF  string `json:"cpf"  validate:"len:11, regex:^[0-9]*$"`
}

func ValidateStudentInfo(student *Student) error {
	if err := validator.Validate(student); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}
	return nil
}
