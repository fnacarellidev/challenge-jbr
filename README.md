# Desafio Consulta Processual

## Backend API

### Tecnologias Utilizadas Na Construção Do Backend

- **Golang**
- **Postgres-14**
- **SQLC**: Gera codigo Go a partir de um schema com tipagem segura para interagir com o banco.
- **Testcontainers (Testes)**: Sobe um conteiner postgres-14 para testar tanto a API como as Queries SQL.

### Endpoints
- **/register_court_case:** Registra um novo processo judicial, recebe toda a informação no corpo da requisição (application/json).
- **/fetch_court_case/:cnj:** Recupera um processo judicial a partir de um CNJ, parâmetro passado no URL da requisição.
- **/healthcheck:** Utilizado pelo docker-compose.yaml para garantir que o serviço do GraphQL só esteja disponível quando o backend estiver pronto.

### Testes Realizados

**Testes realizados na API (backend/tests/api_test.go):**
- Recuperar um processo já registrado (Alice X Bob) e toda a informação a respeito do processo.
- Recuperar um processo já registrado (Michael X Sarah) e toda a informação a respeito do processo.
- Recuperar um processo que não existe.
- Criar um processo que já esta registrado (CNJ duplicado).
- Criar um processo passando informações que não são aceitas pela API.
- Criar um processo com um campo faltando (CNJ).

**Testes realizados diretamente nas queries SQL (backend/tests/sqlc_queries_test.go):**
- Criar um processo.
- Recuperar um processo já registrado.

### Como Executar Os Testes Do Backend

Para executar os testes a partir da raiz do projeto:
```bash
$ go test ./backend/tests -v
```
