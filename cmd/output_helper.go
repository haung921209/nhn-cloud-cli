package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/jmespath/go-jmespath"
	"gopkg.in/yaml.v3"
)

func printOutput(data interface{}) error {
	processed, err := applyJMESPathQuery(data)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		return outputJSON(processed)
	case "yaml":
		return outputYAML(processed)
	default:
		if query != "" {
			return outputJSON(processed)
		}
		return outputTable(processed)
	}
}

func applyJMESPathQuery(data interface{}) (interface{}, error) {
	if query == "" {
		return data, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var jsonInterface interface{}
	if err := json.Unmarshal(jsonData, &jsonInterface); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	result, err := jmespath.Search(query, jsonInterface)
	if err != nil {
		return nil, fmt.Errorf("JMESPath query error: %w", err)
	}

	if result == nil {
		return []interface{}{}, nil
	}

	return result, nil
}

func outputJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func outputYAML(data interface{}) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	return enc.Encode(data)
}

func outputTable(data interface{}) error {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return outputSliceTable(v)
	case reflect.Map:
		return outputMapTable(v)
	case reflect.Struct:
		return outputStructTable(v)
	default:
		fmt.Println(data)
		return nil
	}
}

func outputSliceTable(v reflect.Value) error {
	if v.Len() == 0 {
		fmt.Println("No results found.")
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Interface {
		first = first.Elem()
	}

	if first.Kind() == reflect.Map {
		return outputMapSliceTable(v)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	if first.Kind() == reflect.Struct {
		headers := getStructHeaders(first.Type())
		fmt.Fprintln(w, strings.Join(headers, "\t"))

		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Interface {
				item = item.Elem()
			}
			values := getStructValues(item)
			fmt.Fprintln(w, strings.Join(values, "\t"))
		}
	} else {
		for i := 0; i < v.Len(); i++ {
			fmt.Fprintln(w, fmt.Sprintf("%v", v.Index(i).Interface()))
		}
	}

	return w.Flush()
}

func outputMapSliceTable(v reflect.Value) error {
	if v.Len() == 0 {
		fmt.Println("No results found.")
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Interface {
		first = first.Elem()
	}

	var headers []string
	for _, key := range first.MapKeys() {
		headers = append(headers, fmt.Sprintf("%v", key.Interface()))
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Interface {
			item = item.Elem()
		}
		var values []string
		for _, key := range first.MapKeys() {
			val := item.MapIndex(key)
			if val.IsValid() {
				values = append(values, fmt.Sprintf("%v", val.Interface()))
			} else {
				values = append(values, "")
			}
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	return w.Flush()
}

func outputMapTable(v reflect.Value) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")

	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		fmt.Fprintf(w, "%v\t%v\n", key.Interface(), val.Interface())
	}

	return w.Flush()
}

func outputStructTable(v reflect.Value) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		name = strings.Split(name, ",")[0]

		val := v.Field(i)
		fmt.Fprintf(w, "%s:\t%v\n", name, val.Interface())
	}

	return w.Flush()
}

func getStructHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		name = strings.Split(name, ",")[0]
		headers = append(headers, strings.ToUpper(name))
	}
	return headers
}

func getStructValues(v reflect.Value) []string {
	var values []string
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		val := v.Field(i)
		str := fmt.Sprintf("%v", val.Interface())
		if len(str) > 50 {
			str = str[:47] + "..."
		}
		values = append(values, str)
	}
	return values
}
