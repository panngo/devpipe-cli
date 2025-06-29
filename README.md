# DevPipe

DevPipe é uma ferramenta de túnel reverso que permite expor aplicações locais à internet de forma segura e simples, utilizando um domínio customizado como `*.devpipe.cloud`.

## ✨ Funcionalidades

- Exponha rapidamente qualquer porta local pela internet
- Subdomínio dedicado para cada sessão (ex: `abc123-3000.devpipe.cloud`)
- Integração com HTTPS
- Binários disponíveis para múltiplas plataformas
- **🆕 Sistema de Reconexão Segura**: UUID persistente + chave de segurança para reconexões seguras
- **🆕 Configuração Persistente**: Salva automaticamente UUID e chave de segurança
- **🆕 Limpeza Automática**: Remove configurações inválidas automaticamente
- **🆕 Proxy Transparente**: Repassa todos os headers e dados exatamente como recebidos
- **🆕 Suporte Completo a Métodos HTTP**: Todos os métodos padrão que navegadores aceitam
- **🆕 Suporte a CORS**: Headers CORS automáticos para requisições cross-origin
- **Túneis persistentes**: O endereço do túnel é mantido mesmo após reconexões
- **Reconexão automática**: Reconecta automaticamente quando a conexão é perdida
- **Heartbeat**: Detecta conexões perdidas rapidamente através de pings regulares
- **Tratamento robusto de erros**: Não falha quando há erros individuais nas requisições
- **Logs informativos**: Feedback claro sobre o status da conexão e reconexões

## 🌐 Métodos HTTP Suportados

O DevPipe agora suporta **todos os métodos HTTP padrão** que navegadores podem enviar:

### Métodos Principais
- **GET**: Requisições de leitura (sem corpo)
- **POST**: Envio de dados (com corpo)
- **PUT**: Atualização completa de recursos (com corpo)
- **DELETE**: Remoção de recursos (pode ter corpo)
- **PATCH**: Atualização parcial de recursos (com corpo)

### Métodos Especiais
- **HEAD**: Requisições sem corpo (apenas headers)
- **OPTIONS**: Requisições CORS preflight (suporte automático)
- **TRACE**: Debugging de requisições
- **CONNECT**: Túneis HTTP

### Tratamento Específico por Método

#### GET, HEAD, OPTIONS, TRACE
- Sem corpo de requisição
- HEAD retorna apenas headers (Content-Length: 0)
- OPTIONS inclui headers CORS automáticos

#### POST, PUT, PATCH
- Com corpo de requisição
- Content-Length calculado automaticamente
- Headers de Content-Type preservados

#### DELETE
- Pode ter corpo (opcional)
- Tratamento flexível baseado na presença de dados

### Validação e Erros
- **Métodos não suportados**: Retorna 405 Method Not Allowed
- **Caminhos vazios**: Retorna 400 Bad Request
- **Headers inválidos**: Tratamento automático e correção

## 🔐 Sistema de Reconexão Segura

O DevPipe agora suporta reconexão segura com UUID persistente e chave de segurança, permitindo que você mantenha a mesma URL mesmo quando a conexão cair, com autenticação adicional para maior segurança.

### Como Funciona

1. **Primeira Conexão**: Servidor gera UUID único + chave de segurança
2. **Reconexão**: Cliente deve fornecer UUID + chave de segurança válida
3. **Validação**: Servidor verifica a chave antes de autorizar reconexão
4. **Limpeza Automática**: Túneis inativos são removidos após 1 hora

### Primeira Conexão

Quando você se conecta pela primeira vez, o servidor gera um UUID único e uma chave de segurança:

```bash
./devpipe -port 3000
```

O servidor responde com UUID, chave de segurança e ID completo do túnel:
```json
{
  "tunnel": "abc123-def456-789-3000",
  "uuid": "abc123-def456-789",
  "key": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"
}
```

### Reconexão Segura

Para reconectar com o mesmo UUID, o cliente automaticamente fornece a chave de segurança:

```bash
./devpipe -port 3000  # Reconecta automaticamente com UUID salvo
```

O servidor validará a chave e manterá a mesma URL: `abc123-def456-789-3000.devpipe.cloud`

### Gerenciamento de Configuração

