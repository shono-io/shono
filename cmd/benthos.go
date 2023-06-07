/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/benthosdev/benthos/v4/public/service/servicetest"
	"github.com/spf13/cobra"

	_ "github.com/benthosdev/benthos/v4/public/components/all"
	_ "github.com/shono-io/shono/systems/storage/arangodb"
)

// benthosCmd represents the benthos command
var benthosCmd = &cobra.Command{
	Use:   "benthos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fa := []string{"benthos"}
		fa = append(fa, args...)

		servicetest.RunCLIWithArgs(context.Background(), fa...)
	},
}

func init() {
	rootCmd.AddCommand(benthosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benthosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benthosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
