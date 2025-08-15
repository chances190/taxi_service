package models

import (
	"regexp"
	"strings"
	"time"

	"taxi_service/internal/apperrors"
)

// StatusMotorista representa os possíveis status de um motorista
type StatusMotorista string

const (
	StatusAguardandoAprovacao StatusMotorista = "aguardando_documentos"
	StatusDocumentosAnalise   StatusMotorista = "documentos_em_analise"
	StatusAprovado            StatusMotorista = "aprovado"
	StatusRejeitado           StatusMotorista = "documentos_rejeitados"
	StatusAguardandoExclusao  StatusMotorista = "aguardando_exclusao"
	StatusAtivo               StatusMotorista = "ativo"
	StatusEncerrado           StatusMotorista = "encerrado"
)

// CategoriaCNH representa as categorias de CNH
type CategoriaCNH string

const (
	CategoriaA  CategoriaCNH = "A"
	CategoriaB  CategoriaCNH = "B"
	CategoriaC  CategoriaCNH = "C"
	CategoriaD  CategoriaCNH = "D"
	CategoriaE  CategoriaCNH = "E"
	CategoriaAB CategoriaCNH = "AB"
	CategoriaAC CategoriaCNH = "AC"
	CategoriaAD CategoriaCNH = "AD"
	CategoriaAE CategoriaCNH = "AE"
)

// Motorista representa um motorista no sistema
type Motorista struct {
	ID             string          `json:"id" validate:"required"`
	Nome           string          `json:"nome" validate:"required,min=2,max=100"`
	DataNascimento time.Time       `json:"data_nascimento" validate:"required"`
	CPF            string          `json:"cpf" validate:"required"`
	CNH            string          `json:"cnh" validate:"required,len=11"`
	CategoriaCNH   CategoriaCNH    `json:"categoria_cnh" validate:"required"`
	ValidadeCNH    time.Time       `json:"validade_cnh" validate:"required"`
	PlacaVeiculo   string          `json:"placa_veiculo" validate:"required"`
	ModeloVeiculo  string          `json:"modelo_veiculo" validate:"required,min=3,max=100"`
	Telefone       string          `json:"telefone" validate:"required"`
	Email          string          `json:"email" validate:"required,email"`
	Senha          string          `json:"senha" validate:"required,min=8"`
	Status         StatusMotorista `json:"status"`
	FotoPerfil     string          `json:"foto_perfil"`
	CriadoEm       time.Time       `json:"criado_em"`
	AtualizadoEm   time.Time       `json:"atualizado_em"`
	Documentos     []Documento     `json:"documentos"`
}

// Documento representa um documento enviado pelo motorista
type Documento struct {
	ID             string    `json:"id"`
	TipoDocumento  string    `json:"tipo_documento"` // CNH, CRLV, selfie_cnh
	CaminhoArquivo string    `json:"caminho_arquivo"`
	Formato        string    `json:"formato"`
	Tamanho        int64     `json:"tamanho"`
	Status         string    `json:"status"`
	CriadoEm       time.Time `json:"criado_em"`
}

// Status de documentos
const (
	DocumentoStatusPendente = "pendente"
	DocumentoStatusAprovado = "aprovado"
)

// ValidarCPF valida o formato e dígitos verificadores do CPF
func ValidarCPF(cpf string) error {
	// Exigir somente dígitos
	if !regexp.MustCompile(`^\d{11}$`).MatchString(cpf) {
		return apperrors.ErrCPFInvalido
	}

	// Verifica se não são todos dígitos iguais
	primeiro := cpf[0]
	todosIguais := true
	for i := 1; i < 11; i++ {
		if cpf[i] != primeiro {
			todosIguais = false
			break
		}
	}
	if todosIguais {
		return apperrors.ErrCPFInvalido
	}

	// Calcula primeiro dígito verificador
	soma := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		soma += digit * (10 - i)
	}
	primeiroDigito := (soma * 10) % 11
	if primeiroDigito == 10 {
		primeiroDigito = 0
	}

	if int(cpf[9]-'0') != primeiroDigito {
		return apperrors.ErrCPFInvalido
	}

	// Calcula segundo dígito verificador
	soma = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		soma += digit * (11 - i)
	}
	segundoDigito := (soma * 10) % 11
	if segundoDigito == 10 {
		segundoDigito = 0
	}

	if int(cpf[10]-'0') != segundoDigito {
		return apperrors.ErrCPFInvalido
	}
	return nil
}

