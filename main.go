package main

import (
    "fmt"
    "net/http"
    "gocv.io/x/gocv"
)

func main() {
    // Open the default camera device
    webcam, err := gocv.OpenVideoCapture(0)
    if err != nil {
        fmt.Printf("Error opening camera: %v\n", err)
        return
    }
    defer webcam.Close()

    // Start capturing frames from the camera
    frame := gocv.NewMat()
    defer frame.Close()

    // Set up an HTTP handler that streams the video
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Set the response headers
        w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

        // Continuously read frames from the camera and write them to the response
        for {
            if ok := webcam.Read(&frame); !ok {
                fmt.Println("Device closed")
                break
            }

            // Convert the frame to a JPEG image and write it to the response
            img, _ := gocv.IMEncode(".jpg", frame)
            w.Write([]byte("--frame\r\n"))
            w.Header().Set("Content-Type", "image/jpeg")
            w.Header().Set("Content-Length", fmt.Sprint(img.Len()))
            w.Write([]byte("\r\n"))
            w.Write(img.GetBytes())
            w.Write([]byte("\r\n"))
        }
    })

    // Start the HTTP server
    fmt.Println("Listening on port 4000...")
    http.ListenAndServe(":4000", nil)
}
