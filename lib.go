package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

// HideTextInImage implements a basic Least Significant Bit (LSB) steganography
// algorithm. The LSB algorithm hides the secret text by replacing the least
// significant bits of the red channel in the image's pixels with the bits of
// the secret message. The changes are minimal and imperceptible to the human
// eye, allowing the text to remain hidden.
func HideTextInImage(input image.Image, message []byte) (image.Image, error) {
	bounds := input.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Note that the hidden text should be
	// terminated with a null byte (0) for the extraction to work correctly.
	data := append(message, 0)
	if len(data)*8 > width*height {
		return input, fmt.Errorf("input image is too small to hide the text")
	}

	imgWithText := image.NewRGBA(bounds)
	counter := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := input.At(x, y)
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

	return imgWithText, nil
}

// HideTextInImageFile loads a PNG image from a file and executes
// HideTextInImage on it, saving it to outputFile. Note that the hidden text
// should be terminated with a null byte (0) for the extraction to work
// correctly.
func HideTextInImageFile(inputFile string, msg []byte, outputFile string) error {
	img, imgType, err := loadImage(inputFile)
	if err != nil {
		return err
	}

	imgWithText, err := HideTextInImage(img, msg)
	if err != nil {
		return err
	}

	if strings.ToLower(imgType) != "png" {
		return fmt.Errorf("only PNG images are supported")
	}

	err = saveImage(outputFile, imgWithText)
	if err != nil {
		return err
	}

	return nil
}

// loadImage loads an image from filename.
func loadImage(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open input file %v: %w", filename, err)
	}

	defer file.Close()

	img, imgType, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decodoe image from input file %v: %w", filename, err)
	}

	return img, imgType, nil
}

// saveImage saves an image to filename.
func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for saving image to input file %v: %w", filename, err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to encode image to input file %v: %w", filename, err)
	}

	return nil
}

func ExtractTextFromImage(img image.Image) []byte {
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

	return data.Bytes()
}

func ExtractTextFromImageFile(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to open input file %v for text extraction: %w", inputFile, err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image from file %v for text extraction: %w", inputFile, err)
	}

	message := ExtractTextFromImage(img)

	return string(message), nil
}
