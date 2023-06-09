package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/artdarek/go-unzip"
	log "github.com/sirupsen/logrus"
)

func init() {
	// initialize update.log file and set log output to file
	file, err := os.OpenFile(path.Join(getDefaultDagFolderPath(), "update.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// Update type contains the update processes mandatory data
type Update struct {
	clientRPC          *rpc.Client
	downloadURL        string
	shiFolderPath      *string
	oldShiwaBinaryPath *string
	newVersion         *string
	triggerUpdate      *bool
}

// Signal is used for IPC with ShiwaWallet
type Signal struct {
	PID int
	Msg string
}

type unzippedContents struct {
	newShiwaBinaryPath string
	updateBinaryPath   string
}

func main() {
	var update Update

	update.downloadURL = "https://github.com/grvlle/constellation_wallet/releases/download"

	// ShiwaWallet provides the below data when an update is triggered
	update.shiFolderPath = flag.String("init_shi_path", getDefaultDagFolderPath(), "Enter the directory path to shi folder")
	update.oldShiwaBinaryPath = flag.String("init_shiwa_path", "", "Enter the directory path where the shiwa binary is located")
	update.newVersion = flag.String("new_version", "", "Enter the new semantic version. E.g 1.2.3")
	update.triggerUpdate = flag.Bool("upgrade", false, "Upgrade shiwa wallet binary")
	flag.Parse()

	update.Run()

}

// Run is the main method that runs the full update.
func (u *Update) Run() {
	var err error

	// Clean up old update artifacts
	err = u.CleanUp()
	if err != nil {
		log.Fatalf("Unable to clear previous local state: %v", err)
	}

	// Create a TCP connection to localhost on port 36866
	u.clientRPC, err = rpc.DialHTTP("tcp", "localhost:36866")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	log.Infof("Successfully established RPC connection with Shiwa Wallet")
	defer u.clientRPC.Close()

	zippedArchive, err := u.DownloadAppBinary()
	if err != nil {
		log.Fatalf("Unable to download v%s of Shiwa Wallet: %v", *u.newVersion, err)
	}

	ok, err := u.VerifyChecksum(zippedArchive)
	if err != nil || !ok {
		log.Fatalf("Checksum missmatch. Corrupted download: %v", err)
	}

	contents, err := unzipArchive(zippedArchive, *u.shiFolderPath)
	if err != nil {
		log.Fatalf("Unable to unzip contents: %v", err)
	}

	err = u.BackupApp()
	if err != nil {
		log.Fatalf("Unable to Backup Shiwa Wallet: %v", err)
	}

	err = u.TerminateAppService()
	if err != nil {
		log.Errorf("Unable to terminate Shiwa Wallet: %v", err)
		err = u.RestoreBackup()
		if err != nil {
			log.Fatal("Unable to restore backup: %v", err)
		}
	}

	err = u.ReplaceAppBinary(contents)
	if err != nil {
		log.Errorf("Unable to overwrite old installation: %v", err)
		err = u.RestoreBackup()
		if err != nil {
			log.Fatalf("Unable to restore backup: %v", err)
		}
	}

	err = u.LaunchAppBinary()
	if err != nil {
		log.Errorf("Unable to start up Shiwa after update: %v", err)
		err = u.RestoreBackup()
		if err != nil {
			log.Fatalf("Unable to restore backup: %v", err)
		}
	}

	err = u.CleanUp()
	if err != nil {
		log.Fatalf("Unable to clear previous local state: %v", err)
	}

}

// DownloadAppBinary downloads the latest Shiwa Wallet zip from github releases and returns the path to it
func (u *Update) DownloadAppBinary() (string, error) {

	filename := "shiwawallet.zip"
	osBuild, _ := getUserOS() // returns linux, windows, darwin or unsupported as well as the file extension (e.g ".exe")

	if osBuild == "unsupported" {
		return "", fmt.Errorf("the OS is not supported")
	}

	url := u.downloadURL + "/v" + *u.newVersion + "-" + osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/shiwawallet.zip
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(*u.shiFolderPath, filename)
	err := downloadFile(url, filePath)
	if err != nil {
		return "", fmt.Errorf("unable to download remote checksum: %v", err)
	}

	return filePath, nil
}

// VerifyChecksum takes a file path and will check the file sha256 checksum against the checksum included
// in the downlaod returns false if there's a missmatch
func (u *Update) VerifyChecksum(filePathZip string) (bool, error) {
	// Download checksum
	filename := "checksum.sha256"
	osBuild, _ := getUserOS() // returns linux, windows, darwin or unsupported as well as the file extension (e.g ".exe")

	if osBuild == "unsupported" {
		return false, fmt.Errorf("the OS is not supported")
	}

	url := u.downloadURL + "/v" + *u.newVersion + "-" + osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/checksum.sha256
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(*u.shiFolderPath, filename)
	err := downloadFile(url, filePath)
	if err != nil {
		return false, fmt.Errorf("unable to download remote checksum: %v", err)
	}

	// Read the contents of the downloaded file (remoteChecksum)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(content), "\n")
	remoteChecksum := lines[0]
	log.Infof("Remote file checksum: %v", remoteChecksum)

	// Collect the checksum of the downloaded zip (localChecksum)
	f, err := os.Open(filePathZip)
	if err != nil {
		return false, err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return false, err
	}
	localChecksum := hex.EncodeToString(hasher.Sum(nil))
	log.Infof("Local file checksum: %v", localChecksum)

	return remoteChecksum == localChecksum, nil
}

// TerminateAppService will send an RPC to shiwawallet to terminate the application
func (u *Update) TerminateAppService() error {
	sig := Signal{0, "Terminate Shiwa Wallet. Begining Update..."}
	var response Signal

	err := u.clientRPC.Call("RPCEndpoints.ShutDown", sig, &response)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(response.PID)
	if err != nil {
		log.Panicf("%v", err)
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	return nil
}

// BackupApp backs up the old binary in case of a failed update.
func (u *Update) BackupApp() error {
	_, fileExt := getUserOS()

	// Create backup folder in ~/.shi
	err := os.Mkdir(*u.shiFolderPath+"/backup", 0755)
	if err != nil {
		return fmt.Errorf("unable to create backup folder: %v", err)
	}

	// Copy the old Shiwa Wallet binary into ~/.shi/backup/
	err = copy(*u.oldShiwaBinaryPath, *u.shiFolderPath+"/backup/shiwawallet"+fileExt)
	if err != nil {
		return fmt.Errorf("unable to backup old Shiwa installation: %v", err)
	}

	return nil
}

// ReplaceAppBinary replaces the old update module and shiwa binary with the newly downloaded ones.
func (u *Update) ReplaceAppBinary(contents *unzippedContents) error {
	// Replace old shiwa binary with the new one
	_, fileExt := getUserOS()
	err := copy(contents.newShiwaBinaryPath, *u.oldShiwaBinaryPath)
	if err != nil {
		for i := 5; i > 0; i-- {
			time.Sleep(time.Duration(i) * time.Second)
			err = copy(contents.newShiwaBinaryPath, *u.oldShiwaBinaryPath)
			if err == nil {
				break
			} else if err != nil && i == 0 {
				return fmt.Errorf("unable to overwrite old shiwa binary: %v", err)
			}
		}
	}
	// Replace old update binary with the new one
	if fileExists(contents.updateBinaryPath) {
		err = copy(contents.updateBinaryPath, *u.shiFolderPath+"/update"+fileExt)
		if err != nil {
			return fmt.Errorf("unable to copy update binary to .shi folder: %v", err)
		}
	}
	return nil
}

// LaunchAppBinary executes the new shiwa binary
func (u *Update) LaunchAppBinary() error {
	cmd := exec.Command(*u.oldShiwaBinaryPath)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to execute run command for Shiwa Wallet: %v", err)
	}
	return nil
}

// RestoreBackup restores the backed up files
func (u *Update) RestoreBackup() error {

	log.Infoln("Restoring Backup...")

	// Copy the old Shiwa Wallet binary from ~/.shi/backup/ to the old path
	_, fileExt := getUserOS()
	err := copy(*u.shiFolderPath+"/backup/shiwawallet"+fileExt, *u.oldShiwaBinaryPath)
	if err != nil {
		return fmt.Errorf("unable to overwrite old shiwa binary: %v", err)
	}

	// Copy update binary from ~/.shi/backup/update -> ~/.shi/update
	if fileExists(*u.shiFolderPath + "/backup/update" + fileExt) {
		err = copy(*u.shiFolderPath+"/backup/update"+fileExt, *u.shiFolderPath+"/update"+fileExt)
		if err != nil {
			return fmt.Errorf("unable to copy update binary to .shi folder: %v", err)
		}
	}

	log.Infoln("Backup successfully restored.")

	return nil

}

// CleanUp removes uneccesary artifacts from the update process
func (u *Update) CleanUp() error {

	if fileExists(*u.shiFolderPath + "/shiwawallet.zip") {
		err := os.Remove(*u.shiFolderPath + "/shiwawallet.zip")
		if err != nil {
			return err
		}
	}
	if fileExists(*u.shiFolderPath + "/backup") {
		err := os.RemoveAll(*u.shiFolderPath + "/backup")
		if err != nil {
			return err
		}
	}

	if fileExists(*u.shiFolderPath + "/new_build") {
		err := os.RemoveAll(*u.shiFolderPath + "/new_build")
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(url, filePath string) error {

	tmpFilePath := filePath + ".tmp"
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	out.Close() // Close file to rename

	if err = os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}
	return nil
}

func getDefaultDagFolderPath() string {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Unable to detect UserHomeDir: %v", err)
		return ""
	}
	return userDir + "/.shi"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}

func copy(src string, dst string) error {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Unzips archive to dstPath, returns path to wallet binary
func unzipArchive(zippedArchive, dstPath string) (*unzippedContents, error) {

	uz := unzip.New(zippedArchive, dstPath+"/"+"new_build/")
	err := uz.Extract()
	if err != nil {
		return nil, err
	}
	_, fileExt := getUserOS()

	contents := &unzippedContents{
		newShiwaBinaryPath: dstPath + "/" + "new_build/shiwawallet" + fileExt,
		updateBinaryPath:   dstPath + "/" + "new_build/update" + fileExt,
	}

	return contents, err
}

// getUserOS returns the users OS as well as the file extension of executables for said OS
func getUserOS() (string, string) {
	var osBuild string
	var fileExt string

	switch os := runtime.GOOS; os {
	case "darwin":
		osBuild = "darwin"
		fileExt = ""
	case "linux":
		osBuild = "linux"
		fileExt = ""
	case "windows":
		osBuild = "windows"
		fileExt = ".exe"
	default:
		osBuild = "unsupported"
		fileExt = ""
	}

	return osBuild, fileExt
}
