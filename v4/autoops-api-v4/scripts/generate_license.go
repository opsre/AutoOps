package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type LicensePayload struct {
	CustomerName string  `json:"customerName"`
	MachineCode  string  `json:"machineCode"`
	LicenseType  string  `json:"licenseType"`
	IssuedAt     string  `json:"issuedAt"`
	ExpireAt     *string `json:"expireAt"`
}

func normalizePEM(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	return strings.ReplaceAll(value, "\\n", "\n")
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("private key is not valid PEM")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		key, ok := keyAny.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
		return key, nil
	}
	return nil, fmt.Errorf("failed to parse RSA private key")
}

func main() {
	var (
		privateKeyPath = flag.String("private-key", "", "path to RSA private key PEM")
		customerName   = flag.String("customer", "", "customer name")
		machineCode    = flag.String("machine-code", "", "machine code from deployed instance")
		licenseType    = flag.String("license-type", "", "half_year | annual | permanent")
		issuedAt       = flag.String("issued-at", "", "issuedAt time, format 2006-01-02 15:04:05")
		expireAt       = flag.String("expire-at", "", "expireAt time, format 2006-01-02 15:04:05")
	)
	flag.Parse()

	if *privateKeyPath == "" || *customerName == "" || *machineCode == "" || *licenseType == "" {
		fmt.Fprintln(os.Stderr, "missing required flags")
		flag.Usage()
		os.Exit(1)
	}

	switch *licenseType {
	case "half_year", "annual", "permanent":
	default:
		fmt.Fprintln(os.Stderr, "license-type must be one of: half_year, annual, permanent")
		os.Exit(1)
	}

	issuedAtValue := strings.TrimSpace(*issuedAt)
	if issuedAtValue == "" {
		issuedAtValue = time.Now().Format("2006-01-02 15:04:05")
	}
	if _, err := time.ParseInLocation("2006-01-02 15:04:05", issuedAtValue, time.Local); err != nil {
		fmt.Fprintf(os.Stderr, "invalid issued-at: %v\n", err)
		os.Exit(1)
	}

	var expireAtValue *string
	if *licenseType != "permanent" {
		if strings.TrimSpace(*expireAt) == "" {
			fmt.Fprintln(os.Stderr, "expire-at is required for half_year and annual licenses")
			os.Exit(1)
		}
		expireAtRaw := strings.TrimSpace(*expireAt)
		if _, err := time.ParseInLocation("2006-01-02 15:04:05", expireAtRaw, time.Local); err != nil {
			fmt.Fprintf(os.Stderr, "invalid expire-at: %v\n", err)
			os.Exit(1)
		}
		expireAtValue = &expireAtRaw
	}

	privateKey, err := loadPrivateKey(*privateKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load private key failed: %v\n", err)
		os.Exit(1)
	}

	payload := LicensePayload{
		CustomerName: *customerName,
		MachineCode:  strings.TrimSpace(*machineCode),
		LicenseType:  *licenseType,
		IssuedAt:     issuedAtValue,
		ExpireAt:     expireAtValue,
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal payload failed: %v\n", err)
		os.Exit(1)
	}

	hash := sha256.Sum256(payloadRaw)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "sign payload failed: %v\n", err)
		os.Exit(1)
	}

	licenseKey := base64.StdEncoding.EncodeToString(payloadRaw) + "." + base64.StdEncoding.EncodeToString(signature)
	fmt.Println(licenseKey)
}
