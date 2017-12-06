package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func usage() {
	fmt.Println("app <action> <params> ...")
}

var db *sql.DB

func createSchema() error {
	var err error
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY NOT NULL,
	name VARCHAR(100),
	hobby VARCHAR(100)
)`)
	return err
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	action := os.Args[1]

	var err error
	db, err = sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal(err)
	}

	err = createSchema()
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "create":
		if len(os.Args) < 4 {
			usage()
			os.Exit(1)
		}
		name := os.Args[2]
		hobby := os.Args[3]
		res, err := db.Exec(`INSERT INTO users(name, hobby) values(?, ?)`, name, hobby)
		if err != nil {
			log.Fatal(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User '%s' created, id: %d\n", name, id)

	case "read":
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}
		idString := os.Args[2]
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Fatal(err)
		}

		var name, hobby string

		err = db.QueryRow(`SELECT name, hobby FROM users WHERE id = ?`, id).Scan(&name, &hobby)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("* %d %s: %s\n", id, name, hobby)
	case "update":
		if len(os.Args) < 4 {
			usage()
			os.Exit(1)
		}
		idString := os.Args[2]
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Fatal(err)
		}
		hobby := os.Args[3]

		_, err = db.Exec(`UPDATE users SET hobby = ? WHERE id = ?`, hobby, id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User %d updated.\n", id)

	case "delete":
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}
		idString := os.Args[2]
		id, err := strconv.Atoi(idString)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User %d deleted.\n", id)
	case "list":
		if len(os.Args) < 2 {
			usage()
			os.Exit(1)
		}

		rows, err := db.Query(`SELECT id, name, hobby FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var id int
		var name, hobby string
		for rows.Next() {
			err = rows.Scan(&id, &name, &hobby)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("%d %s: %s\n", id, name, hobby)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
