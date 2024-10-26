Interface Gráfica com Fyne e Comunicação com CLP Siemens
Este projeto demonstra como criar uma interface gráfica utilizando o framework Fyne e integrar a comunicação com um CLP Siemens usando a biblioteca gos7 em Golang.

Requisitos
Go 1.16 ou superior

Biblioteca Fyne

Biblioteca gos7

Instalação
Clone o repositório:

sh

Copiar
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
Instale as dependências:

sh

Copiar
go get fyne.io/fyne/v2
go get github.com/robinson/gos7
Uso
Configure a conexão com o CLP:

No código, ajuste o endereço IP, rack e slot do CLP Siemens conforme necessário:

go

Copiar
handler := gos7.NewTCPClientHandler("192.168.0.1", 0, 1)
Execute a aplicação:

sh

Copiar
go run main.go
Interface Gráfica:

A interface gráfica exibirá os dados de RPM, temperatura e vibração do motor. Use o botão "Atualizar Dados" para ler os dados do CLP e atualizar a interface.

Código
go

Copiar
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
Explicação
Importação de Pacotes
fmt e log: Utilizados para formatar e registrar mensagens de texto e de erro.

time: Define o tempo de timeout para a conexão com o CLP.

fyne: Pacote para criar a interface gráfica, que inclui widgets e containers.

gos7: Biblioteca que permite conectar e ler dados de um CLP Siemens.

Função lerDadosCLP
Esta função executa a leitura dos dados do CLP.

Parâmetros:

client: O cliente que conecta com o CLP.

dbNumber: O número do Data Block (DB) no CLP.

start: A posição inicial do dado.

size: O tamanho do dado.

Processo:

Cria um buffer de bytes para armazenar o dado recebido.

Usa AGReadDB para ler o bloco de dados do CLP e armazena no buffer.

Caso não haja erros, converte o valor do buffer em float e o retorna.

Data Block (DB)
No contexto dos Controladores Lógicos Programáveis (CLPs) Siemens, um Data Block (DB) é uma área de memória usada para armazenar dados. Cada DB pode conter diferentes tipos de dados, como inteiros, floats, strings, entre outros. Os DBs são utilizados para armazenar variáveis que podem ser lidas e escritas durante a operação do CLP. No código acima, o dbNumber especifica qual Data Block está sendo acessado, enquanto start e size determinam a posição inicial e o tamanho dos dados a serem lidos.

Função Principal (main)
Aqui é onde a aplicação é configurada e a interface gráfica é criada.

Configuração da Conexão com o CLP
NewTCPClientHandler: Inicializa a conexão com o CLP, usando o endereço IP, rack, e slot.

Timeouts: Define tempos de espera para a conexão e quando o CLP está inativo.

Conectar ao CLP e Criar Cliente
Connect: Estabelece a conexão com o CLP.

defer handler.Close(): Fecha a conexão ao final do programa.

NewClient: Cria um cliente para interagir com o CLP.

Interface Gráfica com Fyne
app.New(): Inicializa uma nova aplicação Fyne.

NewWindow: Cria uma nova janela de aplicativo com o título especificado.

Criação de Widgets (Labels e Botão)
Cria três Labels que exibem informações de RPM, temperatura e vibração, iniciando com o texto “Aguardando dados...”.

Botão de Atualização
widget.NewButton: Cria um botão com o texto "Atualizar Dados".

lerDadosCLP: Chama a função lerDadosCLP para ler RPM, temperatura e vibração do CLP e armazena os valores em variáveis.

SetText: Atualiza os textos dos Labels com os valores lidos.

dialog.ShowInformation: Exibe uma janela de diálogo informando que os dados foram atualizados.

Layout da Janela e Exibição
NewVBox: Cria um layout vertical para organizar o título, botão e Labels.

SetContent: Define o layout como o conteúdo da janela.

Resize: Define o tamanho inicial da janela.

ShowAndRun: Exibe a janela e inicia a aplicação.
