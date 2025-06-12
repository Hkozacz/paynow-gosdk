package main

type PayloadToSign struct {
	Headers    map[string]string      `json:"headers"`
	Parameters map[string]interface{} `json:"parameters"`
	Body       string                 `json:"body"`
}
