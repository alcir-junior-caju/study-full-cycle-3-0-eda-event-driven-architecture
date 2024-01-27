# Curso Full Cycle 3.0 - Módulo EDA Event Driven Architecture

<div>
    <img alt="Criado por Alcir Junior [Caju]" src="https://img.shields.io/badge/criado%20por-Alcir Junior [Caju]-%23f08700">
    <img alt="License" src="https://img.shields.io/badge/license-MIT-%23f08700">
</div>

---

## Descrição

O Curso Full Cycle é uma formação completa para fazer com que pessoas desenvolvedoras sejam capazes de trabalhar em projetos expressivos sendo capazes de desenvolver aplicações de grande porte utilizando de boas práticas de desenvolvimento.

---

## Desafio

Desenvolva um microsserviço em sua linguagem de preferência que seja capaz de receber via Kafka os eventos gerados pelo microsserviço "Wallet Core" e persistir no banco de dados os balances atualizados para cada conta.

Crie um endpoint: `/balances/{account_id}` que exibe o balance atualizado.

Requisitos para entrega:
- [x] Tudo deve rodar via Docker / Docker-compose
- [x] Com um único docker-compose up -d todos os microsserviços, incluindo o da wallet core precisam estar disponíveis para que possamos fazer a correção.
- [x] Não esqueça de rodar migrations e popular dados fictícios em ambos bancos de dados (wallet core e o microsserviço de balances) de forma automática quando os serviços subirem.
- [x] Gere o arquivo ".http" para realizarmos as chamadas em seu microsserviço da mesma forma que fizemos no microsserviço "wallet core"
- [x] Disponibilize o microsserviço na porta: 3003.

---

## Instruções

Para fazer o `build` de todos serviços no diretório raiz executo o comando `docker compose up --build` com isso as imagens serão baixadas e executadas.

Após todos containers estiverem em rodando, inicie as aplicações:

### Serviço de Wallet

Execute o comando `docker compose exec wallet-app bash` após acessar o container execute o comando `go run cmd/walletcore/main.go` ele irá rodas as migrations e subir o sevidor na porta `8080`.

### Serviço de Balances

Execute o comando `docker compose exec balances-app bash` após acessar o container execute o comando `go run cmd/balances/main.go` ele irá rodas as migrations e subir o sevidor na porta `3003`.

Os dois serviços já tem seu arquivo de `api/client.http` já com os `IDs` corretos, mas nada impede te criar novos registros e usar os mesmos.

---

## Repositório Pai
https://github.com/alcir-junior-caju/study-full-cycle-3-0

---

## Visualizar o projeto na IDE:

Para quem quiser visualizar o projeto na IDE clique no teclado a tecla `ponto`, esse recurso do GitHub é bem bacana

---

### O que é eventos?
- Situações que ocorreram no passado;
- Normalmente deixa efeitos colaterais:
    - Ex.: A porta do carro abriu, Ligar a luz interna;
- Ele pode trabalhar de forma internalizada no software ou externalizada;
- Domain Events (Eventos de domínio):
    - Mudança no estado interno da aplicação / regra de negócios -> Ex.: agregados;

### Tipos de eventos
1. Event Notification: Forma curta de comunicação, ex.: `{ "order": 1, "status": "approved" }`;
2. Event Carried State Transfer: Formato completo para trafegar as informações, ex.: `stream de dados` `{ "order": 1, "status": "approved", "products": [{ ... }, { ... }], "tax": "1%", "client": "John Doe" }`;
3. Event Sourcing: Armazenamento dos eventos baseado em uma linha do tempo, possibilidade de replay;

### Event Colaboration
- Metodo tradicional, precisa consultar informações:
    - John Doe compra um produto -> Checar Estoque do produto comprado -> Muda o catálogo -> Emite nota -> Separação -> ... -> ...;
- Event Colaboration, parte do princípio que já tem as informações para uma determinada ação:
    - John Doe compra um produto;
    - Estoque mudou;
    - Cor Mudou;
    - Nota foi emitida;
    - Erro aconteceu;
    - Produto mudou a descrição;

### CQRS (Command Query Responsibility Segregation) + Event Sourcing
- CQS x CQRS
    - Nível de granularidade;
- Comando: Intenção de mudança do usuário e não tem retorno:
    - Criar Produto
- Consulta: Foca apenas no dado;
- Pode se ter bancos de dados separados para Escrita e Leitura:
    - Escrita: `{ "id": 1, "product": 2, "category": 3 }` -> dispara um evento para atualizar o bando de leitura;
    - Leitura: `{ "id": 1, "product": "Carrinho", "categoria": "Brinquedo" }` -> Informação já consolidada;
_ CQRS + Event Sorcing: Não é uma regra mas é recomendado, pois a consolidação dos dados pode ser feita utilizando o Event Sorcing;

### Elementos táticos de um contexto de eventos
- Evento (Carregar dados);
- Operações que serão executadas quando um evento é chamado;
- Gerenciador dos nossos eventos / operações;
    - Registrar os eventos e suas operações;
    - Despachar / Fire no evento para que suas operações sejam executadas;
