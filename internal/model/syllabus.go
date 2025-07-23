package model

type Syllabus struct {
	Topics []string `json:"topics"`
}

type Schedule struct {
	Week  string
	Topic string
}
