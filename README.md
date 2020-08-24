# puREST

## Golang REST API boilerplate with authentication using PASETO tokens, RBAC authorization, PostgreSQL and Swagger for API docs

### :warning: **_This project is WiP!_**

---
![Go](https://github.com/padurean/puREST/workflows/Go/badge.svg)

### **1. (Re)Generate Swagger docs**

See [Swaggo example for net/http](https://github.com/swaggo/http-swagger)

Then run the following command after any API changes:
`swag init --dir . -g cmd/server/main.go`

---

### **2. Specify the app env at runtime**

- Bash: `PUREST_ENV=test go run main.go`

- Fish: `env PUREST_ENV=test go run main.go`

For all supported values for `PUREST_ENV` see the suffix of _**.env.&lt;env&gt;**_ files:
`development` (**default**), `test` and `production`.

### **3. Run**

Create PostgreSQL database and user:

```console
CREATE USER purest_user WITH ENCRYPTED PASSWORD 'purest_pass';
CREATE DATABASE purest_db;
GRANT ALL PRIVILEGES ON DATABASE purest_db TO purest_user;
```

Run it:

`go run cmd/server/main.go`

or build it and then run it:

```console
go build -o puREST ./cmd/server/main.go
./puREST
```

The built-in Swagger UI can be accessed at: <http://localhost:8000/swagger/>
