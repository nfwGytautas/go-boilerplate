package migrator

import (
	"fmt"
	"io/fs"
	"os"
)

// Migration represents a single migration for a database
type Migration struct {
	Version int    // Version of the migration
	Name    string // Optional name of the migration
	Content string // The actual 'code' of the migration i.e. SQL
}

// FilterFn is called by the loader to parse out migration information
// if a non-nil value is return the loader will also read the file contents
type FilterFn func(filename string) *Migration

// Loader provides configurable migration loading functionality
type Loader struct {
	migrations []Migration
	filter     FilterFn
}

// NewLoader create a new Loader instance
func NewLoader(filter FilterFn) Loader {
	return Loader{
		migrations: make([]Migration, 0),
		filter:     filter,
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

		filename := file.Name()

		m := l.filter(filename)
		if m != nil {
			migrationSQL, err := fs.ReadFile(afs, filename)
			if err != nil {
				return fmt.Errorf("failed to read migration file (%s): %w", filename, err)
			}
			m.Content = string(migrationSQL)

			l.migrations = append(l.migrations, *m)
		}
	}

	return nil
}
