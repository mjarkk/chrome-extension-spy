package chrome

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

// Launch launches chrome without any user provile
// TODO: Fix the long time before command exits after closing google chrome
func Launch(extsPath string, chromeType string, forceClose chan struct{}) error {
	tempDir, err := ioutil.TempDir("", "chrome-data")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	var cmd *exec.Cmd
	go func() {
		<-forceClose
		fmt.Println("Trying to stop chrome process")
		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println("Cloud not stop chrome, killing program")
			os.Exit(1)
		}()
		cmd.Process.Kill()
	}()
	files, err := ioutil.ReadDir(extsPath)
	if err != nil {
		return err
	}
	dirs := []string{}
	for _, folderItem := range files {
		dirs = append(dirs, path.Join(extsPath, folderItem.Name()))
	}
	allExts := strings.Join(dirs, ",")
	cmd = exec.Command(
		ChromeLocation(chromeType),
		"--user-data-dir="+tempDir, // set the data dir in this case a empty dir to make sure chrome starts fully clean
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-features=site-per-process",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",
		"--metrics-recording-only",
		"--no-first-run",
		"--safebrowsing-disable-auto-update",
		"--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
		"--load-extension=\""+allExts+"\"", // Add all extensions to load
		"http://localhost:8080/",           // Open chrome with this url
	)
	_, err = cmd.Output()
	return err
}

// ChromeLocation returns a valid launch command for chrome
func ChromeLocation(chromeType string) string {
	if runtime.GOOS == "windows" {
		switch chromeType {
		case "google-chrome", path.Join("Google", "Chrome"):
			// normal google chrome
			return "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
		case "chromium":
			// chromium
			return "C:\\Users\\mark\\AppData\\Local\\Chromium\\Application\\chrome.exe"
		case "google-chrome-dev", path.Join("Google", "Chrome-dev"):
			// google chrome dev
			return "C:\\Program Files (x86)\\Google\\Chrome Dev\\Application\\chrome.exe"
		case "google-chrome-beta", path.Join("Google", "Chrome-beta"):
			// google chrome beta
			return "C:\\Program Files (x86)\\Google\\Chrome Beta\\Application\\chrome.exe"
		case "google-chrome-unstable", path.Join("Google", "Chrome-unstable"):
			// google chrome unstable
			return "C:\\Program Files (x86)\\Google\\Chrome Unstable\\Application\\chrome.exe"
		case "google-chrome-canary", path.Join("Google", "Chrome-canary"):
			// google chrome canary
			return "C:\\Users\\mark\\AppData\\Local\\Google\\Chrome SxS\\Application\\chrome.exe"
		default:
			return chromeType
		}
	}
	return chromeType
}
