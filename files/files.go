package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func OpenFile(filepath string) {
	content, _ := os.ReadFile(filepath)

	fmt.Println(string(content))
}

func WriteFile(filepath string) {
	content := []byte("Hello, Gophers!")

	os.WriteFile(filepath, content, 0644)
}

func OpenFileWithStream(filepath string) {
	file, _ := os.Open(filepath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	err := scanner.Err()

	if err != nil {
		log.Fatalf("GAVNO")
	}
}

func WriteFileWithStream(filepath string) {
	file, _ := os.Create(filepath)

	defer file.Close()

	bufWriter := bufio.NewWriter(file)

	bufWriter.WriteString("content123")

	bufWriter.Flush()
}
