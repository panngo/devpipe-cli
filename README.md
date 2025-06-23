# DevPipe

DevPipe é uma ferramenta de túnel reverso que permite expor aplicações locais à internet de forma segura e simples, utilizando um domínio customizado como `*.devpipe.cloud`.

## ✨ Funcionalidades

- Exponha rapidamente qualquer porta local pela internet
- Subdomínio dedicado para cada sessão (ex: `abc123-3000.devpipe.cloud`)
- Integração com HTTPS
- Binários disponíveis para múltiplas plataformas

## 🚀 Começando

### Pré-requisitos

- Go 1.22+
- Docker (opcional para execução em container)
- Conta no [devpipe.cloud](https://devpipe.cloud) (futuro)

### Instalação via binário

```bash
curl -sL https://devpipe.cloud/install.sh | bash

Compilação manual

git clone https://github.com/panngo/devpipe-cli.git
cd devpipe-cli
go build -o devpipe src/main.go

🧪 Exemplo de uso

devpipe --port 3000

Acesse então:

https://<uuid>-3000.devpipe.cloud

🐳 Executando via Docker

docker run --rm -p 3000:3000 devpipe/devpipe-server

🛠 Arquitetura
	•	Cliente (devpipe-cli): Inicia conexão WebSocket com o servidor e encaminha as requisições
	•	Servidor (devpipe-server): Proxy reverso que gerencia conexões e subdomínios via Traefik

📦 Publicação

Imagens Docker disponíveis em:

ghcr.io/seu-usuario/devpipe-server

📄 Licença

Este projeto está licenciado sob a MIT License.

⸻

Desenvolvido com ❤️ por @panngo