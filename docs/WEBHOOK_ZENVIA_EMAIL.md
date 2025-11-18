# Webhook Zenvia Email API

## Endpoint

**POST** `/webhook/zenvia/email`

## Descrição

Este endpoint processa eventos de status de email recebidos do webhook Zenvia. Ele atualiza automaticamente o status de mensagens de email no banco de dados e registra ocorrências quando emails são rejeitados ou não entregues.

## Origem

Convertido da procedure WinDev `pgEventoZenvia_Email()`.

## Fluxo de Processamento

1. Recebe o payload JSON do webhook Zenvia
2. Valida o tipo de callback (deve ser `MESSAGE_STATUS`)
3. Valida o formato do ID da mensagem (deve ser numérico)
4. Mapeia o código de status Zenvia para status interno
5. Busca dados do email no banco de dados (tabela `emailmensagem`)
6. Atualiza registros nas tabelas `logsapi` e `logsapihistorico`
7. Se o email foi rejeitado (status 124), cria uma ocorrência na tabela `ocorrencias`
8. Registra toda a requisição na tabela `WSREQUISICOES`

## Headers

| Header | Tipo | Obrigatório | Descrição |
|--------|------|-------------|-----------|
| Content-Type | string | Sim | Deve ser `application/json` |
| Authorization | string | Opcional | Token de autorização (registrado mas não validado) |

## Payload JSON

```json
{
  "type": "MESSAGE_STATUS",
  "message": {
    "to": "destinatario@example.com",
    "externalId": "12345"
  },
  "messageStatus": {
    "code": "DELIVERED",
    "description": "Email entregue com sucesso",
    "causes": [
      {
        "reason": "INVALID_EMAIL",
        "details": "Endereço inválido"
      }
    ]
  }
}
```

### Campos do Payload

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| type | string | Sim | Tipo de callback. Deve ser `MESSAGE_STATUS` (case-insensitive) |
| message.to | string | Sim | Endereço de email do destinatário |
| message.externalId | string | Sim | ID da mensagem (emsgcodigo). Deve ser numérico |
| messageStatus.code | string | Sim | Código do status da mensagem |
| messageStatus.description | string | Sim | Descrição do status (máx. 1000 caracteres) |
| messageStatus.causes | array | Não | Array de causas (usado para status REJECTED/NOT_DELIVERED) |
| messageStatus.causes[].reason | string | Não | Motivo da rejeição |
| messageStatus.causes[].details | string | Não | Detalhes da rejeição |

## Mapeamento de Status

O sistema mapeia os códigos de status Zenvia para status internos:

| Status Zenvia | Status Interno | Tag | Descrição |
|---------------|----------------|-----|-----------|
| SENT | 121 | AgendadoProvedor | Mensagem agendada/enviada ao provedor |
| DELIVERED | 122 | Entregue | Email entregue ao destinatário |
| READ | 123 | Aberto | Email foi aberto |
| CLICKED | 123 | Aberto | Email teve um link clicado |
| REJECTED | 124 | NãoEntregue | Email rejeitado (cria ocorrência) |
| NOT_DELIVERED | 124 | NãoEntregue | Email não entregue (cria ocorrência) |

## Respostas

### Sucesso (200 OK)

```json
{
  "success": true,
  "message": "Webhook processado com sucesso"
}
```

### Tipo de Mensagem Não Processado (200 OK)

```json
{
  "success": true,
  "message": "Tipo de mensagem não processado: other_type"
}
```

### ID Inválido (200 OK)

```json
{
  "success": true,
  "message": "ID não corresponde ao formato esperado"
}
```

### Mensagem Não Encontrada (200 OK)

```json
{
  "success": true,
  "message": "Mensagem não encontrada no banco de dados"
}
```

### Erro ao Ler Body (400 Bad Request)

```json
{
  "success": false,
  "message": "Erro ao ler body da requisição"
}
```

### JSON Inválido (400 Bad Request)

