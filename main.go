// table writes csv type formatted input in tabular format

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func write(f *os.File, records [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, record := range records {
		for j, cell := range record {
			fmt.Fprintf(w, strings.TrimSpace(cell))
			if j < len(cell)-1 {
				fmt.Fprintf(w, "\t")
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "table: invalid argument(s): '%s'\n", strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}
	reader := csv.NewReader(os.Stdin)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "table: %s", err)
	}
	write(os.Stdout, records)
}
