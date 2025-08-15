package controllers

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"taxi_service/internal/apperrors"
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
		return apperrors.ErrCampoObrigatorio
	}
	motorista, err := c.motoristaService.CadastrarMotorista(request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Cadastro realizado com sucesso", "motorista": resumoMotorista(motorista)})
}

// (Removidos endpoints JSON de upload individual e em lote para simplificação)

// UploadDocumentosArquivos POST /api/documents/:id/upload/files (multipart)
// Espera campos de formulário: files[] (até 3) e para cada arquivo um campo tipo_{index} com valores CNH|CRLV|selfie_cnh
func (c *MotoristaController) UploadDocumentosArquivos(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	form, err := ctx.MultipartForm()
	if err != nil || form == nil {
		return apperrors.ErrCampoObrigatorio
	}
	files := form.File["files"]
	if len(files) == 0 {
		return apperrors.ErrNenhumDocumentoEnviado
	}
	if len(files) > 3 {
		return apperrors.ErrLimiteArquivosExcedido
	}

	// criar diretório do motorista
	baseDir := filepath.Join("data", motoristaID)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return apperrors.ErrFalhaCriarDiretorio
	}

	var uploadRequests []services.UploadDocumentoRequest
	for idx, fh := range files {
		// tipo vem de campo tipo_0, tipo_1...
		tipoField := "tipo_" + strconv.Itoa(idx)
		tipos := form.Value[tipoField]
		if len(tipos) == 0 {
			return apperrors.ErrCampoObrigatorio
		}
		tipo := tipos[0]
		ext := filepath.Ext(fh.Filename)
		destino := filepath.Join(baseDir, tipo+ext)
		if err := ctx.SaveFile(fh, destino); err != nil {
			return apperrors.ErrFalhaSalvarArquivo
		}
		uploadRequests = append(uploadRequests, services.UploadDocumentoRequest{
			TipoDocumento:  tipo,
			CaminhoArquivo: destino,
			Formato:        ext[1:],
			Tamanho:        fh.Size,
		})
	}

	if err := c.motoristaService.UploadDocumentosLote(motoristaID, uploadRequests); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"message": "Arquivos enviados", "quantidade": len(uploadRequests)})
}

// DownloadDocumento GET /api/documents/:id/file/:tipo
func (c *MotoristaController) DownloadDocumento(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")
	tipo := ctx.Params("tipo")

	motorista, err := c.motoristaService.BuscarMotorista(motoristaID)
	if err != nil {
		return err
	}
	for _, doc := range motorista.Documentos {
		if doc.TipoDocumento == tipo {
			// Verifica se arquivo existe
			if _, err := os.Stat(doc.CaminhoArquivo); err == nil {
				return ctx.SendFile(doc.CaminhoArquivo)
			}
			return apperrors.ErrArquivoNaoEncontrado
		}
	}
	return apperrors.ErrDocumentoNaoEncontrado
}

// BuscarMotorista GET /api/motoristas/:id
func (c *MotoristaController) BuscarMotorista(ctx *fiber.Ctx) error {
	motorista, err := c.motoristaService.BuscarMotorista(ctx.Params("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"motorista": detalhesMotorista(motorista)})
}

// FotoPerfil GET /api/profile/:id/photo
func (c *MotoristaController) FotoPerfil(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")
	motorista, err := c.motoristaService.BuscarMotorista(motoristaID)
	if err != nil {
		return err
	}
	if motorista.FotoPerfil == "" {
		return apperrors.ErrFotoNaoEncontrada
	}
	if _, err := os.Stat(motorista.FotoPerfil); err != nil {
		return apperrors.ErrArquivoNaoEncontrado
	}
	return ctx.SendFile(motorista.FotoPerfil)
}

// (Removido endpoint de validação automática)

// AprovarMotorista PUT /api/motoristas/:id/aprovar
func (c *MotoristaController) AprovarMotorista(ctx *fiber.Ctx) error {
	motoristaID := ctx.Params("id")

	if err := c.motoristaService.AprovarMotorista(motoristaID); err != nil {
		return err
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
		return apperrors.ErrCampoObrigatorio
	}

	if request.Motivo == "" {
		return apperrors.ErrCampoObrigatorio
	}

	if err := c.motoristaService.RejeitarMotorista(motoristaID, request.Motivo); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "Motorista rejeitado",
	})
}

