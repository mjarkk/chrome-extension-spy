package firefox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// CreateEmptyProfile creates a empty firefox profile
func (f *FF) CreateEmptyProfile() {
	// firefox --profile /tmp/testFF about:addons https://mkopenga.com
	if f.Err() {
		return
	}

	profile, err := ioutil.TempDir("", "FFProfile")
	if err != nil {
		f.HasErr = err
		return
	}
	f.TmpDirs.Profile = profile

	var cmd *exec.Cmd

	go func() {
		fmt.Println("Starting firefox in background to create a empty user profile, waiting 5 seconds before closing firefox...")
		time.Sleep(time.Second * 5)
		cmd.Process.Signal(os.Kill)
	}()

	cmd = exec.Command(
		f.LaunchCMD,
		"--profile", profile,
		"-headless",
		"about:addons", "http://localhost:8080",
	)
	cmd.Output()
}

// Launch launches firefox
func (f *FF) Launch(kill chan struct{}) {
	var cmd *exec.Cmd

	go func() {
		<-kill
		fmt.Println("Closing firefox")
		cmd.Process.Signal(os.Kill)
	}()

	cmd = exec.Command(
		f.LaunchCMD,
		"--profile", f.TmpDirs.Profile,
		"-headless",
		"about:addons", "http://localhost:8080#isFF",
	)
	cmd.Output()
}
