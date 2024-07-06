package watercolor

import "bytes"

type ImageFormat int

const (
	ImageFormatUnknown ImageFormat = iota
	ImageFormatJPEG
	ImageFormatPNG
	ImageFormatAnimatedPNG
	ImageFormatWebP
	ImageFormatAnimatedWebP
	ImageFormatGIF
	ImageFormatAVIF
	ImageFormatAnimatedAVIF
)

func isJPEG(data *[]byte) bool {
	return bytes.HasPrefix(*data, []byte("\xff\xd8\xff"))

}

func isPNG(data *[]byte) bool {
	return bytes.HasPrefix(*data, []byte("\x89PNG\r\n\x1a\n"))
}

func isAnimatedPNG(data *[]byte) bool {
	// APNGのシグネチャをチェック
	return len(*data) > 41 && bytes.Equal((*data)[37:41], []byte("acTL"))
}

func isWebP(data *[]byte) bool {
	if len(*data) < 12 {
		return false
	}

	return bytes.HasPrefix(*data, []byte("RIFF")) && bytes.Equal((*data)[8:12], []byte("WEBP"))
}

func isAnimatedWebP(data *[]byte) bool {
	// Animated WebPの場合、ファイルの0x1Eから0x22がANIMになる
	return len(*data) > 0x22 && string((*data)[0x1E:0x22]) == "ANIM"
}

func isGIF(data *[]byte) bool {
	return bytes.HasPrefix(*data, []byte("GIF87a")) || bytes.HasPrefix(*data, []byte("GIF89a"))
}

func isAVIF(data *[]byte) bool {
	return len(*data) > 12 && bytes.Equal((*data)[4:12], []byte("ftypavif"))
}

func detectImageFormat(data *[]byte) ImageFormat {
	if isJPEG(data) {
		return ImageFormatJPEG
	}

	if isPNG(data) {
		if isAnimatedPNG(data) {
			return ImageFormatAnimatedPNG
		}
		return ImageFormatPNG
	}

	if isWebP(data) {
		if isAnimatedWebP(data) {
			return ImageFormatAnimatedWebP
		}
		return ImageFormatWebP
	}

	if isGIF(data) {
		return ImageFormatGIF
	}

	if isAVIF(data) {
		return ImageFormatAVIF
	}

	return ImageFormatUnknown
}
