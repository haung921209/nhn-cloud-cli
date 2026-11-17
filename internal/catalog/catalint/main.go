// Command catalint validates catalog files without building the CLI, then
// checks every catalog in their directories together (name, alias, endpoint
// and operation clashes):
//
//	go run ./internal/catalog/catalint internal/catalog/data/sms.json ...
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/haung921209/nhn-cloud-cli/internal/catalog"
)

func main() {
	failed := false
	dirs := map[string]bool{}
	for _, path := range os.Args[1:] {
		dirs[filepath.Dir(path)] = true
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
	for dir := range dirs {
		services, err := catalog.LoadFS(os.DirFS(dir), ".")
		if err != nil {
			fmt.Println(dir+":", err)
			failed = true
			continue
		}
		for _, e := range catalog.ValidateAll(services) {
			fmt.Println(dir+":", e)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
