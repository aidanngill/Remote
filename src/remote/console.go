package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Task struct {
	Status string `json:"status"`
	Error  string `json:"error_code"`
	ID     int    `json:"task_id"`
	Title  string `json:"title"`
}

type TaskProgress struct {
	Total       string `json:"length_total"`
	Transferred string `json:"transferred_total"`
}

var (
	reHex   = regexp.MustCompile(`(0[xX][0-9a-fA-F]+)`)
	devices []net.IP
)

func createConsoleRoute(endpoint string) string {
	return fmt.Sprintf("http://%s:12800/api/%s", consoleBoxAddress.Text(), endpoint)
}

func createHostRoute(endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", devices[serverBoxDevices.Selected()], 8080, endpoint)
}

func createPackage(filePath string) (Package, error) {
	if filePath == "" {
		return Package{}, errors.New("the file path cannot be empty")
	}

	if !strings.HasSuffix(strings.ToLower(filePath), ".pkg") {
		return Package{}, errors.New("only valid package files ending with the PKG extension may be used")
	}

	stat, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return Package{}, errors.New("the given file no longer exists")
	}

	h := md5.New()
	h.Write([]byte(filePath))

	return Package{
		ID:       hex.EncodeToString(h.Sum(nil)),
		Path:     filePath,
		Row:      uint16(len(packages)),
		Size:     uint64(stat.Size()),
		Progress: 0,
	}, nil
}

func createTask(pkg Package) (Task, error) {
	url := createConsoleRoute("install")

	defaultTask := Task{
		ID:     0,
		Status: "fail",
	}

	data := map[string]interface{}{
		"type":     "direct",
		"packages": []string{createHostRoute(pkg.ID)},
	}

	body, err := json.Marshal(data)

	if err != nil {
		return defaultTask, err
	}

	res, err := client.Post(url, "application/json", bytes.NewBuffer([]byte(body)))

	if err != nil {
		return defaultTask, err
	}

	if res.StatusCode != 200 {
		return defaultTask, fmt.Errorf("status code != 200, was %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return defaultTask, errors.New("failed to read the response body")
	}

	processed := reHex.ReplaceAll([]byte(body), []byte(`"$1"`))

	var task Task
	json.Unmarshal(processed, &task)

	if task.Status == "fail" {
		ec := "the console failed to start the task"

		switch task.Error {
		case "0x80990015":
			ec = "the given title ID already exists"
		}

		return defaultTask, errors.New(ec)
	}

	return task, nil
}

func getTaskProgress(task Task) (int64, error) {
	url := createConsoleRoute("get_task_progress")

	data := map[string]interface{}{
		"task_id": task.ID,
	}

	body, err := json.Marshal(data)

	if err != nil {
		return 0, err
	}

	res, err := client.Post(url, "application/json", bytes.NewBuffer([]byte(body)))

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, err
	}

	processed := reHex.ReplaceAll([]byte(body), []byte(`"$1"`))

	var progress TaskProgress
	json.Unmarshal(processed, &progress)

	transferred, err := strconv.ParseInt(progress.Transferred, 0, 64)

	if err != nil {
		panic("invalid hex value in console response")
	}

	total, err := strconv.ParseInt(progress.Total, 0, 64)

	if err != nil {
		panic("invalid hex value in console response")
	}

	if total == 0 {
		return 0, nil
	}

	x := new(big.Float).SetInt64(transferred)
	y := new(big.Float).SetInt64(total)
	z := new(big.Float).Quo(x, y)

	f, _ := new(big.Float).Mul(z, big.NewFloat(100.0)).Int64()

	return f, nil
}
