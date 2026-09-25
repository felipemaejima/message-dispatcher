# message-dispatcher

Um **processador de filas extensível** escrito em Go para despachar mensagens por múltiplos serviços de comunicação.

A ideia central do projeto é a **plugabilidade**: tanto o broker de filas quanto os serviços de disparo são implementações intercambiáveis. Novos brokers (Kafka, SQS, etc.) e novos serviços de disparo (WhatsApp, SMS, Slack, etc.) podem ser adicionados sem alterar o núcleo da aplicação.

Atualmente implementados:
- **Consumer (broker):** RabbitMQ
- **Notifier (disparo):** Resend (e-mail)

---

## Como funciona

```
[ Provider: Consumer ]  →  [ Application ]  →  [ Provider: Notifier ]
   (ex: RabbitMQ)            dispatcher.go        (ex: Resend, WhatsApp,
                             roteia pelo            SMS, Slack...)
                             campo "channel"
```

1. O consumer (provider) se conecta ao broker e fica aguardando mensagens
2. Ao receber uma mensagem, repassa para o `dispatcher` na camada de aplicação
3. O dispatcher lê o campo `channel` do payload e roteia para o notifier correspondente
4. O notifier (provider) executa a entrega pelo serviço configurado

---

## Estrutura do projeto

O projeto segue a **Arquitetura Hexagonal (Ports & Adapters)**: o domínio e a lógica de aplicação são completamente isolados das implementações externas. Novos brokers ou serviços de disparo são adicionados apenas criando um novo provider — sem alterar nada no núcleo.

```
message-dispatcher/
├── docker-compose.yml                      # Ambiente local (RabbitMQ)
├── Makefile                                # Comandos de desenvolvimento
├── cmd/
│   └── main.go                             # Entrypoint da aplicação
├── internal/
│   ├── application/
│   │   └── dispatcher.go                   # Orquestra o roteamento das mensagens
│   ├── domain/
│   │   └── notification.go                 # Entidade central do domínio
│   ├── ports/
│   │   ├── consumer.go                     # Interface: broker de filas
│   │   └── notifier.go                     # Interface: serviço de disparo
│   └── providers/
│       ├── consumers/
│       │   └── rabbitmq/
│       │       ├── consumer.go             # Implementação: RabbitMQ
│       │       └── consumer_integration_test.go
│       ├── logger/
│       │   └── logger.go                   # Logger da aplicação
│       └── notifiers/
│           └── resend/
│               ├── notifier.go             # Implementação: Resend (e-mail)
│               └── notifier_integration_test.go
├── storage/
│   └── logs/                               # Logs gerados pela aplicação
├── .env.example
├── payload.json                            # Exemplo de payload para testes
├── go.mod
└── go.sum
```

### Camadas

| Camada | Responsabilidade |
|---|---|
| `domain` | Entidades e regras de negócio puras, sem dependências externas |
| `ports` | Interfaces (contratos) que definem o que consumers e notifiers devem implementar |
| `application` | Orquestra o fluxo: recebe a notificação e despacha para o notifier correto |
| `providers` | Implementações concretas dos ports (RabbitMQ, Resend, etc.) |

---

## Tecnologias

| Tecnologia | Uso |
|---|---|
| [Go](https://golang.org) | Linguagem principal |
| [RabbitMQ (amqp091-go)](https://github.com/rabbitmq/amqp091-go) | Consumer provider (implementação atual) |
| [Resend](https://github.com/resend/resend-go) | Notifier provider de e-mail (implementação atual) |
| [godotenv](https://github.com/joho/godotenv) | Carregamento de variáveis de ambiente |

---

## Como executar

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) com Compose
- `make` (`sudo apt install make`)
- Conta no [Resend](https://resend.com) com domínio verificado

Não é necessário ter Go instalado: o toolchain roda em container.

### 1. Clone o repositório

```bash
git clone https://github.com/felipemaejima/message-dispatcher.git
cd message-dispatcher
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

Preencha `RESEND_API_KEY` e `EMAIL_FROM`. Os valores de RabbitMQ já vêm apontando
para o `docker-compose.yml`.

### 3. Suba o ambiente e rode

```bash
make up    # sobe o RabbitMQ e declara a fila
make run   # roda o dispatcher
```

O dispatcher ficará aguardando mensagens na fila configurada.

Em outro terminal:

```bash
make publish   # publica o payload.json na fila
make peek      # mostra quantas mensagens estão na fila
```

### Comandos disponíveis

`make help` lista todos. Os principais:

| Comando | Descrição |
|---|---|
| `make up` / `make down` | Sobe / derruba o RabbitMQ |
| `make run` | Roda o dispatcher |
| `make publish` | Publica `payload.json` na fila |
| `make peek` | Conta mensagens na fila |
| `make logs` | Logs do RabbitMQ |
| `make test` | Testes unitários (sem rede) |
| `make test-integration` | Testes de integração (envia e-mail real) |
| `make build` | Compila o binário local |
| `make shell` | Shell no container Go |

A UI de gerenciamento do RabbitMQ fica em http://localhost:15672 (`app` / `app`).

---

## Formato do Payload

As mensagens publicadas na fila devem seguir o seguinte formato JSON:

```json
{
  "channel": "resend",
  "message": {
    "to": [
      "destinatario@exemplo.com"
    ],
    "subject": "Assunto da mensagem",
    "body": "<h1>Conteúdo HTML</h1>"
  }
}
```

### Campos

| Campo | Tipo | Descrição |
|---|---|---|
| `channel` | `string` | Nome do notifier de destino (ex: `resend`, `whatsapp`) |
| `message.to` | `array<string>` | Lista de destinatários |
| `message.subject` | `string` | Assunto (específico por notifier) |
| `message.body` | `string` | Corpo da mensagem (suporta HTML) |

O campo `channel` é o identificador que o dispatcher usa para rotear a mensagem ao notifier correto.

---

## Extensibilidade

Para adicionar uma nova implementação, basta criar um novo provider respeitando o contrato definido nos `ports` e registrá-lo no dispatcher.

### Novo notifier (serviço de disparo)

Crie uma pasta em `internal/providers/notifiers/<nome>/` implementando a interface `ports/notifier.go`. Exemplos de notifiers que podem ser adicionados:

- WhatsApp
- SMS
- Slack
- Telegram
- Push Notification

### Novo consumer (broker de filas)

Crie uma pasta em `internal/providers/consumers/<nome>/` implementando a interface `ports/consumer.go`. Exemplos de consumers que podem ser adicionados:

- Apache Kafka
- AWS SQS
- Redis Streams
- Google Pub/Sub

---

## Variáveis de Ambiente

| Variável | Descrição |
|---|---|
| `ENVIRONMENT` | Ambiente de execução (`development`, `production`) |
| `RESEND_API_KEY` | Chave de API do Resend |
| `EMAIL_FROM` | Domínio verificado para envio de e-mails (ex: `email@seudominio.com`) |
| `RABBITMQ_BROKER_URL` | URL de conexão com o RabbitMQ (formato AMQP) |
| `MESSAGES_QUEUE` | Nome da fila a ser consumida |

---

## Licença

Este projeto está sob a licença MIT.

---

Desenvolvido por [Felipe Maejima](https://github.com/felipemaejima)
