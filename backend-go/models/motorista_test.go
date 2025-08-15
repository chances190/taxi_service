package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidarCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"CPF válido com pontos", "11144477735", true},
		{"CPF válido sem pontos", "11144477735", true},
		{"CPF inválido - dígitos incorretos", "123.456.789-99", false},
		{"CPF inválido - todos iguais", "111.111.111-11", false},
		{"CPF inválido - poucos dígitos", "123.456.789", false},
		{"CPF inválido - muitos dígitos", "123.456.789-000", false},
		{"CPF com letras", "abc.def.ghi-jk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarCPF(tt.cpf)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidarCNH(t *testing.T) {
	tests := []struct {
		name     string
		cnh      string
		expected bool
	}{
		{"CNH válida", "12345678901", true},
		{"CNH com menos dígitos", "1234567890", false},
		{"CNH com mais dígitos", "123456789012", false},
		{"CNH com letras", "1234567890a", false},
		{"CNH vazia", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarCNH(tt.cnh)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidarPlaca(t *testing.T) {
	tests := []struct {
		name     string
		placa    string
		expected bool
	}{
		{"Placa formato antigo válida", "ABC1234", true},
		{"Placa formato Mercosul válida", "ABC1D23", true},
		{"Placa formato antigo minúscula", "abc1234", true},
		{"Placa inválida - poucos caracteres", "ABC123", false},
		{"Placa inválida - muitos caracteres", "ABC12345", false},
		{"Placa inválida - formato incorreto", "AB1234C", false},
		{"Placa vazia", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarPlaca(tt.placa)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidarTelefone(t *testing.T) {
	tests := []struct {
		name     string
		telefone string
		expected bool
	}{
		{"Telefone celular válido com parênteses", "(11) 99999-9999", true},
		{"Telefone celular válido sem formatação", "11999999999", true},
		{"Telefone fixo válido", "(11) 3333-3333", true},
		{"Telefone fixo sem formatação", "1133333333", true},
		{"Telefone inválido - poucos dígitos", "119999999", false},
		{"Telefone inválido - muitos dígitos", "1199999999999", false},
		{"Telefone com letras", "11abcdefghi", false},
		{"Telefone vazio", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarTelefone(tt.telefone)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidarEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Email válido", "joao.silva@email.com", true},
		{"Email válido com números", "user123@test.com", true},
		{"Email válido com subdomínio", "user@mail.test.com", true},
		{"Email inválido - sem @", "email_invalido", false},
		{"Email inválido - sem domínio", "user@", false},
		{"Email inválido - sem usuário", "@domain.com", false},
		{"Email inválido - domínio sem TLD", "user@domain", false},
		{"Email vazio", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidarIdade(t *testing.T) {
	agora := time.Now()

	tests := []struct {
		name           string
		dataNascimento time.Time
		expectError    bool
		errorMessage   string
	}{
		{
			"Idade válida - 25 anos",
			agora.AddDate(-25, 0, 0),
			false,
			"",
		},
		{
			"Idade válida - exatamente 18 anos",
			agora.AddDate(-18, 0, 0),
			false,
			"",
		},
		{
			"Idade inválida - 17 anos",
			agora.AddDate(-17, 0, 0),
			true,
			"motorista deve ter pelo menos 18 anos",
		},
		{
			"Idade inválida - 10 anos",
			agora.AddDate(-10, 0, 0),
			true,
			"motorista deve ter pelo menos 18 anos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarIdade(tt.dataNascimento)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidarValidadeCNH(t *testing.T) {
	agora := time.Now()

	tests := []struct {
		name         string
		validadeCNH  time.Time
		expectError  bool
		errorMessage string
	}{
		{
			"CNH válida - vence em 1 ano",
			agora.AddDate(1, 0, 0),
			false,
			"",
		},
		{
			"CNH válida - vence hoje",
			agora,
			false,
			"",
		},
		{
			"CNH vencida - venceu ontem",
			agora.AddDate(0, 0, -1),
			true,
			"CNH vencida. Renove sua CNH para prosseguir",
		},
		{
			"CNH vencida - venceu há 1 ano",
			agora.AddDate(-1, 0, 0),
			true,
			"CNH vencida. Renove sua CNH para prosseguir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarValidadeCNH(tt.validadeCNH)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidarForcaSenha(t *testing.T) {
	tests := []struct {
		name          string
		senha         string
		expectedForce string
		expectError   bool
		errorMessage  string
	}{
		{
			"Senha fraca - muito curta",
			"123456",
			"Fraca",
			true,
			"senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo",
		},
		{
			"Senha fraca - sem maiúscula",
			"minhasenha123!",
			"Fraca",
			true,
			"senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo",
		},
		{
			"Senha fraca - sem minúscula",
			"MINHASENHA123!",
			"Fraca",
			true,
			"senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo",
		},
		{
			"Senha fraca - sem número",
			"MinhaSenha!",
			"Fraca",
			true,
			"senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo",
		},
		{
			"Senha fraca - sem símbolo",
			"MinhaSenha123",
			"Fraca",
			true,
			"senha deve ter pelo menos 8 caracteres, incluindo maiúscula, minúscula, número e símbolo",
		},
		{
			"Senha média - todos critérios mas curta",
			"MinhaS1!",
			"Média",
			false,
			"",
		},
		{
			"Senha forte",
			"MinhaSenh@123",
			"Forte",
			false,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			force, err := ValidarForcaSenha(tt.senha)

			assert.Equal(t, tt.expectedForce, force)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidarDocumento(t *testing.T) {
	tests := []struct {
		name         string
		formato      string
		tamanho      int64
		expectError  bool
		errorMessage string
	}{
		{
			"Documento válido - JPG 2MB",
			"JPG",
			2 * 1024 * 1024,
			false,
			"",
		},
		{
			"Documento válido - PNG 1.5MB",
			"PNG",
			int64(1.5 * 1024 * 1024),
			false,
			"",
		},
		{
			"Documento válido - PDF 1MB",
			"PDF",
			1 * 1024 * 1024,
			false,
			"",
		},
		{
			"Formato inválido - TXT",
			"TXT",
			1 * 1024 * 1024,
			true,
			"formato não suportado. Use JPG, PNG ou PDF",
		},
		{
			"Arquivo muito grande - 6MB",
			"JPG",
			6 * 1024 * 1024,
			true,
			"arquivo muito grande. Tamanho máximo: 5MB",
		},
		{
			"Formato minúsculo válido",
			"jpg",
			1 * 1024 * 1024,
			false,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarDocumento(tt.formato, tt.tamanho)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
