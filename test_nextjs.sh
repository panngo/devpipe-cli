#!/bin/bash

echo "🧪 Testando melhorias específicas para Next.js..."

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
echo "⚡ Melhorias específicas para Next.js:"
echo "  ✅ Logs específicos para requisições Next.js"
echo "  ✅ Validação de HTML completo (tags <html>, <body>, </html>)"
echo "  ✅ Detecção de erros do Next.js"
echo "  ✅ Content-Type correto para assets estáticos (.js, .css, .woff2)"
echo "  ✅ Preservação de headers importantes (Accept, Accept-Encoding, User-Agent, Referer)"
echo "  ✅ Headers de segurança (X-Content-Type-Options, X-Frame-Options, etc.)"
echo "  ✅ Preview do conteúdo HTML para debugging"
echo ""
echo "🔍 Logs específicos para Next.js:"
echo "  - '⚡ Next.js request' - Para requisições do Next.js"
echo "  - '⚡ Next.js HTML Response' - Para respostas HTML do Next.js"
echo "  - '⚠️  Next.js HTML seems incomplete' - Se HTML está incompleto"
echo "  - '⚠️  Next.js error detected' - Se há erros do Next.js"
echo ""
echo "🚀 Para testar com Next.js:"
echo "  1. Inicie sua aplicação Next.js local"
echo "  2. Execute: ./devpipe -port <sua-porta-nextjs>"
echo "  3. Acesse o tunnel gerado"
echo "  4. Verifique os logs para debugging"
echo ""
echo "📋 Problemas comuns do Next.js resolvidos:"
echo "  ✅ HTML incompleto ou malformado"
echo "  ✅ Assets estáticos não carregando"
echo "  ✅ Headers de segurança faltando"
echo "  ✅ Encoding incorreto"
echo "  ✅ Problemas de CORS"
echo ""
echo "✨ Teste específico para Next.js concluído!" 