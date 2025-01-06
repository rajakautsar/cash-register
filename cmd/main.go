package main

import (
	"cash-register/internal"
	"cash-register/internal/models"
	"cash-register/internal/register"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var reg *register.Register
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/add-user", authMiddleware(addUserHandler, "admin"))
	http.HandleFunc("/add-item", authMiddleware(addItemHandler, "admin"))
	http.HandleFunc("/remove-item", authMiddleware(removeItemHandler, "admin"))
	http.HandleFunc("/calculate-total", authMiddleware(calculateTotalHandler, "employee"))
	http.HandleFunc("/sales-report", authMiddleware(salesReportHandler, "admin"))
	http.HandleFunc("/add-sale", authMiddleware(addSaleHandler, "admin"))

	// Print the server port
	fmt.Printf("Server is running at http://localhost:%s\n", port)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the Cash Register System API"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Here you should check the credentials against your database
	// For simplicity, we are using hardcoded values
	var user models.User
	if creds.Username == "admin" && creds.Password == "password" {
		user = models.User{ID: "1", Username: "admin", Password: "password", Role: "admin"}
	} else if creds.Username == "employee" && creds.Password == "password" {
		user = models.User{ID: "2", Username: "employee", Password: "password", Role: "employee"}
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(100 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := map[string]string{"message": "Login successful", "token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func authMiddleware(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Role != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate user data
	if user.ID == "" || user.Username == "" || user.Password == "" || user.Role == "" {
		http.Error(w, "Missing user data", http.StatusBadRequest)
		return
	}

	// Save user to database
	err = reg.SaveUserToDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "User added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var item models.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate item data
	if item.ID == "" || item.Name == "" || item.Price <= 0 {
		http.Error(w, "Missing item data", http.StatusBadRequest)
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

func salesReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		http.Error(w, "Missing period parameter", http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	now := time.Now()

	switch period {
	case "daily":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 1)
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	case "yearly":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
	default:
		http.Error(w, "Invalid period", http.StatusBadRequest)
		return
	}

	sales, err := reg.GetSalesByDate(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

func addSaleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var sales []models.Sale
	err := json.NewDecoder(r.Body).Decode(&sales)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, sale := range sales {
		sale.ID = uuid.New().String()
		sale.Date = time.Now()
		err := reg.AddSale(sale)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := map[string]string{"message": "Sales added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
