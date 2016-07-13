package main

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path"
	"github.com/gorilla/handlers"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/nfnt/resize"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse URL (N.B. IMAGEROOT environment variable)
	requestedImage := r.URL.Path[1:]
	requestedFile := path.Join(os.Getenv("IMAGEROOT"), requestedImage)

	// Get image
	fp, err := os.Open(requestedFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s for %s\n", err, requestedFile)
		w.WriteHeader(404)
		return
	}

	// Decode image
	img, imgFormat, err := image.Decode(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s decoding %s\n", err, requestedFile)
		w.WriteHeader(400)
		return
	}
	fmt.Fprintf(os.Stdout, "imgFormat is %s\n", imgFormat)

	// Resize image
	img = resize.Resize(uint(img.Bounds().Max.X / 2), 0, img, resize.Lanczos3)

	// Encode image and stream
	w.Header().Add("content-type", "image/jpeg")
	jpeg.Encode(w, img, nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// log requests in Apache CLF format to a file that will rollover to a
	// new one every 50 MB
	logger := handlers.LoggingHandler(&lumberjack.Logger{
		Filename:   "imageshrinker-access.log",
		MaxSize:    50, // megabytes
		MaxBackups: 100000,
		MaxAge:     91, //days
	}, mux)

	http.ListenAndServe(":8080", logger)
}
