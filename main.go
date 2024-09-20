package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func searchDir(dirPath string) []WasmFile {
	wasmFiles := []WasmFile{}
	files, err := os.ReadDir(dirPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "node_modules" {
			wasmFiles = append(wasmFiles, searchDir(dirPath+file.Name()+"/")...)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if filepath.Ext(file.Name()) == ".wasm" {
				wasmFiles = append(wasmFiles, WasmFile{FilePath: dirPath + file.Name(), FileName: file.Name()})
			}
		}
	}
	return wasmFiles
}

type WasmFile struct {
	FilePath string
	FileName string
}

func main() {
	wasmFiles := searchDir("../")

	data := map[string][]WasmFile{
		"WasmFiles": wasmFiles,
	}

	gameHandler := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))

		templ.Execute(w, data)
	}

	http.HandleFunc("/", gameHandler)

	for _, wasmFile := range data["WasmFiles"] {
		fmt.Println(wasmFile.FileName)
		http.HandleFunc("/"+wasmFile.FileName, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/wasm")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			log.Printf("Serving WASM file %s to %s...\n", wasmFile.FilePath, r.RemoteAddr)
			http.ServeFile(w, r, wasmFile.FilePath)
		})
	}

	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./wasm_exec.js")
	})
	http.HandleFunc("/main.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./main.html")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
