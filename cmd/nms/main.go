// ********* ESTE FICHEIRO É SÓ USADO PARA TESTAR O PARSE *********

package main

import (
	"fmt"
	a "nms/pkg/utils"
	"os"
	"path/filepath"
)

func main() {
	// Caminho para o ficheiro JSON
	filePath := filepath.Join("..", "..", "configs", "settings.json") // ajuste para o caminho do seu ficheiro JSON

	// Leitura do ficheiro JSON
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[ERROR] Unable to read file:", err)
		os.Exit(1)
	}

	// Parse do JSON
	tasks := a.ParseDataFromJson(data)

	// Exemplo de uso dos dados parseados
	for _, task := range tasks {
		fmt.Printf("Task: %+v\n", task)
	}
}

// ********* ESTE FICHEIRO É SÓ USADO PARA TESTAR O PARSE *********
