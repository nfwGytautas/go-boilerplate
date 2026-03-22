package migrator

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// Migration represents a single migration for a database
type Migration struct {
	Version int    // Version of the migration
	Name    string // Optional name of the migration
	Content string // The actual 'code' of the migration i.e. SQL
}

// Loader provides configurable migration loading functionality
type Loader struct {
	migrations []Migration
	format     string
}

func NewLoader() Loader {
	return Loader{
		make([]Migration, 0),
		format: "",
	}
}

// FromDirectory read the specified directory for migrations and load them into the internal buffer
func (l *Loader) FromDirectory(dir string) error {
	dirFS := os.DirFS(dir)
	return l.FromFS(dirFS)
}

// FromFS loads the migrations from an embedded filesystem, does not recurse into subdirectories
func (l *Loader) FromFS(afs fs.FS) error {
	files, err := fs.ReadDir(afs, ".")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		m, err := l.parseMigrationFile(afs, file)
		if err != nil {
			return err
		}

		if m != nil {
			migrations = append(migrations, *m)
		}
	}

	return nil
}

func (l *Loader) parseMigrationFile(afs fs.FS, file fs.DirEntry) (m *Migration, err error) {
	// File name format: <id>_<name>.sql
	fileName := file.Name()

	// Ignore non sql files
	if !strings.HasSuffix(fileName, ".sql") {
		return
	}

	// Ignore fixture files
	if strings.HasSuffix(fileName, ".fixture.sql") {
		return
	}

	fileName = strings.TrimSuffix(fileName, ".sql")

	m = &Migration{}

	// Find the first underscore to separate the ID from the name
	underscoreIndex := strings.Index(fileName, "_")
	if underscoreIndex == -1 {
		return m, fmt.Errorf("invalid migration file name format for file: %s, (the format is <id>_<name>.sql)", fileName)
	}

	versionStr := fileName[:underscoreIndex]
	m.Name = fileName[underscoreIndex+1:]

	m.Version, err = strconv.Atoi(versionStr)
	if err != nil {
		return m, fmt.Errorf("failed to parse migration file version (%s): %w", fileName, err)
	}

	if m.Version <= 0 {
		return m, fmt.Errorf("migration file version cannot be <= 0: %s", fileName)
	}

	migrationSQL, err := fs.ReadFile(afs, fileName+".sql")
	if err != nil {
		return m, fmt.Errorf("failed to read migration file (%s): %w", fileName, err)
	}
	m.Content = string(migrationSQL)

	return m, nil
}
