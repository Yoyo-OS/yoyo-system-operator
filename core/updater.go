package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/yoyo-os/yso/settings"
)

var (
	checkLogPath       = "/var/log/yso-check.log"
	firstSetupDonePath = "/etc/yoyo-first-setup-done"
	abrootLockPath     = "/tmp/abroot-transactions.lock"
)

// AreABRootTransactionsLocked checks if there are any abroot transactions
// currently running
func AreABRootTransactionsLocked() bool {
	_, err := os.Stat(abrootLockPath)

	return err == nil
}

// NeedUpdate checks if the system needs to be updated according to the latest
// update log compared to the YSO configuation
func NeedUpdate() bool {
	if AreABRootTransactionsLocked() {
		fmt.Println("ABRoot transactions are currently locked. Skipping update check.")
		return false
	}

	if !okToUpdate() {
		fmt.Println("Deferring update until first-setup is complete.")
		return false
	}

	res := false
	currentTime := time.Now()
	schedule := settings.GetConfigValue("updates.schedule")
	latestCheck := getLatestCheck()
	if latestCheck == nil {
		fmt.Println("No previous update check found. Triggering update.")
		return true
	}

	switch schedule {
	case "daily":
		if currentTime.Sub(*latestCheck).Hours() >= 24 {
			res = true
		}
	case "weekly":
		if currentTime.Sub(*latestCheck).Hours() >= 168 {
			res = true
		}
	case "monthly":
		if currentTime.Sub(*latestCheck).Hours() >= 720 {
			res = true
		}
	}

	status, _, _ := HasUpdates()
	if status {
		res = true
	}

	return res
}

// HasUpdates checks if the system has updates available
func HasUpdates() (bool, []string, error) {
	aptBinary := "/usr/bin/___apt___"
	if _, err := os.Stat(aptBinary); err != nil {
		aptBinary = "/usr/bin/apt" // Fallback to apt only for development purposes
		if _, err := os.Stat(aptBinary); err != nil {
			return false, nil, err
		}
	}

	update_cmd := exec.Command(aptBinary, "update")
	update_cmd.Stdin = os.Stdin
	update_cmd.Stdout = os.Stdout
	update_cmd.Stderr = os.Stderr
	if err := update_cmd.Run(); err != nil {
		return false, nil, err
	}

	list_cmd := exec.Command(aptBinary, "list", "--upgradable")
	list_cmd.Env = os.Environ()
	list_cmd.Env = append(list_cmd.Env, "LANG=en_US.UTF-8")
	output, err := list_cmd.Output()
	if err != nil {
		return false, nil, err
	}

	packages := strings.Split(string(output), "\n")
	if len(packages) <= 2 {
		return false, nil, nil
	}

	list_updates := []string{}
	for _, pkg := range packages[1 : len(packages)-1] { // First and last lines are not packages
		cols := strings.Split(pkg, " ")
		pkg_name, pkg_newver, pkg_oldver := cols[0], cols[1], cols[5]
		pkg_name = strings.Split(pkg_name, "/")[0]  // Remove source info
		pkg_oldver = pkg_oldver[:len(pkg_oldver)-1] // Remove trailing "]"

		list_updates = append(list_updates, fmt.Sprintf("  - %s\t%s -> %s", pkg_name, pkg_oldver, pkg_newver))
	}

	return true, list_updates, nil
}
func okToUpdate() bool {
	// check yso logs for previous run
	if _, err := os.Stat(checkLogPath); os.IsNotExist(err) {
		// yso has never run, check if first-setup has been done
		if _, err := os.Stat(firstSetupDonePath); os.IsNotExist(err) {
			// first setup has not completed. don't try to update
			return false
		}

	}
	return true
}

// getLatestCheck returns the latest check time from the log file, it also
// write the current time to the log file if it doesn't exist
func getLatestCheck() *time.Time {
	var latestCheck time.Time

	if _, err := os.Stat(checkLogPath); os.IsNotExist(err) {
		latestCheck = time.Now()
		writeLatestCheck(latestCheck)
		return &latestCheck
	}

	file, err := os.Open(checkLogPath)
	if err != nil {
		fmt.Println("Failed to open update log file:", err)
		return nil
	}
	defer file.Close()

	content, err := os.ReadFile(checkLogPath)
	if err != nil {
		fmt.Println("Failed to read update log file:", err)
		return nil
	}

	timeString := string(content)
	timeArr := strings.Split(timeString, " ")
	latestCheck, err = time.Parse("2006-01-02 15:04:05", timeArr[0]+" "+timeArr[1])
	if err != nil {
		fmt.Println("Failed to parse update log file:", err)
		return nil
	}

	return &latestCheck
}

// writeLatestCheck writes the latest check time to the log file
func writeLatestCheck(t time.Time) error {
	file, err := os.OpenFile(checkLogPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, t)
	if err != nil {
		return err
	}

	return nil
}

// TryUpdate tries to update the system via ABRoot
func TryUpdate(force bool) error {
	if !okToUpdate() {
		fmt.Println("Deferring update until first-setup is done.")
		return nil

	}
	if !force && !SmartUpdate() {
		fmt.Println("Smart update detected device is being used, skipping update.")
		return nil
	}

	if os.Getenv("YSO_DEBUG_SMARTUPDATE") != "" {
		fmt.Println("Smart update is in debug mode, skipping update.")
		return nil
	}

	writeLatestCheck(time.Now())

	file, err := os.Create("/tmp/" + time.Now().Format("20060102150405") + "-script")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("#!/bin/bash\napt update && apt upgrade -y && apt autoremove -y")
	if err != nil {
		return err
	}

	cmd := exec.Command("abroot", "exec", "--force-run", "sh", file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	SendNotification("Update", "System updated successfully, restart to apply changes.")
	return nil
}

// SmartUpdate checks if the device is currently being used, then returns true
// if the device is not being used
func SmartUpdate() bool {
	if settings.GetConfigValue("updates.smart") == false {
		fmt.Println("Smart update is disabled, skipping check.")
		return true
	}

	fmt.Println("Starting smart update check...")

	commonChecks := GetCommonChecks()

	// battery check
	if commonChecks.IsLaptop && commonChecks.LowBattery {
		fmt.Printf("Low battery detected, skipping update.")
		return false
	}

	// internet usage check
	if commonChecks.HighInternetUsage {
		if commonChecks.InternetUsage == 0 {
			fmt.Println("Something goes wrong during internet usage check, skipping update.")
			return false
		}
		fmt.Printf("Connection being used (%d), skipping update.", commonChecks.InternetUsage)
		return false
	}

	// ram check
	if commonChecks.HighMemoryUsage {
		if commonChecks.MemoryUsage == 0 {
			fmt.Println("Something goes wrong during memory usage check, skipping update.")
			return false
		}
		fmt.Printf("Memory usage detected (%d), skipping update.", commonChecks.MemoryUsage)
		return false
	}

	// cpu check
	if commonChecks.HighCPUUsage {
		if commonChecks.CPUUsage == 0 {
			fmt.Println("Something goes wrong during CPU usage check, skipping update.")
			return false
		}
		fmt.Printf("CPU usage detected (%d), skipping update.", commonChecks.CPUUsage)
		return false
	}

	// cpu temp check
	if commonChecks.CPUTemp >= 70 {
		fmt.Printf("High CPU temperature detected (%d), skipping update.", commonChecks.CPUTemp)
		return false
	}

	return true
}
