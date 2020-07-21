package transcoder

import (
	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"net/http"
	"strconv"
)

type WebP struct {
	imgsz		uint
	wpOptions *webp.Options
}

func NewWebP(imgsz uint, quality int) *WebP {
	return &WebP{
		imgsz: imgsz,
		wpOptions: &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		},
	}
}

func (t *WebP) Transcode(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) error {
	img, err := webp.Decode(r)
	if err != nil {
		return err
	}

	img = resize.Thumbnail(t.imgsz, t.imgsz, img, resize.NearestNeighbor)

	wpOptions := t.wpOptions
	qualityString := headers.Get("X-Compy-Quality")
	if qualityString != "" {
		if quality, err := strconv.Atoi(qualityString); err != nil {
			return err
		} else {
			wpOptions.Quality = float32(quality)
		}
	}

	if err = webp.Encode(w, img, wpOptions); err != nil {
		return err
	}
	return nil
}

func (t *WebP) TranscodeHead(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) {
	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
	}
}