```bash
# Limpar configuração salva (força nova conexão)
./devpipe -clear-config

# Verificar configuração salva
cat ~/.devpipe/tunnel.json
```

## 🚀 Começando

### Pré-requisitos

- Go 1.22+
- Docker (opcional para execução em container)
- Conta no [devpipe.cloud](https://devpipe.cloud) (futuro)

### Instalação via binário

```bash
curl -sL https://devpipe.cloud/install.sh | bash
```

### Compilação manual

```bash
git clone https://github.com/panngo/devpipe-cli.git
cd devpipe-cli
go build -o devpipe
```

## 🧪 Exemplo de uso

```bash
# Primeira conexão (gera novo UUID)
./devpipe -port 3000

# Reconexão (usa UUID salvo)
./devpipe -port 3000

# Limpar configuração e forçar nova conexão
./devpipe -clear-config -port 3000
```

Acesse então:
```
https://<uuid>-3000.devpipe.cloud
```

## 🧪 Testando

Execute os scripts de teste para verificar todas as funcionalidades:

```bash
# Testar reconexão segura
./test_reconnection.sh

# Testar métodos HTTP
./test_http_methods.sh

# Testar todas as funcionalidades
make test-all
```

### Testes Disponíveis
- ✅ **Reconexão Segura**: UUID e chave de segurança
- ✅ **Métodos HTTP**: Todos os métodos padrão
- ✅ **CORS**: Headers automáticos para cross-origin
- ✅ **Next.js**: Otimizações específicas
- ✅ **Swagger**: Suporte completo
- ✅ **Concorrência**: Múltiplas requisições simultâneas

## 🐳 Executando via Docker

```bash
docker run --rm -p 3000:3000 devpipe/devpipe-server
```

## 🛠 Arquitetura

- **Cliente (devpipe-cli)**: Inicia conexão WebSocket com o servidor e encaminha as requisições
- **Servidor (devpipe-server)**: Proxy reverso que gerencia conexões e subdomínios via Traefik
- **🆕 Sistema de Configuração**: Gerencia UUID e chave de segurança persistentes
- **🆕 Validação HTTP**: Suporte completo a todos os métodos HTTP

## 📦 Publicação

Imagens Docker disponíveis em:
```
ghcr.io/seu-usuario/devpipe-server
```

## 📄 Licença

Este projeto está licenciado sob a MIT License.

---

Desenvolvido com ❤️ por @panngo

# DevPipe CLI

Um cliente CLI para o DevPipe que permite expor serviços locais através de túneis seguros.

## Funcionalidades

- **🆕 Sistema de Reconexão Segura**: UUID persistente + chave de segurança para reconexões seguras
- **🆕 Configuração Persistente**: Salva automaticamente UUID e chave de segurança
- **🆕 Limpeza Automática**: Remove configurações inválidas automaticamente
- **🆕 Suporte Completo a Métodos HTTP**: Todos os métodos padrão que navegadores aceitam
- **🆕 Suporte a CORS**: Headers CORS automáticos para requisições cross-origin
- **🆕 Otimização para Swagger**: Carregamento inicial melhorado com headers adequados
- **Túneis persistentes**: O endereço do túnel é mantido mesmo após reconexões
- **Reconexão automática**: Reconecta automaticamente quando a conexão é perdida
- **Heartbeat**: Detecta conexões perdidas rapidamente através de pings regulares
- **Tratamento robusto de erros**: Não falha quando há erros individuais nas requisições
- **Logs informativos**: Feedback claro sobre o status da conexão e reconexões

## Instalação

```bash
# Clone o repositório
git clone https://github.com/panngo/devpipe-cli.git
cd devpipe-cli

# Compile o projeto
go build -o devpipe

# Ou use o script de instalação
./install.sh
```

## Uso

```bash
# Expor a porta padrão (3000)
./devpipe

# Expor uma porta específica
./devpipe -port 8080

# Limpar configuração salva
./devpipe -clear-config

# Limpar configuração e expor porta específica
./devpipe -clear-config -port 8080
```

## 🆕 Melhorias de Segurança

### Sistema de Reconexão Segura
- **UUID Persistente**: Cada túnel tem um UUID único que é mantido entre reconexões
- **Chave de Segurança**: Chave de 32 bytes para autenticar reconexões
- **Validação Automática**: Servidor valida a chave antes de autorizar reconexão
- **Configuração Persistente**: UUID e chave são salvos automaticamente em `~/.devpipe/tunnel.json`
- **Limpeza Automática**: Configurações inválidas são removidas automaticamente

### Gerenciamento de Configuração
- **Salvamento Automático**: Configuração é salva após primeira conexão bem-sucedida
- **Carregamento Automático**: Configuração é carregada automaticamente em reconexões
- **Limpeza Manual**: Use `-clear-config` para remover configuração salva
- **Fallback Seguro**: Se a configuração for inválida, cria nova conexão automaticamente

### Logs de Segurança
- **🔐 Secure Reconnection**: Indica quando está usando reconexão segura
- **🔑 UUID**: Exibe o UUID da conexão atual
- **💾 Configuration Saved**: Confirma quando a configuração foi salva
- **🗑️ Configuration Cleared**: Confirma quando a configuração foi removida

## 🆕 Suporte Completo a Métodos HTTP

### Métodos Suportados
O cliente agora suporta **todos os métodos HTTP padrão**:

- **GET**: Requisições de leitura (sem corpo)
- **POST**: Envio de dados (com corpo)
- **PUT**: Atualização completa de recursos (com corpo)
- **DELETE**: Remoção de recursos (pode ter corpo)
- **PATCH**: Atualização parcial de recursos (com corpo)
- **HEAD**: Requisições sem corpo (apenas headers)
- **OPTIONS**: Requisições CORS preflight (suporte automático)
- **TRACE**: Debugging de requisições
- **CONNECT**: Túneis HTTP

### Tratamento Específico

#### Métodos sem Corpo (GET, HEAD, OPTIONS, TRACE)
- Requisições criadas sem corpo
- HEAD retorna apenas headers (Content-Length: 0)
- OPTIONS inclui headers CORS automáticos

#### Métodos com Corpo (POST, PUT, PATCH, DELETE)
- Corpo incluído quando presente
- Content-Length calculado automaticamente
- Headers de Content-Type preservados

### Validação e Erros
- **Métodos não suportados**: Retorna 405 Method Not Allowed
- **Caminhos vazios**: Retorna 400 Bad Request
- **Headers inválidos**: Tratamento automático e correção

### Suporte a CORS
- **Headers automáticos**: Access-Control-Allow-* headers para OPTIONS
- **Métodos permitidos**: Todos os métodos HTTP suportados
- **Headers permitidos**: Content-Type, Authorization, etc.
- **Cache**: 24 horas para requisições preflight


### Reconexão Automática
- Quando a conexão WebSocket é perdida, o cliente tenta reconectar automaticamente
- Até 5 tentativas de reconexão com delay exponencial
- Mantém o mesmo tunnel ID quando possível usando UUID e chave de segurança

### Heartbeat
- Envia pings a cada 30 segundos para manter a conexão ativa
- Detecta conexões perdidas em até 35 segundos
- Reinicia automaticamente após reconexão

### Tratamento de Erros
- Erros individuais nas requisições não afetam a conexão principal
- Respostas de erro apropriadas para requisições malformadas
- Logs detalhados para debugging


```
devpipe-cli/
├── client/          # Lógica do cliente
├── ws/             # Gerenciamento de conexões WebSocket
├── ui/             # Interface do usuário
├── config/         # 🆕 Gerenciamento de configuração persistente
├── main.go         # Ponto de entrada
└── README.md       # Documentação
```

## Protocolo de Comunicação

### Registro Inicial
```json
{
  "action": "register",
  "port": "3000"
}
```

### Registro com Reconexão Segura
```json
{
  "action": "register",
  "port": "3000",
  "uuid": "abc123-def456-789",
  "key": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"
}
```

### Resposta do Servidor
```json
{
  "tunnel": "abc123-def456-789-3000",
  "uuid": "abc123-def456-789",
  "key": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"
}
```

### Heartbeat
```json
{
  "action": "ping"
}
```

## Logs

O cliente fornece logs informativos sobre:
- Status da conexão inicial
- Tentativas de reconexão segura
- Sucesso/falha na reconexão
- Mudanças no tunnel ID
- Erros de requisição
- Heartbeat status
- **🆕 Status de segurança**: UUID, chave de segurança, configuração
- **🆕 Métodos HTTP**: Logs específicos para cada método
- **🆕 CORS**: Logs de requisições preflight

## Troubleshooting

### Problemas com Reconexão Segura
Se a reconexão segura falhar:

1. **Limpe a configuração**: `./devpipe -clear-config`
2. **Verifique os logs**: Procure por mensagens de erro de autenticação
3. **Nova conexão**: O cliente criará automaticamente uma nova conexão
4. **Verifique o arquivo de configuração**: `cat ~/.devpipe/tunnel.json`

### Problemas com Métodos HTTP
Se algum método HTTP não funcionar:

1. **Verifique os logs**: Procure por "❌ Unsupported HTTP method"
2. **Teste com curl**: Use o script `test_http_methods.sh`
3. **Verifique o servidor**: Certifique-se de que o servidor local suporta o método
4. **CORS**: Para requisições cross-origin, verifique os headers CORS

### Problemas com HTML/Swagger
Se o HTML do Swagger ou outras páginas web aparecem estranhas:

1. **Verifique os logs**: Procure por logs com `🔍 Swagger request` ou `📄 HTML Response`
2. **Encoding**: O cliente agora adiciona automaticamente `charset=utf-8` se necessário
3. **Headers**: Todos os headers importantes são preservados
4. **Cache**: Headers de cache são mantidos para melhor performance

### Logs de Debugging
```bash
# Para ver logs detalhados de requisições Swagger
./devpipe -port 3000 2>&1 | grep -E "(🔍|📄|🔧)"

# Para ver logs de segurança
./devpipe -port 3000 2>&1 | grep -E "(🔐|🔑|💾|🗑️)"

# Para ver logs de métodos HTTP
./devpipe -port 3000 2>&1 | grep -E "(🌐|📦)"

# Para ver todos os logs
./devpipe -port 3000
```

## 🆕 Modo Proxy Transparente

### Como Funciona
O DevPipe agora funciona como um **proxy transparente**, repassando todos os headers e dados exatamente como recebidos, sem modificações ou interpretações:

```go
// PROXY MODE: Copy ALL headers exactly as received (transparent proxy)
for k, v := range req.Headers {
    // Only skip headers that Go's HTTP client manages automatically
    if strings.ToLower(k) == "host" {
        // Skip Host header, it will be defined automatically by Go
        continue
    }
    if strings.ToLower(k) == "connection" {
        // Skip Connection header, it will be managed by Go
        continue
    }
    if strings.ToLower(k) == "transfer-encoding" {
        // Skip Transfer-Encoding, let Go handle it
        continue
    }
    // Copy all other headers exactly as received
    httpReq.Header.Set(k, v)
}
```

### Headers Preservados
- **Todos os headers de requisição** são repassados exatamente como recebidos
- **Todos os headers de resposta** são preservados sem modificação
- **Content-Length** é calculado automaticamente pelo Go
- **Host, Connection, Transfer-Encoding** são gerenciados pelo Go HTTP client

### Benefícios
1. **Compatibilidade Total**: Funciona com qualquer aplicação sem modificações
2. **Headers Preservados**: Todos os headers customizados são mantidos
3. **Dados Intactos**: Corpo das requisições e respostas não são modificados
4. **Performance**: Sem overhead de processamento de headers
5. **Simplicidade**: Comportamento previsível e transparente

### Casos de Uso
- **APIs REST**: Todos os headers de autenticação preservados
- **Aplicações Web**: Headers de sessão e cookies mantidos
- **Microserviços**: Headers de rastreamento e metadados preservados
- **APIs GraphQL**: Headers de autorização e contexto mantidos
- **WebSockets**: Headers de upgrade preservados
- **Uploads de Arquivo**: Headers de multipart mantidos

### Exemplo de Headers Preservados
```bash
# Headers de requisição preservados
Authorization: Bearer token123
X-API-Key: abc123def456
X-Request-ID: req-12345
X-Forwarded-For: 192.168.1.1
User-Agent: CustomApp/1.0
Accept: application/json
Content-Type: application/json

# Headers de resposta preservados
Content-Type: application/json; charset=utf-8
Cache-Control: no-cache
ETag: "abc123"
X-Powered-By: Express
Set-Cookie: session=abc123; HttpOnly
```