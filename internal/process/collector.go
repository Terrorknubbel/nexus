package process

import (
	"bufio"
	"nexus/pkg/models"
	"os"
	"strconv"
	"strings"
)

func Collect() []models.ProcessInfo {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	processes := []models.ProcessInfo{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}

		process := models.ProcessInfo{}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			process.PID = -1
		} else {
			process.PID = pid
		}

		commPath := "/proc/" + entry.Name() + "/comm"
		if f, err := os.Open(commPath); err != nil {
			process.Name = "unknown"
		} else {
			defer f.Close()
			r := bufio.NewReader(f)
			line, _ := r.ReadString('\n')
			process.Name = strings.TrimSpace(line)
		}

		StatePath := "/proc/" + entry.Name() + "/status"
		if f, err := os.Open(StatePath); err != nil {
			println("Cannot open", StatePath, ":", err.Error())
			process.State = "unknown"
		} else {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			process.State = "unknown" // Default
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "State:") {
					// Beispiel: "State:\tS (sleeping)"
					// Wir extrahieren den Text in den Klammern
					if idx := strings.Index(line, "("); idx != -1 {
						if idx2 := strings.Index(line[idx:], ")"); idx2 != -1 {
							process.State = line[idx+1 : idx+idx2]
						}
					}
					break
				}
			}
		}

		processes = append(processes, process)
	}

	return processes
}
