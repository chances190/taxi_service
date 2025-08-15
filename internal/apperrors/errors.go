package apperrors

import "github.com/gofiber/fiber/v2"

// Error representa um erro padronizado da aplicação.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *Error) Error() string { return e.Message }

func New(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, Status: status}
}

// Erros de domínio (mensagens em pt-BR para o cliente)
var (
	ErrMotoristaNaoEncontrado   = New("motorista.nao_encontrado", "motorista não encontrado", fiber.StatusNotFound)
	ErrCPFJaCadastrado          = New("motorista.cpf_ja_cadastrado", "CPF já cadastrado", fiber.StatusConflict)
	ErrCNHJaCadastrada          = New("motorista.cnh_ja_cadastrada", "CNH já cadastrada", fiber.StatusConflict)
	ErrEmailJaCadastrado        = New("motorista.email_ja_cadastrado", "e-mail já cadastrado", fiber.StatusConflict)
	ErrSenhasNaoConferem        = New("motorista.senhas_nao_conferem", "senhas não conferem", fiber.StatusBadRequest)
	ErrSenhaAtualIncorreta      = New("motorista.senha_atual_incorreta", "senha atual incorreta", fiber.StatusUnauthorized)
	ErrDocumentoTipoInvalido    = New("documento.tipo_invalido", "tipo de documento inválido", fiber.StatusBadRequest)
	ErrDocumentoDuplicadoBatch  = New("documento.duplicado_batch", "tipo de documento duplicado na mesma requisição", fiber.StatusBadRequest)
	ErrNenhumDocumentoEnviado   = New("documento.nenhum_enviado", "nenhum documento enviado", fiber.StatusBadRequest)
	ErrDocumentosObrigPendentes = New("documento.obrigatorios_pendentes", "documentos obrigatórios pendentes", fiber.StatusBadRequest)
	ErrSenhaFraca               = New("senha.fraca", "senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo", fiber.StatusBadRequest)
	ErrCampoObrigatorio         = New("validation.campo_obrigatorio", "campo obrigatório ausente", fiber.StatusBadRequest)
	ErrCPFInvalido              = New("validation.cpf_invalido", "CPF inválido", fiber.StatusBadRequest)
	ErrCNHInvalida              = New("validation.cnh_invalida", "CNH deve ter 11 dígitos", fiber.StatusBadRequest)
	ErrEmailInvalido            = New("validation.email_invalido", "formato de email inválido", fiber.StatusBadRequest)
	ErrTelefoneInvalido         = New("validation.telefone_invalido", "formato de telefone inválido", fiber.StatusBadRequest)
	ErrPlacaInvalida            = New("validation.placa_invalida", "formato de placa inválido", fiber.StatusBadRequest)
	ErrMotoristaMenorIdade      = New("validation.menor_idade", "motorista deve ter pelo menos 18 anos", fiber.StatusBadRequest)
	ErrCNHVencida               = New("validation.cnh_vencida", "CNH vencida. Renove sua CNH para prosseguir", fiber.StatusBadRequest)
	ErrDocumentoFormatoInvalido = New("validation.documento_formato", "formato não suportado. Use JPG, PNG ou PDF", fiber.StatusBadRequest)
	ErrDocumentoMuitoGrande     = New("validation.documento_tamanho", "arquivo muito grande. Tamanho máximo: 5MB", fiber.StatusBadRequest)
	ErrDataNascimentoInvalida   = New("validation.data_nascimento_formato", "formato de data de nascimento inválido. Use DD/MM/AAAA", fiber.StatusBadRequest)
	ErrValidadeCNHInvalida      = New("validation.validade_cnh_formato", "formato de validade da CNH inválido. Use DD/MM/AAAA", fiber.StatusBadRequest)
	ErrFotoFormatoInvalido      = New("validation.foto_formato", "formato de foto não suportado. Use JPG, JPEG, PNG ou WEBP", fiber.StatusBadRequest)
	ErrFotoMuitoGrande          = New("validation.foto_tamanho", "foto muito grande. Tamanho máximo: 5MB", fiber.StatusBadRequest)
	ErrLimiteArquivosExcedido   = New("upload.limite_arquivos", "limite de arquivos excedido", fiber.StatusBadRequest)
	ErrFalhaCriarDiretorio      = New("infra.criar_diretorio", "falha ao criar diretório", fiber.StatusInternalServerError)
	ErrFalhaSalvarArquivo       = New("infra.salvar_arquivo", "falha ao salvar arquivo", fiber.StatusInternalServerError)
	ErrDocumentoNaoEncontrado   = New("documento.nao_encontrado", "documento não encontrado", fiber.StatusNotFound)
	ErrArquivoNaoEncontrado     = New("arquivo.nao_encontrado", "arquivo não encontrado", fiber.StatusNotFound)
	ErrFotoNaoEncontrada        = New("foto.nao_encontrada", "foto não encontrada", fiber.StatusNotFound)
)

// HTTPStatus retorna status adequado.
func HTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Status
	}
	return fiber.StatusInternalServerError
}

// ToPayload garante retorno padronizado.
func ToPayload(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return New("internal.erro", "erro interno", fiber.StatusInternalServerError)
}
