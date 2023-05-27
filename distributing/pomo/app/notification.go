//go:build !containers && !disable_notification
// +build !containers,!disable_notification

package app

import "pragprog.com/rggo/distributing/notify"

func send_notification(msg string) {
	n := notify.New("Pomodoro", msg, notify.SeverityNormal)

	n.Send()
}
