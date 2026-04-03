# Sistema RBAC + ABAC

Sistema híbrido de controle de acesso que combina **RBAC** (Controle de Acesso Baseado em Perfis) e **ABAC** (Controle de Acesso Baseado em Atributos), com backend em Go/Gin e frontend em React/Vite.

## Conceito

- **RBAC** define *quem pode fazer o quê*: usuários recebem perfis, perfis contêm permissões.
- **ABAC** define *quando e de onde* essas permissões são válidas: políticas de atributos são vinculadas a perfis e avaliadas a cada requisição.

### Tipos de Política ABAC

| Tipo      | Descrição                                       | Exemplo de valor        |
|-----------|------------------------------------------------|-------------------------|
| `horario` | Restringe acesso a uma faixa de horário         | `08:00-18:00`           |
| `ip`      | Restringe acesso a uma lista de IPs/CIDRs       | `127.0.0.1,::1,10.0.0.0/8` |

Se um perfil **não** possui políticas ABAC, o acesso é irrestrito por atributos (apenas RBAC é avaliado). Se possui, **todas** as políticas devem ser satisfeitas.

## Estrutura do Projeto

```
backend/   — API em Go (Gin, JWT, armazenamento em memória)
frontend/  — SPA em React (Vite, CSS puro)
```

## Executando

### Backend

```bash
cd backend
go run .
```

API inicia em `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Abre em `http://localhost:5173`.

## Dados Iniciais

A aplicação inicia com:

| Entidade   | Nome         | Detalhes                              |
|------------|-------------|---------------------------------------|
| Usuário    | `admin`     | senha: `admin123`, perfil: admin      |
| Perfil     | `admin`     | todas as permissões                   |
| Perfil     | `viewer`    | `user:read`, `role:read` + políticas ABAC |
| Política   | `horário comercial` | tipo: `horario`, valor: `08:00-18:00` |
| Política   | `rede interna`      | tipo: `ip`, valor: `127.0.0.1,::1`   |
| Permissão  | `user:read` |                                       |
| Permissão  | `user:write`|                                       |
| Permissão  | `role:read` |                                       |
| Permissão  | `role:write`|                                       |

## Exemplos de Uso

### Login
```bash
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Listar usuários (protegido — requer `user:read`)
```bash
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>"
```

### Criar um perfil
```bash
curl -s -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"editor"}'
```

### Atribuir perfis a um usuário
```bash
curl -s -X PUT http://localhost:8080/api/users/<user_id>/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"role_ids":["<role_id>"]}'
```

### Atribuir permissões a um perfil
```bash
curl -s -X PUT http://localhost:8080/api/roles/<role_id>/permissions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"permission_ids":["<perm_id>"]}'
```

### Criar uma política ABAC
```bash
curl -s -X POST http://localhost:8080/api/policies \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"horário noturno","type":"horario","value":"18:00-23:59"}'
```

### Atribuir políticas ABAC a um perfil
```bash
curl -s -X PUT http://localhost:8080/api/roles/<role_id>/policies \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"policy_ids":["<policy_id>"]}'
```

## Rotas

| Método | Rota                            | Permissão     |
|--------|---------------------------------|---------------|
| POST   | `/login`                        | pública       |
| GET    | `/api/users`                    | `user:read`   |
| POST   | `/api/users`                    | `user:write`  |
| PUT    | `/api/users/:id/roles`          | `role:write`  |
| GET    | `/api/roles`                    | `role:read`   |
| POST   | `/api/roles`                    | `role:write`  |
| PUT    | `/api/roles/:id/permissions`    | `role:write`  |
| PUT    | `/api/roles/:id/policies`       | `role:write`  |
| GET    | `/api/permissions`              | `role:read`   |
| POST   | `/api/permissions`              | `role:write`  |
| GET    | `/api/policies`                 | `role:read`   |
| POST   | `/api/policies`                 | `role:write`  |

> **Nota:** Todas as rotas `/api/*` (exceto `/api/me`) passam pelo middleware ABAC. Se o usuário possui políticas que não são satisfeitas (ex: fora do horário permitido), a requisição é negada com 403.
