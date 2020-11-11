// table writes csv type formatted input in tabular format

package main

import (
	"encoding/csv"
	goflag "flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
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
	delimiterRune := parseArgs()
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "table: invalid argument(s): '%s'\n", strings.Join(flag.Args(), " "))
		os.Exit(1)
	}
	reader := csv.NewReader(os.Stdin)
	reader.Comma = delimiterRune
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "table: %v\n", err)
	}
	write(os.Stdout, records)
}
