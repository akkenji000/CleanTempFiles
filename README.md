# CleanTempFiles

# 🧽 Faxineira Japonesa

Um utilitário de limpeza leve e eficiente para Windows, construído em **Go (Golang)**. O foco é automatizar a limpeza de arquivos temporários que desenvolvedores e gamers costumam acumular.

## ✨ Funcionalidades
- **Desenvolvimento:** Limpa caches de pacotes `NPM` e `NuGet`.
- **Gaming:** Limpa logs do `League of Legends`.
- **Sistema:** Flush DNS, limpeza de `Temp`, `Prefetch` e esvaziamento da Lixeira.
- **Interface Gráfica:** Desenvolvida com a biblioteca `walk`, garantindo um visual nativo do Windows.
- **Performance:** Uso de Goroutines para garantir que a UI não trave durante a exclusão dos arquivos.

## 🛠️ Como executar
Se você é apenas um usuário, baixe o executável na aba [Releases](link-da-sua-release).

Se você é desenvolvedor:
1. Tenha o Go instalado.
2. Gere o arquivo de recursos: `go run github.com/akavel/rsrc@latest -manifest app.manifest -o rsrc.syso`
3. Compile: `go build -ldflags="-H=windowsgui -s -w" -o FaxineiraJaponesa.exe`

## 👨‍💻 Autor
Desenvolvido por **Kenji**.
