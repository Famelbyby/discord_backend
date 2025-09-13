package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const DefaultSessionLifetime = 10 * time.Hour

var DefaultSessionLifetimeString = strconv.Itoa(int(DefaultSessionLifetime.Seconds()))

type ErrorWrapper struct {
	Err string `json:"err"`
}

func WriteError(w http.ResponseWriter, err string) {
	errWrapper := ErrorWrapper{Err: err}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json, _ := json.Marshal(errWrapper)
	w.Write(json)
}
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	inputFile.Close()

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}
	return nil
}
func Remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}
func DateTimeToSecondsForDb(date string) (string, error) {
	time, err := time.Parse("02.01.2006 15:04:05", date)
	time = time.AddDate(-55, 0, 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(time.Unix()), nil
}
