package util

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateTables(db *sqlx.DB) {
	sqlDir := "../devops"

	// Create migrations table if it doesn't exist
	createMigrationsTable := `
	CREATE TABLE IF NOT EXISTS migrations (
		filename TEXT PRIMARY KEY,
		executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createMigrationsTable); err != nil {
		log.Printf("Error creating migrations table: %v", err)
		return
	}

	// Read executed files list from migrations table
	executedFiles := []string{}
	err := db.Select(&executedFiles, "SELECT filename FROM migrations ORDER BY filename ASC")
	if err != nil {
		log.Printf("Error reading migrations table: %v", err)
		return
	}

	// Get all SQL files sorted by name
	var sqlFiles []string
	err = filepath.WalkDir(sqlDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".sql") {
			sqlFiles = append(sqlFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error reading SQL directory: %v", err)
		return
	}

	// Sort files by name
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Printf("No SQL files found in %s", sqlDir)
		return
	}

	log.Printf("Found %d SQL file(s)", len(sqlFiles))

	executedCount := 0
	skippedCount := 0

	// Execute SQL files that haven't been executed yet
	for _, sqlFile := range sqlFiles {
		filename := filepath.Base(sqlFile)

		if contains(executedFiles, filename) {
			log.Printf("SKIPPED: %s (already executed)", filename)
			skippedCount++
			continue
		}

		log.Printf("EXECUTING: %s", filename)

		// Read and execute SQL file
		content, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Printf("ERROR: Failed to read %s: %v", filename, err)
			continue
		}

		log.Printf("--- Content of %s ---", filename)
		log.Printf("%s", string(content))
		log.Printf("--- End of %s ---", filename)

		// Execute the SQL
		if _, err := db.Exec(string(content)); err != nil {
			log.Printf("ERROR: Failed to execute %s: %v", filename, err)
			continue
		}

		// Mark as executed in migrations table
		if err := markExecuted(db, filename); err != nil {
			log.Printf("ERROR: Failed to mark %s as executed: %v", filename, err)
		} else {
			log.Printf("SUCCESS: %s executed successfully", filename)
			executedCount++
		}
	}

	log.Printf("=== SUMMARY ===")
	log.Printf("Total SQL files found: %d", len(sqlFiles))
	log.Printf("Files executed: %d", executedCount)
	log.Printf("Files skipped: %d", skippedCount)
}

func markExecuted(db *sqlx.DB, filename string) error {
	_, err := db.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
	return err
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
