# Configuração CORS (Cross-Origin Resource Sharing)

## Visão Geral

O WSICRMREST implementa suporte completo para CORS (Cross-Origin Resource Sharing), permitindo que aplicações web de diferentes origens façam requisições à API de forma segura e controlada.

## Por Que CORS é Necessário?

Por padrão, navegadores web implementam a **Same-Origin Policy** (Política de Mesma Origem), que impede que scripts em uma página web façam requisições para um domínio diferente do que serviu a página.

### Exemplo de Problema sem CORS:

```
Frontend: https://app.example.com
API:      https://api.example.com

❌ Erro: O navegador bloqueia a requisição por violação de CORS
```

### Como CORS Resolve:

O servidor API adiciona headers HTTP especiais informando ao navegador quais origens estão autorizadas a fazer requisições.

## Configuração

### Arquivo: `dbinit.ini`

Adicione a seção `[CORS]` no arquivo `dbinit.ini`:

```ini
[CORS]
# Origens permitidas separadas por vírgula (deixe vazio para permitir todas - desenvolvimento)
# Exemplo: AllowedOrigins=https://app.example.com,https://admin.example.com
AllowedOrigins=

# Métodos HTTP permitidos
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS

# Headers permitidos (aceita tanto grant-type quanto Grant_type)
AllowedHeaders=Origin,Content-Type,Content-Length,Accept-Encoding,Authorization,Grant_type,X-CSRF-Token

# Permite credenciais (cookies, authorization headers)
AllowCredentials=true

# Tempo de cache do preflight em segundos (12 horas = 43200)
MaxAge=43200
```

## Parâmetros de Configuração

### 1. AllowedOrigins

**Lista de origens permitidas** (separadas por vírgula).

**Modo Desenvolvimento (Permissivo):**
```ini
AllowedOrigins=
```
- ✅ Permite **TODAS** as origens (`Access-Control-Allow-Origin: *`)
- ⚠️ **Use apenas em desenvolvimento**

**Modo Produção (Restritivo):**
```ini
AllowedOrigins=https://app.example.com,https://admin.example.com
```
- ✅ Permite apenas as origens listadas
- 🔒 **Recomendado para produção**

**Exemplo com múltiplas origens:**
```ini
AllowedOrigins=https://crm.example.com,https://painel.example.com,https://app.example.com
```

### 2. AllowedMethods

**Métodos HTTP permitidos** nas requisições CORS.

**Valor Padrão:**
```ini
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
```

**Customização:**
```ini
# Apenas leitura
AllowedMethods=GET,OPTIONS

# Leitura e escrita
AllowedMethods=GET,POST,PUT,DELETE,OPTIONS
```

### 3. AllowedHeaders

**Headers HTTP permitidos** nas requisições.

**Valor Padrão:**
```ini
AllowedHeaders=Origin,Content-Type,Content-Length,Accept-Encoding,Authorization,Grant_type,X-CSRF-Token
```

**Headers Importantes:**

| Header | Descrição |
|--------|-----------|
| `Origin` | Identifica a origem da requisição |
| `Content-Type` | Tipo de conteúdo (application/json) |
| `Authorization` | Token Bearer JWT |
| `Grant_type` | Para geração de token (compatibilidade WinDev) |
| `X-CSRF-Token` | Proteção contra CSRF |

**Nota:** O header `Grant_type` está incluído para manter compatibilidade com o sistema WinDev original, que usa underscore ao invés de hífen.

### 4. AllowCredentials

**Permite envio de credenciais** (cookies, headers de autenticação).

**Valor Padrão:**
```ini
AllowCredentials=true
```

**Opções:**
- `true` ou `1` - Permite credenciais (recomendado para autenticação JWT)
- `false` ou `0` - Não permite credenciais

**Importante:** Se `AllowCredentials=true`, você **não pode** usar `AllowedOrigins=*`. Deve especificar origens exatas.

### 5. MaxAge

**Tempo de cache do preflight** em segundos.

**Valor Padrão:**
```ini
MaxAge=43200  # 12 horas
```

**Exemplos:**
```ini
MaxAge=3600   # 1 hora
MaxAge=86400  # 24 horas
MaxAge=43200  # 12 horas (recomendado)
```

Quanto maior o valor, menos requisições OPTIONS (preflight) o navegador faz, melhorando a performance.

## Headers HTTP Adicionados pelo Middleware

