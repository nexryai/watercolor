package watercolor

import (
	"os"
	"testing"
)

func TestIsJPEG(t *testing.T) {
	jpegData, err := os.ReadFile("testdata/pic.jpg")
	if err != nil {
		t.Error(err)
	}

	pngData, err := os.ReadFile("testdata/pic.png")
	if err != nil {
		t.Error(err)
	}

	if !isJPEG(&jpegData) {
		t.Error("isJPEG() failed: not detected")
	} else if isJPEG(&pngData) {
		t.Error("isJPEG() failed: incorrect detection")
	} else {
		t.Log("isJPEG() passed")
	}
}

func TestIsPNG(t *testing.T) {
	pngData, err := os.ReadFile("testdata/pic.png")
	if err != nil {
		t.Error(err)
	}

	jpegData, err := os.ReadFile("testdata/pic.jpg")
	if err != nil {
		t.Error(err)
	}

	if !isPNG(&pngData) {
		t.Error("isPNG() failed: not detected")
	} else if isPNG(&jpegData) {
		t.Error("isPNG() failed: incorrect detection")
	} else {
		t.Log("isPNG() passed")
	}
}

func TestIsAnimatedPNG(t *testing.T) {
	data, err := os.ReadFile("testdata/apng1.png")
	if err != nil {
		t.Error(err)
	}

	if !isAnimatedPNG(&data) {
		t.Error("isAnimatedPNG() failed")
	} else {
		t.Log("isAnimatedPNG() passed")
	}
}
