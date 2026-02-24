package model

import (
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Student representa a entidade principal do domínio.
//
// As tags `json` controlam a serialização/deserialização JSON.
// As tags `validate` usam o pacote validator.v2 para validação declarativa.
//
// Decisão: mover o model para internal/model/ porque ele é específico deste projeto,
// não precisa ser exposto para outros módulos externos.
type Student struct {
	gorm.Model
	Name string `json:"name" validate:"nonzero"`
	RG   string `json:"rg" validate:"len=9,regexp=^[0-9]*$"`
	CPF  string `json:"cpf" validate:"len=11,regexp=^[0-9]*$"`
}

// Validate verifica se os campos do Student atendem às regras definidas nas tags.
//
// Decisão: usar método no struct em vez de função solta (ValidateStudentInfo).
// Isso torna a chamada mais idiomática: student.Validate() em vez de ValidateStudentInfo(&student).
func (s *Student) Validate() error {
	return validator.Validate(s)
}
