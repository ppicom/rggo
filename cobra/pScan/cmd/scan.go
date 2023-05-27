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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		portsStr, err := cmd.Flags().GetStringSlice("ports")
		if err != nil {
			return err
		}

		udp, err := cmd.Flags().GetBool("udp")
		if err != nil {
			return err
		}

		open, err := cmd.Flags().GetBool("open")
		if err != nil {
			return err
		}

		closed, err := cmd.Flags().GetBool("closed")
		if err != nil {
			return err
		}

		ttl, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return err
		}

		if open && closed {
			return fmt.Errorf("cannot use --open and --closed at the same time")
		}

		ports := []int{}

		for _, s := range portsStr {
			if r, ok := ran(s); ok {
				ports = append(ports, getAll(r)...)
				continue
			}

			p, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			ports = append(ports, p)
		}

		return scanAction(os.Stdout, hostsFile, ports, udp, open, closed, ttl)
	},
}

func ran(s string) ([]int, bool) {
	slice := strings.Split(s, "-")

	if len(slice) != 2 {
		return []int{}, false
	}

	fStr, lStr := slice[0], slice[1]

	f, err := strconv.Atoi(fStr)
	if err != nil {
		return []int{}, false
	}

	l, err := strconv.Atoi(lStr)
	if err != nil {
		return []int{}, false
	}

	if f > l {
		return []int{}, false
	}

	return []int{f, l}, true
}

func getAll(ran []int) []int {
	ports := []int{}

	for i := ran[0]; i <= ran[1]; i++ {
		ports = append(ports, i)
	}

	return ports
}

func scanAction(out io.Writer, hostsFile string, ports []int, udp, open, closed bool, ttl time.Duration) error {
	if err := validate(ports, udp); err != nil {
		return err
	}

	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results, _ := scan.Run(hl, ports, udp, open, closed, ttl)

	return printResults(out, results)
}

func validate(ports []int, udp bool) error {
	if udp {
		return nil
	}
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("%d: port is out of range", p)
		}
	}

	return nil
}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if r.NotFound {
			message += " Host not found\n\n"
			continue
		}

		message += fmt.Sprintln()

		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprint(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringSliceP("ports", "p", []string{"22", "80", "443"}, "ports to scan")
	scanCmd.Flags().Bool("udp", false, "use UDP instead of TCP")
	scanCmd.Flags().Bool("open", false, "show only open ports")
	scanCmd.Flags().Bool("closed", false, "show only closed ports")
	scanCmd.Flags().DurationP("timeout", "t", 0, "set a timeout for the scan")
}
