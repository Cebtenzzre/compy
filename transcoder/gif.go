package transcoder

import (
	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
	"image/gif"
	"net/http"
)

type Gif struct {
	wpOptions  *webp.Options
}

func NewGif(quality int) *Gif {
	return &Gif{
		wpOptions: &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		},
	}
}

func (t *Gif) Transcode(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) error {
	if !SupportsWebP(headers) {
		return w.ReadFrom(r)
	}

	img, err := gif.Decode(r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "image/webp")
	if err = webp.Encode(w, img, t.wpOptions); err != nil {
		return err
	}
	return nil
}

func (t *Gif) TranscodeHead(w *proxy.ResponseWriter, r *proxy.ResponseReader, headers http.Header) {
	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
	}
}
