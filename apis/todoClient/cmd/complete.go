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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:          "complete <id>",
	Short:        "Mark an item as completed",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")

		return completeAction(os.Stdout, apiRoot, args[0])
	},
}

func completeAction(out io.Writer, apiRoot, arg string) error {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("%w: item id must be a number", ErrNotNumber)
	}

	if err := completeItem(apiRoot, id); err != nil {
		return err
	}

	return printComplete(out, id)
}

func printComplete(out io.Writer, id int) error {
	_, err := fmt.Fprintf(out, "Item number %d marked as completed.\n", id)
	return err
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
