package main

import "time"

type Event struct {
	Time         time.Time
	Id           int
	CompetitorId int
	ExtraParams  []string
}

type Competitor struct {
	Id                int
	StartTime         time.Time
	ActualStartTime   time.Time
	LapStartTime      time.Time
	LapsCompleted     int
	InPenalty         bool
	PenaltyStart      time.Time
	PenaltyLaps       int
	TotalPenaltyTime  time.Duration
	CurrentFiringLine int
	FiringSessions    []FiringSession
	Status            string
	DSQ               bool
	Hits              int
	EndTime           time.Time
	LastEventTime     time.Time
}

type FiringSession struct {
	Lap   int
	Line  int
	Hits  map[int]bool
	Start time.Time
	End   time.Time
}
