package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type bitwardenExport struct {
	Folders []bitwardenFolder `json:"folders"`
	Items   []bitwardenItem   `json:"items"`
}

type bitwardenFolder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type bitwardenItem struct {
	ID       string         `json:"id"`
	FolderID string         `json:"folderId"`
	Type     int            `json:"type"`
	RePrompt int            `json:"reprompt"`
	Name     string         `json:"name"`
	Favorite bool           `json:"favorite"`
	Login    bitwardenLogin `json:"login"`
}

type bitwardenLogin struct {
	URIs     []bitwardenLoginURI `json:"uris"`
	UserName string              `json:"username"`
	Password string              `json:"password"`
}

type bitwardenLoginURI struct {
	Match string `json:"match"`
	URI   string `json:"uri"`
}

func main() {
	// site,username,password
	var csvFilePath string
	flag.StringVar(&csvFilePath, "csv", "data.csv", "Path to file *.csv with import data")
	flag.Parse()

	records, err := readCsv(csvFilePath)
	if err != nil {
		panic(err)
	}

	items := make([]bitwardenItem, 0, 400)
	folderID := uuid.New().String()

	for _, r := range records {
		item := bitwardenItem{
			ID:       uuid.New().String(),
			FolderID: folderID,
			Type:     1,
			RePrompt: 0,
			Name:     r[0],
			Login: bitwardenLogin{
				URIs:     []bitwardenLoginURI{bitwardenLoginURI{URI: r[0]}},
				UserName: r[1],
				Password: r[2],
			},
		}
		items = append(items, item)
	}

	export := bitwardenExport{
		Folders: []bitwardenFolder{bitwardenFolder{ID: folderID, Name: "Main"}},
		Items:   items,
	}

	b, err := json.Marshal(export)
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Stdout.Write(b)
}

func readCsv(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	reader.Comment = '#'

	result, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}