// AtualizarPerfil PUT /api/profile/:id
func (c *MotoristaController) AtualizarPerfil(ctx *fiber.Ctx) error {
	var body struct{ Telefone, Email string }
	if err := ctx.BodyParser(&body); err != nil {
		return apperrors.ErrCampoObrigatorio
	}
	m, err := c.motoristaService.AtualizarPerfil(ctx.Params("id"), body.Telefone, body.Email)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Perfil atualizado com sucesso", "motorista": detalhesMotorista(m)})
}

// AlterarSenha PUT /api/profile/:id/password
func (c *MotoristaController) AlterarSenha(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var body struct{ SenhaAtual, NovaSenha, Confirmacao string }
	if err := ctx.BodyParser(&body); err != nil {
		return apperrors.ErrCampoObrigatorio
	}
	if err := c.motoristaService.AlterarSenha(id, body.SenhaAtual, body.NovaSenha, body.Confirmacao); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Senha alterada com sucesso"})
}

// UploadFotoPerfil POST /api/profile/:id/photo multipart campo 'foto'
func (c *MotoristaController) UploadFotoPerfil(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	fh, err := ctx.FormFile("foto")
	if err != nil {
		return apperrors.ErrCampoObrigatorio
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	formato := strings.TrimPrefix(ext, ".")
	baseDir := filepath.Join("data", id, "profile")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return apperrors.ErrFalhaCriarDiretorio
	}
	destino := filepath.Join(baseDir, "foto"+ext)
	if err := ctx.SaveFile(fh, destino); err != nil {
		return apperrors.ErrFalhaSalvarArquivo
	}
	if err := c.motoristaService.UploadFotoPerfil(id, destino, formato, fh.Size); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Foto de perfil atualizada com sucesso", "caminho": destino})
}

// SolicitarExclusao POST /api/profile/:id/request-deletion
func (c *MotoristaController) SolicitarExclusao(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.motoristaService.SolicitarExclusao(id); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Solicitação de exclusão registrada"})
}

// ConfirmarExclusao POST /api/profile/:id/confirm-deletion
func (c *MotoristaController) ConfirmarExclusao(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.motoristaService.ConfirmarExclusao(id); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Sua conta foi encerrada e será excluida permanentemente em 72h"})
}

// VerificarForcaSenha POST /api/motoristas/verificar-senha
func (c *MotoristaController) VerificarForcaSenha(ctx *fiber.Ctx) error {
	var req struct {
		Senha string `json:"senha"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return apperrors.ErrCampoObrigatorio
	}
	forca, err := c.motoristaService.VerificarForcaSenha(req.Senha)
	resp := fiber.Map{"forca": forca}
	if err != nil {
		resp["message"] = err.Error()
	}
	return ctx.JSON(resp)
}

// (Removido endpoint utilitário de validação de upload; validação ocorre ao salvar)

// LoginMotorista POST /api/auth/login
func (c *MotoristaController) LoginMotorista(ctx *fiber.Ctx) error {
	var req struct{ Email, Senha string }
	if err := ctx.BodyParser(&req); err != nil {
		return apperrors.ErrCampoObrigatorio
	}
	m, err := c.motoristaService.LoginMotorista(req.Email, req.Senha)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{"message": "Login realizado com sucesso", "motorista": resumoMotorista(m)})
}

// --- helpers de serialização ---
func resumoMotorista(m *models.Motorista) fiber.Map {
	return fiber.Map{"id": m.ID, "nome": m.Nome, "email": m.Email, "status": m.Status}
}

func detalhesMotorista(m *models.Motorista) fiber.Map {
	fotoURL := ""
	if m.FotoPerfil != "" {
		fotoURL = "/api/profile/" + m.ID + "/photo"
	}
	return fiber.Map{
		"id":              m.ID,
		"nome":            m.Nome,
		"email":           m.Email,
		"telefone":        m.Telefone,
		"cpf":             m.CPF,
		"cnh":             m.CNH,
		"categoria_cnh":   m.CategoriaCNH,
		"validade_cnh":    m.ValidadeCNH,
		"status":          m.Status,
		"modelo_veiculo":  m.ModeloVeiculo,
		"placa_veiculo":   m.PlacaVeiculo,
		"criado_em":       m.CriadoEm,
		"documentos":      m.Documentos,
		"foto_perfil_url": fotoURL,
	}
}
