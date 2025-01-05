package register

import "cash-register/internal/models"

type Register struct {
    items map[string]models.Item
}

func NewRegister() *Register {
    return &Register{
        items: make(map[string]models.Item),
    }
}

func (r *Register) AddItem(item models.Item) {
    r.items[item.ID] = item
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