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
	profile.Get("/:id", motoristaController.BuscarMotorista) // Buscar motorista

	// Rotas de documentos
	documents := apiGroup.Group("/documents")
	documents.Post("/:id/upload", motoristaController.UploadDocumento)     // Upload de documentos
	documents.Post("/:id/validate", motoristaController.ValidarDocumentos) // Validar documentos
	documents.Put("/:id/approve", motoristaController.AprovarMotorista)    // Aprovar motorista
	documents.Put("/:id/reject", motoristaController.RejeitarMotorista)    // Rejeitar motorista

	// Rotas utilitárias
	utils := apiGroup.Group("/utils")
	utils.Post("/check-password", motoristaController.VerificarForcaSenha)     // Verificar força da senha
	utils.Post("/validate-upload", motoristaController.ValidarDocumentoUpload) // Validar upload de documento
}
