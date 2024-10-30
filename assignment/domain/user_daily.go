package domain

import "github.com/glossd/pokergloss/gomq/mqtable"

type UserDaily struct {
	UserID string `bson:"_id"`
	DailyID
	Assignments []*UserAssignment
	Version     int64
}

func NewUserDaily(userID string, d *Daily) *UserDaily {
	var uas []*UserAssignment
	for _, a := range d.Assignments {
		if a != nil {
			uas = append(uas, NewUserAssignment(*a))
		}
	}
	return &UserDaily{
		UserID:      userID,
		DailyID:     d.Day,
		Assignments: uas,
	}
}

type UserAssignment struct {
	Assignment
	CurrentCount int64
	isChanged    bool
	isDone       bool
}

func NewUserAssignment(a Assignment) *UserAssignment {
	return &UserAssignment{
		Assignment:   a,
		CurrentCount: 0,
	}
}

func (ua *UserAssignment) IsDone() bool {
	return ua.CurrentCount >= ua.Assignment.Count
}

func (ua *UserAssignment) Inc() {
	ua.CurrentCount++
	ua.isChanged = true
	if ua.CurrentCount == ua.Assignment.Count {
		ua.isDone = true
	}
}

func (u *UserDaily) Update(p *mqtable.Player, ge *mqtable.GameEnd) (changed bool) {
	for _, ua := range u.Assignments {
		if ua.IsDone() {
			continue
		}
		matched := ua.matchGameEnd(p, ge)
		if matched {
			ua.Inc()
			changed = true
		}
	}
	return changed
}

func (u *UserDaily) UpdateForTournament(p *mqtable.Player, te *mqtable.TournamentEnd) (changed bool) {
	for _, ua := range u.Assignments {
		if ua.IsDone() {
			continue
		}
		matched := ua.matchTournamentEnd(p, te)
		if matched {
			ua.Inc()
			changed = true
		}
	}
	return changed
}

func (u *UserDaily) GetChanged() (assignments []*UserAssignment) {
	for _, a := range u.Assignments {
		if a.isChanged {
			assignments = append(assignments, a)
		}
	}
	return
}
