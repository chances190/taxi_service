package services

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockSMTPServer is a minimal SMTP mock for tests
type MockSMTPServer struct {
	listener  net.Listener
	port      int
	messages  []ReceivedMessage
	mu        sync.Mutex
	running   bool
	responses map[string]string
}

type ReceivedMessage struct {
	From    string
	To      []string
	Subject string
	Body    string
	Headers map[string]string
}

func NewMockSMTPServer() *MockSMTPServer {
	return &MockSMTPServer{
		messages:  []ReceivedMessage{},
		responses: map[string]string{},
	}
}

func (s *MockSMTPServer) Start() error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	s.listener = ln
	s.port = ln.Addr().(*net.TCPAddr).Port
	s.running = true
	go s.acceptLoop()
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (s *MockSMTPServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *MockSMTPServer) GetPort() int { return s.port }

func (s *MockSMTPServer) GetMessages() []ReceivedMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]ReceivedMessage, len(s.messages))
	copy(cp, s.messages)
	return cp
}

func (s *MockSMTPServer) ClearMessages() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = []ReceivedMessage{}
}

func (s *MockSMTPServer) SetResponse(cmd, resp string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responses[strings.ToUpper(cmd)] = resp
}

func (s *MockSMTPServer) acceptLoop() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.running {
				return
			}
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *MockSMTPServer) handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	w.WriteString("220 localhost ESMTP Mock\r\n")
	w.Flush()

	var msg ReceivedMessage
	msg.Headers = map[string]string{}
	var dataMode bool
	var dataLines []string

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if dataMode {
			if line == "." {
				dataMode = false
				s.parseMsg(&msg, strings.Join(dataLines, "\n"))
				s.mu.Lock()
				s.messages = append(s.messages, msg)
				s.mu.Unlock()
				w.WriteString("250 OK\r\n")
				w.Flush()
				msg = ReceivedMessage{Headers: map[string]string{}}
				dataLines = nil
			} else {
				dataLines = append(dataLines, line)
			}
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])
		s.mu.Lock()
		resp, ok := s.responses[cmd]
		s.mu.Unlock()
		if ok {
			w.WriteString(resp + "\r\n")
			w.Flush()
			continue
		}
		switch cmd {
		case "HELO", "EHLO":
			w.WriteString("250-localhost\r\n250-AUTH PLAIN LOGIN\r\n250 OK\r\n")
		case "AUTH":
			w.WriteString("235 Authentication successful\r\n")
		case "MAIL":
			if len(parts) > 1 && strings.HasPrefix(parts[1], "FROM:") {
				msg.From = strings.Trim(parts[1][5:], "<>")
			}
			w.WriteString("250 OK\r\n")
		case "RCPT":
			if len(parts) > 1 && strings.HasPrefix(parts[1], "TO:") {
				msg.To = append(msg.To, strings.Trim(parts[1][3:], "<>"))
			}
			w.WriteString("250 OK\r\n")
		case "DATA":
			w.WriteString("354 End data with <CRLF>.<CRLF>\r\n")
			dataMode = true
		case "QUIT":
			w.WriteString("221 Bye\r\n")
			w.Flush()
			return
		default:
			w.WriteString("500 Unknown\r\n")
		}
		w.Flush()
	}
}

func (s *MockSMTPServer) parseMsg(msg *ReceivedMessage, content string) {
	lines := strings.Split(content, "\n")
	bodyStart := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			bodyStart = i + 1
			break
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			k := strings.TrimSpace(line[:idx])
			v := strings.TrimSpace(line[idx+1:])
			msg.Headers[k] = v
			if strings.ToLower(k) == "subject" {
				msg.Subject = v
			}
		}
	}
	if bodyStart < len(lines) {
		msg.Body = strings.Join(lines[bodyStart:], "\n")
	}
}

