package broker

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func waitForPort(port int, connectChan chan bool) {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			conn.Close()
			connectChan <- true
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func resolveDeviceId(r *http.Request) (string, error) {
	// try to get the device ID from the request headers
	deviceId := r.Header.Get("X-Qbee-Device-Id")

	// if no device ID is provided, we try to use the hostname
	if deviceId == "" {
		// We cannot use the public key digest here as there's a 63 character limit
		// for DNS name labels. We need to check if we can use device UUID instead,
		hostParts := strings.Split(r.Host, ".")
		deviceId = hostParts[0]
	}

	if deviceId == "" {
		return "", fmt.Errorf("no device ID provided")
	}

	return deviceId, nil
}
