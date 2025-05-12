package flags

import (
	"fmt"
	"slices"
	"strings"
)

type Database string

// These are all the current databases supported. If you want to add one, you
// can simply copy and paste a line here. Do not forget to also add it into the
// AllowedDBDrivers slice too!
const (
	Postgres Database = "postgres"
	Redis    Database = "redis"
	None     Database = "none"
)

var AllowedDBDrivers = []string{string(Postgres), string(Redis), string(None)}

func (f Database) String() string {
	return string(f)
}

func (f *Database) Type() string {
	return "Database"
}

func (f *Database) Set(value string) error {
	// Contains isn't available in 1.20 yet
	// if AllowedDBDrivers.Contains(value) {

	if slices.Contains(AllowedDBDrivers, value) {
		// Explicit type assertion to convert string to Database
		*f = Database(value)
		return nil
	}
	return fmt.Errorf("Database to use. Allowed values: %s", strings.Join(AllowedDBDrivers, ", "))
}
