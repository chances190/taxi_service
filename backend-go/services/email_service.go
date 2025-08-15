package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

// EmailService define a interface para envio de emails
type EmailService interface {
	EnviarEmailConfirmacao(email, nome string) error
	EnviarEmailRecebimentoDocumentos(email, nome string) error
	EnviarEmailAprovacao(email, nome string) error
	EnviarEmailRejeicao(email, nome, motivo string) error
}

// SMTPEmailService implementação real usando SMTP
type SMTPEmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// EmailConfig configuração para o serviço de email
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewSMTPEmailService cria uma nova instância do serviço SMTP
func NewSMTPEmailService(config EmailConfig) *SMTPEmailService {
	return &SMTPEmailService{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
	}
}

// NewSMTPEmailServiceFromEnv cria uma instância usando variáveis de ambiente
func NewSMTPEmailServiceFromEnv() *SMTPEmailService {
	port, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))

	return &SMTPEmailService{
		host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		port:     port,
		username: getEnvOrDefault("SMTP_USERNAME", ""),
		password: getEnvOrDefault("SMTP_PASSWORD", ""),
		from:     getEnvOrDefault("SMTP_FROM", "noreply@taxiservice.com"),
	}
}

// enviarEmail método interno para enviar emails
func (s *SMTPEmailService) enviarEmail(to, subject, body string) error {
	// Se as credenciais não estiverem configuradas, não enviar (modo desenvolvimento)
	if s.username == "" || s.password == "" {
		return fmt.Errorf("credenciais SMTP não configuradas")
	}

	// Configurar autenticação
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Montar a mensagem
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, subject, body)

	// Conectar e enviar
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Para conexões TLS (porta 587)
	if s.port == 587 {
		return s.enviarComTLS(addr, auth, to, msg)
	}

	// Para conexões SSL (porta 465) ou sem criptografia (porta 25)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

// enviarComTLS envia email usando conexão TLS
func (s *SMTPEmailService) enviarComTLS(addr string, auth smtp.Auth, to, msg string) error {
	// Conectar ao servidor
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao servidor SMTP: %w", err)
	}
	defer conn.Close()

	// Iniciar TLS
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("erro ao iniciar TLS: %w", err)
	}

	// Autenticar
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("erro na autenticação: %w", err)
	}

	// Definir remetente
	if err = conn.Mail(s.from); err != nil {
		return fmt.Errorf("erro ao definir remetente: %w", err)
	}

	// Definir destinatário
	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("erro ao definir destinatário: %w", err)
	}

	// Enviar dados
	writer, err := conn.Data()
	if err != nil {
		return fmt.Errorf("erro ao iniciar envio de dados: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("erro ao escrever mensagem: %w", err)
	}

	return nil
}

// EnviarEmailConfirmacao envia email de confirmação de cadastro
func (s *SMTPEmailService) EnviarEmailConfirmacao(email, nome string) error {
	subject := "Cadastro realizado com sucesso - Taxi Service"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Bem-vindo ao Taxi Service!</h2>
			<p>Olá <strong>%s</strong>,</p>
			<p>Seu cadastro foi realizado com sucesso! Agora você precisa enviar seus documentos para aprovação.</p>
			<p>Em breve você receberá instruções sobre o próximo passo.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Taxi Service</p>
		</body>
		</html>
	`, nome)

	return s.enviarEmail(email, subject, body)
}

// EnviarEmailRecebimentoDocumentos envia email de confirmação de recebimento de documentos
func (s *SMTPEmailService) EnviarEmailRecebimentoDocumentos(email, nome string) error {
	subject := "Documentos recebidos - Taxi Service"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Documentos Recebidos</h2>
			<p>Olá <strong>%s</strong>,</p>
			<p>Recebemos seus documentos e eles estão sendo analisados por nossa equipe.</p>
			<p>O processo de análise pode levar até 2 dias úteis.</p>
			<p>Você será notificado assim que a análise for concluída.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Taxi Service</p>
		</body>
		</html>
	`, nome)

	return s.enviarEmail(email, subject, body)
}

// EnviarEmailAprovacao envia email de aprovação do cadastro
func (s *SMTPEmailService) EnviarEmailAprovacao(email, nome string) error {
	subject := "Parabéns! Seu cadastro foi aprovado - Taxi Service"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>🎉 Cadastro Aprovado!</h2>
			<p>Olá <strong>%s</strong>,</p>
			<p><strong>Parabéns!</strong> Seu cadastro foi aprovado e você já pode começar a trabalhar como motorista.</p>
			<p>Acesse o aplicativo e comece a receber corridas agora mesmo!</p>
			<br>
			<p>Bem-vindo à família Taxi Service!</p>
			<br>
			<p>Atenciosamente,<br>Equipe Taxi Service</p>
		</body>
		</html>
	`, nome)

	return s.enviarEmail(email, subject, body)
}

// EnviarEmailRejeicao envia email de rejeição do cadastro
func (s *SMTPEmailService) EnviarEmailRejeicao(email, nome, motivo string) error {
	subject := "Documentos rejeitados - Taxi Service"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Documentos Rejeitados</h2>
			<p>Olá <strong>%s</strong>,</p>
			<p>Infelizmente seus documentos foram rejeitados por nossa equipe.</p>
			<p><strong>Motivo:</strong> %s</p>
			<p>Você pode corrigir os problemas identificados e reenviar seus documentos.</p>
			<p>Se tiver dúvidas, entre em contato conosco.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Taxi Service</p>
		</body>
		</html>
	`, nome, motivo)

	return s.enviarEmail(email, subject, body)
}

// getEnvOrDefault obtém variável de ambiente ou retorna valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
