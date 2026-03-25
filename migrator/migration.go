package migrator

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
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
	err        error
}

// NewLoader create a new Loader instance
func NewLoader(filter FilterFn) Loader {
	return Loader{
		migrations: make([]Migration, 0),
		filter:     filter,
		err:        nil,
	}
}

// FromDirectory read the specified directory for migrations and load them into the internal buffer
func (l *Loader) FromDirectory(dir string) *Loader {
	dirFS := os.DirFS(dir)
	return l.FromFS(dirFS)
}

// FromFS loads the migrations from an embedded filesystem, does not recurse into subdirectories
func (l *Loader) FromFS(afs fs.FS) *Loader {
	if l.err != nil {
		return l
	}

	files, err := fs.ReadDir(afs, ".")
	if err != nil {
		l.err = fmt.Errorf("failed to read directory: %w", err)
		return l
	}

	if len(files) == 0 {
		return l
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
				l.err = fmt.Errorf("failed to read migration file (%s): %w", filename, err)
				return l
			}
			m.Content = string(migrationSQL)

			l.migrations = append(l.migrations, *m)
		}
	}

	return l
}

// Result return the loaded migrations and an error if it occured
func (l *Loader) Result() ([]Migration, error) {
	// Sort the migrations by their version
	sort.Slice(l.migrations, func(i, j int) bool {
		return l.migrations[i].Version < l.migrations[j].Version
	})

	return l.migrations, l.err
}
