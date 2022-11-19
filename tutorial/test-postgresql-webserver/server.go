package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var (
	DB          *sql.DB
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_SSLMODE  = os.Getenv("DB_SSLMODE")
)

func init() {
	initDB()
}

func main() {
	app := fiber.New()

	app.Post("/seed", seedHandler)

	app.Get("/select", selectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

// TODO: fix this
func seedHandler(c *fiber.Ctx) error {
	log.Println("Seeding database")
	// create table
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS test (id SERIAL PRIMARY KEY, name TEXT)")
	if err != nil {
		return err
	}

	// insert data
	_, err = DB.Exec("INSERT INTO test (name) VALUES ('test')")
	if err != nil {
		return err
	}

	return c.SendString("seeded")
}
func selectHandler(c *fiber.Ctx) error {
	// select data and return as json
	rows, err := DB.Query("SELECT * FROM test")
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		result = append(result, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}
	return c.JSON(result)
}

func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
