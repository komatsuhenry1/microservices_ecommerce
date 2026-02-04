package utils

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
)

func RunMigrations(db *sql.DB, migrationsFS embed.FS) error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		fmt.Println("Running migration:", name)

		b, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if _, err := db.Exec(string(b)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}
	}

	return nil
}
