package sqlc

import (
	_ "embed"
)

//go:embed files/sqlc.yaml.tmpl
var SqlcConfigTemplate []byte

//go:embed files/001_initial_schema.sql.tmpl
var InitialSchemaTemplate []byte

//go:embed files/users.sql.tmpl
var UsersQueriesTemplate []byte

//go:embed files/posts.sql.tmpl
var PostsQueriesTemplate []byte

// SqlcTemplater provides methods for SQLC-related templates
type SqlcTemplater interface {
	Config() []byte
	InitialSchema() []byte
	UsersQueries() []byte
	PostsQueries() []byte
}

// SqlcTemplate implements SqlcTemplater
type SqlcTemplate struct{}

func (s SqlcTemplate) Config() []byte {
	return SqlcConfigTemplate
}

func (s SqlcTemplate) InitialSchema() []byte {
	return InitialSchemaTemplate
}

func (s SqlcTemplate) UsersQueries() []byte {
	return UsersQueriesTemplate
}

func (s SqlcTemplate) PostsQueries() []byte {
	return PostsQueriesTemplate
}
