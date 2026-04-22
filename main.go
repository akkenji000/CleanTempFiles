package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var mw *walk.MainWindow
	var rbCompleta, rbPersonalizada *walk.RadioButton
	var cbDNS, cbDev, cbLogs, cbTemp, cbLixo *walk.CheckBox

	// Função para habilitar/desabilitar checkboxes conforme o modo
	togglePersonalizado := func() {
		isCustom := rbPersonalizada.Checked()
		cbDNS.SetEnabled(isCustom)
		cbDev.SetEnabled(isCustom)
		cbLogs.SetEnabled(isCustom)
		cbTemp.SetEnabled(isCustom)
		cbLixo.SetEnabled(isCustom)
	}

	executarLimpeza := func() {
		isCompleta := rbCompleta.Checked()

		// 1. Execução em Goroutine para não travar a interface
		go func() {
			var totalBytes int64

			if isCompleta || cbDNS.Checked() {
				flushDNS()
			}

			if isCompleta || cbDev.Checked() {
				userProfile := os.Getenv("USERPROFILE")
				totalBytes += cleanAndMeasure(filepath.Join(userProfile, "AppData", "Local", "npm-cache"))
				totalBytes += cleanAndMeasure(filepath.Join(userProfile, ".nuget", "packages"))
			}

			if isCompleta || cbLogs.Checked() {
				// Caminho padrão do LoL
				totalBytes += cleanAndMeasure(`C:\Riot Games\League of Legends\Logs`)
			}

			if isCompleta || cbTemp.Checked() {
				systemRoot := os.Getenv("SystemRoot")
				totalBytes += cleanAndMeasure(os.TempDir())
				totalBytes += cleanAndMeasure(filepath.Join(systemRoot, "Temp"))
				totalBytes += cleanAndMeasure(filepath.Join(systemRoot, "Prefetch"))
			}

			if isCompleta || cbLixo.Checked() {
				emptyRecycleBin()
			}

			mb := float64(totalBytes) / 1024 / 1024

			// 2. Sincroniza com a thread da UI para exibir o resultado
			mw.Synchronize(func() {
				walk.MsgBox(mw, "Sucesso",
					fmt.Sprintf("Faxina concluída!\n\nEspaço liberado: %.2f MB\nLixeira e DNS processados.", mb),
					walk.MsgBoxIconInformation)
			})
		}()
	}

	// Construção da Janela
	err := MainWindow{
		AssignTo: &mw,
		Title:    "Faxineira Japonesa",
		MinSize:  Size{Width: 500, Height: 300},
		Size:     Size{Width: 550, Height: 350},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Escolha o modo de manutenção:", Font: Font{Bold: true}},

			GroupBox{
				Title:  "Modo de Operação",
				Layout: HBox{},
				Children: []Widget{
					RadioButton{
						AssignTo:  &rbCompleta,
						Text:      "Limpeza Completa",
						OnClicked: togglePersonalizado,
					},
					RadioButton{
						AssignTo:  &rbPersonalizada,
						Text:      "Personalizada",
						OnClicked: togglePersonalizado,
					},
				},
			},

			GroupBox{
				Title:  "Módulos do Sistema",
				Layout: VBox{},
				Children: []Widget{
					CheckBox{AssignTo: &cbDNS, Text: "Limpar Cache de DNS (Rede)", Enabled: false},
					CheckBox{AssignTo: &cbDev, Text: "Limpar Caches de Dev (NPM, NuGet)", Enabled: false},
					CheckBox{AssignTo: &cbLogs, Text: "Limpar Logs de Jogos (League of Legends)", Enabled: false},
					CheckBox{AssignTo: &cbTemp, Text: "Limpar Arquivos Temporários", Enabled: false},
					CheckBox{AssignTo: &cbLixo, Text: "Esvaziar Lixeira", Enabled: false},
				},
			},

			VSpacer{},
			PushButton{
				Text:      "Executar Limpeza Agora",
				OnClicked: executarLimpeza,
			},
		},
	}.Create()

	if err != nil {
		panic(err)
	}

	// Configura o estado inicial antes de exibir a janela
	rbCompleta.SetChecked(true)
	togglePersonalizado()

	mw.Run()
}

// --- FUNÇÕES DE SUPORTE ---

func flushDNS() {
	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Run()
}

// Esvazia a lixeira via API oficial (atualiza o ícone e evita erros de CMD)
func emptyRecycleBin() {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	emptyBin := shell32.NewProc("SHEmptyRecycleBinW")
	// 7 = SHERB_NOCONFIRMATION | SHERB_NOPROGRESSUI | SHERB_NOSOUND
	emptyBin.Call(0, 0, 7)
}

// Soma o espaço apenas se deletar, e limpa pastas vazias
func cleanAndMeasure(dir string) int64 {
	var total int64

	// Primeira passada: Apaga arquivos e soma tamanho
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size := info.Size()
			// Tenta remover; se conseguir, adiciona ao total
			if err := os.Remove(path); err == nil {
				total += size
			}
		}
		return nil
	})

	// Segunda passada: Tenta remover subpastas que ficaram vazias
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && path != dir {
			// os.Remove em pasta só funciona se ela estiver vazia
			os.Remove(path)
		}
		return nil
	})

	return total
}