```json
{
  "success": false,
  "message": "JSON inválido"
}
```

## Operações no Banco de Dados

### Consulta de Email

```sql
SELECT e.emsgcodigo, e.clicodigo, l.logsapiid
FROM emailmensagem e
INNER JOIN logsapi l ON e.emsgcodigo = l.emsgcodigo
WHERE e.Emsgapimsgid = :externalId
AND l.logsapitipmensagem = 1
```

### Atualização de Logs API

```sql
UPDATE logsapi
SET logsapistatus = :status,
    logsapiretorno = :retorno,
    logsapidtaatualizacao = :data
WHERE logsapiid = :id
```

### Inserção de Histórico

```sql
INSERT INTO logsapihistorico(
    logsapihisid,
    logsapiid,
    logsapihisdta,
    logsapihisstatus,
    logsapihisdescricao,
    logsapitag
) VALUES(
    seq_logsapihisid.NEXTVAL,
    :logsapiid,
    :data,
    :status,
    :descricao,
    :tag
)
```

### Inserção de Ocorrência (Status 125)

```sql
-- Primeiro busca dados do cliente e emails da extensão
SELECT c.clinome, c.clicpfcnpj, ce.CliExtEmail2, ce.CliExtEmail3
FROM clientes c
INNER JOIN Clientesextensao ce ON c.clicodigo = ce.clicodigo
WHERE c.clicodigo = :clicodigo

-- Obtém próximo ID usando MAX+1
SELECT NVL(MAX(OcoCod), 0) + 1 FROM Ocorrencia

-- Insere ocorrência se o email corresponder
INSERT INTO Ocorrencia(
    OcoCod,
    EntCod,
    CliCod,
    CodPltCod,
    TocNum,
    OcoTip,
    OcoCliNon,
    OcoDsc,
    OcoUsrSol,
    OcoSolDta,
    UsrAlt,
    DatCad
) VALUES(
    :max_ococod_plus_1,
    :organizacao_codigo,
    :clicodigo,
    :cpf_cnpj,
    721,
    2,
    :cli_nome,
    'Email inválido. Não foi possível o envio de mensagem para esse email, favor preencher o email corretamente.',
    'WebHookSendGrid',
    SYSDATE,
    'WebHookSendGrid',
    SYSDATE
)

-- Limpa o campo de email inconsistente
UPDATE Clientesextensao
SET CliExtEmail2 = ''  -- ou CliExtEmail3, dependendo do email
WHERE clicodigo = :clicodigo
```

**Nota:** O sistema verifica se o email inconsistente corresponde a `CliExtEmail2` ou `CliExtEmail3` antes de criar a ocorrência.

## Logging

Todas as requisições são registradas na tabela `WSREQUISICOES` com:

- UUID único da requisição
- Método: `POST`
- Endpoint: `/webhook/zenvia/email`
- Headers (exceto Authorization)
- Payload completo
- Resposta JSON
- Duração da requisição
- Código de retorno HTTP

Logs estruturados são gravados em arquivo com informações detalhadas:
- Início/fim do processamento
- Dados da mensagem (email, evento, messageId)
- Resultados das consultas e operações no banco
- Warnings e erros

## Exemplos de Uso

### Exemplo 1: Email Entregue

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "DELIVERED",
      "description": "Email entregue com sucesso"
    }
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Webhook processado com sucesso"
}
```

### Exemplo 2: Email Rejeitado

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "invalido@example.com",
      "externalId": "54322"
    },
    "messageStatus": {
      "code": "REJECTED",
      "description": "Email rejeitado pelo servidor",
      "causes": [
        {
          "reason": "INVALID_EMAIL",
          "details": "Endereço de email inválido ou inexistente"
        }
      ]
    }
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Webhook processado com sucesso"
}
```

