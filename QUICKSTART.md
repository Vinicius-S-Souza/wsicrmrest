# 🚀 Quick Start - WSICRMREST

Guia rápido para colocar o webservice REST em funcionamento.

## ⚡ Início Rápido (5 minutos)

### 1. Pré-requisitos

- ✅ Go 1.21+ instalado
- ✅ Oracle Client instalado
- ✅ Variável `LD_LIBRARY_PATH` configurada
- ✅ Arquivo `tnsnames.ora` configurado
- ✅ Acesso ao banco de dados Oracle

### 2. Setup Automático

```bash
# Clone ou navegue até o diretório
cd /home/vinicius/projetos/wsicrmrest

# Execute o setup
make setup
```

### 3. Configurar Credenciais

Edite o arquivo `dbinit.ini`:

```ini
[database]
driver = 2
tns_name = SEU_TNS_NAME
username = seu_usuario
password = sua_senha

[application]
version = Ver 1.26.4.27
version_date = 2025-10-16T11:55:00
environment = production
port = 8080
```

**📖 Documentação completa das variáveis:** [docs/GLOBAL_VARIABLES.md](docs/GLOBAL_VARIABLES.md)

**⚠️ IMPORTANTE sobre Configurações:**
- **Credenciais JWT** (gsKey, gsIss, etc): Valores **fixos no código**, NÃO configuráveis
- **Dados da Organização**: Carregados da tabela `ORGANIZADOR`, NÃO do dbinit.ini
- **Database e Application**: Configuráveis via dbinit.ini

### 4. Criar Tabelas no Oracle (OBRIGATÓRIO)

**⚠️ A tabela ORGANIZADOR é OBRIGATÓRIA para o sistema funcionar!**

Execute o script SQL em `docs/DATABASE_SCHEMA.md` no seu banco Oracle:

```bash
sqlplus seu_usuario/sua_senha@tns_name @docs/create_tables.sql
```

Ou copie e execute manualmente os comandos CREATE TABLE do arquivo.

### 5. Inserir Dados Obrigatórios

**Primeiro, insira o ORGANIZADOR (OBRIGATÓRIO):**

```sql
-- ORGANIZADOR (obrigatório - sistema não inicia sem isso)
INSERT INTO ORGANIZADOR (
    ORGCODIGO,
    ORGNOME,
    ORGCNPJ,
    ORGCODLOJAMATRIZ,
    ORGCODISGA
) VALUES (
    1,
    'Minha Empresa',
    '12345678000199',
    1,
    1001
);
COMMIT;
```

**Depois, insira uma aplicação de teste:**

```sql
-- Aplicação de teste
INSERT INTO WSAPLICACOES (
    WSAPLCLIENTID,
    WSAPLCLIENTSECRET,
    WSAPLIJWTEXPIRACAO,
    WSAPLSCOPO,
    WSAPLSTATUS,
    WSAPLNOME
) VALUES (
    'CLIENTE1234567890',
    'a1234567890b',
    86400,  -- 24 horas
    1,      -- Escopo: clientes
    1,      -- Status: ativo
    'Aplicacao Teste'
);

COMMIT;
```

### 6. Executar o Servidor

```bash
# Compilar e executar
make dev

# OU executar separadamente
make build
make run
```

O servidor iniciará na porta **8080** (padrão).

### 7. Testar as APIs

#### Opção A: Script Automático

```bash
# Edite o script se necessário (CLIENT_ID e CLIENT_SECRET)
make test-api
```

#### Opção B: cURL Manual

**Teste de conexão:**
```bash
curl http://localhost:8080/connect/v1/wsteste
```

**Gerar token:**
```bash
# Criar Basic Auth
AUTH=$(echo -n "CLIENTE1234567890:a1234567890b" | base64)

# Requisitar token
curl -X GET http://localhost:8080/connect/v1/token \
  -H "Authorization: Basic $AUTH" \
  -H "Grant_type: client_credentials"
```

---

## 📋 Comandos Make Disponíveis

```bash
make help           # Mostra todos os comandos disponíveis
make setup          # Configuração inicial do projeto
make deps           # Baixa dependências
make build          # Compila o projeto
make run            # Executa o servidor
make dev            # Compila e executa
make test           # Executa testes
make test-api       # Testa as APIs
make clean          # Remove arquivos de build
make fmt            # Formata o código
make vet            # Verifica o código
make check          # Formata e verifica
```

---

## 🔍 Verificar se Está Funcionando

### 1. Logs

