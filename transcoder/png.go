package transcoder

import (
	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"image/png"
	"net/http"
)

type Png struct {
	imgsz		uint
	wpOptions  *webp.Options
}

func NewPng(imgsz uint, quality int) *Png {
	return &Png{
		imgsz: imgsz,
		wpOptions: &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		},
	}
}

func (t *Png) Transcode(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) error {
	img, err := png.Decode(r)
	if err != nil {
		return err
	}

	img = resize.Thumbnail(t.imgsz, t.imgsz, img, resize.NearestNeighbor)

	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
		if err = webp.Encode(w, img, t.wpOptions); err != nil {
			return err
		}
	} else {
		if err = png.Encode(w, img); err != nil {
			return err
		}
	}
	return nil
}

func (t *Png) TranscodeHead(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) {
	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
	}
}
