package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KTemka1234/Yadro-Impulse-2025-2/config"
	"github.com/KTemka1234/Yadro-Impulse-2025-2/logger"
)

func main() {
	myLogger := logger.InitLogger("./logs/competiotion.log")
	defer logger.Cleanup()

	config := config.ParseConfig("./config/config.json")
	events, err := parseEvents("events")
	if err != nil {
		myLogger.Fatalf("Cannot parse events file: %v", err)
	}
	myLogger.Println("Events file succsessfully parsed!")
	
	comps := make(map[int]*Competitor)
	var outEvents []Event
	for _, event := range events {
		handleEvent(event, config, comps, &outEvents, myLogger)
	}
}

func parseEvents(filename string) ([]Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file with filename %s: %v", filename, err)
	}

	lines := strings.Split(string(data), "\n")
	var events []Event
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "]", 2)
		if len(parts) != 2 {
			continue
		}

		timeStr := strings.TrimSpace(parts[0][1:])
		t, err := time.Parse(config.TimeMillis, timeStr)
		if err != nil {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 2 {
			continue
		}

		eventId, _ := strconv.Atoi(fields[0])
		competitorId, _ := strconv.Atoi(fields[1])
		params := fields[2:]
		event := Event{t, eventId, competitorId, params}
		events = append(events, event)
	}
	return events, nil
}

func getCompetitor(comps map[int]*Competitor, id int) *Competitor {
	if c, ok := comps[id]; ok {
		return c
	}

	c := &Competitor{Id: id}
	comps[id] = c
	return c
}
