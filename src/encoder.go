package src

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Encoder struct {
	File_content		[]int		`json:"content"`
	File_size			int			`json:"size"`
	Index				int			`json:"Start_Index"`
	Header				BmpHeader	`json:"HEADER"`
	DIBheader			DIBHeader	`json:"DIBHEADER"`
}

func New_encoder(json_file string) *Encoder {
	info := &Encoder{ }

	dir, e := filepath.Abs(json_file)

	_, e = os.Stat(dir)

	if e != nil {
		log.Fatal(e)
	}

	file, err := os.Open(dir)

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(info)

	if err != nil {
		log.Fatal(err)
	}

	return info
}