package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

// runNativeREPL runs a basic interactive shell using the provided database connection
func runNativeREPL(db *sql.DB, dbName string, promptLabel string) {
	fmt.Printf("Connected to database '%s' (Native Mode)\n", dbName)
	fmt.Println("Type SQL commands ending with semicolon ';'. Type 'exit', 'quit' or '\\q' to exit.")
	fmt.Println("--------------------------------------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	var queryBuffer strings.Builder
	prompt := fmt.Sprintf("%s> ", promptLabel)

	fmt.Printf("%s> ", promptLabel)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Handle exit commands
		if trimmed == "exit" || trimmed == "quit" || trimmed == "\\q" {
			fmt.Println("Bye")
			break
		}

		if trimmed == "" {
			if queryBuffer.Len() == 0 {
				fmt.Print(prompt)
			} else {
				fmt.Printf("%s-> ", strings.Repeat(" ", len(promptLabel)))
			}
			continue
		}

		// Append to buffer
		queryBuffer.WriteString(line)
		queryBuffer.WriteString(" ")

		// Check for terminator
		if strings.HasSuffix(trimmed, ";") {
			fullQuery := queryBuffer.String()
			queryBuffer.Reset()

			// Simple timing
			start := time.Now()
			executeNativeQuery(db, fullQuery)
			duration := time.Since(start)
			fmt.Printf("(Time: %v)\n", duration)

			fmt.Print(prompt)
		} else {
			// Continuation prompt
			fmt.Printf("%s-> ", strings.Repeat(" ", len(promptLabel)))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

// executeNativeQuery executes a single query and prints results in a table format
func executeNativeQuery(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer rows.Close()

	// get column names
	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("Error getting columns: %v\n", err)
		return
	}

	// Prepare result buffer
	// We'll read all rows to format them nicely (simple tabwriter)
	// For huge result sets, this simple REPL might not be ideal, but good for admin tasks.

	// Using tabwriter for alignment
	// (Note: standard tabwriter doesn't always handle wide columns perfectly but better than raw tabs)
	// We will just print naive tab separated for now or use the logic from connect.go

	// Reuse the printing logic from connect.go (we will move it here or duplicate slightly for REPL independence)

	fmt.Println("--------------------------------------------------")
	for _, c := range cols {
		fmt.Printf("%s\t", c)
	}
	fmt.Println("\n--------------------------------------------------")

	rowValues := make([]interface{}, len(cols))
	rowPointers := make([]interface{}, len(cols))
	for i := range rowValues {
		rowPointers[i] = &rowValues[i]
	}

	rowCount := 0
	for rows.Next() {
		rowCount++
		err := rows.Scan(rowPointers...)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		for _, val := range rowValues {
			if b, ok := val.([]byte); ok {
				fmt.Printf("%s\t", string(b))
			} else {
				fmt.Printf("%v\t", val)
			}
		}
		fmt.Println("")
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("(%d rows)\n", rowCount)
}
