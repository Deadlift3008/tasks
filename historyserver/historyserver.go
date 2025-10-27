package historyserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Request() {
	bts, _ := json.Marshal(map[string]string{
		"key": "value",
	})

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(bts))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("REQUEST PARASHA")
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	fmt.Println(string(responseBody))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(map[string]bool{
		"success": true,
	})

	w.Write(body)
}

func Server() {
	http.HandleFunc("/path", Handler)

	http.ListenAndServe(":8080", nil)
}
