# Desafio Consulta Processual

O objetivo desse desafio é construir um sistema de consulta processual, onde terá um banco de dados
já populado com alguns processos falsos, e recuperar esses processos por meio de uma API GraphQL.
Também é possível cadastrar novos processos pelo frontend.

## Requisitos

- Docker
- Go >= 1.22

## 🚀 Começando

O frontend da aplicação estará disponível na porta 3000, para rodar o projeto basta executar:
```bash
$ docker compose up
```

## Modelo De Dados

![database schema](https://github.com/fnacarellidev/challenge-jbr/blob/main/assets/schema.png)

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
- Recuperar um processo que foi registrado.
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

## GraphQL API

### Objetos 
```
type CaseUpdateInput {
    update_date: DateTime!
    update_details: String!
}

type CourtCase {
    cnj: String!
    plaintiff: String!
    defendant: String!
    court_of_origin: String!
    start_date: String!
    updates: [CaseUpdateInput]!
}
```

### Operações Disponíveis
- **court_case(cnj: String!, court_of_origin: String!):** Query que recupera o processo a partir de um CNJ e Tribunal.
- **new_court_case(cnj: String!, plaintiff: String!, defendant: String!, court_of_origin: String!, start_date: DateTime!, updates: \[CaseUpdateInput\])**: Mutation que cria um processo judicial com os parâmetros passados.

### Testes Realizados

**Testes realizados na operação court_case:**
- Recuperar todas as informações de um caso específico (Alice X Bob)
- Recuperar apenas Autor e Réu de um caso específico (Alice X Bob)
- Recuperar apenas Autor de um caso específico (Alice X Bob)
- Recuperar um caso que não foi registrado.
- Recuperar apenas as atualizações de um caso específico (Alice X Bob)

**Testes realizados na operação new_court_case:**
- Criar um processo.
- Criar um processo com CNJ já cadastrado.
- Criar um processo com CNJ vazio.
- Criar um processo passando tipos incorretos.

### Como Executar Os Testes Do GraphQL

Para executar os testes a partir da raiz do projeto:
```bash
$ go test ./graphql-api/tests -v
```
