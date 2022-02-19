package media

import (
	"io"
	"image"
	"image/png"
	_ "image/jpeg"
	_ "image/gif"
	"golang.org/x/image/draw"
)

type PreviewOptions struct {
	Width  int
	Height int
}

func GeneratePreview(src io.ReadCloser, out io.Writer, opts *PreviewOptions) error {
		sourceImage, _, err := image.Decode(src)
		if err != nil {
			return err
		}
		newImage := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
		draw.ApproxBiLinear.Scale(
			newImage,
			newImage.Rect,
			sourceImage,
			sourceImage.Bounds(),
			draw.Over,
			nil,
		)
		return png.Encode(out, newImage)
}
