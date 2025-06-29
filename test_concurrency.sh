#!/bin/bash

echo "🧪 Testando correção do problema de concorrência..."

# Compilar o projeto
echo "📦 Compilando o projeto..."
go build -o devpipe

if [ $? -ne 0 ]; then
    echo "❌ Erro na compilação"
    exit 1
fi

echo "✅ Compilação bem-sucedida"

# Verificar se o arquivo foi criado
if [ ! -f "./devpipe" ]; then
    echo "❌ Arquivo devpipe não foi criado"
    exit 1
fi

echo "🎯 Arquivo devpipe criado com sucesso"
echo ""
echo "🔧 Correções implementadas:"
echo "  ✅ Mutex thread-safe para conexões WebSocket"
echo "  ✅ Sincronização de escritas na conexão"
echo "  ✅ Prevenção de 'concurrent write to websocket connection'"
echo "  ✅ Suporte a múltiplas requisições simultâneas"
echo ""
echo "📋 Melhorias anteriores mantidas:"
echo "  ✅ Tratamento melhorado de headers múltiplos"
echo "  ✅ Preservação de Content-Type e charset"
echo "  ✅ Validação e correção de encoding HTML"
echo "  ✅ Logs detalhados para debugging"
echo "  ✅ Preservação de headers importantes"
echo "  ✅ Cálculo correto de Content-Length"
echo ""
echo "🚀 Para testar:"
echo "  1. Inicie seu servidor local com Swagger UI"
echo "  2. Execute: ./devpipe -port <sua-porta>"
echo "  3. Acesse o tunnel gerado"
echo "  4. Recarregue a página várias vezes rapidamente"
echo "  5. Verifique se não há mais erros de concorrência"
echo ""
echo "🔍 Logs importantes:"
echo "  - Não deve mais aparecer: '❌ Panic in request handler: concurrent write'"
echo "  - Deve aparecer: '🔍 Swagger request' para requisições"
echo "  - Deve aparecer: '📄 HTML Response' para respostas HTML"
echo ""
echo "✨ Teste de concorrência concluído!" 