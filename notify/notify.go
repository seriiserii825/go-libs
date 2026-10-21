package notify

import "os/exec"

// Send показывает уведомление на рабочем столе через notify-send.
func Send(title, message string) error {
	return exec.Command("notify-send", "--app-name", title, title, message).Run()
}
