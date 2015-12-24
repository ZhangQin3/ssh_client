package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

type password string

func (p password) Password(user string) (password string, err error) {
	return string(p), nil
}

func main() {
	server := "10.89.255.1"
	// fmt.Print("Remote host? (Default=localhost): ")
	// server := scanConfig()
	// if server == "" {
	// 	server = "localhost"
	// }

	// fmt.Print("UserName? (Default=root): ")
	// user := scanConfig()
	// if user == "" {
	// 	user = "root"
	// }
	// fmt.Print("Password?: ")
	// password := scanConfig()

	// config := &ssh.ClientConfig{
	// 	User: user,
	// 	Auth: []ssh.AuthMethod{
	// 		// SSHAgent(),
	// 		ssh.Password(password),
	// 	},
	// }

	config := &ssh.ClientConfig{
		User: "ethan",
		Auth: []ssh.AuthMethod{
			// SSHAgent(),
			ssh.Password("xxxxx"),
		},
	}
	conn, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	// session.Stderr = os.Stderr
	in, _ := session.StdinPipe()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("Xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}
	// Determine is running on Windows
	var isWindows = os.PathSeparator == '\\' && os.PathListSeparator == ';'

	// Accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		if isWindows && len(str) > 1 {
			// Due to the windows ends a line with \r\n, so remove \n away
			fmt.Fprint(in, str[:len(str)-1])
		} else {
			fmt.Fprint(in, str)
		}
	}
}

func scanConfig() string {
	config, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	config = strings.TrimSpace(config)
	return config
}
