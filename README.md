# auth-proxy

Protects your backends using JWT + LDAP

## Quickstart

Bring up LDAP server and test backend (terminal 1):

> you should have Docker installed

```bash
cd integration
docker-compose up -d
go run server.go
```

> Test user `jdoe` is available

Run auth-proxy (terminal 2):
```bash
go run cmd/auth-proxy -c configs/config.yaml
```

Test:
```bash
❯ curl -H "Authorization: Bearer eyJh..." http://localhost:8989/hello?abc=xyz -i
HTTP/1.1 200 OK
Content-Length: 291
Content-Type: text/plain; charset=utf-8
Date: Tue, 07 Jul 2020 08:01:36 GMT

Hello!
You called /hello
Headers: map[Accept:[*/*] Accept-Encoding:[gzip] Authorization:[Bearer eyJh...] User-Agent:[curl/7.68.0] X-Forwarded-For:[127.0.0.1] X-Username:[jdoe]]
Username: jdoe
```
---
To issue a JWT token use:
```bash
go run utils/issue.go --username jdoe --duration 5m
```

Now you can send requests to port `8989`

Server expects a token in `Authorization: Bearer <token>` header

## How it works

1. auth-proxy gets request
2. It first _authenticates_ using JSON web token
3. Then it _authorizes_ the request using username from JWT payload via LDAP
4. Last, it checks the Access Control Lists to limit access for specific routes
5. Request is proxied to backend, username from JWT is written to header (`X-Username` by default)

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

## Build

```bash
make build
./auth-proxy
```

## Task

> Разработка сервиса авторизации для доступа к API веб-приложения с применением JWT-токена

Разработать сервис (backend – PHP/Go/Java), позволяющий:
- [x] Проксировать вызовы к целевому API
- [x] Реализовать целевые вызовы API с настраиваемыми параметрами времени жизни JWT-токена
- [x] Позволяющий создать конфигурацию доступных методов для вызова для каждого клиента
- [x] Логировать действия по вызову API


Для проверки реализации требуется также подготовить:

- [x] Целевой API сервер
- [x] Тестового клиента (curl)

- [x] > Авторизация и аутентификация проводится с применением JWT-токенов (jwt.io), а также с использованием интеграции с каталогом LDAP.
