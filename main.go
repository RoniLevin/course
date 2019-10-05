package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	s := http.NewServeMux()

	s.HandleFunc("/upload", uploadHandler)
	s.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		indexFile, err := ioutil.ReadFile("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(indexFile)
	})

	handler := handlers.LoggingHandler(os.Stdout, s)
	http.ListenAndServe(":3000", handler)
	//http.ListenAndServe(":3000", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := getChunk(r); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getChunk(r *http.Request) error {
	// Part 1: Chunk Number
	// Part 2: Chunk Size
	// Part 3: Current Chunk Size
	// Part 4: Total Size
	// Part 5: Identifier
	// Part 6: File Name
	// Part 7: Relative Path
	// Part 8: Total Chunks
	// Part 9: Chunk Data
	reader, err := r.MultipartReader()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	//part 1: Chunk Number
	part, err := reader.NextPart()
	if err != nil {
		return err
	}
	if _, err = io.Copy(buf, part); err != nil {
		return err
	}
	chunkNumber := buf.String()
	buf.Reset()

	//move through unused parts 2-5
	for i := 2; i < 6; i++ {
		_, err = reader.NextPart()
		if err != nil {
			return err
		}
	}

	//part 6: File Name
	part, err = reader.NextPart()
	if err != nil {
		return err
	}
	if _, err = io.Copy(buf, part); err != nil {
		return err
	}

	var f *os.File
	if chunkNumber == "1" {
		f, err = os.OpenFile(buf.String(),
			os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.OpenFile(buf.String(),
			os.O_APPEND|os.O_WRONLY, 0644)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	buf.Reset()

	//move through unused parts 7-8
	for i := 7; i < 9; i++ {
		_, err = reader.NextPart()
		if err != nil {
			return err
		}
	}

	//part 9: Chunk Data
	part, err = reader.NextPart()
	if err != nil {
		return err
	}

	if _, err = io.Copy(buf, part); err != nil {
		return err
	}

	return nil
}
