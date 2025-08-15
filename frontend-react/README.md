# Taxi Service Frontend (React)

Frontend React moderno (Mobile-first) com Material UI, tema dark e integração com a API Go.

## Tecnologias
- React 18 + TypeScript
- Vite
- Material UI (MUI v6)
- React Query (@tanstack)
- React Hook Form + Yup
- Axios

## Scripts
- `npm run dev` - ambiente de desenvolvimento
- `npm run build` - build produção
- `npm run preview` - preview da build

## Variáveis de Ambiente
Crie `.env`:
```
VITE_API_URL=http://localhost:3000
```

## Estrutura de Rotas
| Frontend | API | Descrição |
|----------|-----|-----------|
| POST /api/auth/register | /register (page) | Cadastro motorista |
| POST /api/auth/login | /login (page) | Login motorista |
| GET /api/profile/:id | /profile/:id | Perfil motorista |
| POST /api/documents/:id/upload | /documents/:id/upload | Upload docs |
| POST /api/documents/:id/validate | (ação) | Validar docs |
| PUT /api/documents/:id/approve | (ação) | Aprovar motorista |
| PUT /api/documents/:id/reject | (ação) | Rejeitar motorista |

## TODOs Futuro
- Armazenar token (Auth)
- Proteção de rotas
- Feedback granular de validação
- Testes unitários e2e

## Execução
Instalar dependências:
```
npm install
```
Iniciar:
```
npm run dev
```
