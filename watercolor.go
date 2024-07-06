package watercolor

type TargetFormat int

const (
	TargetFormatWebP TargetFormat = iota
	TargetFormatAVIF
)

type TargetImage struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
	Format    TargetFormat
}
