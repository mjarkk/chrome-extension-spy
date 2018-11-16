package chrome

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Launch launches chrome without any user provile
// TODO: Fix the long time before command exits after closing google chrome
func Launch(extsPath string, launchCommand string, forceClose chan struct{}) error {
	tempDir, err := ioutil.TempDir("", "chrome-data")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	var cmd *exec.Cmd
	go func() {
		<-forceClose
		fmt.Println("killing chrome process")
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
		launchCommand,
		"--user-data-dir="+tempDir,
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
		"--load-extension=\""+allExts+"\"",
	)
	_, err = cmd.Output()
	return err
}