func TestSMTPEmailService(t *testing.T) {
	// Teste usando configuração manual
	t.Run("Criar serviço com configuração manual", func(t *testing.T) {
		config := EmailConfig{
			Host:     "smtp.gmail.com",
			Port:     587,
			Username: "test@gmail.com",
			Password: "password",
			From:     "noreply@taxiservice.com",
		}

		service := NewSMTPEmailService(config)

		assert.Equal(t, "smtp.gmail.com", service.host)
		assert.Equal(t, 587, service.port)
		assert.Equal(t, "test@gmail.com", service.username)
		assert.Equal(t, "password", service.password)
		assert.Equal(t, "noreply@taxiservice.com", service.from)
	})

	// Teste usando variáveis de ambiente
	t.Run("Criar serviço a partir de variáveis de ambiente", func(t *testing.T) {
		// Configurar variáveis de ambiente temporárias
		os.Setenv("SMTP_HOST", "smtp.test.com")
		os.Setenv("SMTP_PORT", "465")
		os.Setenv("SMTP_USERNAME", "test@test.com")
		os.Setenv("SMTP_PASSWORD", "testpass")
		os.Setenv("SMTP_FROM", "test@taxiservice.com")

		defer func() {
			// Limpar variáveis de ambiente após o teste
			os.Unsetenv("SMTP_HOST")
			os.Unsetenv("SMTP_PORT")
			os.Unsetenv("SMTP_USERNAME")
			os.Unsetenv("SMTP_PASSWORD")
			os.Unsetenv("SMTP_FROM")
		}()

		service := NewSMTPEmailServiceFromEnv()

		assert.Equal(t, "smtp.test.com", service.host)
		assert.Equal(t, 465, service.port)
		assert.Equal(t, "test@test.com", service.username)
		assert.Equal(t, "testpass", service.password)
		assert.Equal(t, "test@taxiservice.com", service.from)
	})

	// Teste com valores padrão quando variáveis não existem
	t.Run("Usar valores padrão quando variáveis não existem", func(t *testing.T) {
		service := NewSMTPEmailServiceFromEnv()

		assert.Equal(t, "smtp.gmail.com", service.host)
		assert.Equal(t, 587, service.port)
		assert.Equal(t, "noreply@taxiservice.com", service.from)
		// username e password podem estar vazios se não configurados
	})

	// Teste de envio de email sem credenciais deve retornar erro
	t.Run("Envio de email sem credenciais retorna erro", func(t *testing.T) {
		config := EmailConfig{
			Host:     "smtp.gmail.com",
			Port:     587,
			Username: "", // Sem credenciais
			Password: "",
			From:     "noreply@taxiservice.com",
		}

		service := NewSMTPEmailService(config)

		err := service.EnviarEmailConfirmacao("test@example.com", "João Silva")
		assert.Error(t, err)

		err = service.EnviarEmailRecebimentoDocumentos("test@example.com", "João Silva")
		assert.Error(t, err)

		err = service.EnviarEmailAprovacao("test@example.com", "João Silva")
		assert.Error(t, err)

		err = service.EnviarEmailRejeicao("test@example.com", "João Silva", "Documentos com problemas de qualidade")
		assert.Error(t, err)
	})

}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			"Retorna valor padrão quando variável não existe",
			"TEST_VAR_NOT_EXISTS",
			"default",
			"",
			"default",
		},
		{
			"Retorna valor da variável quando existe",
			"TEST_VAR_EXISTS",
			"default",
			"environment_value",
			"environment_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailFunctionalityWithMockSMTPServer(t *testing.T) {
	// Criar e iniciar o servidor SMTP mock
	mockServer := NewMockSMTPServer()
	err := mockServer.Start()
	assert.NoError(t, err)
	defer mockServer.Stop()

	// Configuração para usar o servidor mock
	config := EmailConfig{
		Host:     "127.0.0.1",
		Port:     mockServer.GetPort(),
		Username: "test@test.com",
		Password: "password123",
		From:     "noreply@taxiservice.com",
	}

	t.Run("Envio de email de confirmação com sucesso", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		err := service.EnviarEmailConfirmacao("joao@example.com", "João Silva")
		assert.NoError(t, err)

		// Aguardar um pouco para a mensagem ser processada
		time.Sleep(50 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 1)

		msg := messages[0]
		assert.Equal(t, "noreply@taxiservice.com", msg.From)
		assert.Equal(t, []string{"joao@example.com"}, msg.To)
		assert.Contains(t, msg.Subject, "Cadastro realizado com sucesso")
		assert.Contains(t, msg.Body, "João Silva")
		assert.Contains(t, msg.Body, "Bem-vindo ao Taxi Service!")
	})

	t.Run("Envio de email de recebimento de documentos", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		err := service.EnviarEmailRecebimentoDocumentos("maria@example.com", "Maria Santos")
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 1)

		msg := messages[0]
		assert.Equal(t, "noreply@taxiservice.com", msg.From)
		assert.Equal(t, []string{"maria@example.com"}, msg.To)
		assert.Contains(t, msg.Subject, "Documentos recebidos")
		assert.Contains(t, msg.Body, "Maria Santos")
		assert.Contains(t, msg.Body, "Recebemos seus documentos")
	})

	t.Run("Envio de email de aprovação", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		err := service.EnviarEmailAprovacao("carlos@example.com", "Carlos Lima")
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 1)

		msg := messages[0]
		assert.Equal(t, "noreply@taxiservice.com", msg.From)
		assert.Equal(t, []string{"carlos@example.com"}, msg.To)
		assert.Contains(t, msg.Subject, "cadastro foi aprovado")
		assert.Contains(t, msg.Body, "Carlos Lima")
		assert.Contains(t, msg.Body, "🎉 Cadastro Aprovado!")
	})

	t.Run("Envio de email de rejeição", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		motivo := "Documentos com baixa qualidade de imagem"
		err := service.EnviarEmailRejeicao("ana@example.com", "Ana Costa", motivo)
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 1)

		msg := messages[0]
		assert.Equal(t, "noreply@taxiservice.com", msg.From)
		assert.Equal(t, []string{"ana@example.com"}, msg.To)
		assert.Contains(t, msg.Subject, "Documentos rejeitados")
		assert.Contains(t, msg.Body, "Ana Costa")
		assert.Contains(t, msg.Body, motivo)
	})

	t.Run("Múltiplos emails em sequência", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		// Enviar vários emails
		err1 := service.EnviarEmailConfirmacao("user1@example.com", "User 1")
		err2 := service.EnviarEmailRecebimentoDocumentos("user2@example.com", "User 2")
		err3 := service.EnviarEmailAprovacao("user3@example.com", "User 3")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)

		time.Sleep(100 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 3)

		// Verificar que todos os emails foram recebidos
		recipients := make([]string, 0)
		for _, msg := range messages {
			recipients = append(recipients, msg.To[0])
		}

		assert.Contains(t, recipients, "user1@example.com")
		assert.Contains(t, recipients, "user2@example.com")
		assert.Contains(t, recipients, "user3@example.com")
	})

	t.Run("Verificação de formato HTML dos emails", func(t *testing.T) {
		service := NewSMTPEmailService(config)
		mockServer.ClearMessages()

		err := service.EnviarEmailConfirmacao("test@example.com", "Test User")
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		messages := mockServer.GetMessages()
		assert.Len(t, messages, 1)

		msg := messages[0]
		assert.Contains(t, msg.Headers["Content-Type"], "text/html")
		assert.Contains(t, msg.Body, "<html>")
		assert.Contains(t, msg.Body, "<body>")
		assert.Contains(t, msg.Body, "</body>")
		assert.Contains(t, msg.Body, "</html>")
	})
}

