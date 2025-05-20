![tests ci badge](https://github.com/fundacao-lusofona/luso-wiki/actions/workflows/ci.yml/badge.svg)
# Luso Wiki (Aberto a sugestões para nomes)

## Tecnologias utilizadas:
- [Go](https://go.dev/) `v1.23.4`
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
