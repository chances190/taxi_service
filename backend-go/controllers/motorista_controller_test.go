package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"taxi_service/models"
	"taxi_service/services"
)

// MockMotoristaService mock do serviço de motorista
type MockMotoristaService struct {
	mock.Mock
}

func (m *MockMotoristaService) CadastrarMotorista(request services.CadastroMotoristaRequest) (*models.Motorista, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaService) ValidarDadosCadastro(request services.CadastroMotoristaRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockMotoristaService) UploadDocumento(motoristaID string, request services.UploadDocumentoRequest) error {
	args := m.Called(motoristaID, request)
	return args.Error(0)
}

func (m *MockMotoristaService) ValidarDocumentos(motoristaID string) error {
	args := m.Called(motoristaID)
	return args.Error(0)
}

func (m *MockMotoristaService) AprovarMotorista(motoristaID string) error {
	args := m.Called(motoristaID)
	return args.Error(0)
}

func (m *MockMotoristaService) RejeitarMotorista(motoristaID string, motivo string) error {
	args := m.Called(motoristaID, motivo)
	return args.Error(0)
}

func (m *MockMotoristaService) BuscarMotorista(id string) (*models.Motorista, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaService) VerificarForcaSenha(senha string) (string, error) {
	args := m.Called(senha)
	return args.String(0), args.Error(1)
}

