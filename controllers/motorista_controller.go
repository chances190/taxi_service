package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"taxi_service/models"
	"taxi_service/services"
)

// MotoristaController gerencia as rotas relacionadas a motoristas
type MotoristaController struct {
	motoristaService services.MotoristaService
}

// NewMotoristaController cria uma nova instância do controller
func NewMotoristaController(motoristaService services.MotoristaService) *MotoristaController {
	return &MotoristaController{
		motoristaService: motoristaService,
	}
}

// CadastrarMotorista POST /api/motoristas
func (c *MotoristaController) CadastrarMotorista(ctx *fiber.Ctx) error {
	var request services.CadastroMotoristaRequest

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	motorista, err := c.motoristaService.CadastrarMotorista(request)
	if err != nil {
		// Determinar o status code baseado no tipo de erro
		statusCode := fiber.StatusBadRequest

		switch err.Error() {
		case "CPF já cadastrado", "CNH já cadastrada", "Email já cadastrado":
			statusCode = fiber.StatusConflict
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Cadastro realizado com sucesso",
		"motorista": fiber.Map{
			"id":     motorista.ID,
			"nome":   motorista.Nome,
			"email":  motorista.Email,
			"status": motorista.Status,
		},
	})
}

// UploadDocumento POST /api/motoristas/:id/documentos
func (c *MotoristaController) UploadDocumento(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	var request services.UploadDocumentoRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	err := c.motoristaService.UploadDocumento(motoristaID, request)
	if err != nil {
		statusCode := fiber.StatusBadRequest

		if err.Error() == "Motorista não encontrado" {
			statusCode = fiber.StatusNotFound
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Documentos enviados com sucesso",
	})
}

// BuscarMotorista GET /api/motoristas/:id
func (c *MotoristaController) BuscarMotorista(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	motorista, err := c.motoristaService.BuscarMotorista(motoristaID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Motorista não encontrado",
		})
	}

	return ctx.JSON(fiber.Map{
		"motorista": fiber.Map{
			"id":             motorista.ID,
			"nome":           motorista.Nome,
			"email":          motorista.Email,
			"telefone":       motorista.Telefone,
			"status":         motorista.Status,
			"modelo_veiculo": motorista.ModeloVeiculo,
			"placa_veiculo":  motorista.PlacaVeiculo,
			"criado_em":      motorista.CriadoEm,
			"documentos":     motorista.Documentos,
		},
	})
}

// ValidarDocumentos POST /api/motoristas/:id/validar-documentos
func (c *MotoristaController) ValidarDocumentos(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	err := c.motoristaService.ValidarDocumentos(motoristaID)
	if err != nil {
		statusCode := fiber.StatusBadRequest

		if err.Error() == "Motorista não encontrado" {
			statusCode = fiber.StatusNotFound
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Documentos validados com sucesso",
	})
}

// AprovarMotorista PUT /api/motoristas/:id/aprovar
func (c *MotoristaController) AprovarMotorista(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	err := c.motoristaService.AprovarMotorista(motoristaID)
	if err != nil {
		statusCode := fiber.StatusBadRequest

		if err.Error() == "Motorista não encontrado" {
			statusCode = fiber.StatusNotFound
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Motorista aprovado com sucesso",
	})
}

// RejeitarMotorista PUT /api/motoristas/:id/rejeitar
func (c *MotoristaController) RejeitarMotorista(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	var request struct {
		Motivo string `json:"motivo" validate:"required"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	if request.Motivo == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Motivo é obrigatório",
		})
	}

	err := c.motoristaService.RejeitarMotorista(motoristaID, request.Motivo)
	if err != nil {
		statusCode := fiber.StatusBadRequest

		if err.Error() == "Motorista não encontrado" {
			statusCode = fiber.StatusNotFound
		}

		return ctx.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Motorista rejeitado",
	})
}

// VerificarForcaSenha POST /api/motoristas/verificar-senha
func (c *MotoristaController) VerificarForcaSenha(ctx *fiber.Ctx) error {
	var request struct {
		Senha string `json:"senha" validate:"required"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	force, err := c.motoristaService.VerificarForcaSenha(request.Senha)

	response := fiber.Map{
		"forca": force,
	}

	if err != nil {
		response["message"] = err.Error()
	}

	return ctx.JSON(response)
}

// ValidarDocumentoUpload validação de upload de arquivo
func (c *MotoristaController) ValidarDocumentoUpload(ctx *fiber.Ctx) error {
	var request struct {
		Formato string `json:"formato" validate:"required"`
		Tamanho string `json:"tamanho" validate:"required"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	tamanho, err := strconv.ParseInt(request.Tamanho, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tamanho inválido",
		})
	}

	if err := models.ValidarDocumento(request.Formato, tamanho); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Documento válido",
	})
}

// LoginMotorista POST /api/auth/login
func (c *MotoristaController) LoginMotorista(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email" validate:"required,email"`
		Senha string `json:"senha" validate:"required"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}

	motorista, err := c.motoristaService.LoginMotorista(request.Email, request.Senha)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Login realizado com sucesso",
		"motorista": fiber.Map{
			"id":    motorista.ID,
			"nome":  motorista.Nome,
			"email": motorista.Email,
		},
	})
}
