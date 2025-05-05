package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	c "github.com/KTemka1234/Yadro-Impulse-2025-2/config"
	l "github.com/KTemka1234/Yadro-Impulse-2025-2/logger"
)

type EventHandler func(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger)

var eventHandlers = map[int]EventHandler{
	1:  handleRegistration,
	2:  handleStartTimeSet,
	3:  handleCompetitorOnStart,
	4:  handleCompStart,
	5:  handleCompOnFireRange,
	6:  handleTargetHit,
	7:  handleCompLeftFireRange,
	8:  handleCompOnPenalty,
	9:  handleCompLeftPenalty,
	10: handleCompLapEnd,
	11: handleCompDNF,
	32: handleCompDSQ,
	33: handleCompFinish,
}

// Incoming event handler
func handleEvent(event Event, config c.Config, comps map[int]*Competitor, out *[]Event, logger *log.Logger) {
	comp := getCompetitor(comps, event.CompetitorId)
	if handler, exists := eventHandlers[event.Id]; exists {
		handler(event, comp, config, out, logger)
	}
	comp.LastEventTime = event.Time
}

// 1. A competitor registered
func handleRegistration(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.Status = "Registered"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 2. Competitor's start time was set by a draw
func handleStartTimeSet(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	if len(event.ExtraParams) == 0 {
		panic("Incoming event error: no startTime value specified")
	}

	startTime, err := time.Parse("15:04:05.000", event.ExtraParams[0])
	if err != nil {
		panic("Can't parse startTime in events: " + err.Error())
	}
	comp.StartTime = startTime
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id, comp.StartTime.Format(c.TimeMillis))
}

// 3. A competitor is on the start line
func handleCompetitorOnStart(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.Status = "OnStartLine"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 4. A competitor has started
func handleCompStart(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	durations := strings.Split(config.StartDelta, ":")
	if len(durations) < 3 {
		panic("Invalid config: startDelta must have hh:mm:ss time format")
	}

	strDuration := durations[0] + "h" + durations[1] + "m" + durations[2] + "s"
	startDelta, err := time.ParseDuration(strDuration)
	if err != nil {
		panic("Can't parse config startDelta time: " + err.Error())
	}

	maxStartTime := comp.StartTime.Add(startDelta)

	// DSQ
	if event.Time.After(maxStartTime) {
		dsqEvent := Event{event.Time, 32, comp.Id, nil}
		handleCompDSQ(dsqEvent, comp, config, nil, logger)
		*out = append(*out, dsqEvent)
		return
	}

	comp.ActualStartTime = event.Time
	comp.LapStartTime = event.Time
	comp.Status = "Racing"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 5. A competitor is on the firing range
func handleCompOnFireRange(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	if len(event.ExtraParams) == 0 {
		return
	}
	firingLine, err := strconv.Atoi(event.ExtraParams[0])
	if err != nil {
		panic("Incoming event error: can't parse firingRange number")
	}

	comp.FiringSessions = append(comp.FiringSessions, FiringSession{
		Lap:   comp.LapsCompleted + 1,
		Line:  firingLine,
		Start: event.Time,
		Hits:  make(map[int]bool),
	})
	comp.Status = fmt.Sprintf("Shooting firing line #%d", firingLine)
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id, firingLine)
}

// 6. Target hit
func handleTargetHit(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	if len(comp.FiringSessions) == 0 {
		panic("Logic error: competitor does not have shooting sessions")
	}
	if len(event.ExtraParams) == 0 {
		panic("Incoming event error: no target hit number specified")
	}

	target, err := strconv.Atoi(event.ExtraParams[0])
	if err != nil {
		panic("Incoming event error: can't parse hit target number")
	}
	session := &comp.FiringSessions[len(comp.FiringSessions)-1]

	if !session.Hits[target] {
		session.Hits[target] = true
		comp.Hits++
	} else {
		panic(fmt.Sprintf("Competitor %d has already hit target %d", comp.Id, target))
	}
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), target, comp.Id)
}

// 7. A competitor left the firing range
func handleCompLeftFireRange(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	if len(comp.FiringSessions) == 0 {
		panic("Logic error: a competitor does not have any firing sessions")
	}

	session := &comp.FiringSessions[len(comp.FiringSessions)-1]
	session.End = event.Time

	// Target miss
	missed := 5 - len(session.Hits) // 5 targets on firing line
	if missed > 0 {
		comp.PenaltyLaps += missed // 1 penalty lap by 1 target miss
		comp.Status = "InPenalty"
	}
	comp.Status = "Racing"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 8. A competitor entered the penalty laps
func handleCompOnPenalty(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.PenaltyStart = event.Time
	comp.Status = "InPenalty"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 9. A competitor left the penalty laps
func handleCompLeftPenalty(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	if !comp.PenaltyStart.IsZero() {
		comp.TotalPenaltyTime += event.Time.Sub(comp.PenaltyStart)
		comp.PenaltyStart = time.Time{}
	}
	comp.Status = "Racing"
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 10. A competitor ended the main lap
func handleCompLapEnd(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.LapsCompleted++

	// Finish state
	if comp.LapsCompleted == config.Laps {
		outEvent := Event{event.Time, 33, comp.Id, nil}
		handleCompFinish(outEvent, comp, config, nil, logger)
		*out = append(*out, outEvent)
		comp.EndTime = event.Time
		comp.Status = "Finished"
		return
	}

	comp.LapStartTime = event.Time
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 11. A competitor DNF
func handleCompDNF(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.Status = "DNF"
	if len(event.ExtraParams) > 0 {
		msg := strings.Join(event.ExtraParams, " ")
		comp.Status += ": " + msg
		logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id, msg)
	} else {
		logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
	}
	comp.EndTime = event.Time
}

// 32. A competitor DSQ
func handleCompDSQ(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.DSQ = true
	comp.Status = "Disqualified (Not Started)"
	comp.EndTime = event.Time
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}

// 33. A competitor has finished
func handleCompFinish(event Event, comp *Competitor, config c.Config, out *[]Event, logger *log.Logger) {
	comp.Status = "Finished"
	comp.EndTime = event.Time
	logger.Printf(l.EventsLog[event.Id], event.Time.Format(c.TimeMillis), comp.Id)
}
