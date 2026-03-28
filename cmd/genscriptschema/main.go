// Command genscriptschema generates JSON Schema from Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"

	"github.com/plexusone/w3pilot/script"
)

func main() {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	schema := r.Reflect(&script.Script{})
	schema.ID = "https://github.com/plexusone/w3pilot/script/w3pilot-script.schema.json"
	schema.Title = "WebPilot Test Script"
	schema.Description = "Schema for WebPilot browser automation test scripts"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
