package printer

import (
	"bufio"
	"gb_go_best_practics/homework/sqlcsv/call"
	"gb_go_best_practics/homework/sqlcsv/crypto"
	"os"
)

func Print(input string) error {
	var output string
	res, err := call.Execute(input)
	if err != nil {
		return err
	}

	for _, r := range res {
		s := crypto.EncodeCSV(r)
		output += s
	}

	writer := bufio.NewWriter(os.Stdout)
	// Error return value of `writer.WriteString` is not checked (errcheck)
	_, err = writer.WriteString(output + "\n")
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
