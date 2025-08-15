package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"taxi_service/models"
)

// MockMotoristaRepository is a mock implementation of repositories.MotoristaRepository
type MockMotoristaRepository struct {
	mock.Mock
}

func (m *MockMotoristaRepository) Criar(motorista *models.Motorista) error {
	args := m.Called(motorista)
	return args.Error(0)
}

func (m *MockMotoristaRepository) BuscarPorID(id string) (*models.Motorista, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaRepository) BuscarPorEmail(email string) (*models.Motorista, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaRepository) BuscarPorCPF(cpf string) (*models.Motorista, error) {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaRepository) BuscarPorCNH(cnh string) (*models.Motorista, error) {
	args := m.Called(cnh)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Motorista), args.Error(1)
}

func (m *MockMotoristaRepository) Atualizar(motorista *models.Motorista) error {
	args := m.Called(motorista)
	return args.Error(0)
}

func (m *MockMotoristaRepository) Deletar(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMotoristaRepository) ListarTodos() ([]*models.Motorista, error) {
	args := m.Called()
	return args.Get(0).([]*models.Motorista), args.Error(1)
}

// MockEmailService is a mock email service for testing
type MockEmailService struct {
	mock.Mock
	emailsEnviados []EmailEnviado
}

type EmailEnviado struct {
	Para    string
	Assunto string
	Corpo   string
}

func (m *MockEmailService) EnviarEmailConfirmacao(email, nome string) error {
	args := m.Called(email, nome)
	m.emailsEnviados = append(m.emailsEnviados, EmailEnviado{
		Para:    email,
		Assunto: "Cadastro realizado com sucesso - Taxi Service",
		Corpo:   fmt.Sprintf("Olá %s, seu cadastro foi realizado com sucesso!", nome),
	})
	return args.Error(0)
}

func (m *MockEmailService) EnviarEmailRecebimentoDocumentos(email, nome string) error {
	args := m.Called(email, nome)
	m.emailsEnviados = append(m.emailsEnviados, EmailEnviado{
		Para:    email,
		Assunto: "Documentos recebidos - Taxi Service",
		Corpo:   fmt.Sprintf("Olá %s, recebemos seus documentos.", nome),
	})
	return args.Error(0)
}

func (m *MockEmailService) EnviarEmailAprovacao(email, nome string) error {
	args := m.Called(email, nome)
	m.emailsEnviados = append(m.emailsEnviados, EmailEnviado{
		Para:    email,
		Assunto: "Parabéns! Seu cadastro foi aprovado - Taxi Service",
		Corpo:   fmt.Sprintf("Olá %s, seu cadastro foi aprovado!", nome),
	})
	return args.Error(0)
}

func (m *MockEmailService) EnviarEmailRejeicao(email, nome, motivo string) error {
	args := m.Called(email, nome, motivo)
	m.emailsEnviados = append(m.emailsEnviados, EmailEnviado{
		Para:    email,
		Assunto: "Documentos rejeitados - Taxi Service",
		Corpo:   fmt.Sprintf("Olá %s, seus documentos foram rejeitados. Motivo: %s", nome, motivo),
	})
	return args.Error(0)
}

func (m *MockEmailService) ObterEmailsEnviados() []EmailEnviado {
	return m.emailsEnviados
}

func (m *MockEmailService) LimparEmails() {
	m.emailsEnviados = []EmailEnviado{}
}

