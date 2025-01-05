package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("Starting Cash Register System...")
    
    // Initialize your cash register system here

    // Example: Start a simple HTTP server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to the Cash Register System!")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}