# DevPipe CLI GitHub Actions

## Pipelines Disponíveis

### 🚀 Robust Release (`robust-release.yml`)

**Pipeline principal para criar releases do DevPipe CLI.**

#### **Como usar:**
1. Vá para **Actions** no GitHub
2. Selecione **"Robust Release"**
3. Clique em **"Run workflow"**
4. Configure os parâmetros:
   - **Version**: Versão do release (ex: `v1.0.2`)
   - **Prerelease**: Se é um prerelease
   - **Draft**: Se deve criar como draft

#### **Funcionalidades:**
- ✅ **Build para todas as plataformas**: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64, arm64)
- ✅ **Upload automático** dos binários para o release
- ✅ **Verificação** de downloads após o release
- ✅ **Suporte a draft** e prerelease
- ✅ **Deleção automática** de releases existentes
- ✅ **Logs detalhados** para debug

#### **Parâmetros:**
| Parâmetro | Descrição | Padrão | Obrigatório |
|-----------|-----------|--------|-------------|
| `version` | Versão do release (ex: v1.0.2) | v1.0.2 | ✅ |
| `prerelease` | Se é um prerelease | false | ❌ |
| `draft` | Se deve criar como draft | false | ❌ |

#### **Exemplos de uso:**

**Release normal:**
```bash
# No GitHub Actions UI:
Version: v1.0.2
Prerelease: false
Draft: false
```

**Prerelease:**
```bash
# No GitHub Actions UI:
Version: v1.1.0-beta.1
Prerelease: true
Draft: false
```

**Draft para revisão:**
```bash
# No GitHub Actions UI:
Version: v1.0.3
Prerelease: false
Draft: true
```

#### **Binários gerados:**
- `devpipe-linux-amd64` - Linux x86_64
- `devpipe-linux-arm64` - Linux ARM64
- `devpipe-darwin-amd64` - macOS Intel
- `devpipe-darwin-arm64` - macOS Apple Silicon
- `devpipe-windows-amd64.exe` - Windows x86_64
- `devpipe-windows-arm64.exe` - Windows ARM64

#### **URLs de download:**
Após o release, os binários estarão disponíveis em:
```
https://github.com/panngo/devpipe-cli/releases/download/{version}/devpipe-{platform}
```

#### **Troubleshooting:**
- **Erro 403**: Verifique se o workflow tem permissão `contents: write`
- **Release já existe**: O workflow deleta automaticamente releases existentes
- **Build falha**: Verifique se todos os arquivos `.go` estão no mesmo package

---

## Histórico de Mudanças

### v1.0.2
- ✅ **Limpeza**: Removidas 4 pipelines redundantes
- ✅ **Consolidação**: Mantida apenas `robust-release.yml`
- ✅ **Melhorias**: Adicionado suporte a draft releases
- ✅ **Documentação**: README completo para a pipeline

### Pipelines Removidas (Redundantes):
- ❌ `quick-release.yml` - Funcionalidade incorporada no robust-release
- ❌ `release.yml` - Muito complexa, matrix desnecessária
- ❌ `simple-release.yml` - Muito similar ao robust-release
- ❌ `deploy.yml` - Só build, sem release 