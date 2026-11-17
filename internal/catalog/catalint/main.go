// Command catalint validates catalog files without building the CLI:
//
//	go run ./internal/catalog/catalint internal/catalog/data/sms.json ...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
)

func main() {
	failed := false
	for _, path := range os.Args[1:] {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(path+":", err)
			failed = true
			continue
		}
		var s catalog.Service
		dec := json.NewDecoder(strings.NewReader(string(data)))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&s); err != nil {
			fmt.Println(path+":", err)
			failed = true
			continue
		}
		errs := catalog.Validate(s)
		for _, e := range errs {
			fmt.Println(path+":", e)
		}
		if len(errs) > 0 {
			failed = true
		} else {
			fmt.Printf("%s: ok (%d operations)\n", path, len(s.Operations))
		}
	}
	if failed {
		os.Exit(1)
	}
}
