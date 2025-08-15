package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

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

// NewMotoristaService cria uma nova instância do serviço
func NewMotoristaService(motoristaRepo repositories.MotoristaRepository, emailService EmailService) MotoristaService {
	return &MotoristaServiceImpl{
		motoristaRepo: motoristaRepo,
		emailService:  emailService,
	}
}

// CadastrarMotorista realiza o cadastro de um novo motorista
func (s *MotoristaServiceImpl) CadastrarMotorista(request CadastroMotoristaRequest) (*models.Motorista, error) {
	// Validar dados de entrada
	if err := s.ValidarDadosCadastro(request); err != nil {
		return nil, err
	}

	// Verificar se senhas coincidem
	if request.Senha != request.ConfirmacaoSenha {
		return nil, errors.New("senhas não conferem")
	}

	// Verificar se já existe motorista com mesmo CPF, CNH ou email
	if _, err := s.motoristaRepo.BuscarPorCPF(request.CPF); err == nil {
		return nil, errors.New("CPF já cadastrado")
	}

	if _, err := s.motoristaRepo.BuscarPorCNH(request.CNH); err == nil {
		return nil, errors.New("CNH já cadastrada")
	}

	if _, err := s.motoristaRepo.BuscarPorEmail(request.Email); err == nil {
		return nil, errors.New("e-mail já cadastrado")
	}

	// Parsear datas
	dataNascimento, err := time.Parse("02/01/2006", request.DataNascimento)
	if err != nil {
		return nil, errors.New("formato de data de nascimento inválido. Use DD/MM/AAAA")
	}

	validadeCNH, err := time.Parse("02/01/2006", request.ValidadeCNH)
	if err != nil {
		return nil, errors.New("formato de validade da CNH inválido. Use DD/MM/AAAA")
	}

	// Validar idade
	if err := models.ValidarIdade(dataNascimento); err != nil {
		return nil, err
	}

	// Validar validade da CNH
	if err := models.ValidarValidadeCNH(validadeCNH); err != nil {
		return nil, err
	}

	// Criar motorista
	motorista := &models.Motorista{
		ID:             uuid.New().String(),
		Nome:           request.Nome,
		DataNascimento: dataNascimento,
		CPF:            limparString(request.CPF),
		CNH:            limparString(request.CNH),
		CategoriaCNH:   models.CategoriaCNH(request.CategoriaCNH),
		ValidadeCNH:    validadeCNH,
		PlacaVeiculo:   strings.ToUpper(strings.TrimSpace(request.PlacaVeiculo)),
		ModeloVeiculo:  request.ModeloVeiculo,
		Telefone:       limparString(request.Telefone),
		Email:          strings.ToLower(strings.TrimSpace(request.Email)),
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
	// Validar campos obrigatórios
	if strings.TrimSpace(request.Nome) == "" {
		return errors.New("nome é obrigatório")
	}
	if strings.TrimSpace(request.CPF) == "" {
		return errors.New("CPF é obrigatório")
	}
	if strings.TrimSpace(request.CNH) == "" {
		return errors.New("CNH é obrigatória")
	}
	if strings.TrimSpace(request.Email) == "" {
		return errors.New("e-mail é obrigatório")
	}
	if strings.TrimSpace(request.Senha) == "" {
		return errors.New("senha é obrigatória")
	}
	if strings.TrimSpace(request.Telefone) == "" {
		return errors.New("telefone é obrigatório")
	}
	if strings.TrimSpace(request.PlacaVeiculo) == "" {
		return errors.New("placa do veículo é obrigatória")
	}

	// Validar formatos
	if !models.ValidarCPF(request.CPF) {
		return errors.New("CPF inválido")
	}

	if !models.ValidarCNH(request.CNH) {
		return errors.New("CNH deve ter 11 dígitos")
	}

	if !models.ValidarEmail(request.Email) {
		return errors.New("formato de email inválido")
	}

	if !models.ValidarTelefone(request.Telefone) {
		return errors.New("formato de telefone inválido")
	}

	if !models.ValidarPlaca(request.PlacaVeiculo) {
		return errors.New("formato de placa inválido")
	}

	// Validar força da senha
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

	// Normalizar tipo de documento (aceitar variações da selfie)
	tipo := strings.TrimSpace(strings.ToUpper(request.TipoDocumento))
	if strings.Contains(tipo, "SELFIE") { // aceita "SELFIE COM CNH" etc
		request.TipoDocumento = "selfie_cnh"
	}
	if tipo == "CRLV" { // já correto, apenas manter caixa
		request.TipoDocumento = "CRLV"
	}
	if tipo == "CNH" { // manter
		request.TipoDocumento = "CNH"
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
		return errors.New("tipo de documento inválido")
	}

	// Buscar motorista
	motorista, err := s.motoristaRepo.BuscarPorID(motoristaID)
	if err != nil {
		return errors.New("motorista não encontrado")
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
				Status:         "pendente",
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
		Status:         "pendente",
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
		motorista.Status = models.StatusDocumentosAnalise
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
		return errors.New("nenhum documento enviado")
	}
	// Map para evitar tipos duplicados na mesma requisição
	vistos := map[string]bool{}
	for _, r := range requests {
		if r.TipoDocumento == "" {
			return errors.New("tipo de documento vazio")
		}
		lower := strings.ToLower(r.TipoDocumento)
		if strings.Contains(lower, "selfie") {
			r.TipoDocumento = "selfie_cnh"
		}
		if vistos[r.TipoDocumento] {
			return fmt.Errorf("tipo de documento duplicado: %s", r.TipoDocumento)
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
	motorista, err := s.motoristaRepo.BuscarPorID(motoristaID)
	if err != nil {
		return errors.New("motorista não encontrado")
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
			return errors.New("documentos obrigatórios pendentes")
		}
	}

	// Marcar todos documentos como aprovados
	for i := range motorista.Documentos {
		motorista.Documentos[i].Status = "aprovado"
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
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return nil, errors.New("motorista não encontrado")
	}
	if telefone != "" {
		if !models.ValidarTelefone(telefone) {
			return nil, errors.New("Formato de telefone inválido.")
		}
		m.Telefone = limparString(telefone)
	}
	if email != "" {
		if !models.ValidarEmail(email) {
			return nil, errors.New("Formato de email inválido.")
		}
		m.Email = strings.ToLower(strings.TrimSpace(email))
	}
	m.AtualizadoEm = time.Now()
	if err := s.motoristaRepo.Atualizar(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AlterarSenha altera senha com validações
func (s *MotoristaServiceImpl) AlterarSenha(id, senhaAtual, novaSenha, confirmacao string) error {
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return errors.New("motorista não encontrado")
	}
	if strings.TrimSpace(senhaAtual) == "" {
		return errors.New("Senha atual é obrigatória.")
	}
	if m.Senha != senhaAtual {
		return errors.New("Senha atual incorreta.")
	}
	if strings.TrimSpace(novaSenha) == "" {
		return errors.New("Nova senha é obrigatória.")
	}
	if novaSenha != confirmacao {
		return errors.New("Nova senha e confirmação não correspondem.")
	}
	if _, err := models.ValidarForcaSenha(novaSenha); err != nil {
		return err
	}
	m.Senha = novaSenha
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// UploadFotoPerfil salva caminho da foto (arquivo já salvo pelo controller)
func (s *MotoristaServiceImpl) UploadFotoPerfil(id string, caminho string, formato string, tamanho int64) error {
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return errors.New("motorista não encontrado")
	}
	formatoU := strings.ToUpper(formato)
	permitidos := map[string]bool{"JPG": true, "JPEG": true, "PNG": true, "WEBP": true}
	if !permitidos[formatoU] {
		return errors.New("Formato não suportado. Use JPG, PNG ou WEBP")
	}
	if tamanho > 5*1024*1024 {
		return errors.New("Foto muito grande. Tamanho máximo: 5MB")
	}
	m.FotoPerfil = caminho
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// SolicitarExclusao marca status aguardando_exclusao
func (s *MotoristaServiceImpl) SolicitarExclusao(id string) error {
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return errors.New("motorista não encontrado")
	}
	m.Status = models.StatusAguardandoExclusao
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// ConfirmarExclusao marca status encerrado
func (s *MotoristaServiceImpl) ConfirmarExclusao(id string) error {
	m, err := s.motoristaRepo.BuscarPorID(id)
	if err != nil {
		return errors.New("motorista não encontrado")
	}
	m.Status = models.StatusEncerrado
	m.AtualizadoEm = time.Now()
	return s.motoristaRepo.Atualizar(m)
}

// RejeitarMotorista rejeita um motorista com motivo
func (s *MotoristaServiceImpl) RejeitarMotorista(motoristaID string, motivo string) error {
	motorista, err := s.motoristaRepo.BuscarPorID(motoristaID)
	if err != nil {
		return errors.New("motorista não encontrado")
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
		return nil, errors.New("e-mail não encontrado")
	}

	if motorista.Senha != senha {
		return nil, errors.New("senha incorreta")
	}

	return motorista, nil
}

// limparString remove caracteres especiais de strings como CPF, CNH e telefone
func limparString(s string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(s, "")
}
