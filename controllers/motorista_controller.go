package controllers

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"

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

// (Removidos endpoints JSON de upload individual e em lote para simplificação)

// UploadDocumentosArquivos POST /api/documents/:id/upload/files (multipart)
// Espera campos de formulário: files[] (até 3) e para cada arquivo um campo tipo_{index} com valores CNH|CRLV|selfie_cnh
func (c *MotoristaController) UploadDocumentosArquivos(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	form, err := ctx.MultipartForm()
	if err != nil || form == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "multipart inválido"})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nenhum arquivo enviado"})
	}
	if len(files) > 3 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "máximo de 3 arquivos"})
	}

	// criar diretório do motorista
	baseDir := filepath.Join("data", motoristaID)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "falha ao criar diretório"})
	}

	var uploadRequests []services.UploadDocumentoRequest
	for idx, fh := range files {
		// tipo vem de campo tipo_0, tipo_1...
		tipoField := "tipo_" + strconv.Itoa(idx)
		tipos := form.Value[tipoField]
		if len(tipos) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "campo " + tipoField + " ausente"})
		}
		tipo := tipos[0]
		ext := filepath.Ext(fh.Filename)
		destino := filepath.Join(baseDir, tipo+ext)
		if err := ctx.SaveFile(fh, destino); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "falha ao salvar arquivo"})
		}
		uploadRequests = append(uploadRequests, services.UploadDocumentoRequest{
			TipoDocumento:  tipo,
			CaminhoArquivo: destino,
			Formato:        ext[1:],
			Tamanho:        fh.Size,
		})
	}

	if err := c.motoristaService.UploadDocumentosLote(motoristaID, uploadRequests); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "motorista não encontrado" {
			status = fiber.StatusNotFound
		}
		return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Arquivos enviados", "quantidade": len(uploadRequests)})
}

// DownloadDocumento GET /api/documents/:id/file/:tipo
func (c *MotoristaController) DownloadDocumento(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")
	tipo := ctx.Params("tipo")

	motorista, err := c.motoristaService.BuscarMotorista(motoristaID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "motorista não encontrado"})
	}
	for _, doc := range motorista.Documentos {
		if doc.TipoDocumento == tipo {
			// Verifica se arquivo existe
			if _, err := os.Stat(doc.CaminhoArquivo); err == nil {
				return ctx.SendFile(doc.CaminhoArquivo)
			}
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "arquivo não encontrado"})
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "documento não encontrado"})
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

// (Removido endpoint de validação automática)

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

// (Removido endpoint utilitário de validação de upload; validação ocorre ao salvar)

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
