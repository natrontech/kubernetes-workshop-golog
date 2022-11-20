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
	seeded      = false
)

type Test struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	if !seeded {
		_, err = DB.Exec("DROP TABLE IF EXISTS test")
		if err != nil {
			log.Fatal(err)
		}

		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS test (id SERIAL, name VARCHAR(255))")
		if err != nil {
			log.Fatal(err)
		}

		seedData := []Test{
			{Name: "test1"},
			{Name: "test2"},
			{Name: "test3"},
		}

		for _, v := range seedData {
			_, err = DB.Exec("INSERT INTO test (name) VALUES ($1)", v.Name)
			if err != nil {
				log.Fatal(err)
			}
		}

		seeded = true
	}
}

func main() {
	app := fiber.New()

	app.Get("/select", selectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
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

func initDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return err
	}
	return nil
}
