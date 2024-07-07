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

	if !isWebP(webpBytes) {
		t.Error("output is not webp")
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

	if !isAVIF(webpBytes) {
		t.Error("output is not AVIF")
	}

	os.Create("testout/pic.avif")
	os.WriteFile("testout/pic.avif", *webpBytes, 0644)
}

func TestProcessStaticImagePngToWebP(t *testing.T) {
	pngData, err := os.ReadFile("testdata/pic.png")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&pngData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   80,
		Format:    TargetFormatWebP,
	})
	if err != nil {
		t.Error(err)
	}

	if !isWebP(webpBytes) {
		t.Error("output is not webp")
	}

	os.Create("testout/pic.webp")
	os.WriteFile("testout/pic.webp", *webpBytes, 0644)
}

func TestProcessStaticImagePngToAVIF(t *testing.T) {
	pngData, err := os.ReadFile("testdata/pic.png")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&pngData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   60,
		Format:    TargetFormatAVIF,
	})
	if err != nil {
		t.Error(err)
	}

	if !isAVIF(webpBytes) {
		t.Error("output is not AVIF")
	}

	os.Create("testout/pic.avif")
	os.WriteFile("testout/pic.avif", *webpBytes, 0644)
}

func TestProcessStaticImageAvifToWebP(t *testing.T) {
	avifData, err := os.ReadFile("testdata/pic.avif")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&avifData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   80,
		Format:    TargetFormatWebP,
	})
	if err != nil {
		t.Error(err)
	}

	if !isWebP(webpBytes) {
		t.Error("output is not webp")
	}

	os.Create("testout/pic.webp")
	os.WriteFile("testout/pic.webp", *webpBytes, 0644)
}

func TestProcessStaticImageIcoToWebP(t *testing.T) {
	icoData, err := os.ReadFile("testdata/pic.ico")
	if err != nil {
		t.Error(err)
	}

	webpBytes, err := ProcessStaticImage(&icoData, &TargetImage{
		MaxWidth:  1920,
		MaxHeight: 9999,
		Quality:   80,
		Format:    TargetFormatWebP,
	})
	if err != nil {
		t.Error(err)
	}

	if !isWebP(webpBytes) {
		t.Error("output is not webp")
	}

	os.Create("testout/pic.webp")
	os.WriteFile("testout/pic.webp", *webpBytes, 0644)
}
