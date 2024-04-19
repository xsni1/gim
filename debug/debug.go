package debug

import (
	"fmt"
	"os"
)

func Debug(s string) {
	// >>>
	// no idea how to debug this shit lmao
	f, err := os.OpenFile("/dev/pts/8", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err opening dev", err)
		os.Exit(1)
	}
	defer f.Close()
    f.Write([]byte(s))
	// >>>
}
