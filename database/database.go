package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// User struct represents a user in the database
type Shape struct {
	ID     int
	Width  float32
	Height float32
	Area   *float32
}

// Init initializes the SQLite database
func Init(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	fmt.Println("Database initialized.")
}

// GetUsers retrieves all users from the database
func GetShapes() {
	query := `SELECT id, width, height, area FROM shape`
	rows, err := DB.Query(query)
	if err != nil {
		log.Fatalf("Error querying users: %v", err)
	}
	defer rows.Close()

	fmt.Println("Shapes in the database:")
	for rows.Next() {
		var id int
		var width float32
		var height float32
		var area *float32 = nil
		err := rows.Scan(&id, &width, &height, &area)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		fmt.Printf("ID: %d, Width: %.2f, Height: %.2f\n", id, width, height)
	}
}

func UpdateShapeArea(id int) {
	updateQuery := `UPDATE shape SET area = width * height WHERE id = ?`
	_, err := DB.Exec(updateQuery, id)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}
	fmt.Printf("Updated area for shape with id  %d\n", id)
}

func GetShapeByID(id int) (*Shape, error) {
	query := "SELECT id, width, height, area FROM shape WHERE id = ?"

	var shape Shape

	// QueryRow for a single result
	err := DB.QueryRow(query, id).Scan(&shape.ID, &shape.Width, &shape.Height, &shape.Area)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no shape found with id %d", id)
	} else if err != nil {
		return nil, err
	}

	return &shape, nil
}
