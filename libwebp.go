package watercolor

/*
   #cgo LDFLAGS: -lwebp -lwebpmux
   #cgo darwin pkg-config: libwebp
   #include <webp/encode.h>
   #include <webp/mux.h>
*/
import "C"
import (
	"errors"
	"image"
	"unsafe"
)

var (
	ErrEncoderReturnedUnknownError = errors.New("error occurred while encoding WebP")
)

func rgbaToWebP(srcRgba *image.RGBA, quality int) (*[]byte, error) {
	if srcRgba == nil {
		return nil, ErrImageIsNil
	}

	qualityFloat := float64(quality)

	var cOutput *C.uint8_t
	outPtr := (**C.uint8_t)(unsafe.Pointer(&cOutput))

	length := C.WebPEncodeRGBA((*C.uint8_t)(unsafe.Pointer(&srcRgba.Pix[0])), C.int(srcRgba.Bounds().Dx()), C.int(srcRgba.Bounds().Dy()), C.int(srcRgba.Stride), C.float(qualityFloat), outPtr)
	if length == 0 {
		return nil, ErrEncoderReturnedUnknownError
	}

	// Convert the C array to a Go byte slice
	var webpBytes []byte

	/*
		sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&webpBytes)))
		sliceHeader.Cap = int(length)
		sliceHeader.Len = int(length)
		sliceHeader.Data = uintptr(unsafe.Pointer(cOutput))
	*/

	webpBytes = C.GoBytes(unsafe.Pointer(cOutput), C.int(length))

	// Free the memory allocated by WebP
	C.WebPFree(unsafe.Pointer(cOutput))

	return &webpBytes, nil
}
