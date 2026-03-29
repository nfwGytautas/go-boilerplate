package migrator

import (
	"log/slog"
	"strconv"
	"strings"
)

// filters.go contains prebuilt filters that can be used together with the loader
// for some out of the box functionality

// FilterVersionedSQLFiles filters out all non *.sql files and parses out an id and name from
// the file name with the following format `<id>_<name>.sql`
func FilterVersionedSQLFiles(filename string) *Migration {
	filename, valid := strings.CutSuffix(filename, ".sql")

	// Ignore non sql files
	if !valid {
		slog.Debug(
			"filter_versioned_sql_files",
			"filename", filename,
			"message", "dropped non sql file",
		)
		return nil
	}

	// Find the first underscore to separate the ID from the name
	versionStr, name, valid := strings.Cut(filename, "_")
	if !valid {
		slog.Warn(
			"filter_versioned_sql_files",
			"filename", filename,
			"message", "incorrect format",
		)
		return nil
	}

	// Check if the version is a valid number
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		slog.Warn(
			"filter_versioned_sql_files",
			"filename", filename,
			"message", "version not a number",
		)
		return nil
	}

	return &Migration{
		Version: version,
		Name:    name,
	}
}
