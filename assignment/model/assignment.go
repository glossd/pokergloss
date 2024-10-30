package model

import "github.com/glossd/pokergloss/assignment/domain"

type Assignment struct {
	Type domain.AssignmentType `json:"type"`

	Prefix   string      `json:"prefix"`
	Variable interface{} `json:"variable"`
	Suffix   string      `json:"suffix"`

	Prize        int64 `json:"prize"`
	Count        int64 `json:"count"`
	CurrentCount int64 `json:"currentCount"`
}

func ToAssignment(a *domain.UserAssignment) *Assignment {
	if a == nil {
		return nil
	}
	prefix, variable, suffix := a.GetName()
	return &Assignment{
		Type:         a.Type,
		Prefix:       prefix,
		Variable:     variable,
		Suffix:       suffix,
		Prize:        a.GetPrize(),
		Count:        a.Count,
		CurrentCount: a.CurrentCount,
	}
}

func FromUserDaily(ud *domain.UserDaily) (assignments []*Assignment) {
	for _, a := range ud.Assignments {
		assignments = append(assignments, ToAssignment(a))
	}
	return assignments
}
