# 🧩 Twitter API - PT/BR

API de estudo utilizando **GO Lang** com persistência de dados via **MongoDB**, seguindo práticas e teorias da **Arquitetura Hexagonal**


## 🚀 Como Executar

A API possui diversas maneiras de execução:

- `Docker`
- `Makefile(Linux)`
- `docker-exec.bat(Windows)`
- `Linha de comando`

---

## ✅ Pré-requisitos

Antes de executar, certifique-se de ter os seguintes itens instalados:

- [GO Lang](https://go.dev/dl/)
- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/)
- [VSCode](https://code.visualstudio.com/Download)

---

## ℹ️ Informações Importantes

- 🚪 A aplicação roda na porta: `8080`

## 💻 Passo a passo

- Após o download das ferramentas do pré-requisito abra o seu terminal/cmd e execute esse comando para clonar o repositório
    ```
    git clone https://github.com/teamcubation/factory-meli.git
    ```
- Depois de clonado, no mesmo terminal use os comandos em sequência
    ```
    cd factory-meli
    git checkout twitter-yuri
    cd twitter-yuri
    ```
- A configuração já está certa! Agora vamos configurar o **MongoDB**
    1. Dentro da pasta twitter-yuri, crie um novo arquivo chamado `.env`, com as mesmas informações do arquivo `.env.example`
    2. Teremos algumas variáveis dentro do nosso `.env`, como:
        ```
        MONGO_USERNAME=
        MONGO_PASSWORD=
        MONGO_HOST=
        MONGO_APP_NAME=
        ```
        Nessas credenciais, passe as configurações do seu **MongoDB**. Caso não possua uma configuração, utilize a minha como exemplo:
        ```
        MONGO_USERNAME=yuripadlipskas
        MONGO_PASSWORD=teste123
        MONGO_HOST=cluster-twitter.elatoly.mongodb.net
        MONGO_APP_NAME=cluster-twitter
        ```
- Pronto! Agora é só seguir o passo a passo de **execução**.
        

## ⚒️ Execução

#### 🪟 Windows

Pelo **Windows** temos duas opções:

1. Execute o script `docker-exec.bat` na raiz do projeto
    ```bash
    docker-exec.bat
    ```
2. Use o comando do docker compose na raiz do seu projeto
    ```bash
    docker compose --project-name twitter up -d
    ```

#### 🐧 Linux

Pelo **Linux** temos também duas opções:

1. `Makefile`
    - Rode esses comandos no seu terminal na raiz do projeto em sequência
    ```bash
    sudo apt-get update
    sudo apt-get -y install make
    make up
    ```
    - Caso queira derrubar o container, apenas execute
    ```bash
    make down
    ```
    - Para mais informações referentes ao **Makefile** execute
    ```bash
    make help
    ```
2. Use o comando do docker compose na raiz do seu projeto
    ```bash
    docker compose --project-name twitter up -d
    ```    

## 📁 Estrutura do Projeto

```bash
README.md                    # README da raiz do repositório
.gitignore                   
twitter-yuri/
│
├── application/             # Camada de aplicação da arquitetura
├── cmd/                     # Pasta onde contem o main.go
├── core/                    # Camada de domain da arquitetura
├── infrastructure/          # Camada de infra da arquitetura
├── .env.example             # Configurações do seu banco de dados
├── docker-compose.yml       # Docker Compose para facilitar a execução
├── docker-exec.bat          # Script para execução local no Windows
├── Dockerfile               # Dockerfile 
├── go.mod                   # Módulo principal do projeto
├── go.sum                   # Versionamento e integridade dos packages
├── README.md                # Este documento

```

## 📚 Endpoints

Temos no total **seis endpoints** para uso nessa API, que ficam globalizados no arquivo `handler.go`, dividos em três módulos com o prefixo padrão **`http://localhost:8080`**

### User
Endpoints relacionado aos **usuários**
- POST `/user/create`
    - body: 
        ```json
        {
            "name": "<nome do usuário>", //STRING
            "follow_users": "<IDs dos usuários>" //[]STRING
        }
        ```
- PATCH `/user/follow`
    - body: 
        ```json
        {
            "user_id": "<ID usuário>", //STRING
            "target_user_id": "<ID usuário destino>" //STRING
        }
        ```
- PATCH `/user/unfollow`
    - body: 
        ```json
        {
            "user_id": "<ID usuário>", //STRING
            "target_user_id": "<ID usuário destino>" //STRING
        }
        ```

### Tweet
Endpoints relacionado aos **tweets**
- POST `/tweet/create`
    - body: 
        ```json
        {
            "title": "Titulo do tweet", // STRING
            "description": "Descrição do tweet", // STRING
            "author_id": "ID do usuário" // STRING
        }
        ```
- GET `/tweet`
    - URL padrão: http://localhost:8080/tweet?id=ID_USER&page=1&tweetsPerPage=10


### Timeline
Endpoints relacionado com a **timeline**
- GET `/timeline/create`
    - URL padrão: http://localhost:8080/timeline?id=ID_USER&page=1&tweetsPerPage=10


## ⚒️ Ferramentas

  - **GO Lang**
  - **Gin**
  - **MongoDB**
  - **Docker**

