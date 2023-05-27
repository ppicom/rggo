/*
Copyright © 2023 Pere Picó Muntaner
Copyrights apply to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostn>",
	Aliases:      []string{"d"},
	Short:        "Delete host(s) from the list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}

		fmt.Fprintln(out, "Deleted host:", h)
	}

	return hl.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
