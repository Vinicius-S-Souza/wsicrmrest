# Webhook Zenvia SMS API

## Endpoint

**POST** `/webhook/zenvia/sms`

## Descrição

Este endpoint processa eventos de status de SMS recebidos do webhook Zenvia. Ele atualiza automaticamente o status de mensagens SMS no banco de dados e registra ocorrências quando SMS são rejeitados ou não entregues.

## Origem

Convertido da procedure WinDev `pgEventoZenvia_SMS()`.

## Fluxo de Processamento

1. Recebe o payload JSON do webhook Zenvia
2. Valida o tipo de callback (deve ser `MESSAGE_STATUS`)
3. Valida o formato do ID da mensagem (deve ser numérico)
4. Mapeia o código de status Zenvia para status interno
5. Busca dados do SMS no banco de dados (tabela `smsmensagem`)
6. Atualiza registros nas tabelas `logsapi` e `logsapihistorico`
7. Atualiza status na tabela `smsmensagem`
8. Se o SMS foi rejeitado (status 125), cria uma ocorrência na tabela `Ocorrencia` e limpa o celular
9. Registra toda a requisição na tabela `WSREQUISICOES`

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
    "to": "+5511999998888",
    "externalId": "54321"
  },
  "messageStatus": {
    "code": "DELIVERED",
    "description": "SMS entregue com sucesso",
    "causes": [
      {
        "reason": "INVALID_NUMBER",
        "details": "Número inválido"
      }
    ]
  }
}
```

### Campos do Payload

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| type | string | Sim | Tipo de callback. Deve ser `MESSAGE_STATUS` (case-insensitive) |
| message.to | string | Sim | Número de celular do destinatário (formato internacional) |
| message.externalId | string | Sim | ID da mensagem (smscodigo). Deve ser numérico |
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
| DELIVERED | 122 | Entregue | SMS entregue ao destinatário |
| REJECTED | 124 | NãoEntregue | SMS rejeitado |
| NOT_DELIVERED | 124 | NãoEntregue | SMS não entregue |

**Nota:** O código original WinDev verifica status 125 para criação de ocorrência, mas os eventos mapeados geram apenas 121, 122 e 124. O status 125 pode ser implementado em versões futuras para casos específicos de bounce.

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

### Consulta de SMS

```sql
SELECT s.smscodigo, s.clicodigo, l.logsapiid
FROM smsmensagem s
INNER JOIN logsapi l ON s.SMSAPIID = l.emsgcodigo
WHERE s.smscodigo = :smscodigo
AND l.logsapitipmensagem = 2
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
INSERT INTO logsApiHistorico(
    logsApiId,
    logsApiHisSequencial,
    logsApiHisTag,
    logsApiHisDescricao,
    logsapihisstatus,
    logsapihisdata
) VALUES(
    :logsApiId,
    (SELECT NVL(MAX(LogsApiHisSequencial), 0) + 1 FROM LogsApiHistorico WHERE LogsApiId = :logsApiId),
    :plataforma,
    :descricao,
    :status,
    :data
)
```

### Atualização de Status SMS

```sql
UPDATE smsmensagem
SET smsstsenvio = :status
WHERE smscodigo = :smscodigo
```

### Inserção de Ocorrência (Status 125)

```sql
-- Obtém próximo ID usando MAX+1
SELECT NVL(MAX(OcoCod), 0) + 1 FROM Ocorrencia

