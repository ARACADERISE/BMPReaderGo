package main

import (
	"bmpdecoder/src"
	"fmt"
)

func main() {
	decoder := src.NewDecoder("test.bmp")
	_, err := decoder.DecodeHeader()

	src.PrintErr(err)

	_, err = decoder.PickupPixels()

	decoder.Write(err)

	info := src.New_encoder("bmpFileInfo.json")
	fmt.Println(info)
}