O middleware CORS adiciona automaticamente os seguintes headers HTTP:

### Headers Principais:

1. **Access-Control-Allow-Origin**
   - Especifica qual origem está autorizada
   - Valor: `*` (todas) ou origem específica

2. **Access-Control-Allow-Methods**
   - Métodos HTTP permitidos
   - Exemplo: `GET,POST,PUT,PATCH,DELETE,OPTIONS`

3. **Access-Control-Allow-Headers**
   - Headers que podem ser usados na requisição
   - Exemplo: `Origin,Content-Type,Authorization`

4. **Access-Control-Allow-Credentials**
   - Permite envio de cookies e headers de autenticação
   - Valor: `true` ou ausente

5. **Access-Control-Max-Age**
   - Tempo de cache do preflight em segundos
   - Exemplo: `43200`

6. **Access-Control-Expose-Headers**
   - Headers que o navegador pode acessar
   - Valor: `Content-Length`

7. **Vary: Origin**
   - Usado quando origens específicas são permitidas
   - Garante cache correto em CDNs

## Requisições Preflight (OPTIONS)

### O Que É Preflight?

Navegadores fazem uma **requisição OPTIONS** antes da requisição real quando:
- Método não é GET, HEAD ou POST simples
- Usa headers customizados (Authorization, Content-Type: application/json)
- Faz requisições com credenciais

### Como Funciona:

```
1. Navegador envia OPTIONS /api/endpoint
   ├─ Origin: https://app.example.com
   ├─ Access-Control-Request-Method: POST
   └─ Access-Control-Request-Headers: authorization,content-type

2. Servidor responde com headers CORS
   ├─ Access-Control-Allow-Origin: https://app.example.com
   ├─ Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
   ├─ Access-Control-Allow-Headers: authorization,content-type
   └─ HTTP Status: 204 No Content

3. Navegador faz a requisição real
   POST /api/endpoint
   └─ Authorization: Bearer token...
```

### Tratamento no WSICRMREST:

O middleware detecta requisições OPTIONS e responde automaticamente com:
- Status: `204 No Content`
- Headers CORS apropriados
- Sem processamento adicional

## Cenários de Uso

### Cenário 1: Desenvolvimento Local

**Situação:**
- Frontend em `http://localhost:3000`
- API em `http://localhost:8080`

**Configuração:**
```ini
[CORS]
AllowedOrigins=
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
AllowedHeaders=Origin,Content-Type,Authorization,Grant_type
AllowCredentials=true
MaxAge=43200
```

**Resultado:** Permite todas as origens (`*`)

### Cenário 2: Produção com Domínio Único

**Situação:**
- Frontend em `https://app.example.com`
- API em `https://api.example.com`

**Configuração:**
```ini
[CORS]
AllowedOrigins=https://app.example.com
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
AllowedHeaders=Origin,Content-Type,Authorization,Grant_type
AllowCredentials=true
MaxAge=86400
```

**Resultado:** Apenas `https://app.example.com` pode fazer requisições

### Cenário 3: Produção com Múltiplos Domínios

**Situação:**
- Frontend principal: `https://app.example.com`
- Painel admin: `https://admin.example.com`
- Aplicativo mobile web: `https://mobile.example.com`

**Configuração:**
```ini
[CORS]
AllowedOrigins=https://app.example.com,https://admin.example.com,https://mobile.example.com
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
AllowedHeaders=Origin,Content-Type,Authorization,Grant_type
AllowCredentials=true
MaxAge=43200
```

**Resultado:** Apenas os três domínios listados podem fazer requisições

### Cenário 4: API Pública (Sem Autenticação)

**Situação:**
- API pública acessível de qualquer origem
- Sem uso de cookies ou tokens

**Configuração:**
```ini
[CORS]
AllowedOrigins=
AllowedMethods=GET,OPTIONS
AllowedHeaders=Origin,Content-Type
AllowCredentials=false
MaxAge=86400
```

**Resultado:** Qualquer origem pode fazer GET, sem credenciais

## Logs de CORS

### Log de Inicialização

Quando o servidor inicia, os logs mostram a configuração CORS:

**Modo Desenvolvimento:**
```
INFO CORS configurado para permitir TODAS as origens (*) - Modo Desenvolvimento
DEBUG Configurações CORS methods=GET,POST,PUT,PATCH,DELETE,OPTIONS headers=Origin,Content-Type,Authorization credentials=true max_age=43200
```

