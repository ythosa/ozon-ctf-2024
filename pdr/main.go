package main

const (
	sourceName = "traffic.pcapng"
	resultName = "traffic.decoded.png"
)

const (
	width  = 1280
	height = 720
)

func main() {
	parser := newPdrParser()
	reader := newPdrReader(sourceName, parser, 0)
	reader.read()

	encryptedImageBytes := parser.fetchEncryptedImage()
	keyImageBytes := parser.fetchKeyImage()

	decodeImage(encryptedImageBytes, keyImageBytes, resultName)
}
