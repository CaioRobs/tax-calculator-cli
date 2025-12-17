# Tax Calculator (Capital Gains) CLI

CLI para calcular imposto sobre lucro obtido em operações de compra e venda de ações linha a linha via STDIN. Cada linha de entrada é um array JSON de operações; a saída correspondente (STDOUT) é um array JSON com os impostos devidos por operação.

## Como funciona

- Operações suportadas: `buy` (compra) e `sell` (venda).
- Preço médio ponderado: a cada compra, o preço médio ponderado é atualizado considerando o preço médio por ação e a quantidade de ações anteriores, assim como o preço por unidade, a quantidade de ações adquiridas e a nova quantidade total de ações.
- Prejuízo acumulado: prejuízos em vendas são acumulados e usados para abater lucros futuros.
- Tributação: vendas só são tributadas quando o valor total da operação de venda ultrapassa 20000. A alíquota aplicada sobre o lucro líquido (após abater prejuízos acumulados) é de 20%.
- Compras nunca geram imposto; vendas podem gerar imposto conforme as regras acima.

### Contrato de I/O

- Entrada: cada linha no STDIN deve ser um array JSON com objetos de operação no formato:

```json
{ "operation": "buy" | "sell", "unit-cost": float64, "quantity": int }
```

- Saída: para cada linha de entrada, o programa imprime um array JSON de mesmo tamanho contendo objetos:

```json
{ "tax": float64 }
```

### Exemplo

- Entrada:

```json
[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":15.00,"quantity":50},{"operation":"sell","unit-cost":15.00,"quantity":50}]
```

- Saída:

```json
[{"tax":0},{"tax":0},{"tax":0}]
```

No exemplo vemos 3 operações (compra, venda e outra venda) em sequencia. Primeiro na compra o preço medio ponderado (weightedAveragePrice) é calulado e atualizado com o valor 10.00, o numero total de ações (totalSharesQuantity) se torna a quantidade que foi adquirida (100) e o prejuízo acumulado (accumulatedLoss) no momento é 0. Na saída, o primeiro resultado da lista mostra que não houve imposto incidente (por se tratar de uma compra).
Em seguida, uma venda onde o numero total de ações é reduzido para 50 (pois outras 50 foram vendidas), o preço médio se mantém pois não houve outra compra e o prejuízo se mantém em 0 pois foi uma operação lucrativa. Porém o valor total da operação foi menor do que 20000 então não houve imposto incidente. O mesmo acontece para a terceira operação, também de venda.

## Estrutura do projeto

``` shell
.
├── 📁internal
│    ├── 📁cli
│    │    └── processor.go
│    ├── 📁domain
│    │    ├── operation.go
│    │    └── tax_result.go
│    └── 📁services
│         └── tax_calculator.go
├── 📁io
│    ├── expected_output.txt
│    └── input.txt
└─ main.go
```

- `main.go`: ponto de entrada; cria o `Processor` e executa.

- `internal/cli/processor.go`: lê linhas do STDIN, parseia JSON, chama o serviço e imprime a saída.
- `internal/domain/operation.go`: modelo das operaçãos de entrada.
- `internal/domain/tax_result.go`: modelo dos resultados de imposto por operação.
- `internal/services/tax_calculator.go`: regras de negócio (compra e venda, preço médio, prejuízo acumulado, tributação 20% acima de 20k por operação de venda).

- `io/input.txt` / `io/expected_output.txt`: exemplos de entrada/saída.

