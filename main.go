package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func getEnv(key string) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, _ := strconv.Atoi(value); ok {
			return i
		}
	}
	fmt.Println("Number of iterations not set - 'export NUM_RAND=5'")
	os.Exit(0)
	return 0
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func CountDigits(i int) (count int) {
	for i != 0 {

		i /= 10
		count = count + 1
	}
	return count
}

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		stopsecs := getEnv("STOP_SECS")
		sig := <-sigs
		for j := stopsecs; j > 0; j-- {
			fmt.Printf("Stopping in %d seconds\n", j)
			time.Sleep(time.Second)
		}
		fmt.Println()
		fmt.Println(sig)
		os.Exit(127)
	}()

	iterations := getEnv("NUM_RAND")
	for i := 0; i < iterations; i++ {
		key, err := GenerateRandomString(32)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%*d : %s\n", CountDigits(iterations-1), i, key)
		time.Sleep(time.Second)
	}
}
