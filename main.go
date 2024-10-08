package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	steg "github.com/charles-m-knox/steganographics/pkg/steganographics"
)

var (
	inputFile   string
	outputFile  string
	hiddenText  string
	flagAddr    string
	flagCert    string
	flagKey     string
	flagVersion bool
)

var version = "dev" // to be modified at compile time

func parseFlags() {
	flag.StringVar(&flagAddr, "addr", "", "the address (host and port) to listen on, such as 0.0.0.0:29104")
	flag.StringVar(&flagCert, "cert", "", "the cert.pem file to use for TLS - leave blank for no TLS")
	flag.StringVar(&flagKey, "key", "", "the key.pem file to use for TLS - leave blank for no TLS")

	flag.StringVar(&inputFile, "input", "", "input PNG file to encode the secret into")
	flag.StringVar(&outputFile, "output", "", "output file that will contain the encoded secret")
	flag.StringVar(&hiddenText, "secret", "", "the message to encode into the input file")

	const versionHelp = "print version information and exit"

	flag.BoolVar(&flagVersion, "version", false, versionHelp)
	flag.BoolVar(&flagVersion, "v", false, versionHelp)

	flag.Parse()
}

func main() {
	parseFlags()

	if flagVersion {
		os.Stdout.Write([]byte(fmt.Sprintf("%v\n", version)))

		return
	}

	if flagAddr != "" {
		server(flagAddr)
	} else if inputFile != "" && outputFile != "" && hiddenText != "" {
		err := steg.HideTextInImageFile(inputFile, []byte(hiddenText), outputFile)
		if err != nil {
			log.Printf("Error: %s\n", err)
		}

		log.Println("successfully hid text in image, exiting now")
	} else if inputFile != "" && outputFile == "" && hiddenText == "" {
		extractedText, err := steg.ExtractTextFromImageFile(inputFile)
		if err != nil {
			log.Printf("failed to extract text from file %v: %v", inputFile, err.Error())
		}

		os.Stdout.Write([]byte(extractedText))
	} else if inputFile == "" && outputFile == "" && hiddenText != "" {
		img, _, err := image.Decode(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read img from stdin: %v", err.Error())
		}

		output, err := steg.HideTextInImage(img, []byte(hiddenText))
		if err != nil {
			log.Fatalf("failed to read img from stdin: %v", err.Error())
		}

		err = png.Encode(os.Stdout, output)
		if err != nil {
			log.Fatalf("failed to write img to stdout: %v", err.Error())
		}
	} else if inputFile == "" && outputFile == "" && hiddenText == "" {
		img, _, err := image.Decode(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read img from stdin: %v", err.Error())
		}

		os.Stdout.Write(steg.ExtractTextFromImage(img))
	}
}
