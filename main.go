package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"syscall"
	"time"
)

var (
	email          = flag.String("email", "", "Email for the scaleway client")
	chosenOrgName  = flag.String("org", "", "Organization for the scaleway client")
	chosenDuration = flag.Int("duration", 0, "Duration for the API key (in hours, max 8)")
)

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to scw-2fa-init. If you want to discover the command line argument options, use -h or --help.")

	if *email == "" {
		fmt.Print("Enter email: ")
		scanner.Scan()
		*email = scanner.Text()
	}

	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
	}
	password := string(bytePassword)
	fmt.Println()

	fmt.Print("Enter 2FA token: ")
	scanner.Scan()
	token := scanner.Text()

	client, err := NewScalewayClient(*email, password, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	orgs, err := client.ListOrganizations()
	if err != nil {
		fmt.Println(err)
		return
	}

	if *chosenOrgName == "" {
		orgNames := make([]string, 0, len(orgs))
		for name := range orgs {
			orgNames = append(orgNames, name)
		}

		sort.Strings(orgNames) // Optional: sort for consistent ordering

		fmt.Println("Available organizations: ")
		for i, name := range orgNames {
			fmt.Println(i+1, ".", name)
		}

		fmt.Print("Choose an organization by entering the corresponding number: ")
		scanner.Scan()
		chosenNumber, _ := strconv.Atoi(scanner.Text())

		if chosenNumber < 1 || chosenNumber > len(orgNames) {
			fmt.Println("Invalid selection.")
			return
		}

		*chosenOrgName = orgNames[chosenNumber-1]
	}

	if *chosenDuration == 0 {
		fmt.Println("Choose duration for the API key (in hours, max 8): ")
		scanner.Scan()
		*chosenDuration, _ = strconv.Atoi(scanner.Text())

		if *chosenDuration < 1 || *chosenDuration > 8 {
			fmt.Println("Invalid duration, setting to default (1 hour)")
			*chosenDuration = 1 // default option is 1 hour
		}
	}

	if orgKey, ok := orgs[*chosenOrgName]; ok {
		apiKey, err := client.CreateAPIKey(orgKey, time.Duration(*chosenDuration)*time.Hour)
		if err != nil {
			fmt.Println(err)
			return
		}

		scalewayInit(apiKey, orgKey)
	} else {
		fmt.Println("Invalid selection.")
	}
}

const SCW_INIT_COMMAND = "scw config destroy && scw init secret-key=%s access-key=%s organization-id=%s send-telemetry=false install-autocomplete=true"

func scalewayInit(apiKey *ApiKey, organizationId string) {
	command := fmt.Sprintf(SCW_INIT_COMMAND, apiKey.SecretKey, apiKey.AccessKey, organizationId)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Command execution failed", err)
	} else {
		fmt.Println("Command executed successfully")
	}
}
