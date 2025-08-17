package flags

import (
	"fmt"
	"strings"
)

type Database string

// Database options: PostgreSQL and Supabase
const (
	Postgres Database = "postgres"
	Supabase Database = "supabase"
	None     Database = "none"
)

var AllowedDBDrivers = []string{string(Postgres), string(Supabase), string(None)}

func (f Database) String() string {
	return string(f)
}

func (f *Database) Type() string {
	return "Database"
}

func (f *Database) Set(value string) error {
	// Contains isn't available in 1.20 yet
	// if AllowedDBDrivers.Contains(value) {
	for _, database := range AllowedDBDrivers {
		if database == value {
			*f = Database(value)
			return nil
		}
	}

	return fmt.Errorf("Database to use. Allowed values: %s", strings.Join(AllowedDBDrivers, ", "))
}
