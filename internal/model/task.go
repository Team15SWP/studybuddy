package model

import "study_buddy/pkg/constants"

type Statistics struct {
	Easy   int32 `json:"easy"`
	Medium int32 `json:"medium"`
	Hard   int32 `json:"hard"`
	Total  int32 `json:"total"`
}

type GeneratedTask struct {
	ID               int64        `json:"id"`
	UserID           int64        `json:"user_id"`
	TaskName         string       `json:"Task_name"`
	TaskDescription  string       `json:"Task_description"`
	SampleInputCases []InputCases `json:"Sample_input_cases"`
	Hints            Hints        `json:"Hints"`
	Solution         string       `json:"Full_solution"`
	Difficulty       int32        `json:"Difficulty"`
	Solved           int32        `json:"Solved"`
}

type Task struct {
	ID               int64        `json:"id"`
	TaskName         string       `json:"Task_name"`
	TaskDescription  string       `json:"Task_description"`
	SampleInputCases []InputCases `json:"Sample_input_cases"`
	Hints            Hints        `json:"Hints"`
	Difficulty       string       `json:"Difficulty"`
}

func (g *GeneratedTask) ToServer() *Task {
	if g == nil {
		return nil
	}
	return &Task{
		ID:               g.ID,
		TaskName:         g.TaskName,
		TaskDescription:  g.TaskDescription,
		SampleInputCases: g.SampleInputCases,
		Hints:            g.Hints,
		Difficulty:       constants.DifficultyToString(g.Difficulty),
	}
}

type InputCases struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type Hints struct {
	Hint1 string `json:"Hint1"`
	Hint2 string `json:"Hint2"`
	Hint3 string `json:"Hint3"`
}

type Question struct {
	Task     string `json:"task"`
	Code     string `json:"code"`
	Request  bool   `json:"question"`
	Verdict  bool   `json:"correct"`
	Feedback string `json:"feedback"`
}

type Feedback struct {
	Feedback string `json:"feedback"`
}

func (q *Question) ToFeedback() *Feedback {
	if q == nil {
		return nil
	}
	return &Feedback{
		Feedback: q.Feedback,
	}
}
