package main

import (
	"ddns/cmd"
	"github.com/spf13/cobra"
	"os"
)

var Version = "development"

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args, Version)
	cobra.CheckErr(rootCmd.Execute())
}
