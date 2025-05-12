package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Database string

// The current databases supported by Go Symphony
const (
	PostgreSQL Database = "postgresql"
	None       Database = "none"
)

var AllowedDBDrivers = []string{string(PostgreSQL), string(None)}

func (d Database) String() string {
	return string(d)
}

func (d *Database) Type() string {
	return "Database"
}

func (d *Database) Set(value string) error {
	if !slices.Contains(AllowedDBDrivers, value) {
		return fmt.Errorf("Database driver to use. Allowed values: %s", strings.Join(AllowedDBDrivers, ", "))
	}
	// Explicit type assertion to convert string to Database
	*d = Database(value)
	return nil
}