func TestMotoristaController(t *testing.T) {
	// setup fornece um novo app e mockService com as rotas registradas
	setup := func() (*fiber.App, *MockMotoristaService) {
		app := fiber.New()
		mockService := new(MockMotoristaService)
		controller := NewMotoristaController(mockService)
		// Registrar rotas
		app.Post("/api/motoristas", controller.CadastrarMotorista)
		app.Get("/api/motoristas/:id", controller.BuscarMotorista)
		app.Post("/api/motoristas/:id/documentos", controller.UploadDocumento)
		app.Post("/api/motoristas/:id/validar-documentos", controller.ValidarDocumentos)
		app.Put("/api/motoristas/:id/aprovar", controller.AprovarMotorista)
		app.Put("/api/motoristas/:id/rejeitar", controller.RejeitarMotorista)
		app.Post("/api/motoristas/verificar-senha", controller.VerificarForcaSenha)
		app.Post("/api/motoristas/validar-documento", controller.ValidarDocumentoUpload)
		return app, mockService
	}

	t.Run("Cadastrar motorista com sucesso", func(t *testing.T) {
		app, mockService := setup()
		request := services.CadastroMotoristaRequest{
			Nome:             "João Silva",
			DataNascimento:   "15/03/1990",
			CPF:              "123.456.789-09",
			CNH:              "12345678901",
			CategoriaCNH:     "B",
			ValidadeCNH:      "15/03/2030",
			PlacaVeiculo:     "ABC1234",
			ModeloVeiculo:    "Honda Civic 2020",
			Telefone:         "(11) 99999-9999",
			Email:            "joao.silva@email.com",
			Senha:            "MinhaSenh@123",
			ConfirmacaoSenha: "MinhaSenh@123",
		}

		motorista := &models.Motorista{
			ID:     "123",
			Nome:   "João Silva",
			Email:  "joao.silva@email.com",
			Status: models.StatusAguardandoAprovacao,
		}

		mockService.On("CadastrarMotorista", request).Return(motorista, nil)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest("POST", "/api/motoristas", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Cadastro realizado com sucesso", response["message"])
		assert.NotNil(t, response["motorista"])

		mockService.AssertExpectations(t)
	})

	t.Run("Erro - CPF já cadastrado", func(t *testing.T) {
		app, mockService := setup()
		request := services.CadastroMotoristaRequest{
			Nome:             "João Silva",
			DataNascimento:   "15/03/1990",
			CPF:              "123.456.789-09",
			CNH:              "12345678901",
			CategoriaCNH:     "B",
			ValidadeCNH:      "15/03/2030",
			PlacaVeiculo:     "ABC1234",
			ModeloVeiculo:    "Honda Civic 2020",
			Telefone:         "(11) 99999-9999",
			Email:            "joao.silva@email.com",
			Senha:            "MinhaSenh@123",
			ConfirmacaoSenha: "MinhaSenh@123",
		}

		mockService.On("CadastrarMotorista", request).Return(nil, errors.New("CPF já cadastrado"))

		body, _ := json.Marshal(request)
		req := httptest.NewRequest("POST", "/api/motoristas", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Buscar motorista", func(t *testing.T) {
		app, mockService := setup()
		motorista := &models.Motorista{
			ID:    "123",
			Nome:  "João Silva",
			Email: "joao.silva@email.com",
		}

		mockService.On("BuscarMotorista", "123").Return(motorista, nil)

		req := httptest.NewRequest("GET", "/api/motoristas/123", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.NotNil(t, response["motorista"])

		mockService.AssertExpectations(t)
	})

	t.Run("Buscar motorista não encontrado", func(t *testing.T) {
		app, mockService := setup()
		mockService.On("BuscarMotorista", "999").Return(nil, assert.AnError)

		req := httptest.NewRequest("GET", "/api/motoristas/999", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Upload documento", func(t *testing.T) {
		app, mockService := setup()
		uploadRequest := services.UploadDocumentoRequest{
			TipoDocumento:  "CNH",
			CaminhoArquivo: "/uploads/cnh.jpg",
			Formato:        "JPG",
			Tamanho:        2 * 1024 * 1024,
		}

		mockService.On("UploadDocumento", "123", uploadRequest).Return(nil)

		body, _ := json.Marshal(uploadRequest)
		req := httptest.NewRequest("POST", "/api/motoristas/123/documentos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Documentos enviados com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Validar documentos", func(t *testing.T) {
		app, mockService := setup()
		mockService.On("ValidarDocumentos", "123").Return(nil)

		req := httptest.NewRequest("POST", "/api/motoristas/123/validar-documentos", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Documentos validados com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Aprovar motorista", func(t *testing.T) {
		app, mockService := setup()
		mockService.On("AprovarMotorista", "123").Return(nil)

		req := httptest.NewRequest("PUT", "/api/motoristas/123/aprovar", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Motorista aprovado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Rejeitar motorista", func(t *testing.T) {
		app, mockService := setup()
		rejectRequest := map[string]string{
			"motivo": "Documentos com problemas de qualidade",
		}

		mockService.On("RejeitarMotorista", "123", "Documentos com problemas de qualidade").Return(nil)

		body, _ := json.Marshal(rejectRequest)
		req := httptest.NewRequest("PUT", "/api/motoristas/123/rejeitar", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Motorista rejeitado", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Verificar força da senha", func(t *testing.T) {
		app, mockService := setup()
		senhaRequest := map[string]string{
			"senha": "MinhaSenh@123",
		}

		mockService.On("VerificarForcaSenha", "MinhaSenh@123").Return("Forte", nil)

		body, _ := json.Marshal(senhaRequest)
		req := httptest.NewRequest("POST", "/api/motoristas/verificar-senha", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Forte", response["forca"])

		mockService.AssertExpectations(t)
	})

	t.Run("Validar documento upload - válido", func(t *testing.T) {
		app, _ := setup()
		validationRequest := map[string]string{
			"formato": "JPG",
			"tamanho": "2097152", // 2MB
		}

		body, _ := json.Marshal(validationRequest)
		req := httptest.NewRequest("POST", "/api/motoristas/validar-documento", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "Documento válido", response["message"])
	})

	t.Run("Validar documento upload - inválido", func(t *testing.T) {
		app, _ := setup()
		validationRequest := map[string]string{
			"formato": "TXT",
			"tamanho": "1048576", // 1MB
		}

		body, _ := json.Marshal(validationRequest)
		req := httptest.NewRequest("POST", "/api/motoristas/validar-documento", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Contains(t, response["error"], "formato não suportado")
	})
}
