package model

type Package struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Queries map[string]string `json:"queries"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
