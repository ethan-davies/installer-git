package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
)

// Change these variables your liking
const gitHubAPIRepoURL = "https://api.github.com/repos/ethan-davies/myapp" // Make sure you use https://api.github.com/repos/
const gitHubRepoURL = "https://github.com/ethan-davies/myapp"
const name = "myapp" // This will be used to inform the user
const installFileName = "myapp" // The folder where your binary will be installed

// These must be used in the gitub releases (make sure to have file extensions such as .exe)
const windowsBinaryName = "myapp"
const linuxBinaryName = "myapp"
const darwinBinaryName = "myapp-macos"

// WARNING: Editing the following code may break your installation. 
func getInstallDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), installFileName)
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", installFileName)
	default: // Default to Linux-like systems
		return filepath.Join(os.Getenv("HOME"), "." + installFileName)
	}
}

func getBinaryFileName() string {
	if runtime.GOOS == "windows" {
		return installFileName + ".exe"
	}
	return installFileName
}

func downloadFile(url, targetPath string) error {
	fmt.Printf("Downloading file from: %s\n", url)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", response.Status)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}



func addToPath(path string) error {
	pathVar := os.Getenv("PATH")
	if !strings.Contains(pathVar, path) {
		pathVar = path + string(os.PathListSeparator) + pathVar
		err := os.Setenv("PATH", pathVar)
		if err != nil {
			return err
		}

		// Add to PATH using system-specific command
		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("setx", "PATH", pathVar)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		case "linux", "darwin":
			cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo 'export PATH="%s:$PATH"' >> ~/.bashrc`, path))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}

	return nil
}

func waitForKeyPress() {
	fmt.Println("Press Enter to continue...")
	fmt.Scanln() // Wait for Enter key
}

func fetchLatestVersion() (*version.Version, error) {
    releasesURL := fmt.Sprintf("%s/releases/latest", gitHubAPIRepoURL)

    client := http.DefaultClient
    req, err := http.NewRequest("GET", releasesURL, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Accept", "application/vnd.github.v3+json")
    response, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    var releaseData struct {
        TagName string `json:"tag_name"`
    }
    if err := json.Unmarshal(body, &releaseData); err != nil {
        return nil, err
    }

    return version.NewVersion(releaseData.TagName)
}



func main() {
	fmt.Printf("Installing %s...", name)

	installDir := getInstallDir()
	binDir := filepath.Join(installDir, "bin")
	binaryFileName := getBinaryFileName()
	binaryPath := filepath.Join(binDir, binaryFileName)

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		fmt.Println("Error fetching latest version:", err)
		waitForKeyPress()
		return
	}

	fmt.Println("Latest version:", latestVersion)

	// Determine the platform-specific URL for the binary
	var platformURL string
	switch runtime.GOOS {
	case "windows":
		platformURL = fmt.Sprintf("%s/releases/download/v%s/%s", gitHubRepoURL, latestVersion, windowsBinaryName)
	case "linux":
		platformURL = fmt.Sprintf("%s/releases/download/v%s/%s", gitHubRepoURL, latestVersion, linuxBinaryName)
	case "darwin":
		platformURL = fmt.Sprintf("%s/releases/download/v%s/%s", gitHubRepoURL, latestVersion, darwinBinaryName)
	default:
		fmt.Println("Unsupported platform:", runtime.GOOS)
		waitForKeyPress()
		return
	}

	// Create the installation directory if it doesn't exist
	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		fmt.Println("Error creating installation directory:", err)
		waitForKeyPress()
		return
	}

	// Download the binary
	fmt.Printf("Downloading %s binary...", name)
	err = downloadFile(platformURL, binaryPath)
	if err != nil {
		fmt.Println("Error downloading binary:", err)
		waitForKeyPress()
		return
	}

	// Make the binary executable
	fmt.Println("Setting file permissions...")
	err = os.Chmod(binaryPath, 0755)
	if err != nil {
		fmt.Println("Error setting file permissions:", err)
		waitForKeyPress()
		return
	}

	// Add the installation directory to the system's PATH
	fmt.Println("Adding to system PATH...")
	err = addToPath(binDir)
	if err != nil {
		fmt.Println("Error adding to PATH:", err)
		waitForKeyPress()
		return
	}

	// Print success message
	fmt.Printf("%s has been successfully installed to %s.\n", name, installDir)
	waitForKeyPress()
}