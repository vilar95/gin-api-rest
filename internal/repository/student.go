package repository

import (
	"github.com/vilar95/gin-api-rest/internal/model"
	"gorm.io/gorm"
)

// StudentRepository encapsula o acesso ao banco de dados para a entidade Student.
//
// Decisão: separar a camada de acesso a dados dos handlers HTTP.
// Antes, os handlers chamavam database.DB diretamente, misturando responsabilidades.
// Com o repository:
// - Os handlers focam apenas em HTTP (parse request, retornar response).
// - O repository foca em persistência (queries, inserts, updates).
// - Facilita testes: podemos testar handlers com um repository mock/fake.
type StudentRepository struct {
	db *gorm.DB
}

// NewStudentRepository cria uma nova instância do repositório.
// Recebe *gorm.DB como dependência explícita (injeção de dependência).
func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// FindAll retorna todos os estudantes cadastrados.
func (r *StudentRepository) FindAll() ([]model.Student, error) {
	var students []model.Student
	if err := r.db.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

// FindByID busca um estudante pelo ID.
func (r *StudentRepository) FindByID(id string) (model.Student, error) {
	var student model.Student
	if err := r.db.First(&student, id).Error; err != nil {
		return student, err
	}
	return student, nil
}

// FindByCPF busca um estudante pelo CPF.
func (r *StudentRepository) FindByCPF(cpf string) (model.Student, error) {
	var student model.Student
	if err := r.db.Where("cpf = ?", cpf).First(&student).Error; err != nil {
		return student, err
	}
	return student, nil
}

// Create insere um novo estudante no banco.
func (r *StudentRepository) Create(student *model.Student) error {
	return r.db.Create(student).Error
}

// Update salva as alterações de um estudante existente.
func (r *StudentRepository) Update(student *model.Student) error {
	return r.db.Save(student).Error
}

// Delete remove um estudante do banco (soft-delete via gorm.Model).
func (r *StudentRepository) Delete(student *model.Student) error {
	return r.db.Delete(student).Error
}
