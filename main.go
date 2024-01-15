package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

var (
	inputFile  string
	outputFile string
	hiddenText string
	httpServer bool
	port       int
	address    string
)

func init() {
	flag.StringVar(&inputFile, "input", "input.png", "Input PNG file to encode the secret into")
	flag.StringVar(&outputFile, "output", "output.png", "Output file that will contain the encoded secret")
	flag.StringVar(&hiddenText, "secret", "secret message", "The message to encode into the input file")
	flag.BoolVar(&httpServer, "server", false, "If specified, will run an http server that encodes/decodes on demand")
	flag.IntVar(&port, "port", 8080, "The port that the server will listen on")
	flag.StringVar(&address, "addr", "0.0.0.0", "The bind address for the server")
}

func main() {
	flag.Parse()

	if httpServer {
		server(address, port)
	}

	err := hideTextInImage(inputFile, hiddenText, outputFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Successfully hidden text in the image.")
	}

	// Note that the hidden text should be terminated with a null byte (0) for the extraction to work correctly.
	extractedText, err := extractTextFromImage(outputFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Hidden text: ", extractedText)
	}
}

func hideTextInImage(inputFile, hiddenText, outputFile string) error {
	img, imgType, err := loadImage(inputFile)
	if err != nil {
		return err
	}

	if strings.ToLower(imgType) != "png" {
		return fmt.Errorf("only PNG images are supported")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	data := append([]byte(hiddenText), 0)
	if len(data)*8 > width*height {
		return fmt.Errorf("Input image is too small to hide the text")
	}

	imgWithText := image.NewRGBA(bounds)
	counter := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			rgba := color.RGBAModel.Convert(c).(color.RGBA)

			if counter < len(data)*8 {
				b := data[counter/8]
				bit := (b >> uint(counter%8)) & 1
				rgba.R = (rgba.R & 0xFE) | bit
				counter++
			}

			imgWithText.Set(x, y, rgba)
		}
	}

	err = saveImage(outputFile, imgWithText, imgType)
	if err != nil {
		return err
	}

	return nil
}

func loadImage(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, imgType, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, imgType, nil
}

func saveImage(filename string, img image.Image, imgType string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func extractTextFromImage(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var data bytes.Buffer
	var currentByte byte
	bitIndex := 0

	for counter := 0; counter < width*height; counter++ {
		x, y := counter%width, counter/width

		c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
		bit := c.R & 1

		currentByte |= (bit << uint(bitIndex))
		bitIndex++

		if bitIndex == 8 {
			if currentByte == 0 {
				break
			}
			data.WriteByte(currentByte)
			currentByte = 0
			bitIndex = 0
		}
	}

	return data.String(), nil
}
