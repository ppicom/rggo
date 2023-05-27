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
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// movieCmd represents the movie command
var movieCmd = &cobra.Command{
	Use:          "movie <id>",
	Short:        "See details of a movie",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		APIRoot, APIKey := viper.GetString("api-root"), viper.GetString("api-key")

		return movieAction(os.Stdout, APIRoot, APIKey, args[0])
	},
}

func movieAction(out io.Writer, APIRoot, APIKey string, arg string) error {
	ID, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	m, err := getMovie(APIRoot, APIKey, ID)
	if err != nil {
		return err
	}

	return printMovie(out, m)
}

func printMovie(out io.Writer, m movie) error {
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)
	fmt.Fprintf(w, "%s\t\t%s\n", m.Title, m.Tagline)
	fmt.Fprintf(w, "%s\n", m.Overview)
	return w.Flush()
}

func init() {
	rootCmd.AddCommand(movieCmd)
}
