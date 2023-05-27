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
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// trendingCmd represents the trending command
var trendingCmd = &cobra.Command{
	Use:          "trending",
	Short:        "List the trending movies of the week",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot, apiKey := viper.GetString("api-root"), viper.GetString("api-key")

		return listTrendingAction(os.Stdout, apiRoot, apiKey)
	},
}

func listTrendingAction(out io.Writer, apiRoot, apiKey string) error {
	results, err := getTrending(apiRoot, apiKey)
	if err != nil {
		return err
	}

	return printTrending(out, results)
}

func printTrending(out io.Writer, results []trending) error {
	w := tabwriter.NewWriter(out, 14, 2, 0, ' ', 0)

	for i, r := range results {
		fmt.Fprintf(w, "%d.\t%s\t\tID: %d\n", i+1, r.OriginalTitle, r.ID)
		fmt.Fprintf(w, "\t%s\t%f\n", "Vote average:", r.VoteAverage)
	}

	return w.Flush()
}

func init() {
	rootCmd.AddCommand(trendingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trendingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trendingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
