package watercolor

/*
   #cgo LDFLAGS: -lwebp -lwebpmux
   #cgo darwin pkg-config: libwebp
   #include <webp/encode.h>
   #include <webp/mux.h>
*/
import "C"
import (
	"image"
	"unsafe"
)

func rgbaToWebP(srcRgba *image.RGBA, quality int) (*[]byte, error) {
	if srcRgba == nil {
		return nil, ErrImageIsNil
	}

	var cOutput *C.uint8_t
	outPtr := (**C.uint8_t)(unsafe.Pointer(&cOutput))

	C.WebPEncodeRGBA((*C.uint8_t)(unsafe.Pointer(&srcRgba.Pix[0])), C.int(srcRgba.Bounds().Dx()), C.int(srcRgba.Bounds().Dy()), C.int(srcRgba.Stride), C.float(quality), outPtr)

	webpBytes := C.GoBytes(unsafe.Pointer(cOutput), C.int(len(srcRgba.Pix)))

	// Free the memory allocated by WebP
	C.WebPFree(unsafe.Pointer(cOutput))

	return &webpBytes, nil
}
