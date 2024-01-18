package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"gitea.cmcode.dev/cmcode/steganographics/secrets"
)

var (
	inputFile  string
	outputFile string
	hiddenText string
	flagAddr   string
	flagCert   string
	flagKey    string
)

func parseFlags() {
	flag.StringVar(&flagAddr, "addr", "", "the address (host and port) to listen on, such as 0.0.0.0:29104")
	flag.StringVar(&flagCert, "cert", "", "the cert.pem file to use for TLS - leave blank for no TLS")
	flag.StringVar(&flagKey, "key", "", "the key.pem file to use for TLS - leave blank for no TLS")

	flag.StringVar(&inputFile, "input", "", "Input PNG file to encode the secret into")
	flag.StringVar(&outputFile, "output", "", "Output file that will contain the encoded secret")
	flag.StringVar(&hiddenText, "secret", "", "The message to encode into the input file")

	flag.Parse()
}

func main() {
	parseFlags()

	if flagAddr != "" {
		server(flagAddr)
	} else if inputFile != "" && outputFile != "" && hiddenText != "" {
		err := secrets.HideTextInImageFile(inputFile, []byte(hiddenText), outputFile)
		if err != nil {
			log.Printf("Error: %s\n", err)
		}

		log.Println("successfully hid text in image, exiting now")
	} else if inputFile != "" && outputFile == "" && hiddenText == "" {
		extractedText, err := secrets.ExtractTextFromImageFile(inputFile)
		if err != nil {
			log.Printf("failed to extract text from file %v: %v", inputFile, err.Error())
		}

		os.Stdout.Write([]byte(extractedText))
	} else if inputFile == "" && outputFile == "" && hiddenText != "" {
		img, _, err := image.Decode(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read img from stdin: %v", err.Error())
		}

		output, err := secrets.HideTextInImage(img, []byte(hiddenText))
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

		os.Stdout.Write(secrets.ExtractTextFromImage(img))
	}
}
