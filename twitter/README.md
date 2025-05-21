# Twitter

## Tecnologias utilizadas:
- [Go](https://go.dev/) `v1.23.4`
  - Library [uuid](https://github.com/google/uuid) `v1.6.0`
  - Library [godotenv](https://github.com/joho/godotenv) `v1.5.1`
  - Library [pq](https://github.com/lib/pq) `v1.10.9`
  - Library [testify](https://github.com/stretchr/testify) `v1.10.0`
- [PostgreSQL](https://www.postgresql.org/) `v16.6`
- [Docker](https://www.docker.com/)

## Rodando a API:
Essa é uma API que foi construída para ser rodada utilizando o Docker, uma poderosa ferramenta que permite a utilização de imagens de sistemas operacionais inteiros para a execução de código de forma consistente. Especificamente, para rodar essa API, utilizamos o [Docker Compose](https://docs.docker.com/compose/), que permite a utilização de múltiplos containers em um único serviço. Para rodar nosso código, é só seguir os seguintes passos:
1. [Instale e configure](https://docs.docker.com/guides/getting-started/) o Docker.
2. [Instale e configure](https://docs.docker.com/compose/install/) o Docker Compose.
3. Crie um arquivo `.env` na raiz deste repositório e preencha as variáveis conforme o arquivo `.env.example`.
   1. Note que é extremamente importante preencher a variável `ENV` corretamente com o valor `development`, pois preencher ela com o valor `production` gera o executável para deployment sem utilizar o Docker Compose e não permite o hot-reloading da aplicação.
4. Abra um terminal na raiz deste repositório.
5. Execute o comando `docker-compose up --build`.
   1. Note que este comando pode variar dependendo da sua instalação ou distribuição do Linux, utilizando Debian no WSL 2, o meu é `docker compose up --build`.
   2. Também é possível executar isso usando o [Docker Desktop](https://www.docker.com/products/docker-desktop/), mas visto como eu não o utilizo pessoalmente, não deixei isso nessa documentação.
   3. Note que a flag `--build` só é necessária na primeira vez que rodar a imagem, ou caso hajam alterações na image/compose futuramente.
6. Caso você tenha configurado seu ambiente corretamente, você deverá ter 2 containers rodando, um do banco de dados e um da API em si (rode o comando `docker container ls -a` para verificar). Em `development`, o container da API utiliza a ferramenta `Air` e graças as configurações do Docker Compose, o hot-reloading funciona. Em `production`, apenas é gerado um executável único que então é executado pela API, logo o hot-reloading não está disponível. É só acessar as rotas da API no seu browser ou API Client favorito.
7. Caso você esteja em ambiente de desenvolvimento, as migrations do banco de dados não são rodadas automaticamente. No diretório raiz do repositório, rode o comando `make m-up` e elas serão rodadas no seu banco de dados (novamente, não se esqueça de preencher o seu `.env` antes de fazer isso, ou nada vai funcionar). Caso você utilize Windows e não tenha o `Make` instalado para rodar comandos do `Makefile`, eu recomendo a instalação [nesse link](https://gnuwin32.sourceforge.net/packages/make.htm) já que ele é essencial para o desenvolvimento em Go. Se ainda assim você você preferir não instalar, terá que montar a sua string de conexão do PostgreSQL manualmente e rodar o comando `goose -dir sql/schema postgres (PG_CONN_STRING) up` no seu terminal.
8.  Quando tiver terminado a utilização, apenas dê um `ctrl + c` no terminal e execute o comando `docker-compose down` caso deseje deletar os containers. Como a aplicação possui um volume, os dados não serão perdidos no seu banco a não ser que você delete o volume também.
   1. Note que este comando pode variar dependendo da sua instalação ou distribuição do Linux, utilizando Debian no WSL 2, o meu é `docker compose down`.

Excelente! Que bom que todos os testes estão passando. Isso demonstra a solidez e a boa cobertura dos seus serviços.

Vamos adicionar uma seção de testes automatizados ao seu README.md em português, explicando como rodá-los.

## Testes Automatizados
Esta aplicação conta com uma suíte de testes automatizados unitários para garantir a corretude e a qualidade do código. Os testes são implementados utilizando o framework de testes padrão do Go (testing) em conjunto com a biblioteca testify/mock para simular dependências e isolar os componentes a serem testados.

### Como Rodar os Testes
Para executar todos os testes da aplicação, siga estes passos:

1. Navegue até a raiz do seu repositório em um terminal.
2. Execute o comando `go test ./... -v -count=1` ou, se você tiver a ferramenta `make` instalada, você pode rodar `make test`, que é um `alias` para rodar esse exato comando.

### Explicação do comando
- `go test`: O comando principal para executar testes Go.
- `./...`: Este padrão indica que o go test deve procurar e executar todos os testes em todos os pacotes dentro do diretório atual e seus subdiretórios.
- `-v`: (verbose) Esta flag exibe a saída detalhada de cada teste, mostrando quais testes passaram (PASS), quais falharam (FAIL) e qualquer log que você tenha adicionado (como os t.Logf que usamos para depuração).
- `-count=1`: Esta flag é importante para garantir que os testes não sejam executados múltiplas vezes a partir do cache do Go. Isso é crucial quando se trabalha com mocks e estado, prevenindo resultados inconsistentes.

### Exemplo de Saída (Sucesso)
```
goos: linux
goarch: amd64
pkg: github.com/twitter-tq/vinofsteel/internal/services
--- PASS: TestCreateUser (0.00s)
    --- PASS: TestCreateUser/Valid_user_creation (0.00s)
    --- PASS: TestCreateUser/Invalid_email_format (0.00s)
    --- PASS: TestCreateUser/Email_already_exists (0.00s)
    --- PASS: TestCreateUser/Empty_email (0.00s)
    --- PASS: TestCreateUser/Database_error_during_lookup (0.00s)
--- PASS: TestGetUserTimeline (0.00s)
    --- PASS: TestGetUserTimeline/Successful_timeline_retrieval_with_defaults (0.00s)
    --- PASS: TestGetUserTimeline/Successful_timeline_with_custom_params (0.00s)
    --- PASS: TestGetUserTimeline/Invalid_user_ID_in_path (0.00s)
    --- PASS: TestGetUserTimeline/User_not_found (0.00s)
    --- PASS: TestGetUserTimeline/Database_error_during_user_lookup (0.00s)
    --- PASS: TestGetUserTimeline/Invalid_limit_parameter (0.00s)
    --- PASS: TestGetUserTimeline/Invalid_offset_parameter (0.00s)
    --- PASS: TestGetUserTimeline/Invalid_tweet_num_parameter (0.00s)
    --- PASS: TestGetUserTimeline/Database_error_during_timeline_retrieval (0.00s)
    --- PASS: TestGetUserTimeline/Empty_timeline_(no_followed_users_or_no_tweets) (0.00s)
--- PASS: TestFollowUser (0.00s)
    --- PASS: TestFollowUser/Successful_follow (0.00s)
    --- PASS: TestFollowUser/Invalid_request_payload (0.00s)
    --- PASS: TestFollowUser/Invalid_follower_ID_in_body (0.00s)
    --- PASS: TestFollowUser/Follower_not_found (0.00s)
    --- PASS: TestFollowUser/Database_error_during_follower_lookup (0.00s)
    --- PASS: TestFollowUser/Invalid_followed_ID_in_path (0.00s)
    --- PASS: TestFollowUser/Followed_not_found (0.00s)
    --- PASS: TestFollowUser/Database_error_during_followed_lookup (0.00s)
    --- PASS: TestFollowUser/Cannot_follow_yourself (0.00s)
    --- PASS: TestFollowUser/Follow_relationship_already_exists (0.00s)
    --- PASS: TestFollowUser/Database_error_during_FindByIds_check (0.00s)
    --- PASS: TestFollowUser/Database_error_during_Save (0.00s)
--- PASS: TestUnfollowUser (0.00s)
    --- PASS: TestUnfollowUser/Successful_unfollow (0.00s)
    --- PASS: TestUnfollowUser/Invalid_request_payload (0.00s)
    --- PASS: TestUnfollowUser/Invalid_follower_ID_in_body (0.00s)
    --- PASS: TestUnfollowUser/Follower_not_found (0.00s)
    --- PASS: TestUnfollowUser/Database_error_during_follower_lookup (0.00s)
    --- PASS: TestUnfollowUser/Invalid_followed_ID_in_path (0.00s)
    --- PASS: TestUnfollowUser/Followed_not_found (0.00s)
    --- PASS: TestUnfollowUser/Database_error_during_followed_lookup (0.00s)
    --- PASS: TestUnfollowUser/Cannot_unfollow_yourself (0.00s)
    --- PASS: TestUnfollowUser/Cannot_unfollow_a_user_you_don't_follow (0.00s)
    --- PASS: TestUnfollowUser/Database_error_during_FindByIds_check (0.00s)
    --- PASS: TestUnfollowUser/Database_error_during_Delete (0.00s)
    --- PASS: TestUnfollowUser/Follow_relationship_not_found_by_Delete_(return_error) (0.00s)
PASS
ok      github.com/twitter-tq/vinofsteel/internal/services  0.004s
```

## Documentação da API
A API oferece os seguintes endpoints para gerenciar usuários, tweets e relacionamentos de seguidores:

### Usuários

#### Criar um novo usuário
- **URL**: `/users`
- **Método**: `POST`
- **Descrição**: Cria um novo usuário na plataforma
- **Corpo da requisição**:
  ```json
  {
    "email": "string"
  }
  ```
- **Resposta de sucesso**:
  - **Código**: 201 Created
  - **Conteúdo**:
    ```json
    {
      "id": "uuid",
      "created_at": "string",
      "updated_at": "string",
      "email": "string",
    }
    ```
- **Respostas de erro**:
  - **Código**: 400 Bad Request - Se dados estiverem incompletos ou inválidos
  - **Código**: 409 Conflict - Se o username ou email já existir

### Timeline

#### Obter timeline de um usuário
- **URL**: `/timeline/{user_id}`
- **Método**: `GET`
- **Parâmetros de URL**:
  - `user_id`: ID do usuário (UUID)
- **Parâmetros de consulta (query params)**:
  - `limit` (opcional): Número máximo de usuários a retornar (padrão: 10)
  - `offset` (opcional): Índice inicial para paginação (padrão: 0)
  - `tweet_num` (opcional): Número máximo de tweets por usuário seguido (padrão: 10)
- **Descrição**: Retorna os tweets do usuário e das pessoas que ele segue, em ordem cronológica decrescente, com suporte a paginação
- **Resposta de sucesso**:
  - **Código**: 200 OK
  - **Conteúdo**:
    ```json
    [
      {
        "user": {
          "id": "uuid",
          "created_at": "string",
          "updated_at": "string",
          "email": "string",
        },
        "tweets:" [
          {
            "id": "uuid",
            "created_at": "string",
            "updated_at": "string",
            "post": "string",
            "creator_id": "uuid",
          }
        ]
      }
    ]
    ```
- **Respostas de erro**:
  - **Código**: 400 Bad Request - Se algum dos parâmetros de consulta for inválido
  - **Código**: 404 Not Found - Se o usuário não existir

### Tweets

#### Criar um novo tweet
- **URL**: `/tweets/{creator_id}`
- **Método**: `POST`
- **Parâmetros de URL**:
  - `creator_id`: ID do usuário criador (UUID)
- **Descrição**: Cria um novo tweet para o usuário especificado
- **Corpo da requisição**:
  ```json
  {
    "post": "string"
  }
  ```
- **Resposta de sucesso**:
  - **Código**: 201 Created
  - **Conteúdo**:
    ```json
    {
      "id": "uuid",
      "created_at": "string",
      "updated_at": "string",
      "post": "string",
      "creator_id": "uuid",
    }
    ```
- **Respostas de erro**:
  - **Código**: 400 Bad Request - Se o conteúdo estiver vazio ou inválido
  - **Código**: 404 Not Found - Se o usuário não existir

#### Listar tweets de um usuário
- **URL**: `/tweets/{creator_id}`
- **Método**: `GET`
- **Parâmetros de URL**:
  - `creator_id`: ID do usuário criador (UUID)
- **Descrição**: Retorna todos os tweets criados pelo usuário, em ordem cronológica decrescente
- **Resposta de sucesso**:
  - **Código**: 200 OK
  - **Conteúdo**:
    ```json
    [
      {
        "id": "uuid",
        "created_at": "string",
        "updated_at": "string",
        "post": "string",
        "creator_id": "uuid",
      }
    ]
    ```
- **Respostas de erro**:
  - **Código**: 404 Not Found - Se o usuário não existir

### Seguidores

#### Seguir um usuário
- **URL**: `/follows/{followed_id}`
- **Método**: `POST`
- **Parâmetros de URL**:
  - `followed_id`: ID do usuário a seguir (UUID)
- **Descrição**: Cria um relacionamento de seguidor entre o usuário autenticado e o usuário especificado
- **Corpo da requisição**:
  ```json
  {
    "follower_id": "uuid"
  }
  ```
- **Resposta de sucesso**:
  - **Código**: 201 Created
  - **Conteúdo**:
    ```json
    {
      "id": "uuid",
      "created_at": "string",
      "updated_at": "string",
      "follower_id": "uuid",
      "followed_id": "uuid",
    }
    ```
- **Respostas de erro**:
  - **Código**: 400 Bad Request - Se os dados estiverem incompletos
  - **Código**: 404 Not Found - Se algum dos usuários não existir
  - **Código**: 409 Conflict - Se o mesmo id for passado em ambos os parâmetros
  - **Código**: 409 Conflict - Se o relacionamento já existir

#### Deixar de seguir um usuário
- **URL**: `/follows/{followed_id}`
- **Método**: `DELETE`
- **Parâmetros de URL**:
  - `followed_id`: ID do usuário a deixar de seguir (UUID)
- **Descrição**: Remove um relacionamento de seguidor entre o usuário autenticado e o usuário especificado
- **Corpo da requisição**:
  ```json
  {
    "follower_id": "uuid"
  }
  ```
- **Resposta de sucesso**:
  - **Código**: 204 No Content
- **Respostas de erro**:
  - **Código**: 400 Bad Request - Se os dados estiverem incompletos
  - **Código**: 404 Not Found - Se algum dos usuários não existir
  - **Código**: 404 Not Found - Se o relacionamento não existir