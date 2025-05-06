package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/KTemka1234/Yadro-Impulse-2025-2/config"
	"github.com/KTemka1234/Yadro-Impulse-2025-2/logger"
)

func main() {
	os.Mkdir("logs", 0755)
	eventsLogger := logger.InitLogger("./logs/competiotion.log")
	defer logger.Cleanup()

	c := config.ParseConfig("./config/config_xmpl.json")
	events, err := parseEvents("events_xmpl")
	if err != nil {
		eventsLogger.Fatalf("Cannot parse events file: %v", err)
	}
	eventsLogger.Println("Events file succsessfully parsed!")

	comps := make(map[int]*Competitor)
	var outEvents []Event
	for _, event := range events {
		handleEvent(event, c, comps, &outEvents, eventsLogger)
	}

	generateReport("./logs/report.log", comps, c)
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

func generateReport(filename string, comps map[int]*Competitor, config config.Config) {
	resultLogger := logger.InitLogger(filename)
	defer logger.Cleanup()

	// Copy and sort competitors by finish time
	var results []*Competitor
	for _, c := range comps {
		results = append(results, c)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].EndTime.Before(results[j].EndTime)
	})

	resultLogger.Println("\nResult report. Report structure represents below")
	resultLogger.Println("[Total time/status] ID [{lap #1 time, lap avg. speed},...] {total penalty time, penalty avg. speed} hit/shot\n")

	for _, comp := range results {
		var status string
		switch {
		case comp.DSQ:
			status = "[NotStarted]"
		case comp.Status == "NotFinished":
			status = "[NotFinished]"
		default:
			if comp.LapsCompleted >= config.Laps {
				totalTime := comp.EndTime.Sub(comp.ActualStartTime)
				status = fmt.Sprintf("[%s]", formatDuration(totalTime))
			} else {
				status = "[NotFinished]"
			}
		}

		var lapStats strings.Builder
		lapStats.WriteRune('[')
		for i := 0; i < config.Laps; i++ {
			if i < len(comp.LapTimes) {
				lapDuration := formatDuration(comp.LapTimes[i])
				lapAvgSpeed := float64(config.LapLen) / comp.LapTimes[i].Seconds()
				lapStats.WriteString(fmt.Sprintf("{%s, %0.3f}", lapDuration, lapAvgSpeed))
			} else {
				lapStats.WriteString("{,}")
			}

			if i < config.Laps-1 {
				lapStats.WriteString(", ")
			}
		}
		lapStats.WriteRune(']')

		penaltySpeed := 0.0
		if comp.TotalPenaltyTime > 0 {
			penaltySpeed = float64(config.PenaltyLen*comp.PenaltyLaps) / comp.TotalPenaltyTime.Seconds()
		}
		penaltyTime := formatDuration(comp.TotalPenaltyTime)
		penaltyStats := fmt.Sprintf("{%s, %0.3f}", penaltyTime, penaltySpeed)

		hitShot := fmt.Sprintf("%d/%d", comp.Hits, len(comp.FiringSessions)*5)

		resultLogger.Printf("%s %d %s %s %s\n",
			status,
			comp.Id,
			lapStats.String(),
			penaltyStats,
			hitShot,
		)
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	millis := d.Milliseconds() % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}
