### ✅ Explicação do código

Este programa em Go implementa um **servidor web básico** usando o framework **Gin**. Serve como exemplo de como estruturar um microserviço simples com as camadas típicas de _handler_ (controlador), _usecase_ (lógica de negócio) e _entidade_ (domínio).

---

#### 1. **Estrutura geral**

- `main()`:

  - Cria um `usecase` (lógica de negócio).
  - Cria um `handler` (controlador HTTP).
  - Configura o router com várias rotas e inicia o servidor em `localhost:8080`.

---

#### 2. **Rotas HTTP**

- `GET /` → Retorna `"Welcome to the home page!"`.
- `GET /hello` → Retorna `"Hello, world!"`.
- `POST /bye` → Recebe um JSON com a chave `"message"` e devolve como resposta.
- `POST /items` → Recebe um `item` (JSON), salva usando a lógica de negócio.
- `GET /items` → Retorna todos os `item` existentes.

---

#### 3. **Handler**

O `handler`:

- Valida as entradas (`c.BindJSON`).
- Chama a camada de negócio (`usecase`).
- Retorna respostas HTTP com os códigos e mensagens apropriadas.

---

#### 4. **Usecase (lógica de negócio)**

- `saveItem(it item)`: Simula o salvamento de um item.
- `listItems()`: Retorna um mapa vazio de itens (ainda não há armazenamento real).

---

#### 5. **Domínio (`item`)**

A entidade `item` representa um produto ou objeto com vários atributos: `ID`, `Code`, `Title`, `Price`, etc. É uma struct típica com tags JSON para serialização.

---

---

### 🧪 Exercício

#### **Extensão do programa - Desafio para os alunos**

> ✍️ **Objetivo**: Estender este servidor para **adicionar, obter, atualizar e remover itens**, tudo em um **único arquivo Go**.

---

#### 🛠️ Requisitos técnicos

1. **Adicionar armazenamento em memória**

   - Usar um `map[int]item` como banco de dados simulado em memória (substituir o `map` vazio do `listItems` atual).
   - Armazenar os itens que chegam via POST.

2. **Novas rotas a implementar**

   - `GET /items/:id`: Buscar um item pelo seu `ID`.
   - `PUT /items/:id`: Atualizar um item pelo `ID`. O corpo deve conter o item atualizado.
   - `DELETE /items/:id`: Remover um item pelo `ID`.

3. **Validações**

   - Verificar se os dados do JSON são válidos.
   - Garantir que o `item.ID` não esteja duplicado ao fazer um POST.
   - Verificar se o item existe antes de atualizar ou remover.

4. **Resposta**

   - Sempre usar respostas em `JSON` com os códigos HTTP adequados (`200`, `201`, `400`, `404`, etc.).

---

#### 💡 Dicas

- Toda a lógica e estruturas devem estar contidas **em um único arquivo**.
