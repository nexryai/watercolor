package watercolor

import (
	"io"
	"os"
	"testing"
)

func TestApngToWebP(t *testing.T) {
	// Open our animated PNG file
	f, err := os.Open("testdata/apng1.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	apngBytes, err := io.ReadAll(f)
	webp, err := apngToWebP(&apngBytes, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Save to webp file
	outFile, err := os.Create("testout/apng1.webp")
	if err != nil {
		t.Fatal(err)
	}

	_, err = outFile.Write(*webp)
	if err != nil {
		t.Fatal(err)
	}
}