Os logs são gravados em:
```
log/wsicrmrest_YYYY-MM-DD.log
```

Verifique se há erros:
```bash
tail -f log/wsicrmrest_$(date +%Y-%m-%d).log
```

### 2. Health Check

```bash
# Deve retornar informações da organização
curl http://localhost:8080/connect/v1/wsteste | jq
```

Resposta esperada:
```json
{
  "code": "000",
  "organizadorCodigo": 1,
  "organizadorNome": "Minha Empresa",
  "organizadorCnpj": "12345678000199",
  "organizadorLojaMatriz": 1,
  "organizadorCodIsga": 1001,
  "versao": "1.0.0",
  "versaoData": "2025-01-27"
}
```

### 3. Verificar Logs no Banco

```sql
-- Últimas requisições
SELECT
    WSREQDTARECEBE,
    WSREQENDPOINT,
    WSREQMETODO,
    WSREQCODRESPOSTA,
    WSREQDURACAO
FROM WSREQUISICOES
ORDER BY WSREQDTARECEBE DESC
FETCH FIRST 10 ROWS ONLY;

-- Tokens gerados
SELECT
    WSAPLCLIENTID,
    WSLTKDATA,
    WSLTKEXPIRACAO
FROM WSAPLLOGTOKEN
ORDER BY WSLTKDATA DESC
FETCH FIRST 10 ROWS ONLY;
```

---

## 🐛 Problemas Comuns

### Erro: "cannot connect to database"

**Solução:**
```bash
# Verificar se Oracle Client está instalado
echo $LD_LIBRARY_PATH

# Testar conexão manualmente
sqlplus usuario/senha@tns_name

# Verificar tnsnames.ora
cat $ORACLE_HOME/network/admin/tnsnames.ora
```

### Erro: "table or view does not exist" ou "Organizador Não Cadastrado"

**Solução:**
- Execute os scripts SQL de criação das tabelas (ver `docs/DATABASE_SCHEMA.md`)
- **OBRIGATÓRIO:** Insira pelo menos um registro na tabela ORGANIZADOR
- Verifique se o usuário tem permissões adequadas

```sql
-- Verificar se existe organizador
SELECT * FROM ORGANIZADOR WHERE ORGCODIGO > 0;

-- Se não existir, inserir:
INSERT INTO ORGANIZADOR (ORGCODIGO, ORGNOME, ORGCNPJ, ORGCODLOJAMATRIZ, ORGCODISGA)
VALUES (1, 'Minha Empresa', '12345678000199', 1, 1001);
COMMIT;
```

### Erro: "401 Unauthorized"

**Solução:**
- Verifique se o CLIENT_ID existe na tabela WSAPLICACOES
- Confirme se o CLIENT_SECRET está correto
- Verifique se WSAPLSTATUS = 1 (ativo)

### Servidor não inicia

**Solução:**
```bash
# Verificar se arquivo dbinit.ini existe
ls -la dbinit.ini

# Verificar porta ocupada
netstat -tlnp | grep 8080

# Ver logs detalhados
./wsicrmrest 2>&1 | tee server.log
```

---

## 📚 Próximos Passos

1. ✅ Servidor funcionando
2. 📖 Ler `README.md` completo
3. 🗄️ Consultar `docs/DATABASE_SCHEMA.md` para estrutura das tabelas
4. 🔐 Adicionar mais aplicações na tabela WSAPLICACOES
5. 🚀 Implementar novas APIs

---

## 📞 Ajuda

- **README completo:** [README.md](README.md)
- **Estrutura das tabelas:** [docs/DATABASE_SCHEMA.md](docs/DATABASE_SCHEMA.md)
- **Scripts de teste:** [scripts/test_apis.sh](scripts/test_apis.sh)

---

## ✅ Checklist de Validação

Antes de ir para produção, verifique:

- [ ] Banco de dados Oracle configurado e acessível
- [ ] Tabelas criadas (WSAPLICACOES, WSAPLLOGTOKEN, WSREQUISICOES)
- [ ] Aplicações registradas em WSAPLICACOES
- [ ] Chave JWT forte configurada (mínimo 32 caracteres)
- [ ] Logs sendo gravados corretamente
- [ ] Testes de API executados com sucesso
- [ ] Permissões de arquivos adequadas
- [ ] Backup do dbinit.ini configurado
- [ ] Monitoramento configurado
- [ ] Rotação de logs configurada

**Pronto! Seu webservice REST está funcionando! 🎉**
