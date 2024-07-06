package watercolor

import (
	"os"
	"testing"
)

func TestProcessStaticImage(t *testing.T) {
	jpegData, err := os.ReadFile("testdata/pic.jpg")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&jpegData, &TargetImage{MaxWidth: 720, MaxHeight: 720})
	if err != nil {
		t.Error(err)
	}

	os.Create("testout/pic.webp")
	os.WriteFile("testout/pic.webp", *webpBytes, 0644)
}
