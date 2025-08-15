package routes

import (
	"taxi_service/controllers"
	"taxi_service/repositories"
	"taxi_service/services"

	"github.com/gofiber/fiber/v2"
)

func SetupMotoristaRoutes(api fiber.Router) {
	// Inicializar dependências
	motoristaRepo := repositories.NewJSONMotoristaRepository()
	emailService := services.NewSMTPEmailServiceFromEnv()
	motoristaService := services.NewMotoristaService(motoristaRepo, emailService)
	motoristaController := controllers.NewMotoristaController(motoristaService)

	// Grupo de rotas da API
	apiGroup := api.Group("/api")

	// Rotas de autenticação
	auth := apiGroup.Group("/auth")
	auth.Post("/register", motoristaController.CadastrarMotorista) // Cadastro de motorista
	auth.Post("/login", motoristaController.LoginMotorista)        // Login de motorista

	// Rotas de perfil
	profile := apiGroup.Group("/profile")
	profile.Get("/:id", motoristaController.BuscarMotorista)                     // Buscar motorista
	profile.Put("/:id", motoristaController.AtualizarPerfil)                     // Atualizar telefone/email
	profile.Put("/:id/password", motoristaController.AlterarSenha)               // Alterar senha
	profile.Post("/:id/photo", motoristaController.UploadFotoPerfil)             // Upload foto
	profile.Get("/:id/photo", motoristaController.FotoPerfil)                    // Obter foto
	profile.Post("/:id/request-deletion", motoristaController.SolicitarExclusao) // Solicitar exclusão
	profile.Post("/:id/confirm-deletion", motoristaController.ConfirmarExclusao) // Confirmar exclusão

	// Rotas de documentos
	documents := apiGroup.Group("/documents")
	documents.Post("/:id/upload/files", motoristaController.UploadDocumentosArquivos) // Upload múltiplo multipart (arquivos reais)
	documents.Get("/:id/file/:tipo", motoristaController.DownloadDocumento)           // Download/visualização de arquivo
	documents.Put("/:id/approve", motoristaController.AprovarMotorista)               // Aprovar motorista
	documents.Put("/:id/reject", motoristaController.RejeitarMotorista)               // Rejeitar motorista

	// Rotas utilitárias
	utils := apiGroup.Group("/utils")
	utils.Post("/check-password", motoristaController.VerificarForcaSenha) // Verificar força da senha
}
