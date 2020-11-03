// table writes csv type formatted input in tabular format

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func lengths(records [][]string) []int {
	columnLengths := make([]int, len(records))
	for _, record := range records {
		for j, cell := range record {
			cell := strings.TrimSpace(cell)
			if columnLengths[j] < len(cell) {
				columnLengths[j] = len(cell)
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
	reader := csv.NewReader(os.Stdin)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("table: %s", err)
	}
	columnLengths := lengths(records)
	rows := makeRows(records, columnLengths)
	writer := bufio.NewWriter(os.Stdout)
	write(writer, rows)
}