// ValidarCNH valida o formato da CNH
func ValidarCNH(cnh string) error {
	if !regexp.MustCompile(`^\d{11}$`).MatchString(cnh) {
		return apperrors.ErrCNHInvalida
	}
	return nil
}

// ValidarPlaca valida o formato da placa (formato antigo e Mercosul)
func ValidarPlaca(placa string) error {
	placa = strings.ToUpper(strings.TrimSpace(placa))
	placa = strings.ReplaceAll(placa, "-", "")

	// Formato antigo: ABC1234
	formatoAntigo := regexp.MustCompile(`^[A-Z]{3}\d{4}$`)

	// Formato Mercosul: ABC1D23
	formatoMercosul := regexp.MustCompile(`^[A-Z]{3}\d[A-Z]\d{2}$`)

	if formatoAntigo.MatchString(placa) || formatoMercosul.MatchString(placa) {
		return nil
	}
	return apperrors.ErrPlacaInvalida
}

// ValidarTelefone valida o formato do telefone brasileiro
func ValidarTelefone(telefone string) error {
	// Exige somente dígitos e tamanho adequado
	if !regexp.MustCompile(`^[1-9]{2}(9\d{8}|\d{8})$`).MatchString(telefone) {
		return apperrors.ErrTelefoneInvalido
	}
	return nil
}

// ValidarEmail valida o formato do email
func ValidarEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return apperrors.ErrEmailInvalido
	}
	return nil
}

// ValidarIdade verifica se o motorista tem pelo menos 18 anos
func ValidarIdade(dataNascimento time.Time) error {
	agora := time.Now()
	idade := agora.Year() - dataNascimento.Year()

	// Ajusta se ainda não fez aniversário este ano
	if agora.YearDay() < dataNascimento.YearDay() {
		idade--
	}

	if idade < 18 {
		return apperrors.ErrMotoristaMenorIdade
	}

	return nil
}

// ValidarValidadeCNH verifica se a CNH não está vencida
func ValidarValidadeCNH(validadeCNH time.Time) error {
	if validadeCNH.Before(time.Now().Truncate(24 * time.Hour)) {
		return apperrors.ErrCNHVencida
	}
	return nil
}

// ValidarForcaSenha retorna a força da senha e sugestões
func ValidarForcaSenha(senha string) (string, error) {
	if len(senha) < 8 {
		return "Fraca", apperrors.ErrSenhaFraca
	}

	temMaiuscula := regexp.MustCompile(`[A-Z]`).MatchString(senha)
	temMinuscula := regexp.MustCompile(`[a-z]`).MatchString(senha)
	temNumero := regexp.MustCompile(`\d`).MatchString(senha)
	temSimbolo := regexp.MustCompile(`[!@#$%^&*()_+\-=[]{};':"\\|,.<>\/?]`).MatchString(senha)

	criterios := 0
	if temMaiuscula {
		criterios++
	}
	if temMinuscula {
		criterios++
	}
	if temNumero {
		criterios++
	}
	if temSimbolo {
		criterios++
	}

	if criterios < 4 {
		return "Fraca", apperrors.ErrSenhaFraca
	}

	if len(senha) >= 12 && criterios == 4 {
		return "Forte", nil
	}

	return "Média", nil
}

// ValidarDocumento valida um documento enviado
func ValidarDocumento(formato string, tamanho int64) error {
	formatosPermitidos := []string{"JPG", "JPEG", "PNG", "PDF"}
	formatoValido := false

	formato = strings.ToUpper(formato)
	for _, f := range formatosPermitidos {
		if formato == f {
			formatoValido = true
			break
		}
	}

	if !formatoValido {
		return apperrors.ErrDocumentoFormatoInvalido
	}

	// Tamanho máximo: 5MB
	tamanhoMaximo := int64(5 * 1024 * 1024)
	if tamanho > tamanhoMaximo {
		return apperrors.ErrDocumentoMuitoGrande
	}

	return nil
}
