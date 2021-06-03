package src

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"errors"
	"encoding/json"
)

type BmpHeader struct {
	BMPfilesize		int `json:"BMPFileSize"`
	Pixel_start_index	int	`json:"PIXEL_START_INDEX"`
}

type DIBHeader struct {
	Bmp_image_width		int	`json:"image_width"`
	Bmp_image_height	int	`json:"image_height"`
	Bmp_DIB_size		int	`json:"DIB_Size"`
	Color_planes		int	`json:"Color_Planes"`
	Bits_per_pixel		int	`json:"Bits_Per_pixel"`
	Bitmap_size			int	`json:"Bitmap_Size"`
	V_resolution		int	`json:"V_Reso"`
	H_resolution		int	`json:"H_Reso"`
}

type Decoder struct {
	File_content		[]int		`json:"content"`
	File_size			int			`json:"size"`
	Index				int			`json:"Start_Index"`
	Header				BmpHeader	`json:"HEADER"`
	DIBheader			DIBHeader	`json:"DIBHEADER"`
	Type				uint8		`json:"type"`
	PixelArray			[]byte		`json:"PIXEL_ARRAY"`
	DecodedPixelArr		[]int		`json:"DECODED_PIXEL_ARRAY"`
}

func NewDecoder(filename string) Decoder {
	info := Decoder{ }

	// Default value
	info.Header.BMPfilesize = 0

	dir, _ := filepath.Abs(filename)
	file, err := os.Stat(dir)

	if err != nil {
		log.Fatal(err)
	}

	info.File_size = int(file.Size())

	buffer, e := ioutil.ReadFile(file.Name())

	if e != nil {
		log.Fatal(e)
	}

	for i := range buffer {
		info.File_content = append(info.File_content, int(buffer[i]))
	}
	info.Index = 0

	// Go ahead and get the BMP file size.
	info.Header.BMPfilesize = int(info.File_content[2])
	
	// Go ahead and get where the pixel array starts
	info.Header.Pixel_start_index = int(info.File_content[10])
	
	// Get the width and height
	info.DIBheader.Bmp_image_width = int(info.File_content[18])
	info.DIBheader.Bmp_image_height = int(info.File_content[22])

	// Configure whether or not it's rgb or rgba
	index := 0
	length := 0

	for {
		if index == len(info.File_content) - 1 {
			break
		}
		length++
		index++
	}

	if length == 3 * (info.DIBheader.Bmp_image_width * info.DIBheader.Bmp_image_height) {
		info.Type = 0 // rgb
	} else {
		info.Type = 1 // rgba
	}

	return info
}

/*
 * Simple little helper function for error checking in main.go
 */
func PrintErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (info *Decoder) DecodeHeader() (*Decoder, error) {
	var expected_header []int = []int{ 0x42, 0x4D }

	if !(info.File_content[0] == expected_header[0]) && !(info.File_content[1] == expected_header[1]) {
		// Return empty Decoder just in case.
		return &Decoder{}, errors.New(fmt.Sprintf("Invalid start to BMP file: %d", info.File_content[0]))
	}
	info.Index += 14

	info.DIBheader.Bmp_DIB_size = int(info.File_content[info.Index])
	info.Index += 12 // we skip over the 4 bytes for the DIB size, and we already got the height and width. So skip 8 more bytes

	info.DIBheader.Color_planes = int(info.File_content[info.Index])
	info.Index += 2

	info.DIBheader.Bits_per_pixel = int(info.File_content[info.Index])
	info.Index += 10// we just read the current two bytes, the next 4 bytes are for compression(which there is non)

	info.DIBheader.Bitmap_size = int(info.File_content[info.Index])
	info.Index += 4

	// Both are the same.
	info.DIBheader.H_resolution = int(info.File_content[info.Index]) * int(info.File_content[info.Index + 1]) + 10201
	info.DIBheader.V_resolution = info.DIBheader.H_resolution
	info.Index += 12

	if info.Index == info.Header.Pixel_start_index {
		return info, nil
	}

	return &Decoder{}, errors.New(fmt.Sprintf("Invalid ending of header: %d", info.Index))
}

/*
 * Pickup the pixel array within the bmp image.
 *
 * Error if there is a invalid length to the stream, or invalid
 * pixel value(greater than 255, less than 0). There shouldn't be an error,
 * otherwise the bmp image would not render correctly. It's just safe to
 * have error checking whilst decoding a non-compressed image like a
 * BMP image.
 */
func (info *Decoder) PickupPixels() (*Decoder, error) {

	index := info.Header.Pixel_start_index
	var pixel_array []byte

	for {
		if index== len(info.File_content) - 1 {
			break
		}

		if info.File_content[index] < 0 || info.File_content[index] > 0xff {
			return &Decoder{ }, errors.New(fmt.Sprintf("Invalid pixel value: %d", info.File_content[index]))
		}

		pixel_array = append(pixel_array, byte(info.File_content[index]))
		index++
	}

	info.PixelArray = pixel_array
	for i := range info.PixelArray {
		info.DecodedPixelArr = append(info.DecodedPixelArr, int(info.PixelArray[i]))
	}

	return info, nil
}

/* Write function will ONLY write if the last error of the program is nil. Else, we fatally log the error*/
func (info *Decoder) Write(lastError error) {
	if lastError == nil {
		file, err := json.MarshalIndent(info, "", "\t")

		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile("bmpFileInfo.json", file, 0644)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(lastError)
	}
}