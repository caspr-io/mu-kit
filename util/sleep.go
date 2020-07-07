package util

import "time"

// Sleep sleeps for an input number of seconds
func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
