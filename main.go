package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func getDesktopDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[error]: failed to get user home dir path %q: %v", home, err)
	}
	return filepath.Join(home, "Desktop")
}

func getCapturesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[error]: failed to get user home dir path %q: %v", home, err)
	}
	return filepath.Join(home, "Documents", "Captures")
}

func parseFilename(filename string) (year, month, extension string, err error) {
	re := regexp.MustCompile(`^(?:ScreenShot|Screenshot|Screen[\s_]Recording)[\s_](\d{4})-(\d{2})-\d{2}[\s_-](?:(?:at|в)[\s_])?\d{1,2}[\.-]\d{1,2}[\.-]\d{1,2}(?:.(?:AM|PM|am|pm))?(?:.\d+)?\.(\w+)$`)

	matches := re.FindStringSubmatch(filename)
	if len(matches) != 4 {
		return "", "", "", fmt.Errorf("filename does not match expected format: %s", filename)
	}

	return matches[1], matches[2], matches[3], nil
}

func getFiles(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("[error]: failed to read path %q: %v", path, err)
	}

	return entries
}

func createTargetPath(capturesDir, year, month, ext string) (string, error) {
	isVideo := ext == "mov" || ext == "mp4"
	isPict := ext == "png" || ext == "jpeg" || ext == "jpg"

	if !isVideo && !isPict {
		return "", fmt.Errorf("[warn]: unexpected file type: %q", ext)
	}

	typePath := "Screenshots"
	if isVideo {
		typePath = "Screen Recordings"
	}

	targetPath := filepath.Join(capturesDir, typePath, year, month)
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return "", fmt.Errorf("[error]: failed to create %q directory: %v", targetPath, err)
	}

	return targetPath, nil
}

func main() {
	desktopDir := getDesktopDir()
	capturesDir := getCapturesDir()

	for _, entry := range getFiles(desktopDir) {
		filename := entry.Name()
		outFile := filepath.Join(desktopDir, filename)
		log.Println("[info]: processing file: ", outFile)

		if entry.IsDir() {
			log.Printf("[warn]: path %q is dir. Skipping...", outFile)
			continue
		}

		year, month, ext, err := parseFilename(filename)
		if err != nil {
			log.Printf("[warn]: %v. Skipping...", err)
			continue
		}

		target, err := createTargetPath(capturesDir, year, month, ext)
		if err != nil {
			log.Println(err)
			continue
		}

		targetFile := filepath.Join(target, filename)

		if err := os.Rename(outFile, targetFile); err != nil {
			log.Println("[error]: failed to move file: ", outFile)
			continue
		}
	}

	log.Println("[info]: done!")
}
