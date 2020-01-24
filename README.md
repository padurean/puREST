# puREST

## Golang REST API boilerplate with JWT-based authentication, RBAC authorization, PostgreSQL and Swagger for API docs

### :warning: **_This project is WiP!_**

---

### **1. (Re)Generate Swagger docs**

See [Swaggo example for net/http](https://github.com/swaggo/http-swagger)

---

### **2. Specify the app env at runtime**

- Bash: `PUREST_ENV=test go run main.go`

- Fish: `env PUREST_ENV=test go run main.go`

For all supported values for `PUREST_ENV` see the suffix of _**.env.&lt;env&gt;**_ files:
`development` (**default**), `test` and `production`.
