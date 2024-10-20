package agent

import (
    "fmt"
    "net"
    "os"
    "bufio"
)

func open_agent() {
    // Conectar ao servidor
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Erro ao conectar ao servidor:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Conectado ao servidor. Escreve uma mensagem:")

    // Ler e enviar mensagens ao servidor
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("-> ")
        message, _ := reader.ReadString('\n')

        // Enviar mensagem ao servidor
        _, err = conn.Write([]byte(message))
        if err != nil {
            fmt.Println("Erro ao enviar mensagem:", err)
            return
        }

        // Ler a resposta do servidor
        response, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println("Erro ao receber resposta:", err)
            return
        }
        fmt.Print("Resposta do servidor: ", response)
    }
}
