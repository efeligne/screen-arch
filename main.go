package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	errUnexpectedFileType = errors.New("[warn]: unexpected file type")
	errDoesNotMatchFormat = errors.New("[warn]: filename does not match expected format")
)

const (
	expectedMatches = 4
	expectedPattern = `^(?:ScreenShot|Screenshot|Screen[\s_]Recording)[\s_](\d{4})-(\d{2})-\d{2}[\s_-](?:(?:at|в)[\s_])?\d{1,2}[\.-]\d{1,2}[\.-]\d{1,2}(?:.(?:AM|PM|am|pm))?(?:.\d+)?\.(\w+)$` //nolint:lll // can't be split
)

type fileNameParts struct {
	year  string
	month string
	ext   string
}

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

func parseFilename(filename string) (fileNameParts, error) {
	re := regexp.MustCompile(expectedPattern)

	matches := re.FindStringSubmatch(filename)
	if len(matches) != expectedMatches {
		return fileNameParts{}, fmt.Errorf("%w: %s", errDoesNotMatchFormat, filename)
	}

	parts := fileNameParts{
		year:  matches[1],
		month: matches[2],
		ext:   matches[3],
	}

	return parts, nil
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
		return "", fmt.Errorf("%w: %q", errUnexpectedFileType, ext)
	}

	typePath := "Screenshots"
	if isVideo {
		typePath = "Screen Recordings"
	}

	targetPath := filepath.Join(capturesDir, typePath, year, month)

	err := os.MkdirAll(targetPath, 0o755)
	if err != nil {
		return "", fmt.Errorf("[error]: failed to create %q directory: %w", targetPath, err)
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

		parts, err := parseFilename(filename)
		if err != nil {
			log.Printf("[warn]: %v. Skipping...", err)
			continue
		}

		target, err := createTargetPath(capturesDir, parts.year, parts.month, parts.ext)
		if err != nil {
			log.Println(err)
			continue
		}

		targetFile := filepath.Join(target, filename)

		err = os.Rename(outFile, targetFile)
		if err != nil {
			log.Println("[error]: failed to move file: ", outFile)
			continue
		}
	}

	log.Println("[info]: done!")
}
