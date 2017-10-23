// Package v4l, a facade to the Video4Linux video capture interface
// Copyright (C) 2016 Zoltán Korándi <korandi.z@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Command streamcam streams MJPEG from a V4L device over HTTP.
//
// Command Line
//
// Usage:
//   streamcam [flags]
//
// Flags:
//  -a string
//          address to listen on (default ":8080")
//  -d string
//          path to the capture device
//  -f int
//          frame rate
//  -h int
//          image height
//  -l
//          print supported device configs and quit
//  -r
//          reset all controls to default
//  -w int
//          image width
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
)

func main() {
	var (
		d = flag.String("d", "", "path to the capture device")
		w = flag.Int("w", 0, "image width")
		h = flag.Int("h", 0, "image height")
		f = flag.Int("f", 0, "frame rate")
		a = flag.String("a", ":8080", "address to listen on")
		l = flag.Bool("l", false, "print supported device configs and quit")
		r = flag.Bool("r", false, "reset all controls to default")
	)
	flag.Parse()

	if *d == "" {
		devs := v4l.FindDevices()
		if len(devs) != 1 {
			fmt.Fprintln(os.Stderr, "Use -d to select device.")
			for _, info := range devs {
				fmt.Fprintln(os.Stderr, " ", info.Path)
			}
			os.Exit(1)
		}
		*d = devs[0].Path
	}
	fmt.Fprintln(os.Stderr, "Using device", *d)
	cam, err := v4l.Open(*d)
	fatal("Open", err)

	if *l {
		configs, err := cam.ListConfigs()
		fatal("ListConfigs", err)
		fmt.Fprintln(os.Stderr, "Supported device configs:")
		found := false
		for _, cfg := range configs {
			if cfg.Format != mjpeg.FourCC {
				continue
			}
			found = true
			fmt.Fprintln(os.Stderr, " ", cfg2str(cfg))
		}
		if !found {
			fmt.Fprintln(os.Stderr, "  (none)")
		}
		os.Exit(0)
	}

	cfg, err := cam.GetConfig()
	fatal("GetConfig", err)
	cfg.Format = mjpeg.FourCC
	if *w > 0 {
		cfg.Width = *w
	}
	if *h > 0 {
		cfg.Height = *h
	}
	if *f > 0 {
		cfg.FPS = v4l.Frac{uint32(*f), 1}
	}
	fmt.Fprintln(os.Stderr, "Requested config:", cfg2str(cfg))
	err = cam.SetConfig(cfg)
	fatal("SetConfig", err)
	err = cam.TurnOn()
	fatal("TurnOn", err)
	cfg, err = cam.GetConfig()
	fatal("GetConfig", err)
	if cfg.Format != mjpeg.FourCC {
		fmt.Fprintln(os.Stderr, "Failed to set MJPEG format.")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Actual device config:", cfg2str(cfg))

	if *r {
		ctrls, err := cam.ListControls()
		fatal("ListControls", err)
		for _, ctrl := range ctrls {
			cam.SetControl(ctrl.CID, ctrl.Default)
		}
	}

	blankImg := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, blankImg, nil)
	blank = buf.Bytes()

	go handleInterrupt()
	go stream(cam)

	log.Println("Listening on address", *a)
	srv := http.Server{
		Addr:    *a,
		Handler: http.HandlerFunc(serveHTTP),
	}
	err = srv.ListenAndServe()
	fatal("ListenAndServe", err)
}

func cfg2str(cfg v4l.DeviceConfig) string {
	return fmt.Sprintf("%dx%d @ %.4g FPS", cfg.Width, cfg.Height,
		float64(cfg.FPS.N)/float64(cfg.FPS.D))
}

var (
	mu      sync.Mutex
	clients []*client
	stopped bool
	quit    = make(chan int, 1)
	blank   []byte
)

type client struct {
	i  int
	ch chan []byte
}

func newClient() *client {
	mu.Lock()
	defer mu.Unlock()
	if stopped {
		return nil
	}
	clt := &client{
		i:  len(clients),
		ch: make(chan []byte, 1),
	}
	clt.ch <- blank
	clients = append(clients, clt)
	return clt
}

func (clt *client) remove() {
	mu.Lock()
	defer mu.Unlock()
	i := clt.i
	last := len(clients) - 1
	clients[i] = clients[last]
	clients[i].i = i
	clients[last] = nil
	clients = clients[:last]
	clt.i = -1
	if stopped && len(clients) == 0 {
		quit <- 1
	}
}

func handleInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Stopping...")
	mu.Lock()
	stopped = true
	if len(clients) == 0 {
		quit <- 1
	} else {
		for _, clt := range clients {
			close(clt.ch)
		}
	}
	mu.Unlock()
	<-quit
	os.Exit(0)
}

func stream(cam *v4l.Device) {
	for {
		buf, err := cam.Capture()
		if err != nil {
			log.Println("Capture:", err)
			proc, _ := os.FindProcess(os.Getpid())
			proc.Signal(os.Interrupt)
			break
		}
		b := make([]byte, buf.Size())
		buf.ReadAt(b, 0)
		mu.Lock()
		if stopped {
			mu.Unlock()
			break
		}
		for _, clt := range clients {
			select {
			case clt.ch <- b:
			case <-clt.ch:
				clt.ch <- b
			}
		}
		mu.Unlock()
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] New connection\n", r.RemoteAddr)
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Printf("[%s] ResponseWrite is not a Hijacker", r.RemoteAddr)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		log.Printf("[%s] Hijack: %v\n", r.RemoteAddr, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer conn.Close()

	clt := newClient()
	if clt == nil {
		return
	}
	defer clt.remove()

	const B = "45c7pIy0cxa4vWtwGuVuAkbzKAQGpRjz9eyhyHTv"

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n" +
		"Date: " + time.Now().UTC().Format(http.TimeFormat) + "\r\n" +
		"Content-Type: multipart/x-mixed-replace; boundary=" + B + "\r\n" +
		"Cache-Control: no-cache, no-store, max-age=0, must-revalidate\r\n" +
		"Pragma: no-cache\r\n" +
		"\r\n" +
		"--" + B + "\r\n"))
	if err != nil {
		log.Printf("[%s] %v\n", r.RemoteAddr, err)
		return
	}

	for {
		buf := <-clt.ch
		conn.SetDeadline(time.Now().Add(time.Second))
		_, err := conn.Write([]byte("Content-Type: image/jpeg\r\n\r\n"))
		if err != nil {
			log.Printf("[%s] %v\n", r.RemoteAddr, err)
			return
		}
		if buf == nil {
			_, err := conn.Write(blank)
			if err == nil {
				conn.Write([]byte("--" + B + "--\r\n"))
			} else {
				log.Printf("[%s] %v\n", r.RemoteAddr, err)
			}
			log.Printf("[%s] Quitting\n", r.RemoteAddr)
			return
		}
		_, err = conn.Write(buf)
		if err != nil {
			log.Printf("[%s] %v\n", r.RemoteAddr, err)
			return
		}
		_, err = conn.Write([]byte("--" + B + "\r\n"))
		if err != nil {
			log.Printf("[%s] %v\n", r.RemoteAddr, err)
			return
		}
	}
}

func fatal(p string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", p, err)
		os.Exit(1)
	}
}