-- Insere ocorrência
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
    'Celular inválido. Não foi possível o envio de mensagem para esse celular, favor preencher o celular corretamente.',
    'Zenvia',
    SYSDATE,
    'Zenvia',
    SYSDATE
)
```

### Limpeza de Celular Inconsistente

```sql
UPDATE clientes
SET clicelular = ''
WHERE clicodigo = :clicodigo
```

## Logging

Todas as requisições são registradas na tabela `WSREQUISICOES` com:

- UUID único da requisição
- Método: `POST`
- Endpoint: `/webhook/zenvia/sms`
- Headers (exceto Authorization)
- Payload completo
- Resposta JSON
- Duração da requisição
- Código de retorno HTTP

Logs estruturados são gravados em arquivo com informações detalhadas:
- Início/fim do processamento
- Dados da mensagem (celular, evento, messageId)
- Resultados das consultas e operações no banco
- Warnings e erros

## Exemplos de Uso

### Exemplo 1: SMS Entregue

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "DELIVERED",
      "description": "SMS entregue com sucesso"
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
- Atualiza status para 122 (Entregue)
- Registra histórico com tag "Entregue"
- Atualiza campo `smsstsenvio` na tabela `smsmensagem`

### Exemplo 2: SMS Rejeitado

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999997777",
      "externalId": "54322"
    },
    "messageStatus": {
      "code": "REJECTED",
      "description": "SMS rejeitado pela operadora",
      "causes": [
        {
          "reason": "INVALID_NUMBER",
          "details": "Número de telefone inválido ou inexistente"
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
- Inclui detalhes da causa na descrição

### Exemplo 3: SMS Agendado

```bash
curl -X POST "http://localhost:8080/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "54323"
    },
    "messageStatus": {
      "code": "SENT",
      "description": "Mensagem enviada ao provedor"
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
- Atualiza status para 121 (Agendado)
- Registra histórico com tag "AgendadoProvedor"

## Validações

1. **Tipo de Callback:** Apenas `MESSAGE_STATUS` é processado. Outros tipos retornam 200 OK com mensagem informativa.

2. **Formato do ID:** O `externalId` deve ser numérico (regex: `^\d+$`). IDs alfanuméricos são rejeitados.

3. **Status Reconhecido:** Apenas os status listados na tabela de mapeamento são processados. Outros status retornam 200 OK com mensagem informativa.

4. **SMS Existente:** O SMS deve existir nas tabelas `smsmensagem` e `logsapi` com `logsapitipmensagem = 2`. Se não existir, retorna 200 OK com mensagem informativa.

5. **Celular do Cliente:** Para criar ocorrência, o celular no cadastro do cliente deve corresponder ao número que recebeu o SMS.

## Tratamento de Erros

- **Erro de Parsing JSON:** Retorna 400 Bad Request
- **SMS Não Encontrado:** Retorna 200 OK (não é erro, apenas log de warning)
- **Erro de Banco de Dados:** Registrado em log mas retorna 200 OK para não interromper o webhook
- **Status Não Reconhecido:** Retorna 200 OK (não é erro, apenas ignorado)
- **Celular Não Corresponde:** Não cria ocorrência mas continua processamento normal

## Considerações de Performance

- Todas as operações de log no banco (`GravaLogDB`) são executadas em goroutines para não bloquear a resposta
- Conexão com banco usa pool configurado (max 25 conexões, 5 idle)
- Não há validação de autenticação para não impactar performance (webhook de sistema)
- Consulta de sequencial em `logsApiHistorico` é otimizada com subquery

## Segurança

- O endpoint não exige autenticação JWT (é um webhook público)
- O header Authorization é lido mas não validado
- Todos os inputs são sanitizados antes de serem usados em queries SQL
- Queries usam prepared statements onde possível para prevenir SQL injection
- Caracteres especiais (aspas simples) são escapados antes de inserção

## Tabelas Afetadas

1. **logsapi** - Status atualizado
2. **logsapihistorico** - Novo registro de histórico inserido com sequencial auto-incrementado
3. **smsmensagem** - Status de envio (`smsstsenvio`) atualizado
4. **Ocorrencia** - Nova ocorrência inserida (apenas para status 125)
5. **clientes** - Campo `clicelular` limpo (apenas para status 125 com ocorrência)
6. **WSREQUISICOES** - Log da requisição inserido

## Diferenças entre SMS e Email

| Aspecto | SMS | Email |
|---------|-----|-------|
| Tabela principal | `smsmensagem` | `emailmensagem` |
| Campo de lookup | `smscodigo` | `Emsgapimsgid` |
| Join em logsapi | `s.SMSAPIID = l.emsgcodigo` | `e.emsgcodigo = l.emsgcodigo` |
| Tipo mensagem | 2 | 1 |
| Campo status | `smsstsenvio` | N/A |
| Status de bounce | 125 | 124 |
| Tipo ocorrência | 721 | 721 |
| Campo limpo | `clicelular` | `sCliExtEmail2` (ClientesExtensao) |
| Tabela de ocorrência | `Ocorrencia` | `ocorrencias` |

## Relacionamento WinDev

| Elemento WinDev | Elemento Go | Localização |
|----------------|-------------|-------------|
| `pgEventoZenvia_SMS()` | `ZenviaSMSWebhook()` | `internal/handlers/webhook_zenvia_sms.go` |
| `pgInsereLogsAPISms()` | `InsereLogsAPISMS()` | `internal/database/webhook.go:229` |
| `pgInsereOcorrenciaSmsInconsistente()` | `InsereOcorrenciaSmsInconsistente()` | `internal/database/webhook.go:373` |
| `pgInsertLogsApiHistorico()` | `InsertLogsApiHistorico()` | `internal/database/webhook.go:296` |
| `pgSetMsgStatusSMS()` | `SetMsgStatusSMS()` | `internal/database/webhook.go:348` |
| `pgLimpaCelularInconsistente()` | `LimpaCelularInconsistente()` | `internal/database/webhook.go:467` |
| `pgCmdInsert()` / `pgValorMaxCampo()` | MAX(OcoCod)+1 para auto-increment | `internal/database/webhook.go:468-478` |
| `pgImprimirLog()` | `logger.Infow()`/`Errorw()` | Zap logger |
| `SQLExec()`/`SQLFetch()` | `db.QueryRow().Scan()` | `internal/database/webhook.go` |

## Teste

Use o script fornecido para testar todos os cenários:

```bash
./test_webhook_zenvia_sms.sh
```

O script testa:
1. Evento SENT (Agendado)
2. Evento DELIVERED (Entregue)
3. Evento REJECTED (Rejeitado com causa)
4. Evento NOT_DELIVERED (Não entregue com causa)
5. Tipo de callback inválido
6. ID não numérico
7. Status não reconhecido

## Troubleshooting

**SMS não encontrado no banco:**
- Verifique se o `smscodigo` está correto no payload
- Confirme que existe registro em `smsmensagem` com esse código
- Verifique se existe entrada em `logsapi` com `logsapitipmensagem = 2`
- Confirme o join via `s.SMSAPIID = l.emsgcodigo`

**Ocorrência não criada para SMS rejeitado:**
- O código atual só cria ocorrência para status 125 (bounce)
- Status 124 (REJECTED/NOT_DELIVERED) não cria ocorrência automaticamente
- Verifique se o celular do cliente corresponde ao número que recebeu o SMS

**Erro ao atualizar histórico:**
- Verifique se a sequence `seq_ococod` existe no banco
- Confirme que a tabela `logsApiHistorico` tem a estrutura correta
- Verifique permissões de escrita no banco

**Celular não foi limpo:**
- Ocorrência precisa ser criada primeiro (status 125)
- Celular no cadastro deve corresponder ao número do SMS
- Verifique se há erro de banco na operação anterior
