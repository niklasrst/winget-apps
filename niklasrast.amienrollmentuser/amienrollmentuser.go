package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// /help parameter
	if len(os.Args) > 1 && os.Args[1] == "/help" {
		fmt.Println("This tool checks if the current user is the enrollment user for Intune managed Windows clients")
		fmt.Println("Usage: amienrollmentuser")
		return
	}

	// Get Join ID
	regPath := `SYSTEM\CurrentControlSet\Control\CloudDomainJoin\JoinInfo`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.READ)
	if err != nil {
		log.Println("Error opening registry key:", err)
		exitNoPrimaryUser()
	}
	defer k.Close()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Println("Error reading subkey names:", err)
		exitNoPrimaryUser()
	}
	if len(subKeys) == 0 {
		log.Println("No subkeys found")
		exitNoPrimaryUser()
	}
	joinID := subKeys[0]

	// Get Enrollment UPN
	joinInfoPath := fmt.Sprintf(`%s\%s`, regPath, joinID)
	k, err = registry.OpenKey(registry.LOCAL_MACHINE, joinInfoPath, registry.READ)
	if err != nil {
		log.Println("Error opening join info registry key:", err)
		exitNoPrimaryUser()
	}
	defer k.Close()

	puUser, _, err := k.GetStringValue("UserEmail")
	if err != nil {
		log.Println("Error getting UserEmail value:", err)
		exitNoPrimaryUser()
	}

	puUserSid := ""
	if strings.HasPrefix(puUser, "fooUser@") {
		puUser = getCurrentUser()
		if strings.Contains(puUser, `*\`) {
			puUser = puUser[strings.Index(puUser, `\`)+1:]
		}
		puUserSid = getSid(fmt.Sprintf("*\\%s", puUser))
	} else {
		puUserSid = getSid(fmt.Sprintf("azuread\\%s", puUser))
	}

	currentUser := getCurrentUser()
	currentUserSid := getSid(currentUser)

	if currentUserSid == puUserSid {
		fmt.Printf("%s is the enrollment user\n", currentUser)
	} else {
		fmt.Printf("%s is not the enrollment user\n", currentUser)
	}
}

func getCurrentUser() string {
	out, err := exec.Command("powershell", "-Command", "(Get-Process -IncludeUserName -Name explorer | Select-Object UserName -Unique).UserName").Output()
	if err != nil {
		log.Println("Error getting current user:", err)
		exitNoPrimaryUser()
	}
	return strings.TrimSpace(string(out))
}

func getSid(account string) string {
	out, err := exec.Command("powershell", "-Command", fmt.Sprintf("(New-Object System.Security.Principal.NTAccount('%s')).Translate([System.Security.Principal.SecurityIdentifier]).Value", account)).Output()
	if err != nil {
		log.Println("Error getting SID for account:", err)
		exitNoPrimaryUser()
	}
	return strings.TrimSpace(string(out))
}

func exitNoPrimaryUser() {
	fmt.Println("No primary user on this device")
	os.Exit(0)
}