func TestSMTPServerErrorSimulation(t *testing.T) {
	t.Run("Simulação de falha de autenticação", func(t *testing.T) {
		mockServer := NewMockSMTPServer()
		err := mockServer.Start()
		assert.NoError(t, err)
		defer mockServer.Stop()

		// Configurar resposta de erro para AUTH
		mockServer.SetResponse("AUTH", "535 Authentication failed")

		config := EmailConfig{
			Host:     "127.0.0.1",
			Port:     mockServer.GetPort(),
			Username: "test@test.com",
			Password: "wrongpassword",
			From:     "noreply@taxiservice.com",
		}

		service := NewSMTPEmailService(config)
		err = service.EnviarEmailConfirmacao("test@example.com", "Test User")

		// O erro pode vir da autenticação falhada
		// Nota: Dependendo da implementação, pode ser necessário ajustar este teste
		assert.Error(t, err)
	})

	t.Run("Simulação de rejeição de destinatário", func(t *testing.T) {
		mockServer := NewMockSMTPServer()
		err := mockServer.Start()
		assert.NoError(t, err)
		defer mockServer.Stop()

		// Configurar resposta de erro para RCPT
		mockServer.SetResponse("RCPT", "550 User not found")

		config := EmailConfig{
			Host:     "127.0.0.1",
			Port:     mockServer.GetPort(),
			Username: "test@test.com",
			Password: "password123",
			From:     "noreply@taxiservice.com",
		}

		service := NewSMTPEmailService(config)
		err = service.EnviarEmailConfirmacao("invalid@example.com", "Test User")

		assert.Error(t, err)
	})
}

