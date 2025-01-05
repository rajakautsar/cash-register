package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "cash-register/internal"
    "cash-register/internal/register"
    "cash-register/internal/models"
    "github.com/joho/godotenv"
)

var reg *register.Register

func main() {
    fmt.Println("Starting Cash Register System...")

    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize database connection
    err = internal.InitDB()
    if err != nil {
        log.Fatal("Failed to connect to the database!")
    } else {
        fmt.Println("Successfully connected to the database!")
    }

    // Initialize the register
    reg = register.NewRegister(internal.GetDB())

    // Define the port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not specified
    }

    // Define API routes
    http.HandleFunc("/add-item", addItemHandler)
    http.HandleFunc("/remove-item", removeItemHandler)
    http.HandleFunc("/calculate-total", calculateTotalHandler)

    // Print the server port
    fmt.Printf("Server is running at http://localhost:%s\n", port)

    // Start the HTTP server
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func addItemHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var item models.Item
    err := json.NewDecoder(r.Body).Decode(&item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Add item to register
    reg.AddItem(item)

    // Save item to database
    err = reg.SaveItemToDB(item)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"message": "Item added successfully"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func removeItemHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    itemID := r.URL.Query().Get("id")
    if itemID == "" {
        http.Error(w, "Missing item ID", http.StatusBadRequest)
        return
    }
    reg.RemoveItem(itemID)

    response := map[string]string{"message": "Item removed successfully"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func calculateTotalHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    total := reg.CalculateTotal()
    response := map[string]float64{"total": total}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}