package db

import (
	"errors"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func runMigrations(uri string) {
	m, err := migrate.New(migrationFilesURL(), uri)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No database changes to apply")
		} else {
			log.Fatal(err)
		}
	}
}

func migrationFilesURL() string {
	return fmt.Sprintf("file://%s", filepath.Join(conf.ProjectRoot, "messenger", "db", "migrations"))
}
