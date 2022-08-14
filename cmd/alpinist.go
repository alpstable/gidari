package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var filepath string

func main() {
	cmd := &cobra.Command{
		Use:   "alpinest",
		Short: "alpinest is an ETL executable for storing web data for analysis",
	}
	cmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path to the combination file")

	// cmd.AddCommand(operation.CoinbasePro())
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
