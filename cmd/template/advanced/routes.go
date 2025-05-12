package advanced

import (
	_ "embed"
)

//go:embed files/websocket/imports/standard_library.tmpl
var stdLibWebsocketImports []byte

func StdLibWebsocketTemplImportsTemplate() []byte {
	return stdLibWebsocketImports
}
