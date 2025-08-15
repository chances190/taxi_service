package domain

// Centraliza constantes e listas de domínio para evitar strings mágicas.

// Status do motorista
const (
	StatusAguardandoDocumentos = "aguardando_documentos"
	StatusDocumentosAnalise    = "documentos_em_analise"
	StatusAprovado             = "aprovado"
	StatusRejeitado            = "documentos_rejeitados"
	StatusAguardandoExclusao   = "aguardando_exclusao"
	StatusAtivo                = "ativo"
	StatusEncerrado            = "encerrado"
)

// Status de documentos
const (
	DocumentoStatusPendente = "pendente"
	DocumentoStatusAprovado = "aprovado"
)

// Tipos obrigatórios de documentos
var DocumentosObrigatorios = []string{"CNH", "CRLV", "selfie_cnh"}
