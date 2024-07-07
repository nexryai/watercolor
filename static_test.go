package watercolor

import (
	"os"
	"testing"
)

func TestProcessStaticImageJpegToWebP(t *testing.T) {
	jpegData, err := os.ReadFile("testdata/pic.jpg")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&jpegData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   80,
		Format:    TargetFormatWebP,
	})
	if err != nil {
		t.Error(err)
	}

	os.Create("testout/pic.webp")
	os.WriteFile("testout/pic.webp", *webpBytes, 0644)
}

func TestProcessStaticImageJpegToAVIF(t *testing.T) {
	jpegData, err := os.ReadFile("testdata/pic.jpg")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&jpegData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   60,
		Format:    TargetFormatAVIF,
	})
	if err != nil {
		t.Error(err)
	}

	os.Create("testout/pic.avif")
	os.WriteFile("testout/pic.avif", *webpBytes, 0644)
}
