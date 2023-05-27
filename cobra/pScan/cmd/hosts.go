/*
Copyright © 2023 Pere Picó Muntaner
Copyrights apply to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage tho hosts list",
	Long: `Manages the hosts lists for pScan
Add hosts with the add command
Delete hosts with the delete command
List hosts with the list command`,
}

func init() {
	rootCmd.AddCommand(hostsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
