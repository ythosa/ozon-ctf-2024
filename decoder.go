package main

import (
	"crypto/aes"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const (
	decryptionStep = 48
	keyPartSize    = 16
)

func decodeImageBytes(encrypted []byte, key []byte) []byte {
	if len(encrypted) != len(key) {
		log.Fatal("encrypted and keys must be the same length")
	}

	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i += decryptionStep {
		keyPart := key[i : i+keyPartSize]
		aesBlock, err := aes.NewCipher(keyPart)
		if err != nil {
			log.Fatalf("failed to create cipher for block: %v", err)
		}

		for j := 1; j <= decryptionStep/keyPartSize; j++ {
			startIndex := i + (j-1)*keyPartSize
			endIndex := i + j*keyPartSize
			aesBlock.Decrypt(decrypted[startIndex:endIndex], encrypted[startIndex:endIndex])
		}
	}

	return decrypted
}

func decodeImage(encrypted []byte, key []byte, dstpath string) {
	imageData := decodeImageBytes(encrypted, key)

	// Создаем изображение
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Индекс в исходном слайсе imageData
	index := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := imageData[index]
			g := imageData[index+1]
			b := imageData[index+2]
			index += 3
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// Создаем файл для сохранения изображения
	f, err := os.Create(dstpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Сохраняем изображение в формате PNG
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
