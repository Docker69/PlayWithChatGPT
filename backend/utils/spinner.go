package utils

import (
	"fmt"
	"time"
)

func Spinner(done chan bool) {
	fmt.Print("   Waiting for operation to complete ...")
	symbols := []rune("|/-\\") // Array of ASCII characters for the spinner
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	start := time.Now()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if time.Since(start) > 60*time.Second { // Stop the spinner after 60 seconds
				fmt.Println("\rTimeout!")
				close(done)
				return
			}
			for _, s := range symbols { // Iterate over the ascii characters and print them one by one
				fmt.Printf("\r%c ", s) // Use carriage return to replace the previous character with the new one
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
