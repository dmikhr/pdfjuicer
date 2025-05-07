package config

const (
	ImgScaleDownDefault   = 1.0
	ThumbScaleDownDefault = 10.0
	DefaultFilenamePrefix = "page"
	DefaultImgFormat      = "png"
	ThumbnailsDir         = "thumbnails"
)

type Config struct {
	SourcePath string
	SaveDir    string
	Prefix     string
	Postfix    string
	Pages      string
	Image      struct {
		ImgSize      string
		ImgScaleDown float64
		ImgType      string
	}
	Thumb struct {
		CreateThumbnails bool
		ThumbScaleDown   float64
		ThumbnailsSize   string
	}
	WorkersNum  int
	VersionFlag bool
	Quiet       bool
}
