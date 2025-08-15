package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"taxi_service/internal/apperrors"
	"taxi_service/models"
	"taxi_service/repositories"
)

// CadastroMotoristaRequest representa os dados de entrada para cadastro
type CadastroMotoristaRequest struct {
	Nome             string `json:"nome" validate:"required,min=2,max=100"`
	DataNascimento   string `json:"data_nascimento" validate:"required"`
	CPF              string `json:"cpf" validate:"required"`
	CNH              string `json:"cnh" validate:"required"`
	CategoriaCNH     string `json:"categoria_cnh" validate:"required"`
	ValidadeCNH      string `json:"validade_cnh" validate:"required"`
	PlacaVeiculo     string `json:"placa_veiculo" validate:"required"`
	ModeloVeiculo    string `json:"modelo_veiculo" validate:"required,min=3,max=100"`
	Telefone         string `json:"telefone" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	Senha            string `json:"senha" validate:"required,min=8"`
	ConfirmacaoSenha string `json:"confirmacao_senha" validate:"required"`
}

// UploadDocumentoRequest representa os dados para upload de documento
type UploadDocumentoRequest struct {
	TipoDocumento  string `json:"tipo_documento" validate:"required"`
	CaminhoArquivo string `json:"caminho_arquivo" validate:"required"`
	Formato        string `json:"formato" validate:"required"`
	Tamanho        int64  `json:"tamanho" validate:"required"`
}

var documentosObrigatorios = []string{"CNH", "CRLV", "selfie_cnh"}

// MotoristaService define a interface para serviços de motorista
type MotoristaService interface {
	CadastrarMotorista(request CadastroMotoristaRequest) (*models.Motorista, error)
	ValidarDadosCadastro(request CadastroMotoristaRequest) error
	UploadDocumento(motoristaID string, request UploadDocumentoRequest) error
	UploadDocumentosLote(motoristaID string, requests []UploadDocumentoRequest) error
	AprovarMotorista(motoristaID string) error
	RejeitarMotorista(motoristaID string, motivo string) error
	AtualizarPerfil(id string, telefone string, email string) (*models.Motorista, error)
	AlterarSenha(id, senhaAtual, novaSenha, confirmacao string) error
	UploadFotoPerfil(id string, caminho string, formato string, tamanho int64) error
	SolicitarExclusao(id string) error
	ConfirmarExclusao(id string) error
	BuscarMotorista(id string) (*models.Motorista, error)
	VerificarForcaSenha(senha string) (string, error)
	LoginMotorista(email, senha string) (*models.Motorista, error)
}

// MotoristaServiceImpl implementa MotoristaService
type MotoristaServiceImpl struct {
	motoristaRepo repositories.MotoristaRepository
	emailService  EmailService
}

// getMotorista encapsula busca e mapeia erro de not found
func (s *MotoristaServiceImpl) getMotorista(id string) (*models.Motorista, error) {
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return nil, apperrors.ErrMotoristaNaoEncontrado
	}
	return m, nil
}

// NewMotoristaService cria uma nova instância do serviço
func NewMotoristaService(motoristaRepo repositories.MotoristaRepository, emailService EmailService) MotoristaService {
	return &MotoristaServiceImpl{
		motoristaRepo: motoristaRepo,
		emailService:  emailService,
	}
}

// CadastrarMotorista realiza o cadastro de um novo motorista
func (s *MotoristaServiceImpl) CadastrarMotorista(request CadastroMotoristaRequest) (*models.Motorista, error) {
	// Sanitização inicial (remoção de máscara / espaços) antes de qualquer validação
	request.CPF = DigitsOnly(request.CPF)
	request.CNH = DigitsOnly(request.CNH)
	request.Telefone = DigitsOnly(request.Telefone)
	request.PlacaVeiculo = strings.ToUpper(strings.TrimSpace(strings.ReplaceAll(request.PlacaVeiculo, "-", "")))
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Nome = strings.TrimSpace(request.Nome)
	request.ModeloVeiculo = strings.TrimSpace(request.ModeloVeiculo)
	request.Senha = strings.TrimSpace(request.Senha)
	request.ConfirmacaoSenha = strings.TrimSpace(request.ConfirmacaoSenha)

	// Validar dados de entrada
	if err := s.ValidarDadosCadastro(request); err != nil {
		return nil, err
	}

	// Verificar se senhas coincidem
	if request.Senha != request.ConfirmacaoSenha {
		return nil, apperrors.ErrSenhasNaoConferem
	}

	// Verificar se já existe motorista com mesmo CPF, CNH ou email
	if _, err := s.motoristaRepo.BuscarPorCPF(request.CPF); err == nil {
		return nil, apperrors.ErrCPFJaCadastrado
	}

	if _, err := s.motoristaRepo.BuscarPorCNH(request.CNH); err == nil {
		return nil, apperrors.ErrCNHJaCadastrada
	}

	if _, err := s.motoristaRepo.BuscarPorEmail(request.Email); err == nil {
		return nil, apperrors.ErrEmailJaCadastrado
	}

	// Parsear datas
	dataNascimento, err := time.Parse("02/01/2006", request.DataNascimento)
	if err != nil {
		return nil, apperrors.ErrDataNascimentoInvalida
	}

	validadeCNH, err := time.Parse("02/01/2006", request.ValidadeCNH)
	if err != nil {
		return nil, apperrors.ErrValidadeCNHInvalida
	}

	// Validar idade / validade CNH (model já devolve erros estruturados)
	if err := models.ValidarIdade(dataNascimento); err != nil {
		return nil, err
	}
	if err := models.ValidarValidadeCNH(validadeCNH); err != nil {
		return nil, err
	}

	// Criar motorista
	motorista := &models.Motorista{
		ID:             uuid.New().String(),
		Nome:           request.Nome,
		DataNascimento: dataNascimento,
		CPF:            request.CPF,
		CNH:            request.CNH,
		CategoriaCNH:   models.CategoriaCNH(request.CategoriaCNH),
		ValidadeCNH:    validadeCNH,
		PlacaVeiculo:   request.PlacaVeiculo,
		ModeloVeiculo:  request.ModeloVeiculo,
		Telefone:       request.Telefone,
		Email:          request.Email,
		Senha:          request.Senha, // Em produção, seria hasheada
		Status:         models.StatusAguardandoAprovacao,
		CriadoEm:       time.Now(),
		AtualizadoEm:   time.Now(),
		Documentos:     []models.Documento{},
	}

	// Salvar no repositório
	if err := s.motoristaRepo.Criar(motorista); err != nil {
		return nil, fmt.Errorf("erro ao salvar motorista: %w", err)
	}

	// Enviar email de confirmação
	if err := s.emailService.EnviarEmailConfirmacao(motorista.Email, motorista.Nome); err != nil {
		// Log do erro, mas não falha o cadastro
		fmt.Printf("Erro ao enviar email de confirmação: %v\n", err)
	}

	return motorista, nil
}

// ValidarDadosCadastro valida todos os dados de entrada
func (s *MotoristaServiceImpl) ValidarDadosCadastro(request CadastroMotoristaRequest) error {
	// Validar campos obrigatórios (loop evita repetição)
	required := []string{request.Nome, request.CPF, request.CNH, request.Email, request.Senha, request.Telefone, request.PlacaVeiculo}
	for _, v := range required {
		if strings.TrimSpace(v) == "" {
			return apperrors.ErrCampoObrigatorio
		}
	}

	// Validar formatos diretamente (funções agora retornam error)
	if err := models.ValidarCPF(request.CPF); err != nil {
		return err
	}
	if err := models.ValidarCNH(request.CNH); err != nil {
		return err
	}
	if err := models.ValidarEmail(request.Email); err != nil {
		return err
	}
	if err := models.ValidarTelefone(request.Telefone); err != nil {
		return err
	}
	if err := models.ValidarPlaca(request.PlacaVeiculo); err != nil {
		return err
	}

	if _, err := models.ValidarForcaSenha(request.Senha); err != nil {
		return err
	}

	return nil
}

// UploadDocumento adiciona um documento ao motorista
func (s *MotoristaServiceImpl) UploadDocumento(motoristaID string, request UploadDocumentoRequest) error {
	// validar formato e tamanho antes de qualquer acesso ao repo
	if err := models.ValidarDocumento(request.Formato, request.Tamanho); err != nil {
		return err
	}

	// Validar se é um tipo permitido
	permitido := false
	for _, td := range documentosObrigatorios {
		if td == request.TipoDocumento {
			permitido = true
			break
		}
	}
	if !permitido {
		return apperrors.ErrDocumentoTipoInvalido
	}

	// Buscar motorista
	motorista, err := s.getMotorista(motoristaID)
	if err != nil {
		return err
	}

	// Verificar se já existe documento do mesmo tipo
	for i, doc := range motorista.Documentos {
		if doc.TipoDocumento == request.TipoDocumento {
			// Substituir documento existente
			motorista.Documentos[i] = models.Documento{
				ID:             uuid.New().String(),
				TipoDocumento:  request.TipoDocumento,
				CaminhoArquivo: request.CaminhoArquivo,
				Formato:        strings.ToUpper(request.Formato),
				Tamanho:        request.Tamanho,
				Status:         models.DocumentoStatusPendente,
				CriadoEm:       time.Now(),
			}

			motorista.AtualizadoEm = time.Now()
			return s.motoristaRepo.Atualizar(motorista)
		}
	}

	// Adicionar novo documento
	documento := models.Documento{
		ID:             uuid.New().String(),
		TipoDocumento:  request.TipoDocumento,
		CaminhoArquivo: request.CaminhoArquivo,
		Formato:        strings.ToUpper(request.Formato),
		Tamanho:        request.Tamanho,
		Status:         models.DocumentoStatusPendente,
		CriadoEm:       time.Now(),
	}

	motorista.Documentos = append(motorista.Documentos, documento)
	motorista.AtualizadoEm = time.Now()

	// Verificar se todos os documentos obrigatórios foram enviados
	todosEnviados := true
	for _, tipoObrigatorio := range documentosObrigatorios {
		encontrado := false
		for _, doc := range motorista.Documentos {
			if doc.TipoDocumento == tipoObrigatorio {
				encontrado = true
				break
			}
		}
		if !encontrado {
			todosEnviados = false
			break
		}
	}

	// Se todos os documentos foram enviados, mudar status
	if todosEnviados && motorista.Status == models.StatusAguardandoAprovacao {
		motorista.Status = models.StatusDocumentosAnalise // manter compat; domínio duplicado
	}

	if err := s.motoristaRepo.Atualizar(motorista); err != nil {
		return fmt.Errorf("erro ao atualizar motorista: %w", err)
	}

	// Enviar email de confirmação de recebimento
	if todosEnviados {
		if err := s.emailService.EnviarEmailRecebimentoDocumentos(motorista.Email, motorista.Nome); err != nil {
			fmt.Printf("Erro ao enviar email de recebimento: %v\n", err)
		}
	}

	return nil
}

// UploadDocumentosLote faz upload de vários documentos em uma única chamada
func (s *MotoristaServiceImpl) UploadDocumentosLote(motoristaID string, requests []UploadDocumentoRequest) error {
	if len(requests) == 0 {
		return apperrors.ErrNenhumDocumentoEnviado
	}
	// Map para evitar tipos duplicados na mesma requisição
	vistos := map[string]bool{}
	for _, r := range requests {
		if r.TipoDocumento == "" {
			return apperrors.ErrCampoObrigatorio
		}
		lower := strings.ToLower(r.TipoDocumento)
		if strings.Contains(lower, "selfie") {
			r.TipoDocumento = "selfie_cnh"
		}
		if vistos[r.TipoDocumento] {
			return apperrors.ErrDocumentoDuplicadoBatch
		}
		vistos[r.TipoDocumento] = true
		if err := s.UploadDocumento(motoristaID, r); err != nil {
			return err
		}
	}
	return nil
}

// (Removida função de validação automática; aprovação agora é somente manual)

// AprovarMotorista aprova manualmente um motorista
func (s *MotoristaServiceImpl) AprovarMotorista(motoristaID string) error {
	motorista, err := s.getMotorista(motoristaID)
	if err != nil {
		return err
	}

	// Garantir que todos os documentos obrigatórios existem antes de aprovar manualmente
	for _, tipoObrigatorio := range documentosObrigatorios {
		encontrado := false
		for _, doc := range motorista.Documentos {
			if doc.TipoDocumento == tipoObrigatorio {
				encontrado = true
				break
			}
		}
		if !encontrado {
			return apperrors.ErrDocumentosObrigPendentes
		}
	}

	// Marcar todos documentos como aprovados
	for i := range motorista.Documentos {
		motorista.Documentos[i].Status = models.DocumentoStatusAprovado
	}
	motorista.Status = models.StatusAprovado
	motorista.AtualizadoEm = time.Now()

	if err := s.motoristaRepo.Atualizar(motorista); err != nil {
		return fmt.Errorf("erro ao atualizar status do motorista: %w", err)
	}

	return s.emailService.EnviarEmailAprovacao(motorista.Email, motorista.Nome)
}

// AtualizarPerfil atualiza telefone e email
func (s *MotoristaServiceImpl) AtualizarPerfil(id string, telefone string, email string) (*models.Motorista, error) {
	motorista, err := s.getMotorista(id)
	if err != nil {
		return nil, err
	}

	telefone = DigitsOnly(telefone)
	email = strings.ToLower(strings.TrimSpace(email))

	if telefone != "" {
		if err := models.ValidarTelefone(telefone); err != nil {
			return nil, err
		}
		motorista.Telefone = telefone
	}

	if email != "" {
		if err := models.ValidarEmail(email); err != nil {
			return nil, err
		}
		motorista.Email = email
	}

	// Atualizar timestamp e salvar alterações
	motorista.AtualizadoEm = time.Now()
	if err := s.motoristaRepo.Atualizar(motorista); err != nil {
		return nil, fmt.Errorf("erro ao atualizar perfil do motorista: %w", err)
	}

	return motorista, nil
}

// AlterarSenha altera senha com validações
func (s *MotoristaServiceImpl) AlterarSenha(id, senhaAtual, novaSenha, confirmacao string) error {
	motorista, err := s.getMotorista(id)
	if err != nil {
		return err
	}

	senhaAtual = strings.TrimSpace(senhaAtual)
	novaSenha = strings.TrimSpace(novaSenha)
	confirmacao = strings.TrimSpace(confirmacao)

	if senhaAtual == "" {
		return apperrors.ErrCampoObrigatorio
	}
	if motorista.Senha != senhaAtual {
		return apperrors.ErrSenhaAtualIncorreta
	}
	if novaSenha == "" {
		return apperrors.ErrCampoObrigatorio
	}
	if novaSenha != confirmacao {
		return apperrors.ErrSenhasNaoConferem
	}
	if _, err := models.ValidarForcaSenha(novaSenha); err != nil {
		return err
	}

	motorista.Senha = novaSenha
	motorista.AtualizadoEm = time.Now()

	if err := s.motoristaRepo.Atualizar(motorista); err != nil {
		return fmt.Errorf("erro ao alterar senha do motorista: %w", err)
	}

	return nil
}

// UploadFotoPerfil salva caminho da foto (arquivo já salvo pelo controller)
func (s *MotoristaServiceImpl) UploadFotoPerfil(id string, caminho string, formato string, tamanho int64) error {
	m, err := s.getMotorista(id)
	if err != nil {
		return err
	}
	formatoU := strings.ToUpper(formato)
	permitidos := map[string]bool{"JPG": true, "JPEG": true, "PNG": true, "WEBP": true}
	if !permitidos[formatoU] {
		return apperrors.ErrFotoFormatoInvalido
	}
	if tamanho > 5*1024*1024 {
		return apperrors.ErrFotoMuitoGrande
	}
	m.FotoPerfil = caminho
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// SolicitarExclusao marca status aguardando_exclusao
func (s *MotoristaServiceImpl) SolicitarExclusao(id string) error {
	m, err := s.getMotorista(id)
	if err != nil {
		return err
	}
	m.Status = models.StatusAguardandoExclusao
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// ConfirmarExclusao marca status encerrado
func (s *MotoristaServiceImpl) ConfirmarExclusao(id string) error {
	m, err := s.getMotorista(id)
	if err != nil {
		return err
	}
	m.Status = models.StatusEncerrado
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// RejeitarMotorista rejeita um motorista com motivo
func (s *MotoristaServiceImpl) RejeitarMotorista(motoristaID string, motivo string) error {
	motorista, err := s.getMotorista(motoristaID)
	if err != nil {
		return err
	}

	motorista.Status = models.StatusRejeitado
	motorista.AtualizadoEm = time.Now()

	if err := s.motoristaRepo.Atualizar(motorista); err != nil {
		return fmt.Errorf("erro ao atualizar status do motorista: %w", err)
	}

	return s.emailService.EnviarEmailRejeicao(motorista.Email, motorista.Nome, motivo)
}

// BuscarMotorista busca um motorista por ID
func (s *MotoristaServiceImpl) BuscarMotorista(id string) (*models.Motorista, error) {
	return s.motoristaRepo.BuscarPorID(id)
}

// VerificarForcaSenha verifica a força de uma senha
func (s *MotoristaServiceImpl) VerificarForcaSenha(senha string) (string, error) {
	return models.ValidarForcaSenha(senha)
}

// LoginMotorista realiza o login de um motorista
func (s *MotoristaServiceImpl) LoginMotorista(email, senha string) (*models.Motorista, error) {
	motorista, err := s.motoristaRepo.BuscarPorEmail(email)
	if err != nil {
		return nil, apperrors.ErrMotoristaNaoEncontrado
	}
	if motorista.Senha != senha {
		return nil, apperrors.ErrSenhaAtualIncorreta
	}

	return motorista, nil
}

// DigitsOnly remove caracteres especiais de strings como CPF, CNH e telefone
func DigitsOnly(s string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(s, "")
}