func TestCadastrarMotorista(t *testing.T) {
	// helper para criar um request válido
	createValidRequest := func() CadastroMotoristaRequest {
		return CadastroMotoristaRequest{
			Nome:             "João Silva",
			DataNascimento:   "15/03/1990",
			CPF:              "11144477735",
			CNH:              "12345678901",
			CategoriaCNH:     "B",
			ValidadeCNH:      "15/03/2030",
			PlacaVeiculo:     "ABC1234",
			ModeloVeiculo:    "Honda Civic 2020",
			Telefone:         "11999999999",
			Email:            "joao.silva@email.com",
			Senha:            "MinhaSenh@123",
			ConfirmacaoSenha: "MinhaSenh@123",
		}
	}

	t.Run("Successful Registration", func(t *testing.T) {
		mockRepo := new(MockMotoristaRepository)
		mockEmail := new(MockEmailService)
		service := NewMotoristaService(mockRepo, mockEmail)
		request := createValidRequest()

		mockRepo.On("BuscarPorCPF", request.CPF).Return(nil, errors.New("not found"))
		mockRepo.On("BuscarPorCNH", request.CNH).Return(nil, errors.New("not found"))
		mockRepo.On("BuscarPorEmail", request.Email).Return(nil, errors.New("not found"))
		mockRepo.On("Criar", mock.AnythingOfType("*models.Motorista")).Return(nil)
		mockEmail.On("EnviarEmailConfirmacao", request.Email, request.Nome).Return(nil)

		// Execute the method under test
		motorista, err := service.CadastrarMotorista(request)

		// Assertions
		require.NoError(t, err)
		assert.NotEmpty(t, motorista.ID)
		assert.Equal(t, request.Nome, motorista.Nome)
		assert.Equal(t, request.Email, motorista.Email)
		assert.Equal(t, models.StatusAguardandoAprovacao, motorista.Status)

		// Verify email was sent
		emails := mockEmail.ObterEmailsEnviados()
		assert.Len(t, emails, 1)
		assert.Equal(t, request.Email, emails[0].Para)
		assert.Contains(t, emails[0].Assunto, "Cadastro realizado com sucesso")

		// Verify all mock expectations were met
		mockRepo.AssertExpectations(t)
		mockEmail.AssertExpectations(t)
	})

	t.Run("Password Mismatch Error", func(t *testing.T) {
		mockRepo := new(MockMotoristaRepository)
		mockEmail := new(MockEmailService)
		service := NewMotoristaService(mockRepo, mockEmail)
		request := createValidRequest()
		request.CPF = "52998224725"
		request.ConfirmacaoSenha = "MinhaSenh@456"

		_, err := service.CadastrarMotorista(request)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "senhas não conferem")
	})

	t.Run("CPF Already Exists Error", func(t *testing.T) {
		mockRepo := new(MockMotoristaRepository)
		mockEmail := new(MockEmailService)
		service := NewMotoristaService(mockRepo, mockEmail)
		request := createValidRequest()
		request.Email = "maria@email.com"
		request.CNH = "98765432109"

		// somente CPF é chamado antes do erro
		mockRepo.On("BuscarPorCPF", request.CPF).
			Return(&models.Motorista{ID: "existing-id"}, nil)

		_, err := service.CadastrarMotorista(request)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF já cadastrado")

		mockRepo.AssertExpectations(t)
	})

	// Add more test cases as needed
}

