package browserutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func BrowserExecutable(browserType string) string {
	if runtime.GOOS == "windows" {
		switch browserType {
		case "chrome":
			return "chrome.exe"
		case "brave":
			return "brave.exe"
		case "edge":
			return "msedge.exe"
		case "opera":
			return "opera.exe"
		case "firefox":
			return "firefox.exe"
		default:
			return browserType + ".exe"
		}
	}

	switch browserType {
	case "chrome":
		return "google-chrome"
	case "brave":
		return "brave-browser"
	case "edge":
		return "microsoft-edge"
	case "opera":
		return "opera"
	case "firefox":
		return "firefox"
	case "safari":
		return "safari"
	default:
		return browserType
	}
}

func FindBrowser(browserType string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return Windows(browserType)
	case "darwin":
		return Mac(browserType)
	case "linux":
		return Linux(browserType)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func Windows(browserType string) (string, error) {
	var possiblePaths []string

	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	localAppData := os.Getenv("LOCALAPPDATA")
	userProfile := os.Getenv("USERPROFILE")

	switch browserType {
	case "chrome":
		possiblePaths = []string{
			filepath.Join(programFiles, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(programFilesX86, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(localAppData, "Google", "Chrome", "Application", "chrome.exe"),
		}
	case "brave":
		possiblePaths = []string{
			filepath.Join(programFiles, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
			filepath.Join(programFilesX86, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
			filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
		}
	case "edge":
		possiblePaths = []string{
			filepath.Join(programFiles, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(programFilesX86, "Microsoft", "Edge", "Application", "msedge.exe"),
		}
	case "opera":
		possiblePaths = []string{
			filepath.Join(localAppData, "Programs", "Opera", "opera.exe"),
			filepath.Join(userProfile, "AppData", "Local", "Programs", "Opera", "opera.exe"),
			filepath.Join(programFiles, "Opera", "opera.exe"),
			filepath.Join(programFilesX86, "Opera", "opera.exe"),
		}
	case "firefox":
		possiblePaths = []string{
			filepath.Join(programFiles, "Mozilla Firefox", "firefox.exe"),
			filepath.Join(programFilesX86, "Mozilla Firefox", "firefox.exe"),
		}
	default:
		return "", fmt.Errorf("unsupported browser type: %s", browserType)
	}

	if pathExe, err := exec.LookPath(BrowserExecutable(browserType)); err == nil {
		return pathExe, nil
	}

	for _, path := range possiblePaths {
		if fileExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("browser not found in any of the expected locations: %v", possiblePaths)
}

func Mac(browserType string) (string, error) {
	var possiblePaths []string

	switch browserType {
	case "chrome":
		possiblePaths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "brave":
		possiblePaths = []string{
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		}
	case "edge":
		possiblePaths = []string{
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}
	case "opera":
		possiblePaths = []string{
			"/Applications/Opera.app/Contents/MacOS/Opera",
		}
	case "firefox":
		possiblePaths = []string{
			"/Applications/Firefox.app/Contents/MacOS/firefox",
		}
	case "safari":
		possiblePaths = []string{
			"/Applications/Safari.app/Contents/MacOS/Safari",
		}
	default:
		return "", fmt.Errorf("unsupported browser type: %s", browserType)
	}

	if pathExe, err := exec.LookPath(BrowserExecutable(browserType)); err == nil {
		return pathExe, nil
	}

	for _, path := range possiblePaths {
		if fileExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("browser not found in any of the expected locations: %v", possiblePaths)
}

func Linux(browserType string) (string, error) {
	executableName := BrowserExecutable(browserType)

	if pathExe, err := exec.LookPath(executableName); err == nil {
		return pathExe, nil
	}

	// Common Linux installation paths
	var possiblePaths []string
	homeDir, _ := os.UserHomeDir()

	switch browserType {
	case "chrome":
		possiblePaths = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/opt/google/chrome/chrome",
			"/snap/bin/chromium",
		}
	case "brave":
		possiblePaths = []string{
			"/usr/bin/brave-browser",
			"/opt/brave.com/brave/brave",
			"/snap/bin/brave",
		}
	case "edge":
		possiblePaths = []string{
			"/usr/bin/microsoft-edge",
			"/opt/microsoft/msedge/msedge",
		}
	case "opera":
		possiblePaths = []string{
			"/usr/bin/opera",
			"/opt/opera/opera",
			"/snap/bin/opera",
			filepath.Join(homeDir, ".local", "share", "applications", "opera.desktop"),
		}
	case "firefox":
		possiblePaths = []string{
			"/usr/bin/firefox",
			"/opt/firefox/firefox",
			"/snap/bin/firefox",
		}
	default:
		return "", fmt.Errorf("unsupported browser type: %s", browserType)
	}

	for _, path := range possiblePaths {
		if fileExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("browser not found in PATH or common locations: %v", possiblePaths)
}
