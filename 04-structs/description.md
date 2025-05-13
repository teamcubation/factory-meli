1. **Modelar um Pokémon e sua evolução**

   - Crie um `struct` chamado `Pokemon` com os campos: `Name` (string), `Type` (`[]string`), `Level` (int) e `EvolvesTo` (string).
   - Escreva uma função `Evolve(p *Pokemon) error` que:

     1. Verifique se `p.Level >= 16`.
     2. Se for, imprima no console `"¡<Name> evolui para <EvolvesTo>!"` e atualize `p.Name` para o valor de `EvolvesTo`.
     3. Caso contrário, retorne um erro indicando que ainda não pode evoluir.

2. **Gerenciador de time Pokémon**

   - Defina um `struct` chamado `Trainer` com `Name` (string) e `Party` (`[]Pokemon`).
   - Implemente um método de receptor `AddToParty(p Pokemon) error` que:

     1. Adicione `p` ao slice `Party` se o tamanho for menor que 6.
     2. Se a party já tiver 6 Pokémons, retorne um erro informando que o time está cheio.

   - Em `main`, crie um treinador, alguns `Pokemon` e teste adicioná-los, tratando possíveis erros.

3. **Calculadora de dano de movimento**

   - Crie um `struct` `Move` com os campos: `Name` (string), `Power` (int) e `Type` (string).
   - Adicione uma função `CalculateDamage(m Move, targetType string) int` que:

     1. Considere vantagem de tipos: se `m.Type == "Fogo"` e `targetType == "Planta"`, dobre o dano; se `m.Type == "Água"` e `targetType == "Fogo"`, dobre o dano; para demais combinações, o dano é `m.Power`.
     2. Retorne o dano calculado.

   - Em `main`, defina vários movimentos e calcule o dano contra diferentes tipos de Pokémon.

4. **Pokédex baseada em mapas**

   - Defina um tipo `Pokedex` como `map[int]Pokemon`, onde a chave é o `ID` do Pokémon.
   - Escreva as funções:

     1. `AddPokemon(pokedex Pokedex, id int, p Pokemon)`: insere `p` no map.
     2. `GetPokemon(pokedex Pokedex, id int) (Pokemon, bool)`: retorna o Pokémon e um bool indicando se ele existia.
     3. `ListAll(pokedex Pokedex) []Pokemon`: retorna um slice com todos os Pokémons armazenados.

   - Em `main`, crie uma `Pokedex`, adicione vários exemplares, busque um por ID (tratando a ausência) e exiba no console a lista completa.
