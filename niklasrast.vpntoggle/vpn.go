package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// /help parameter
	if len(os.Args) > 1 && os.Args[1] == "/help" {
		fmt.Println("This tool can toggle VPN on or off")
		fmt.Println("Usage: vpntoggle")
		return
	}

	// Check if VPN services are running
	checkServiceCmd := exec.Command("powershell", "Get-Service", "*vpn*", "|", "Where-Object", "{$_.Status -eq 'Running'}")
	output, err := checkServiceCmd.Output()
	if err != nil {
		fmt.Printf("Error checking services: %v\n", err)
		return
	}

	if strings.TrimSpace(string(output)) == "" {
		// Check if any VPN service exists
		checkServiceExistsCmd := exec.Command("powershell", "Get-Service", "*vpn*")
		output, err := checkServiceExistsCmd.Output()
		if err != nil {
			fmt.Printf("Error checking if VPN services exist: %v\n", err)
			return
		}

		if strings.TrimSpace(string(output)) == "" {
			fmt.Println("No VPN service to toggle.")
			os.Exit(0)
		} else {
			// Start VPN services
			startServiceCmd := exec.Command("powershell", "Get-Service", "*vpn*", "|", "Start-Service")
			err := startServiceCmd.Run()
			if err != nil {
				fmt.Printf("Error starting services: %v\n", err)
			} else {
				fmt.Println("VPN services started successfully.")

				// Launch VPN client
				launchCmd := exec.Command("C:\\Program Files (x86)\\Cisco\\Cisco AnyConnect Secure Mobility Client\\vpnui.exe")
				err := launchCmd.Start()
				if err != nil {
					fmt.Printf("Error launching VPN client: %v\n", err)
				} else {
					fmt.Println("VPN Client launched successfully.")
				}
			}
		}
	} else {
		// Stop VPN services
		stopServiceCmd := exec.Command("powershell", "Get-Service", "*vpn*", "|", "Stop-Service")
		err := stopServiceCmd.Run()
		if err != nil {
			fmt.Printf("Error stopping services: %v\n", err)
		} else {
			fmt.Println("VPN services stopped successfully.")
		}

		// Stop VPN processes
		stopProcessCmd := exec.Command("powershell", "Get-Process", "*vpn*", "|", "Stop-Process")
		err = stopProcessCmd.Run()
		if err != nil {
			fmt.Printf("Error stopping processes: %v\n", err)
		} else {
			fmt.Println("VPN processes stopped successfully.")
		}
	}
}
