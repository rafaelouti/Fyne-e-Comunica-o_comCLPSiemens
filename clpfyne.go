package main

import (
	"fmt"
	"log"
	"time"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/robinson/gos7"
)

func lerDadosCLP(client *gos7.Client, dbNumber, start, size int) float64 {
	var buffer = make([]byte, size)
	err := client.AGReadDB(dbNumber, start, size, buffer)
	if err != nil {
		log.Printf("Erro ao ler dados do CLP: %v", err)
		return 0.0
	}
	return gos7.GetFloatAt(buffer, 0)
}

func main() {
	// Configuração de conexão com o CLP
	handler := gos7.NewTCPClientHandler("192.168.0.1", 0, 1)
	handler.Timeout = 5 * time.Second
	handler.IdleTimeout = 5 * time.Second

	// Conecta ao CLP
	err := handler.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar com o CLP: %v", err)
	}
	defer handler.Close()
	client := gos7.NewClient(handler)

	// Criação da interface gráfica
	myApp := app.New()
	myWindow := myApp.NewWindow("Monitoramento CLP Siemens")

	// Botão e labels para exibir dados do CLP
	rpmLabel := widget.NewLabel("RPM: Aguardando dados...")
	tempLabel := widget.NewLabel("Temperatura: Aguardando dados...")
	vibLabel := widget.NewLabel("Vibração: Aguardando dados...")

	// Botão de atualização
	updateButton := widget.NewButton("Atualizar Dados", func() {
		rpm := lerDadosCLP(client, 1, 0, 4)
		temp := lerDadosCLP(client, 1, 4, 4)
		vib := lerDadosCLP(client, 1, 8, 4)
		rpmLabel.SetText(fmt.Sprintf("RPM: %.2f", rpm))
		tempLabel.SetText(fmt.Sprintf("Temperatura: %.2f °C", temp))
		vibLabel.SetText(fmt.Sprintf("Vibração: %.2f G", vib))
		dialog.ShowInformation("Dados Atualizados", "Dados do CLP atualizados com sucesso!", myWindow)
	})

	// Layout da janela
	content := container.NewVBox(
		widget.NewLabel("Monitor de Motor - CLP Siemens"),
		updateButton,
		rpmLabel,
		tempLabel,
		vibLabel,
	)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
