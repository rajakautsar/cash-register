package register

import (
    "cash-register/internal/models"
    "database/sql"
)

type Register struct {
    items map[string]models.Item
    db    *sql.DB
}

func NewRegister(db *sql.DB) *Register {
    return &Register{
        items: make(map[string]models.Item),
        db:    db,
    }
}

func (r *Register) AddItem(item models.Item) {
    r.items[item.ID] = item
    r.items[item.Name]=item
}

func (r *Register) RemoveItem(itemID string) {
    delete(r.items, itemID)
}

func (r *Register) CalculateTotal() float64 {
    total := 0.0
    for _, item := range r.items {
        total += item.Price
    }
    return total
}

func (r *Register) SaveItemToDB(item models.Item) error {
    query := "INSERT INTO items (id, name, price) VALUES (?, ?, ?)"
    _, err := r.db.Exec(query, item.ID, item.Name, item.Price)
    return err
}