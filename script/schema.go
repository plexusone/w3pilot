package script

import _ "embed"

//go:embed w3pilot-script.schema.json
var SchemaJSON []byte

// Schema returns the JSON Schema for WebPilot test scripts.
func Schema() []byte {
	return SchemaJSON
}
