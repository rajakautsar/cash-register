package register

import (
    "cash-register/internal/models"
    "database/sql"
    "time"
)

type Register struct {
    db *sql.DB
}

func NewRegister(db *sql.DB) *Register {
    return &Register{db: db}
}

func (r *Register) SaveUserToDB(user models.User) error {
    query := "INSERT INTO users (id, username, password, role) VALUES (?, ?, ?, ?)"
    _, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.Role)
    return err
}

func (r *Register) AddItem(item models.Item) {
    // Implementasi untuk menambahkan item ke register
}

func (r *Register) SaveItemToDB(item models.Item) error {
    query := "INSERT INTO items (id, name, price) VALUES (?, ?, ?)"
    _, err := r.db.Exec(query, item.ID, item.Name, item.Price)
    return err
}

func (r *Register) RemoveItem(itemID string) {
    // Implementasi untuk menghapus item dari register
}

func (r *Register) CalculateTotal() float64 {
    // Implementasi untuk menghitung total
    return 0.0
}

func (r *Register) SaveSaleToDB(sale models.Sale) error {
    query := "INSERT INTO sales (id, product_id, quantity, total, date) VALUES (?, ?, ?, ?, ?)"
    _, err := r.db.Exec(query, sale.ID, sale.ProductID, sale.Quantity, sale.Total, sale.Date)
    return err
}

func (r *Register) GetSalesByDate(startDate, endDate time.Time) ([]models.Sale, error) {
    query := "SELECT id, product_id, quantity, total, date FROM sales WHERE date BETWEEN ? AND ?"
    rows, err := r.db.Query(query, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sales []models.Sale
    for rows.Next() {
        var sale models.Sale
        var dateStr string
        err := rows.Scan(&sale.ID, &sale.ProductID, &sale.Quantity, &sale.Total, &dateStr)
        if err != nil {
            return nil, err
        }
        sale.Date, err = time.Parse("2006-01-02 15:04:05", dateStr)
        if err != nil {
            return nil, err
        }
        sales = append(sales, sale)
    }
    return sales, nil
}

func (r *Register) AddSale(sale models.Sale) error {
    // Add sale to register
    // Save sale to database
    return r.SaveSaleToDB(sale)
}