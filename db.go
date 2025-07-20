package main

import (
	"database/sql"
	//"fmt"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", PostgresConnString)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		amount NUMERIC(10,2) NOT NULL,
		category TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

func AddExpense(userID int64, amount float64, category string) error {
	query := `INSERT INTO expenses (user_id, amount, category) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, userID, amount, category)
	return err
}

func GetTodayTotal(userID int64) (float64, error) {
	query := `
	SELECT COALESCE(SUM(amount), 0)
	FROM expenses
	WHERE user_id = $1 AND created_at::date = CURRENT_DATE;
	`
	row := db.QueryRow(query, userID)
	var total float64
	err := row.Scan(&total)
	return total, err
}

func GetWeeklyStats(userID int64) (map[string]float64, error) {
	query := `
		SELECT to_char(created_at::date, 'YYYY-MM-DD') as day, SUM(amount)
		FROM expenses
		WHERE user_id = $1 AND created_at >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY day
		ORDER BY day;
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]float64)
	for rows.Next() {
		var date string
		var amount float64
		err := rows.Scan(&date, &amount)
		if err != nil {
			return nil, err
		}
		stats[date] = amount
	}
	return stats, nil
}

