package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"taxi_service/models"
)

// MotoristaRepository define a interface para operações com motoristas
type MotoristaRepository interface {
	Criar(motorista *models.Motorista) error
	BuscarPorID(id string) (*models.Motorista, error)
	BuscarPorEmail(email string) (*models.Motorista, error)
	BuscarPorCPF(cpf string) (*models.Motorista, error)
	BuscarPorCNH(cnh string) (*models.Motorista, error)
	Atualizar(motorista *models.Motorista) error
	Deletar(id string) error
	ListarTodos() ([]*models.Motorista, error)
}

// JSONMotoristaRepository implementa MotoristaRepository usando arquivo JSON
type JSONMotoristaRepository struct {
	filePath string
	mutex    sync.RWMutex
}

// NewJSONMotoristaRepository cria uma nova instância do repositório
func NewJSONMotoristaRepository() *JSONMotoristaRepository {
	return &JSONMotoristaRepository{
		filePath: "./data/motoristas.json",
	}
}

// lerMotoristas lê todos os motoristas do arquivo JSON
func (r *JSONMotoristaRepository) lerMotoristas() ([]*models.Motorista, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Criar diretório se não existir
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Se o arquivo não existir, retornar lista vazia
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*models.Motorista{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var motoristas []*models.Motorista
	if len(data) == 0 {
		return []*models.Motorista{}, nil
	}

	if err := json.Unmarshal(data, &motoristas); err != nil {
		return nil, fmt.Errorf("erro ao deserializar dados: %w", err)
	}

	return motoristas, nil
}

// salvarMotoristas salva todos os motoristas no arquivo JSON
func (r *JSONMotoristaRepository) salvarMotoristas(motoristas []*models.Motorista) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data, err := json.MarshalIndent(motoristas, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar dados: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}

// Criar adiciona um novo motorista
func (r *JSONMotoristaRepository) Criar(motorista *models.Motorista) error {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return err
	}

	// Verificar se já existe motorista com mesmo ID
	for _, m := range motoristas {
		if m.ID == motorista.ID {
			return errors.New("motorista com este ID já existe")
		}
	}

	motoristas = append(motoristas, motorista)
	return r.salvarMotoristas(motoristas)
}

// BuscarPorID busca um motorista por ID
func (r *JSONMotoristaRepository) BuscarPorID(id string) (*models.Motorista, error) {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return nil, err
	}

	for _, motorista := range motoristas {
		if motorista.ID == id {
			return motorista, nil
		}
	}

	return nil, errors.New("motorista não encontrado")
}

// BuscarPorEmail busca um motorista por email
func (r *JSONMotoristaRepository) BuscarPorEmail(email string) (*models.Motorista, error) {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return nil, err
	}

	for _, motorista := range motoristas {
		if motorista.Email == email {
			return motorista, nil
		}
	}

	return nil, errors.New("motorista não encontrado")
}

// BuscarPorCPF busca um motorista por CPF
func (r *JSONMotoristaRepository) BuscarPorCPF(cpf string) (*models.Motorista, error) {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return nil, err
	}

	for _, motorista := range motoristas {
		if motorista.CPF == cpf {
			return motorista, nil
		}
	}

	return nil, errors.New("motorista não encontrado")
}

// BuscarPorCNH busca um motorista por CNH
func (r *JSONMotoristaRepository) BuscarPorCNH(cnh string) (*models.Motorista, error) {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return nil, err
	}

	for _, motorista := range motoristas {
		if motorista.CNH == cnh {
			return motorista, nil
		}
	}

	return nil, errors.New("motorista não encontrado")
}

// Atualizar atualiza um motorista existente
func (r *JSONMotoristaRepository) Atualizar(motorista *models.Motorista) error {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return err
	}

	for i, m := range motoristas {
		if m.ID == motorista.ID {
			motoristas[i] = motorista
			return r.salvarMotoristas(motoristas)
		}
	}

	return errors.New("motorista não encontrado")
}

// Deletar remove um motorista
func (r *JSONMotoristaRepository) Deletar(id string) error {
	motoristas, err := r.lerMotoristas()
	if err != nil {
		return err
	}

	for i, motorista := range motoristas {
		if motorista.ID == id {
			motoristas = append(motoristas[:i], motoristas[i+1:]...)
			return r.salvarMotoristas(motoristas)
		}
	}

	return errors.New("motorista não encontrado")
}

// ListarTodos retorna todos os motoristas
func (r *JSONMotoristaRepository) ListarTodos() ([]*models.Motorista, error) {
	return r.lerMotoristas()
}
