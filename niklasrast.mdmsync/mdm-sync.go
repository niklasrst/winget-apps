package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// /help parameter
	if len(os.Args) > 1 && os.Args[1] == "/help" {
		fmt.Println("This tool can initiate a mdm sync for Microsoft Intune managed Windows clients")
		fmt.Println("Usage: mdmsync")
		return
	}

	// Open the registry key
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Provisioning\OMADM\Accounts`, registry.READ)
	if err != nil {
		log.Fatalf("Failed to open registry key: %v", err)
	}
	defer key.Close()

	// Get the subkey names
	accounts, err := key.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatalf("Failed to read subkey names: %v", err)
	}

	// Check if there are no accounts
	if len(accounts) == 0 {
		fmt.Println("This device is not MDM managed.")
		os.Exit(0)
	}

	// Iterate over the accounts and start the process
	for _, account := range accounts {
		cmd := exec.Command("deviceenroller.exe", "/o", account, "/c", "/b")
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to start process for account %s: %v", account, err)
		} else {
			fmt.Printf("Process started successfully for account %s\n", account)
		}
	}
}