**Ações executadas:**
- Atualiza status para 124 (Não Entregue)
- Registra histórico com tag "NãoEntregue"
- **Nota:** Status 124 não cria ocorrência automaticamente. Apenas status 125 cria ocorrência.

### Exemplo 3: Email Aberto

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
      "externalId": "54323"
    },
    "messageStatus": {
      "code": "READ",
      "description": "Email foi aberto pelo destinatário"
    }
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Webhook processado com sucesso"
}
```

## Validações

1. **Tipo de Callback:** Apenas `MESSAGE_STATUS` é processado. Outros tipos retornam 200 OK com mensagem informativa.

2. **Formato do ID:** O `externalId` deve ser numérico (regex: `^\d+$`). IDs alfanuméricos são rejeitados.

3. **Status Reconhecido:** Apenas os status listados na tabela de mapeamento são processados. Outros status retornam 200 OK com mensagem informativa.

4. **Email Existente:** O email deve existir nas tabelas `emailmensagem` e `logsapi`. Se não existir, retorna 200 OK com mensagem informativa.

## Tratamento de Erros

- **Erro de Parsing JSON:** Retorna 400 Bad Request
- **Email Não Encontrado:** Retorna 200 OK (não é erro, apenas log de warning)
- **Erro de Banco de Dados:** Registrado em log mas retorna 200 OK para não interromper o webhook
- **Status Não Reconhecido:** Retorna 200 OK (não é erro, apenas ignorado)

## Considerações de Performance

- Todas as operações de log no banco (`GravaLogDB`) são executadas em goroutines para não bloquear a resposta
- Conexão com banco usa pool configurado (max 25 conexões, 5 idle)
- Não há validação de autenticação para não impactar performance (webhook de sistema)

## Segurança

- O endpoint não exige autenticação JWT (é um webhook público)
- O header Authorization é lido mas não validado
- Todos os inputs são sanitizados antes de serem usados em queries SQL
- Queries usam prepared statements para prevenir SQL injection

## Tabelas Afetadas

1. **logsapi** - Status atualizado
2. **logsapihistorico** - Novo registro de histórico inserido com sequencial auto-incrementado
3. **EmailMensagem** - Campo `EMsgStsEnvio` atualizado
4. **Ocorrencia** - Nova ocorrência inserida (apenas para status 125)
5. **Clientesextensao** - Campo `CliExtEmail2` ou `CliExtEmail3` limpo (apenas para status 125 com ocorrência)
6. **WSREQUISICOES** - Log da requisição inserido

## Relacionamento WinDev

| Elemento WinDev | Elemento Go | Localização |
|----------------|-------------|-------------|
| `pgEventoZenvia_Email()` | `ZenviaEmailWebhook()` | `internal/handlers/webhook_zenvia_email.go:21` |
| `pgInsereLogsAPI()` | `InsereLogsAPI()` | `internal/database/webhook.go:49` |
| `pgInsereOcorrenciaEmailInconsistente()` | `InsereOcorrenciaEmailInconsistente()` | `internal/database/webhook.go:132` |
| `pgLimpaEmailInconsistente()` | `LimpaEmailInconsistente()` | `internal/database/webhook.go:232` |
| `pgInsertLogsApiHistorico()` | `InsertLogsApiHistorico()` | `internal/database/webhook.go:296` (shared with SMS) |
| `pgSetMsgStatus()` | `SetMsgStatus()` | `internal/database/webhook.go:112` |
| `pgCmdInsert()` / `pgValorMaxCampo()` | MAX(OcoCod)+1 para auto-increment | `internal/database/webhook.go:170-180` |
| `pgImprimirLog()` | `logger.Infow()`/`Errorw()` | Zap logger |
| `SQLExec()`/`SQLFetch()` | `db.QueryRow().Scan()` | `internal/database/webhook.go` |

## Teste

Use o script fornecido para testar todos os cenários:

```bash
./test_webhook_zenvia_email.sh
```

Ou use o Makefile:

```bash
make test-webhook-zenvia
```
