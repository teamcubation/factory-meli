# Projeto Teamcubation - Twitter

Esse é um projeto que está sendo construido baseado nas especificações contidas no arquivo `twitter-TQ.pdf`
Para este, está sendo criado um projeto seguindo uma arquitetura hexagonal com as seguintes com as seguintes tecnologias

---

## Tecnologias

- Golang
- Docker
- Docker Compose
- PostgreSQL
- Bash
- Air

---

## Estrutura "AS-IS" do projeto

```
07-twitter/
    cmd/
        server/
            .air.toml       # Live reload
            .env            # Variáveis de ambiente (ex. banco de dados)
	        main.go         # Inicializar a aplicação
	core/
        models/             # Entidades
	        tweet.go
	        user.go
		ports/              # Interfaces (repos e services/usecases)
			repos/
	            tweet_repository.go
	            user_repository.go
			services/
				tweet_service.go
				user_service.go
    infra/                  # Configurações de Docker Compose
        scripts/            # Scripts Bash do Docker Compose
        docker-compose.yml
    internal/
        adapters/           # Implementações (adapters)
            http/           # Handlers HTTP (delivery)
                handler.go
                router.go
            repositories/   # Persistência de dados
                tweet_memory.go
                user_memory.go
		services/           # Casos de uso (application/service layer)
            tweet_service.go
            user_service.go
    test/                   # Testes unitários
        app/
            user_service_test.go
            tweet_service_test.go
	.gitignore
    go.mod
    go.sum
    README.md
    twitter-TQ.pdf          # Requisitos do projeto
```

**Obs**: O sub-diretório de infra foi criado utilizando como base um repo privado de um dos meus projetos pessoais

- [BeanBuddy](https://github.com/Bean-Buddy/beanbuddy-infra)

---

## Como usar

1. **Clone o repositório e acesse o diretório do projeto:**

```bash
git clone https://github.com/teamcubation/factory-meli.git
git switch exercicios-luiz
cd factory-meli/07-twitter
```

2. **Configure as variáveis de ambiente:**

Edite o arquivo `.env` que está no caminho `cmd/server/`.
Como exemplo, pode utilizar os seguintes dados (ajuste se necessário):

```
POSTGRES_DB=twitterTest
POSTGRES_USER=urubu100
POSTGRES_PASSWORD=urubu100
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

3. **Torne os scripts executáveis no WSL:**

Abra um terminal WSL e que já tenha o Docker instalado.
Navegue até a pasta `infra` e torne os scripts executáveis:

```bash
cd infra
chmod +x scripts/*.sh
```

4. **Suba o banco de dados com Docker Compose:**

Execute o bash de up contido na pasta `scripts/`:

```bash
./scripts/up.sh
```

Caso houver sucesso ou erro, será impresso no terminal.

Mas se quiser, pode verifique se o container está rodando:

```bash
docker ps
```

5. **(Opcional) Acesse o banco de dados via terminal:**

Caso tenha alterado os dados do arquivo `.env`, então aqui também deve ser atualizado:

```bash
docker exec -it twitter-test-db psql -U urubu100 -d twitterTest
```

6. **Execute a aplicação Go com o Air:**

Abra outro terminal na raiz do projeto (`07-twitter`) e execute:

```bash
cd cmd/server
air
```

Se tudo estiver correto, verá o seguinte log no terminal:

```
Sucesso ao se conectar com o Banco de Dados!
```

7. **Parar e remover os containers:**

Para parar e remover os containers do banco, execute o seguinte dentro do terminal WSL:

```bash
./scripts/down.sh
```

---
