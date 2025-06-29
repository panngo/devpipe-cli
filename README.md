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

## 📄 Licença

Este projeto está licenciado sob a MIT License.

---

Desenvolvido com ❤️ por @panngo