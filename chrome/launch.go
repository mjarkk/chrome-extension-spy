package chrome

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Launch launches chrome without any user provile
func Launch(extPath string, closeChrome chan struct{}) error {
	tempDir, err := ioutil.TempDir("", "chrome-data")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	var cmd *exec.Cmd
	go func() {
		<-closeChrome
		fmt.Println("killing chrome process")
		cmd.Process.Kill()
	}()
	cmd = exec.Command(
		"google-chrome-unstable",
		"--user-data-dir="+tempDir,
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-extensions",
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
	)
	_, err = cmd.Output()
	return err
}