Em resumo: main.go incia a aplicação, instanciando um Processor e iniciando a execução. O Processor começa a ler o STDIN, ao reconhecer uma quebra de linha, parseia a entrada para uma lista de structs ([]domain.Operation), executa o cálculo (taxCalculator.Calculate), indicando via ponteiro (&results) a variavel onde serão salvos os resultados ([]domain.TaxResult), transforma a lista de resultados num JSON e lança como string no STDOUT. \
\
O cálculo começa na função taxCalculator.Calculate, onde é criada uma lista temporaria do tamanho da lista de operações recebidas na linha, é alocado um espaço em memória para registrar a quantidade de ações, o preço médio e o prejuízo acumulado (todos com valor inicial 0) e então se começa uma iteração na lista de operações, da primeira em diante.\
Baseado no tipo da operação, compra ou venda, os dados de cada operação e os ponteiros relacionados ao estado de processamento daquela linha são passados para as funções auxiliares buyOperation e sellOperation, respectivamente, que calculam os valores a serem atualizados (quantidade de ações, preço médio e prejuizo), fazem as validações se serão taxadas ou não e calculam o valor do imposto incidente. São registrados os novos valores do estado da aplicação e é retornado (pelas funções auxiliares) um struct de resposta (domain.TaxResult) para ser inserido na lista de resultados daquela linha na posição referente à posição da operação calculada.\
Após todas as operações da linha serem calculadas, os resultados são salvos no espaço alocado anteriormente.

\
Observação: cada arquivo encontrado em /internal (e main.go) possui seu complementar *_test.go contendo testes unitários, omitidos aqui para melhor visualização da estrutura. integration_test.go possui os testes de integração da aplicação.

## Execução com Docker

Build da imagem:

```shell
docker build -t tax-calculator-cli .
```

Rodar o container com STDIN aguardando o input manual:

```shell
docker run -i --rm tax-calculator-cli
```

\
Rodar container com arquivo de entrada contendo os 9 casos descritos no pdf do Code Challenge (pipe via STDIN).

Linux / macOS:

```shell
cat ./io/input.txt | docker run -i --rm tax-calculator-cli
```

Windows (PowerShell):

```shell
Get-Content .\io\input.txt | docker run -i --rm tax-calculator-cli
```

\
Observação: a imagem gerada inclui apenas o binário; o conteúdo do arquivo é enviado pelo STDIN. Note que nesse caso a aplicação termina logo em seguida.

## Testes

Para rodar os testes:

Linux / macOS:

```shell
docker run --rm -v "$(pwd)":/app -w /app golang:1.25 go test ./... -v
```

Windows (PowerShell):

```shell
docker run --rm -v ${PWD}:/app -w /app golang:1.25 go test ./... -v
```

## Makefile (testes e cobertura)

Este projeto inclui um `Makefile` para facilitar a execução de testes e a geração de relatórios de cobertura.

Alvos disponíveis:

- `make test`: executa todos os testes (`go test ./...`).
- `make coverage`: gera `cover.out` e imprime um resumo textual por função (`go tool cover -func`).
- `make cover-html`: gera o relatório HTML `coverage.html` a partir do `cover.out`.
- `make coverage-all`: mede cobertura cruzada entre pacotes com `-coverpkg=./...` e imprime resumo textual.
- `make clean`: remove artefatos de cobertura (`cover.out` e `coverage.html`).

Exemplos de uso (em um ambiente com `make` instalado, como Git Bash/WSL no Windows, ou nativamente no Linux/macOS):

```bash
make test
make coverage
make cover-html
```

### Pré-requisitos para usar o Makefile

Para executar os alvos do `Makefile` (test, coverage, cover-html, coverage-all, clean), garanta que:

- Go está instalado e no PATH
  - Compatível com a versão declarada em `go.mod` (atualmente Go 1.25.x).
  - Verifique com: `go version` e `go tool cover -h` (o `cover` vem junto com o Go).
- GNU Make está instalado e no PATH
  - Verifique com: `make --version`.
  - No Windows, use um ambiente com suporte a `make` (Git Bash ou WSL) ou instale via gerenciador (Chocolatey/Scoop).
- Permissão de leitura/escrita no diretório do projeto
  - Os alvos de cobertura criam/atualizam `cover.out` e `coverage.html` na raiz.

## Desempenho e limites

- Processamento streaming, linha a linha: memória estável para arquivos grandes, desde que processados por linhas.
- Complexidade linear no número de operações por linha.
