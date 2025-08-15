package repositories

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"taxi_service/models"
)

func TestJSONMotoristaRepository(t *testing.T) {
	// Usar arquivo temporário para testes
	tempFile := "./data/test_motoristas.json"

	// Limpar arquivo antes dos testes
	os.Remove(tempFile)
	defer os.Remove(tempFile)

	repo := &JSONMotoristaRepository{
		filePath: tempFile,
	}

	t.Run("Criar motorista", func(t *testing.T) {
		motorista := &models.Motorista{
			ID:             "1",
			Nome:           "João Silva",
			DataNascimento: time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
			CPF:            "12345678909",
			CNH:            "12345678901",
			CategoriaCNH:   models.CategoriaB,
			ValidadeCNH:    time.Date(2030, 3, 15, 0, 0, 0, 0, time.UTC),
			PlacaVeiculo:   "ABC1234",
			ModeloVeiculo:  "Honda Civic 2020",
			Telefone:       "11999999999",
			Email:          "joao.silva@email.com",
			Senha:          "MinhaSenh@123",
			Status:         models.StatusAguardandoAprovacao,
			CriadoEm:       time.Now(),
			AtualizadoEm:   time.Now(),
		}

		err := repo.Criar(motorista)
		assert.NoError(t, err)
	})

	t.Run("Buscar motorista por ID", func(t *testing.T) {
		motorista, err := repo.BuscarPorID("1")
		require.NoError(t, err)
		assert.Equal(t, "1", motorista.ID)
		assert.Equal(t, "João Silva", motorista.Nome)
		assert.Equal(t, "joao.silva@email.com", motorista.Email)
	})

	t.Run("Buscar motorista por email", func(t *testing.T) {
		motorista, err := repo.BuscarPorEmail("joao.silva@email.com")
		require.NoError(t, err)
		assert.Equal(t, "1", motorista.ID)
		assert.Equal(t, "João Silva", motorista.Nome)
	})

	t.Run("Buscar motorista por CPF", func(t *testing.T) {
		motorista, err := repo.BuscarPorCPF("12345678909")
		require.NoError(t, err)
		assert.Equal(t, "1", motorista.ID)
		assert.Equal(t, "João Silva", motorista.Nome)
	})

	t.Run("Buscar motorista por CNH", func(t *testing.T) {
		motorista, err := repo.BuscarPorCNH("12345678901")
		require.NoError(t, err)
		assert.Equal(t, "1", motorista.ID)
		assert.Equal(t, "João Silva", motorista.Nome)
	})

	t.Run("Atualizar motorista", func(t *testing.T) {
		motorista, err := repo.BuscarPorID("1")
		require.NoError(t, err)

		motorista.Nome = "João Silva Santos"
		motorista.Status = models.StatusAprovado

		err = repo.Atualizar(motorista)
		assert.NoError(t, err)

		motoristaAtualizado, err := repo.BuscarPorID("1")
		require.NoError(t, err)
		assert.Equal(t, "João Silva Santos", motoristaAtualizado.Nome)
		assert.Equal(t, models.StatusAprovado, motoristaAtualizado.Status)
	})

	t.Run("Listar todos os motoristas", func(t *testing.T) {
		motoristas, err := repo.ListarTodos()
		require.NoError(t, err)
		assert.Len(t, motoristas, 1)
		assert.Equal(t, "João Silva Santos", motoristas[0].Nome)
	})

	t.Run("Erro ao buscar motorista inexistente", func(t *testing.T) {
		_, err := repo.BuscarPorID("999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "motorista não encontrado")
	})

	t.Run("Erro ao criar motorista com ID duplicado", func(t *testing.T) {
		motorista := &models.Motorista{
			ID:    "1",
			Nome:  "Outro Motorista",
			Email: "outro@email.com",
		}

		err := repo.Criar(motorista)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "motorista com este ID já existe")
	})

	t.Run("Deletar motorista", func(t *testing.T) {
		err := repo.Deletar("1")
		assert.NoError(t, err)

		_, err = repo.BuscarPorID("1")
		assert.Error(t, err)

		motoristas, err := repo.ListarTodos()
		require.NoError(t, err)
		assert.Len(t, motoristas, 0)
	})

	t.Run("Erro ao deletar motorista inexistente", func(t *testing.T) {
		err := repo.Deletar("999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "motorista não encontrado")
	})
}
