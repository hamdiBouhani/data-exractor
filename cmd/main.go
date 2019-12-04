package main

import (
	"fmt"
	"os"

	"sm-connector-be/svc/cmd/consumer"
	"sm-connector-be/svc/cmd/rabbit"
	"sm-connector-be/svc/cmd/start"
	"sm-connector-be/svc/cmd/ws"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "sm-connector",
		Short: "utilities and services",
		Long:  "Top level command for utilities and services of the sm-connecter-be app",
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		start.Cmd,
		consumer.Cmd,
		ws.Cmd,
		rabbit.Cmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
