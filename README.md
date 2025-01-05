# Cash Register Project

This project implements a simple cash register system in Go. It allows users to add and remove items, as well as calculate the total amount for the items in the register.

## Project Structure

```
cash-register
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── register
│   │   ├── register.go  # Register logic and methods
│   │   └── register_test.go # Unit tests for the Register
│   └── models
│       └── item.go      # Item struct definition
├── go.mod                # Module definition and dependencies
└── README.md             # Project documentation
```

## Setup Instructions

1. Clone the repository:
   ```

   git clone <repository-url>
   cd cash-register
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

## Usage Examples

- To add an item:
  ```go
  register.AddItem(Item{ID: "1", Name: "Apple", Price: 0.99})
  ```

- To remove an item:
  ```go
  register.RemoveItem("1")
  ```

- To calculate the total:
  ```go
  total := register.CalculateTotal()
  ```

## Testing

To run the tests, use the following command:
```
go test ./internal/register
```

## License

This project is licensed under the MIT License.