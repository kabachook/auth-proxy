# auth-proxy

> Разработка сервиса авторизации для доступа к API веб-приложения с применением JWT-токена

Разработать сервис (backend – PHP/Go/Java), позволяющий:
- [x] Проксировать вызовы к целевому API
- [x] Реализовать целевые вызовы API с настраиваемыми параметрами времени жизни JWT-токена
- [x] Позволяющий создать конфигурацию доступных методов для вызова для каждого клиента
- [x] Логировать действия по вызову API


Для проверки реализации требуется также подготовить:

- [x] Целевой API сервер
- [ ] Тестового клиента

- [x] > Авторизация и аутентификация проводится с применением JWT-токенов (jwt.io), а также с использованием интеграции с каталогом LDAP.

## Quickstart

Terminal 1:
```bash
cd integration
docker-compose up -d
go run server.go
```

Terminal 2:
```bash
go run cmd/auth-proxy -c configs/config.yaml
```

To issue a JWT token use:
```bash
go run utils/issue.go -username jdoe -duration 5m
```

Now you can send requests to port `8989`

Server expects a token in `Auhtorization: Bearer <token>` header

## Config

Full config available [here](./configs/config.yaml)

```yaml
listen: ":8989" # address to listen

authn:
  jwt:
    secret: "qwerty" # JWT secret
    field: "username" # JSON field in payload which identifies user

authz:
  ldap:
    url: "ldap://localhost:389"
    username: "cn=readonly,dc=example,dc=org"
    password: "readonly"
    base_dn: "dc=example,dc=org"
    filter: "(uid=%s)" # LDAP filter, C-style format string
  acl: # ACL list, only supports allow action, default action is deny
    - action: allow
      users:
        - jdoe
      routes: # Exact matches
        - /hello

backends:
  - name: some_backend
    host: "localhost"
    port: 8080
    scheme: "http"

routes:
- match:
    host: "*" # "*" to match all hosts
  backend: some_backend

proxy:
  header: "X-Username" # identity header for backend
```