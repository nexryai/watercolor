package watercolor

// getResizedImageSizeKeepRatio is a function that returns the scale of the image after resizing while maintaining the aspect ratio.
// if both srcWidth and srcHeight are less than or equal to maxWidth and maxHeight, the original size is returned.
func getResizedImageScaleKeepRatio(srcWidth, srcHeight, maxWidth, maxHeight int) float64 {
	var scale float64

	shouldResize := false
	if srcWidth > maxWidth {
		shouldResize = true
	} else if srcHeight > maxHeight {
		shouldResize = true
	}

	if !shouldResize {
		return 1.0
	} else {
		// 超過が大きい方に合わせる
		widthExcess := srcWidth - maxWidth
		heightExcess := srcHeight - maxHeight

		if widthExcess > heightExcess {
			scale = float64(maxWidth) / float64(srcWidth)
		} else {
			scale = float64(maxHeight) / float64(srcHeight)
		}

		return scale
	}
}
