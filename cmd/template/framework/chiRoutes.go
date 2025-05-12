package framework

import (
	_ "embed"

	"github.com/Tomlord1122/go-symphony/cmd/template/advanced"
)

//go:embed files/routes/chi.go.tmpl
var chiRoutesTemplate []byte

//go:embed files/tests/default-test.go.tmpl
var chiTestHandlerTemplate []byte

// ChiTemplates contains the methods used for building
// an app that uses [github.com/go-chi/chi]
type ChiTemplates struct{}

func (c ChiTemplates) Main() []byte {
	return mainTemplate
}

func (c ChiTemplates) Server() []byte {
	return standardServerTemplate
}

func (c ChiTemplates) Routes() []byte {
	return chiRoutesTemplate
}

func (c ChiTemplates) TestHandler() []byte {
	return chiTestHandlerTemplate
}

func (c ChiTemplates) WebsocketImports() []byte {
	return advanced.StdLibWebsocketTemplImportsTemplate()
}