func TestEmailContentValidation(t *testing.T) {
	mockServer := NewMockSMTPServer()
	err := mockServer.Start()
	assert.NoError(t, err)
	defer mockServer.Stop()

	config := EmailConfig{
		Host:     "127.0.0.1",
		Port:     mockServer.GetPort(),
		Username: "test@test.com",
		Password: "password123",
		From:     "noreply@taxiservice.com",
	}

	service := NewSMTPEmailService(config)

	testCases := []struct {
		name            string
		emailFunc       func() error
		expectedSubject string
		expectedContent []string
	}{
		{
			name: "Email de confirmação",
			emailFunc: func() error {
				return service.EnviarEmailConfirmacao("user@example.com", "Usuario Teste")
			},
			expectedSubject: "Cadastro realizado com sucesso",
			expectedContent: []string{"Bem-vindo ao Taxi Service!", "Usuario Teste", "próximo passo"},
		},
		{
			name: "Email de recebimento de documentos",
			emailFunc: func() error {
				return service.EnviarEmailRecebimentoDocumentos("user@example.com", "Usuario Teste")
			},
			expectedSubject: "Documentos recebidos",
			expectedContent: []string{"Recebemos seus documentos", "Usuario Teste", "2 dias úteis"},
		},
		{
			name: "Email de aprovação",
			emailFunc: func() error {
				return service.EnviarEmailAprovacao("user@example.com", "Usuario Teste")
			},
			expectedSubject: "cadastro foi aprovado",
			expectedContent: []string{"🎉 Cadastro Aprovado!", "Usuario Teste", "família Taxi Service"},
		},
		{
			name: "Email de rejeição",
			emailFunc: func() error {
				return service.EnviarEmailRejeicao("user@example.com", "Usuario Teste", "Documento ilegível")
			},
			expectedSubject: "Documentos rejeitados",
			expectedContent: []string{"Documentos Rejeitados", "Usuario Teste", "Documento ilegível"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer.ClearMessages()

			err := tc.emailFunc()
			assert.NoError(t, err)

			time.Sleep(50 * time.Millisecond)

			messages := mockServer.GetMessages()
			assert.Len(t, messages, 1)

			msg := messages[0]
			assert.Contains(t, msg.Subject, tc.expectedSubject)

			for _, content := range tc.expectedContent {
				assert.Contains(t, msg.Body, content, "Conteúdo esperado não encontrado: %s", content)
			}
		})
	}
}
