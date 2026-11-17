// Command catalint validates catalog files without building the CLI:
//
//	go run ./internal/catalog/catalint internal/catalog/data/sms.json ...
package main

import (
	"fmt"
	"os"

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
		s, err := catalog.Decode(data)
		if err != nil {
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
