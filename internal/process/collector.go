package process

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"nexus/pkg/models"
)

func Collect() []models.Process {
	return collectFrom("/proc")
}

func collectFrom(procRoot string) []models.Process {
	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return nil
	}

	var processes []models.Process
	for _, entry := range entries {
		info, err := readProcessInfo(procRoot, entry)
		if err != nil {
			continue
		}
		processes = append(processes, info)
	}
	return processes
}

func readProcessInfo(procRoot string, entry os.DirEntry) (models.Process, error) {
	pid, err := parsePID(entry)
	if err != nil {
		return models.Process{}, err
	}

	line, err := readStatLine(procRoot, pid)
	if err != nil {
		return models.Process{}, err
	}

	name := extractName(line)
	state := extractState(line)

	return models.Process{PID: pid, Name: name, State: state}, nil
}

func parsePID(entry os.DirEntry) (int, error) {
	if !entry.IsDir() {
		return 0, os.ErrInvalid
	}
	pidStr := entry.Name()
	return strconv.Atoi(pidStr)
}

func readStatLine(procRoot string, pid int) (string, error) {
	path := filepath.Join(procRoot, strconv.Itoa(pid), "stat")
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

func extractName(line string) string {
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start < 0 || end < start+1 {
		return "unknown"
	}
	extracted := line[start+1 : end]
	if extracted == "?" {
		return "unknown"
	}
	return extracted
}

func extractState(line string) string {
	end := strings.LastIndex(line, ")")
	if end < 0 || len(line) <= end+1 {
		return "unknown"
	}

	rest := strings.TrimSpace(line[end+1:])
	fields := strings.Fields(rest)
	if len(fields) == 0 || len(fields[0]) == 0 {
		return "unknown"
	}
	stateChar := fields[0][0]

	switch stateChar {
	case 'R':
		return "running"
	case 'S':
		return "sleeping"
	case 'D':
		return "uninterruptible sleep"
	case 'Z':
		return "zombie"
	case 'T', 't':
		return "stopped"
	case 'X', 'x':
		return "dead"
	case 'K':
		return "wakekill"
	case 'W':
		return "waking"
	case 'P':
		return "parked"
	case 'I':
		return "idle"
	default:
		return string(stateChar)
	}
}
