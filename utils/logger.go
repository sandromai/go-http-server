package utils

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/types"
)

type Logger struct {
	FolderPath     string
	FileName       string
	MessagesPrefix string
}

func (log *Logger) formatMessage(
	message string,
) (string, *types.AppError) {
	formattedMessage := "[" + time.Now().Format(time.DateTime) + "] "

	if log.MessagesPrefix != "" {
		formattedMessage += log.MessagesPrefix + " "
	}

	cleanMessage, appErr := FormatLineBreaks(message)

	if appErr != nil {
		return "", appErr
	}

	spacesRegExp, err := regexp.Compile(`\n +`)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to format spaces.",
		}
	}

	formattedMessage += spacesRegExp.ReplaceAllString(cleanMessage, " ")
	formattedMessage = strings.TrimSpace(formattedMessage) + "\n"

	return formattedMessage, nil
}

func (log *Logger) createFile() (
	string,
	*types.AppError,
) {
	newFileName := ""
	fileExists := true

	for i := 0; i < 20 && fileExists; i++ {
		newFileName = fmt.Sprintf("%v", time.Now().Unix()) + "-" + fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", time.Now().UnixNano())+log.FileName))) + "-" + log.FileName

		if _, err := os.Stat(filepath.Join(log.FolderPath, newFileName)); err != nil && os.IsNotExist(err) {
			fileExists = false
		}
	}

	if newFileName == "" || fileExists {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to generate log file name.",
		}
	}

	fileHandler, err := os.Create(filepath.Join(log.FolderPath, newFileName))

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create log file.",
		}
	}

	fileHandler.Close()

	return newFileName, nil
}

func (log *Logger) getCurrentFilePath() (
	string,
	*types.AppError,
) {
	folderInfo, err := os.Stat(log.FolderPath)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(log.FolderPath, 0755)

			if err != nil {
				return "", &types.AppError{
					StatusCode: 500,
					Message:    "Failed to create log folder.",
				}
			}

			folderInfo, err = os.Stat(log.FolderPath)

			if err != nil {
				return "", &types.AppError{
					StatusCode: 500,
					Message:    "Failed to get log folder data.",
				}
			}
		} else {
			return "", &types.AppError{
				StatusCode: 500,
				Message:    "Failed to get log folder data.",
			}
		}
	}

	if !folderInfo.IsDir() {
		err = os.MkdirAll(log.FolderPath, 0755)

		if err != nil {
			return "", &types.AppError{
				StatusCode: 500,
				Message:    "Failed to create log folder.",
			}
		}
	}

	files, err := os.ReadDir(log.FolderPath)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to open log folder.",
		}
	}

	var currentFile *fs.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileNameParts := strings.Split(file.Name(), "-")

		if fileNameParts[len(fileNameParts)-1] != log.FileName {
			continue
		}

		completeFilePath := filepath.Join(log.FolderPath, file.Name())

		fileInfo, err := os.Stat(completeFilePath)

		if err != nil {
			continue
		}

		fileTime := fileInfo.ModTime().Unix()

		if fileTime <= time.Now().Add(-time.Hour*24*30).Unix() {
			os.Remove(completeFilePath)

			continue
		}

		if currentFile == nil {
			currentFile = &fileInfo

			continue
		}

		currentFileTime := (*currentFile).ModTime().Unix()

		if fileTime > currentFileTime {
			currentFile = &fileInfo
		}
	}

	var currentFileName string
	var appErr *types.AppError

	if currentFile != nil {
		currentFileName = (*currentFile).Name()

		currentFileParts := strings.Split(currentFileName, "-")

		if len(currentFileParts) >= 2 && !strings.Contains(currentFileParts[0], ".") {
			currentFileCreationTime, err := strconv.Atoi(currentFileParts[0])

			if err != nil {
				currentFileName, appErr = log.createFile()

				if appErr != nil {
					return "", appErr
				}
			}

			currentYear, currentMonth, currentDay := time.Now().Date()
			fileYear, fileMonth, fileDay := time.Unix(int64(currentFileCreationTime), 0).Date()

			if currentYear != fileYear || currentMonth != fileMonth || currentDay != fileDay {
				currentFileName, appErr = log.createFile()

				if appErr != nil {
					return "", appErr
				}
			}
		}
	} else {
		currentFileName, appErr = log.createFile()

		if appErr != nil {
			return "", appErr
		}
	}

	return filepath.Join(log.FolderPath, currentFileName), nil
}

func (log *Logger) Save(
	message string,
) *types.AppError {
	currentFilePath, appErr := log.getCurrentFilePath()

	if appErr != nil {
		return appErr
	}

	fileHandler, err := os.OpenFile(
		currentFilePath,
		os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to open log file.",
		}
	}

	defer fileHandler.Close()

	formattedMessage, appErr := log.formatMessage(message)

	if appErr != nil {
		return appErr
	}

	_, err = fileHandler.Write([]byte(
		formattedMessage,
	))

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to save log.",
		}
	}

	return nil
}
