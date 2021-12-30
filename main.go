package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Target struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"target"`
	Source []string `json:"source"`
}

func loadConfiguration(file string) (config Config) {
	configFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return
}

var flagConfig = flag.String("c", "config.json", "Define the path to config.json file.")

func main() {
	flag.Parse()
	config := loadConfiguration(*flagConfig)

	// Add for host key validation if you have it in known_hosts
	// hostKey := getHostKey(config.Target.Host)

	clientConfig := &ssh.ClientConfig{
		User: config.Target.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Target.Password),
		},
		// Remove for host key validation if you have it in known_hosts
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Add for host key validation if you have it in known_hosts
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	conn, err := ssh.Dial("tcp", config.Target.Host+":"+config.Target.Port, clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create new sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for _, f := range config.Source {
		_, filename := filepath.Split(f)

		client.Remove(filename)
		targetFile, err := client.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer targetFile.Close()

		srcFile, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer srcFile.Close()

		// not sure if this overrides existing content in target file
		// so client.Remove() above
		bytes, err := io.Copy(targetFile, srcFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d bytes copied\n", bytes)
	}
}

func getHostKey(host string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
