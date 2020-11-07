// table writes csv type formatted input in tabular format

package main

import (
	"bufio"
	"encoding/csv"
	goflag "flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	flag "github.com/spf13/pflag"
)

func parseArgs() rune {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: table --delimiter ',' < file.csv\n")
		flag.PrintDefaults()
	}
	delimiterPtr := flag.String("delimiter", ",", "delimiter between cells, must be one character or '\\t' for tab separated values")
	flag.Parse()
	delimiterRune, _ := utf8.DecodeRuneInString(*delimiterPtr) // ignore rune size, as we check it above
	delimiterIsTab := *delimiterPtr == `\t`
	if delimiterIsTab {
		delimiterRune = '\t'
	}
	if !delimiterIsTab && len(*delimiterPtr) != 1 {
		fmt.Fprintln(os.Stderr, "delimiter must be one character")
		flag.Usage()
		os.Exit(1)
	}
	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "invalid argument(s):", flag.Args())
		flag.Usage()
		os.Exit(1)
	}
	return delimiterRune
}

func lengths(records [][]string) []int {
	columnLengths := make([]int, len(records[0]))
	for _, record := range records {
		for j, cell := range record {
			cellLen := len(strings.TrimSpace(cell))
			if columnLengths[j] < cellLen {
				columnLengths[j] = cellLen
			}
		}
	}
	return columnLengths
}

func makeRows(records [][]string, columnLengths []int) []string {
	rows := make([]string, len(records))
	for i, record := range records {
		for j, cell := range record {
			record[j] = fmt.Sprintf("%-*s", columnLengths[j], strings.TrimSpace(cell))
		}
		rows[i] = strings.TrimSpace(strings.Join(record, "  "))
	}
	return rows
}

func write(w *bufio.Writer, rows []string) {
	allRows := strings.Join(rows, "\n")
	w.WriteString(allRows)
	w.WriteString("\n")
	w.Flush()
}

func main() {
	delimiterRune := parseArgs()
	reader := csv.NewReader(os.Stdin)
	reader.Comma = delimiterRune
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("table: %s", err)
	}
	columnLengths := lengths(records)
	rows := makeRows(records, columnLengths)
	writer := bufio.NewWriter(os.Stdout)
	write(writer, rows)
}
