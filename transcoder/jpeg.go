package transcoder

import (
	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/pixiv/go-libjpeg/jpeg"
	"net/http"
	"strconv"
)

type Jpeg struct {
	imgsz		uint
	decOptions *jpeg.DecoderOptions
	encOptions *jpeg.EncoderOptions
	wpOptions  *webp.Options
}

func NewJpeg(imgsz uint, quality, wpquality int) *Jpeg {
	return &Jpeg{
		imgsz: imgsz,
		decOptions: &jpeg.DecoderOptions{},
		encOptions: &jpeg.EncoderOptions{
			Quality:        quality,
			OptimizeCoding: true,
		},
		wpOptions: &webp.Options{
			Lossless: false,
			Quality: float32(wpquality),
		},
	}
}

func (t *Jpeg) Transcode(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) error {
	img, err := jpeg.Decode(r, t.decOptions)
	if err != nil {
		return err
	}

	img = resize.Thumbnail(t.imgsz, t.imgsz, img, resize.NearestNeighbor)

	encOptions := t.encOptions
	qualityString := headers.Get("X-Compy-Quality")
	if qualityString != "" {
		if quality, err := strconv.Atoi(qualityString); err != nil {
			return err
		} else {
			encOptions.Quality = quality
		}
	}

	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
		if err = webp.Encode(w, img, t.wpOptions); err != nil {
			return err
		}
	} else {
		if err = jpeg.Encode(w, img, encOptions); err != nil {
			return err
		}
	}
	return nil
}

func (t *Jpeg) TranscodeHead(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) {
	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
	}
}
