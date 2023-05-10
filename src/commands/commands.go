package commands

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"
)

func SetService(db *sql.DB, userID int64, args string, expirationTime time.Time) error {
	parts := strings.Fields(args)
	if len(parts) != 3 {
		return fmt.Errorf("Usage: /set <name> <login> <password>")
	}

	expirationTime = time.Now().Add(1 * time.Minute)
	
	name := parts[0]
	login := parts[1]
	password := parts[2]

	passwordHash := sha256.Sum256([]byte(password))
	passwordHashStr := hex.EncodeToString(passwordHash[:])

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM service WHERE name = $1 AND user_id =$2`, name, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("Failed to check for existing service: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("Service with name %q already exists", name)
	}

	_, err = db.Exec(`INSERT INTO service (user_id, name, login, password, hash, expiration_time) VALUES ($1, $2 ,$3, $4, $5, $6)`, userID, name, login, password, passwordHashStr, expirationTime)
	if err != nil {
		return fmt.Errorf("Failed to insert service: %v", err)
	}

	return nil
}

func GetService(db *sql.DB, userID int64, args string) (login string, password string, err error) {
	parts := strings.Fields(args)
	if len(parts) != 1 {
		return "", "", fmt.Errorf("Usage: /get <name>")
	}

	name := parts[0]

	err = db.QueryRow(`SELECT login, password FROM service WHERE name = $1 AND user_id=$2`, name, userID).Scan(&login, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("Service with name %q not found", name)
		}
		return "", "", fmt.Errorf("Failed to get service: %v", err)
	}
	return login, password, nil
}

func DeleteService(db *sql.DB, userID int64, args string) error {
	name := strings.TrimSpace(args)
	if name == "" {
		return fmt.Errorf("Usage: /delete <name>")
	}

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM service WHERE name = $1 AND user_id = $2`, name, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("Failed to check for existing service: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("Service with name %q not found", name)
	}

	_, err = db.Exec(`DELETE FROM service WHERE name = $1 AND user_id = $2`, name, userID)
	if err != nil {
		return fmt.Errorf("Failed to delete service: %v", err)
	}

	return nil
}

func DeleteExpiredServices(db *sql.DB) {
	rows, err := db.Query(`SELECT id, user_id, name FROM service WHERE expiration_time < NOW()`)
	if err != nil {
		log.Printf("Failed to query expired services: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var userID int64
		var name string

		err := rows.Scan(&id, &userID, &name)
		if err != nil {
			log.Printf("Failed to scan expired service: %v", err)
			continue
		}

		err = DeleteService(db, userID, name)
		if err != nil {
			log.Printf("Failed to delete expired service: %v", err)
			continue
		}

		log.Printf("Deleted expired service: %v", name)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate over expired services: %v", err)
	}
}
