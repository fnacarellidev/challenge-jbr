## Início do Projeto
A construção do projeto acabou sendo mais simples do que eu imaginava, eu não conhecia GraphQL então
fiquei meio confuso em como ele se comunicaria com o backend ao ver o diagrama no desafio, então a
primeira coisa que fui fazer foi pesquisar na pratica como que o GraphQL é implementado pra decidir e
entender como seria aplicado aqui nesse caso, depois de ver alguns posts e trechos de codigo percebi que ele
ia funcionar como um "filtro" pra uma API normal que serve um set completo de informações. Depois que entendi isso
o projeto foi bem tranquilo, já que eu ja tinha uma boa noção de como conteinerizar as aplicações e de como fazer o
backend se comunicar com um banco de dados.

## Implementação

### Backend
Comecei subindo o banco de dados e configurando o backend pra se comunicar
com o banco e fazer as requisições que eu julguei que fossem necessárias, basicamente um SELECT na tabela que registram os casos
e outro SELECT na tabela que linka as atualizações a respeito de um caso com o caso em si. Boa parte do trabalho que teria que ser
manual acabou sendo automatizado pelo [SQLC](https://github.com/sqlc-dev/sqlc), que lê o schema do banco de dados que eu vou interagir
e já gera funções para realizar as queries SQL com tipagem segura.

### GraphQL
Para a minha surpresa essa parte foi bem tranquila, eu achei que teria mais dificuldade por nunca ter mexido com o GraphQL, mas
encontrei bastante conteúdo na internet que me ajudou bastante, inclusive acho que se eu fosse refazer esse projeto eu implementaria
o GraphQL em TypeScript por ter uma documentação mais extensa, por eu ter decidido fazer em Golang a documentação era um pouco mais limitada,
mas também não me arrependo, aprendi bastante coisa.

### Frontend
Acabei deixando o frontend por último, e também foi bem tranquilo já que eu ja tinha feito um projeto em React e também trabalhei bastante
com Vue, então tenho uma boa noção de componentização, fluiu bem essa parte.

### Testes
Depois de ter esse setup básico eu comecei a implementar os testes para o backend e o GraphQL e durante esse processo fui notando alguns
edge-cases que eu havia deixado de lado, algumas tratativas de erro que faltaram e também algumas incoerências nas respostas da minha API
e algumas respostas genéricas demais, daí dei uma melhorada nessa parte. No README.md está melhor detalhado os testes automatizados que foram realizados,
todos utilizando o toolkit [testify](https://github.com/stretchr/testify), mais especificamente o pacote `suite`, principalmente por conta das funcões
de Setup e Teardown dos tests, pra subir os endpoints e o banco de dados ([testcontainers](https://testcontainers.com/)) para utilizar nos testes.

### O que mais eu gostaria de implementar
Acho que uma funcionalidade que eu fiquei pensando bastante em implementar é de ter diversos filtros na aba de pesquisa, filtrar os processos por Autor, por Réu,
e assim por diante. Outra coisa que senti falta é de ter testes end to end, acho que é algo super necessário que dá um diferencial na hora de subir um projeto.
