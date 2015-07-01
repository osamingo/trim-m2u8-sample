package main

import (
	"net/http"
	"strconv"

	"github.com/grafov/m3u8"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const m3u8URL = "TARGET_URL"

func main() {
	goji.Get("/trim", trimM3u8)
	goji.Serve()
}

func trimM3u8(c web.C, w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("sec")

	sec, err := strconv.Atoi(q)
	if err != nil {
		sec = 0
	}

	resp, err := http.Get(m3u8URL)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("faild to get m3u8 file"))
		return
	}
	defer resp.Body.Close()

	p, err := m3u8.NewMediaPlaylist(1024, 1024)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("faild to NewMediaPlaylist"))
		return
	}
	defer p.Close()

	err = p.DecodeFrom(resp.Body, true)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("faild to decode m3u8 file"))
		return
	}

	if sec != 0 {
		p.Segments = append(p.Segments[(sec / int(p.TargetDuration)):])
	}

	w.Write(p.Encode().Bytes())
}

