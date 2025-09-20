package data

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB
const dbPath = "./internal/data/repo.db"

type RepoRow struct {
	ID         int    `json:"id"`
	AccessCode string `json:"access_code"`
	Name       string `json:"name"`
}

func AddRepo(name string) (string, error) {

	accessCode := GenerateAccessCode()

	for range 10 {
		row := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM repo WHERE access_code = ?", accessCode)
		var count int

		err := row.Scan(&count)
		if err != nil {
			count = 1
		}

		if count == 0 {
			break
		}

		accessCode = GenerateAccessCode()
	}

	_, err := db.ExecContext(
		context.Background(),
		"INSERT INTO repo (access_code, name) VALUES (?, ?)",
		accessCode, name,
	)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return accessCode, nil
}

func QueryRepoByCode(code string) (RepoRow, error) {
	row := db.QueryRowContext(context.Background(),`SELECT * FROM repo WHERE access_code=?`, code,)
	var repo RepoRow
	err := row.Scan(&repo.ID, &repo.AccessCode, &repo.Name)

	if err != nil {
		fmt.Println(err)
		return RepoRow{}, nil
	}
	return repo, nil
}

func Init() error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS repo (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			access_code TEXT NOT NULL, 
			name TEXT NOT NULL
		)`,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}