**Modo Produção:**
```
INFO CORS configurado com origens restritas allowed_origins=[https://app.example.com, https://admin.example.com]
DEBUG Configurações CORS methods=GET,POST,PUT,PATCH,DELETE,OPTIONS headers=Origin,Content-Type,Authorization credentials=true max_age=43200
```

## Testando CORS

### 1. Teste Simples com cURL

```bash
# Simular requisição de origem diferente
curl -X OPTIONS http://localhost:8080/connect/v1/token \
  -H "Origin: https://app.example.com" \
  -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: authorization" \
  -v
```

**Resposta Esperada:**
```
< HTTP/1.1 204 No Content
< Access-Control-Allow-Origin: *
< Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
< Access-Control-Allow-Headers: Origin,Content-Type,Authorization,Grant_type
< Access-Control-Allow-Credentials: true
< Access-Control-Max-Age: 43200
```

### 2. Teste com JavaScript no Navegador

Abra o console do navegador em uma página de origem diferente:

```javascript
// Fazer requisição GET para testar CORS
fetch('http://localhost:8080/connect/v1/wsteste', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log('Sucesso:', data))
.catch(error => console.error('Erro CORS:', error));
```

**Sucesso:** Resposta JSON aparece no console
**Erro:** Mensagem de bloqueio CORS

### 3. Teste com Postman

Postman **não** valida CORS (ferramentas de API não são navegadores).

Para testar CORS, use:
- Navegador web (Chrome, Firefox)
- Ferramentas de teste CORS online
- Console do navegador com fetch/XMLHttpRequest

## Problemas Comuns e Soluções

### Erro: "CORS policy: No 'Access-Control-Allow-Origin' header"

**Causa:** Origem não está na lista de `AllowedOrigins`

**Solução:**
```ini
# Adicione a origem na lista
AllowedOrigins=https://app.example.com,https://nova-origem.com
```

### Erro: "Credential is not supported if the CORS header 'Access-Control-Allow-Origin' is '*'"

**Causa:** `AllowCredentials=true` com `AllowedOrigins=` vazio

**Solução:** Especifique origens exatas:
```ini
AllowedOrigins=https://app.example.com
AllowCredentials=true
```

### Erro: "Method PUT is not allowed by Access-Control-Allow-Methods"

**Causa:** Método não está em `AllowedMethods`

**Solução:**
```ini
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
```

### Erro: "Request header Authorization is not allowed"

**Causa:** Header não está em `AllowedHeaders`

**Solução:**
```ini
AllowedHeaders=Origin,Content-Type,Authorization,Grant_type
```

### CORS Funciona no Postman mas não no Navegador

**Causa:** Postman não valida CORS (não é um navegador)

**Solução:** Configure CORS corretamente para navegadores web. Teste no navegador ou com ferramentas específicas de CORS.

## Segurança

### Boas Práticas:

1. **Produção: Sempre especifique origens**
   ```ini
   # ❌ Evite em produção
   AllowedOrigins=

   # ✅ Use origens específicas
   AllowedOrigins=https://app.example.com
   ```

2. **Use HTTPS em produção**
   ```ini
   # ✅ Correto
   AllowedOrigins=https://app.example.com

   # ❌ Inseguro
   AllowedOrigins=http://app.example.com
   ```

3. **Minimize headers permitidos**
   ```ini
   # Apenas o necessário
   AllowedHeaders=Origin,Content-Type,Authorization
   ```

4. **Valide origens com regex** (se necessário)
   - Atualmente não suportado, mas pode ser implementado

5. **Monitore logs de CORS**
   - Verifique tentativas de acesso não autorizadas

## Compatibilidade com WinDev

O header `Grant_type` (com underscore) é mantido para compatibilidade com o sistema WinDev original:

```ini
AllowedHeaders=Origin,Content-Type,Authorization,Grant_type,X-CSRF-Token
```

Este header é usado na geração de tokens JWT no endpoint `/connect/v1/token`.

## Referências

- [MDN - CORS](https://developer.mozilla.org/pt-BR/docs/Web/HTTP/CORS)
- [W3C - CORS Specification](https://www.w3.org/TR/cors/)
- [Enable CORS](https://enable-cors.org/)

## Suporte

Para problemas relacionados a CORS:

1. Verifique os logs de inicialização
2. Teste com cURL para validar headers
3. Use o console do navegador para ver erros específicos
4. Consulte a seção de problemas comuns acima
