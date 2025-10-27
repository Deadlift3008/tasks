package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Data struct {
	id        int    `json:"id"`
	firstName string `json:"first_name"`
}

func main() {
	db, err := sql.Open("postgres", "postgresql://strong_password:user@database-host:5000/my_db")
	if err != nil {
		log.Println("Error opening the database:", err)
		return
	}
	defer db.Close()

	ctx := context.TODO()
	var wg sync.WaitGroup
	data := []Data{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}}
	for _, d := range data {
		go func() {
			wg.Add(1)

			d.firstName = "New " + d.firstName

			query := fmt.Sprintf("INSERT INTO data (name) VALUES ('%s')", d.firstName)
			_, err := db.Exec(query)

			if err != nil {
				log.Println("Error inserting data:", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	f, err := json.Marshal(data)
	if err != nil {
		log.Println("Error inserting data:", err)
		return
	}

	var buf bytes.Buffer
	buf.Write(f)

	req, err := http.NewRequestWithContext(ctx, "GET", "domain.ru/my_path", &buf)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(io.ReadAll(resp.Body))
}
