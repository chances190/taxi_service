package test

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestMotorista(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioMotorista,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/motorista/"},
			Output:   colors.Colored(os.Stdout),
			TestingT: t,
			Strict:   true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func enviouDocumentosComProblemasDeQualidade(arg1 string) error {
	return godog.ErrPending
}


func enviouDocumentosVlidos(arg1 string) error {
	return godog.ErrPending
}


func estouAutenticadoComo(arg1 string) error {
	return godog.ErrPending
}


func estouNaPgina(arg1 string) error {
	return godog.ErrPending
}


func euEstouNaPgina(arg1 string) error {
	return godog.ErrPending
}


func euEstouNaPginaAtravsDeLinkVlidoPara(arg1, arg2 string) error {
	return godog.ErrPending
}


func euVejoAMensagem(arg1 string) error {
	return godog.ErrPending
}


func existeUmMotoristaCadastradoComODadoNoValor(arg1, arg2 string) error {
	return godog.ErrPending
}


func existeUmMotoristaCadastradoComOsDados(arg1 *godog.Table) error {
	return godog.ErrPending
}


func faoUploadDosDocumentosObrigatrios(arg1 *godog.Table) error {
	return godog.ErrPending
}


func noEstouAutenticadoComoMotorista() error {
	return godog.ErrPending
}


func noExisteUmMotoristaCadastradoComOsDados(arg1 *godog.Table) error {
	return godog.ErrPending
}


func oSistemaOuAnalistaVerificaOsDocumentosDe(arg1 string) error {
	return godog.ErrPending
}


func oStatusDe(arg1, arg2 string) error {
	return godog.ErrPending
}


func preenchoOCampoComOValor(arg1, arg2 string) error {
	return godog.ErrPending
}


func preenchoOsCampos(arg1 *godog.Table) error {
	return godog.ErrPending
}


func queEstouNaPginaAtravsDeLinkVlidoPara(arg1, arg2 string) error {
	return godog.ErrPending
}


func queOStatusDe(arg1, arg2 string) error {
	return godog.ErrPending
}


func realizoLoginComEmailESenha(arg1, arg2 string) error {
	return godog.ErrPending
}


func solicitoExcluirMinhaContaPermanentemente() error {
	return godog.ErrPending
}


func submetoOCadastroComOsDados(arg1 *godog.Table) error {
	return godog.ErrPending
}


func submetoOFormulrio() error {
	return godog.ErrPending
}


func tentoFazerUploadDeArquivoComTamanhoEFormato(arg1, arg2 string) error {
	return godog.ErrPending
}


func tentoSalvarAsAlteraes() error {
	return godog.ErrPending
}


func umEmailEnviadoParaComOAssunto(arg1, arg2 string) error {
	return godog.ErrPending
}


func vejoAMensagem(arg1 string) error {
	return godog.ErrPending
}


func vejoAMensagemDeErro(arg1 string) error {
	return godog.ErrPending
}


func vejoMeuDados(arg1 *godog.Table) error {
	return godog.ErrPending
}


func vejoMinhaFotoDePerfil() error {
	return godog.ErrPending
}


func vejoOCampoComValor(arg1, arg2 string) error {
	return godog.ErrPending
}


func InitializeScenarioMotorista(ctx *godog.ScenarioContext) {
		ctx.Step(`^"([^"]*)" enviou documentos com problemas de qualidade$`, enviouDocumentosComProblemasDeQualidade)
		ctx.Step(`^"([^"]*)" enviou documentos válidos$`, enviouDocumentosVlidos)
		ctx.Step(`^estou autenticado como "([^"]*)"$`, estouAutenticadoComo)
		ctx.Step(`^estou na página "([^"]*)"$`, estouNaPgina)
		ctx.Step(`^eu estou na página "([^"]*)" através de link válido para "([^"]*)"$`, euEstouNaPginaAtravsDeLinkVlidoPara)
		ctx.Step(`^eu estou na página "([^"]*)"$`, euEstouNaPgina)
		ctx.Step(`^eu vejo a mensagem "([^"]*)"$`, euVejoAMensagem)
		ctx.Step(`^existe um motorista cadastrado com o dado "([^"]*)" no valor "([^"]*)"$`, existeUmMotoristaCadastradoComODadoNoValor)
		ctx.Step(`^existe um motorista cadastrado com os dados:$`, existeUmMotoristaCadastradoComOsDados)
		ctx.Step(`^faço upload dos documentos obrigatórios:$`, faoUploadDosDocumentosObrigatrios)
		ctx.Step(`^não estou autenticado como motorista$`, noEstouAutenticadoComoMotorista)
		ctx.Step(`^não existe um motorista cadastrado com os dados:$`, noExisteUmMotoristaCadastradoComOsDados)
		ctx.Step(`^o sistema ou analista verifica os documentos de "([^"]*)"$`, oSistemaOuAnalistaVerificaOsDocumentosDe)
		ctx.Step(`^o status de ""([^"]*)" é "([^"]*)"$`, oStatusDe)
		ctx.Step(`^o status de "([^"]*)" é "([^"]*)"$`, oStatusDe)
		ctx.Step(`^o status de "([^"]*)" é "([^"]*)":$`, oStatusDe)
		ctx.Step(`^preencho o campo "([^"]*)" com o valor "([^"]*)"$`, preenchoOCampoComOValor)
		ctx.Step(`^preencho os campos:$`, preenchoOsCampos)
		ctx.Step(`^que estou na página "([^"]*)" através de link válido para "([^"]*)"$`, queEstouNaPginaAtravsDeLinkVlidoPara)
		ctx.Step(`^que o status de ""([^"]*)" é "([^"]*)"$`, queOStatusDe)
		ctx.Step(`^realizo login com email "([^"]*)" e senha "([^"]*)"$`, realizoLoginComEmailESenha)
		ctx.Step(`^solicito excluir minha conta permanentemente$`, solicitoExcluirMinhaContaPermanentemente)
		ctx.Step(`^submeto o cadastro com os dados:$`, submetoOCadastroComOsDados)
		ctx.Step(`^submeto o formulário$`, submetoOFormulrio)
		ctx.Step(`^tento fazer upload de arquivo com tamanho "([^"]*)" e formato "([^"]*)"$`, tentoFazerUploadDeArquivoComTamanhoEFormato)
		ctx.Step(`^tento salvar as alterações$`, tentoSalvarAsAlteraes)
		ctx.Step(`^um email é enviado para "([^"]*)" com o assunto "([^"]*)"$`, umEmailEnviadoParaComOAssunto)
		ctx.Step(`^vejo a mensagem "([^"]*)"$`, vejoAMensagem)
		ctx.Step(`^vejo a mensagem de erro "([^"]*)"$`, vejoAMensagemDeErro)
		ctx.Step(`^vejo meu dados:$`, vejoMeuDados)
		ctx.Step(`^vejo minha foto de perfil$`, vejoMinhaFotoDePerfil)
		ctx.Step(`^vejo o campo "([^"]*)" com valor "([^"]*)"$`, vejoOCampoComValor)
}