func TestUploadDocumento(t *testing.T) {
	// Setup
	mockRepo := new(MockMotoristaRepository)
	mockEmail := new(MockEmailService)
	service := NewMotoristaService(mockRepo, mockEmail)

	// Create a test driver
	testDriverID := uuid.New().String()
	testDriver := &models.Motorista{
		ID:         testDriverID,
		Nome:       "Test Driver",
		Email:      "test@driver.com",
		Status:     models.StatusAguardandoAprovacao,
		Documentos: []models.Documento{},
	}

	// Valid document upload request
	validUploadRequest := UploadDocumentoRequest{
		TipoDocumento:  "CNH",
		CaminhoArquivo: "/uploads/cnh_test.jpg",
		Formato:        "JPG",
		Tamanho:        2 * 1024 * 1024, // 2MB
	}

	t.Run("Successful Document Upload", func(t *testing.T) {
		// Setup mocks
		mockRepo.On("BuscarPorID", testDriverID).Return(testDriver, nil)
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Motorista")).Return(nil)

		// Execute method under test
		err := service.UploadDocumento(testDriverID, validUploadRequest)

		// Assertions
		require.NoError(t, err)
		assert.Len(t, testDriver.Documentos, 1)
		assert.Equal(t, "CNH", testDriver.Documentos[0].TipoDocumento)
		assert.Equal(t, "JPG", testDriver.Documentos[0].Formato)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Upload All Required Documents", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockMotoristaRepository)
		mockEmail = new(MockEmailService)
		service = NewMotoristaService(mockRepo, mockEmail)

		// Update driver with the first document already added
		testDriver.Documentos = []models.Documento{
			{
				TipoDocumento:  "CNH",
				CaminhoArquivo: "/uploads/cnh_test.jpg",
				Formato:        "JPG",
				Status:         "pendente",
			},
		}

		// Setup mock expectations for CRLV upload
		crlvRequest := UploadDocumentoRequest{
			TipoDocumento:  "CRLV",
			CaminhoArquivo: "/uploads/crlv_test.png",
			Formato:        "PNG",
			Tamanho:        int64(1.5 * 1024 * 1024), // 1.5MB
		}

		mockRepo.On("BuscarPorID", testDriverID).Return(testDriver, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Motorista")).Return(nil).Once()

		// Execute first upload
		err := service.UploadDocumento(testDriverID, crlvRequest)
		require.NoError(t, err)

		// Setup mock expectations for selfie upload
		selfieRequest := UploadDocumentoRequest{
			TipoDocumento:  "selfie_cnh",
			CaminhoArquivo: "/uploads/selfie_test.jpg",
			Formato:        "JPG",
			Tamanho:        1 * 1024 * 1024, // 1MB
		}

		// The driver now should have 2 documents
		driverWithTwoDocuments := &models.Motorista{
			ID:     testDriverID,
			Nome:   "Test Driver",
			Email:  "test@driver.com",
			Status: models.StatusAguardandoAprovacao,
			Documentos: []models.Documento{
				{
					TipoDocumento:  "CNH",
					CaminhoArquivo: "/uploads/cnh_test.jpg",
					Formato:        "JPG",
					Status:         "pendente",
				},
				{
					TipoDocumento:  "CRLV",
					CaminhoArquivo: "/uploads/crlv_test.png",
					Formato:        "PNG",
					Status:         "pendente",
				},
			},
		}

		mockRepo.On("BuscarPorID", testDriverID).Return(driverWithTwoDocuments, nil).Once()
		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Motorista")).Return(nil).Once()
		mockEmail.On("EnviarEmailRecebimentoDocumentos", "test@driver.com", "Test Driver").Return(nil)

		// Execute second upload
		err = service.UploadDocumento(testDriverID, selfieRequest)
		require.NoError(t, err)

		// Final driver state with status change
		driverWithAllDocs := &models.Motorista{
			ID:     testDriverID,
			Nome:   "Test Driver",
			Email:  "test@driver.com",
			Status: models.StatusDocumentosAnalise,
			Documentos: []models.Documento{
				{
					TipoDocumento:  "CNH",
					CaminhoArquivo: "/uploads/cnh_test.jpg",
					Formato:        "JPG",
					Status:         "pendente",
				},
				{
					TipoDocumento:  "CRLV",
					CaminhoArquivo: "/uploads/crlv_test.png",
					Formato:        "PNG",
					Status:         "pendente",
				},
				{
					TipoDocumento:  "selfie_cnh",
					CaminhoArquivo: "/uploads/selfie_test.jpg",
					Formato:        "JPG",
					Status:         "pendente",
				},
			},
		}

		// Setup for BuscarMotorista call
		mockRepo.On("BuscarPorID", testDriverID).Return(driverWithAllDocs, nil).Once()

		// Verify status change
		motorista, err := service.BuscarMotorista(testDriverID)
		require.NoError(t, err)
		assert.Equal(t, models.StatusDocumentosAnalise, motorista.Status)
		assert.Len(t, motorista.Documentos, 3)

		// Verify email was sent
		emails := mockEmail.ObterEmailsEnviados()
		assert.GreaterOrEqual(t, len(emails), 1)

		mockRepo.AssertExpectations(t)
		mockEmail.AssertExpectations(t)
	})

	t.Run("File Too Large Error", func(t *testing.T) {
		mockRepo := new(MockMotoristaRepository)
		mockEmail := new(MockEmailService)
		service := NewMotoristaService(mockRepo, mockEmail)

		largeFileRequest := UploadDocumentoRequest{
			TipoDocumento:  "CNH",
			CaminhoArquivo: "/uploads/cnh_large.jpg",
			Formato:        "JPG",
			Tamanho:        6 * 1024 * 1024, // 6MB - too large
		}

		err := service.UploadDocumento(testDriverID, largeFileRequest)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "arquivo muito grande")
	})

	t.Run("Invalid Format Error", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockMotoristaRepository)
		mockEmail = new(MockEmailService)
		service = NewMotoristaService(mockRepo, mockEmail)

		invalidFormatRequest := UploadDocumentoRequest{
			TipoDocumento:  "CNH",
			CaminhoArquivo: "/uploads/cnh.txt",
			Formato:        "TXT",
			Tamanho:        1 * 1024 * 1024,
		}

		// No need to mock repository calls since validation should fail first

		err := service.UploadDocumento(testDriverID, invalidFormatRequest)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "formato não suportado")
	})

	t.Run("Document Validation", func(t *testing.T) {
		// Reset mocks
		mockRepo = new(MockMotoristaRepository)
		mockEmail = new(MockEmailService)
		service = NewMotoristaService(mockRepo, mockEmail)
		mockEmail.LimparEmails()

		driverWithAllDocs := &models.Motorista{
			ID:     testDriverID,
			Nome:   "Test Driver",
			Email:  "test@driver.com",
			Status: models.StatusDocumentosAnalise,
			Documentos: []models.Documento{
				{
					TipoDocumento:  "CNH",
					CaminhoArquivo: "/uploads/cnh_test.jpg",
					Formato:        "JPG",
					Status:         "pendente",
				},
				{
					TipoDocumento:  "CRLV",
					CaminhoArquivo: "/uploads/crlv_test.png",
					Formato:        "PNG",
					Status:         "pendente",
				},
				{
					TipoDocumento:  "selfie_cnh",
					CaminhoArquivo: "/uploads/selfie_test.jpg",
					Formato:        "JPG",
					Status:         "pendente",
				},
			},
		}

		mockRepo.On("BuscarPorID", testDriverID).Return(driverWithAllDocs, nil).Once()

		// The driver with approved documents
		approvedDriver := &models.Motorista{
			ID:     testDriverID,
			Nome:   "Test Driver",
			Email:  "test@driver.com",
			Status: models.StatusAprovado,
			Documentos: []models.Documento{
				{
					TipoDocumento:  "CNH",
					CaminhoArquivo: "/uploads/cnh_test.jpg",
					Formato:        "JPG",
					Status:         "aprovado",
				},
				{
					TipoDocumento:  "CRLV",
					CaminhoArquivo: "/uploads/crlv_test.png",
					Formato:        "PNG",
					Status:         "aprovado",
				},
				{
					TipoDocumento:  "selfie_cnh",
					CaminhoArquivo: "/uploads/selfie_test.jpg",
					Formato:        "JPG",
					Status:         "aprovado",
				},
			},
		}

		mockRepo.On("Atualizar", mock.AnythingOfType("*models.Motorista")).Return(nil).Once()
		mockEmail.On("EnviarEmailAprovacao", "test@driver.com", "Test Driver").Return(nil)

		// Execute validation
		err := service.ValidarDocumentos(testDriverID)
		require.NoError(t, err)

		// Setup for status check
		mockRepo.On("BuscarPorID", testDriverID).Return(approvedDriver, nil).Once()

		// Verify status change
		motorista, err := service.BuscarMotorista(testDriverID)
		require.NoError(t, err)
		assert.Equal(t, models.StatusAprovado, motorista.Status)

		// Verify documents are approved
		for _, doc := range motorista.Documentos {
			assert.Equal(t, "aprovado", doc.Status)
		}

		// Verify approval email was sent
		emails := mockEmail.ObterEmailsEnviados()
		assert.Len(t, emails, 1)
		assert.Contains(t, emails[0].Assunto, "aprovado")

		mockRepo.AssertExpectations(t)
		mockEmail.AssertExpectations(t)
	})
}
