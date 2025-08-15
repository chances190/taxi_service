# Taxi Service (Driver Prototype)

Frontend: React + Vite, testes com Cypress
Backend: Go + Fiber, testes com Godog

Features: Cadastro e Manutenção de Motorista, Recuperação de Conta por E-Mail

## Quick Start
Backend:
```bash
go mod download
go build -o taxi_service
./taxi_service
```
Frontend:
```bash
cd frontend-react
npm install
npm run dev
```
Base URL API: http://localhost:3000

## Rotas

| Método  | Rota                                      | Descrição                              |
|---------|-------------------------------------------|----------------------------------------|
| POST    | /api/auth/register                        | Registro de usuário                    |
| POST    | /api/auth/login                           | Login de usuário                       |
| GET     | /api/profile/:id                          | Obter perfil do usuário                |
| PUT     | /api/profile/:id                          | Atualizar perfil do usuário            |
| PUT     | /api/profile/:id/password                 | Alterar senha do usuário               |
| POST    | /api/profile/:id/photo                    | Enviar foto de perfil                  |
| GET     | /api/profile/:id/photo                    | Obter foto de perfil                   |
| POST    | /api/profile/:id/request-deletion         | Solicitar exclusão de perfil           |
| POST    | /api/profile/:id/confirm-deletion         | Confirmar exclusão de perfil           |
| POST    | /api/documents/:id/upload/files           | Enviar arquivos de documentos          |
| GET     | /api/documents/:id/file/:tipo             | Obter arquivo de documento específico  |
| PUT     | /api/documents/:id/approve                | Aprovar documento                      |
| PUT     | /api/documents/:id/reject                 | Rejeitar documento                     |
| POST    | /api/utils/check-password                 | Verificar senha                        |
| GET     | /health                                   | Verificar saúde da aplicação           |

## Próximos Passos

* Testes E2E com Cypress
* Autenticação com token

