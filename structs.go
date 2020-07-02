package main

// Challenges describes the challenges as returned from the /acpi/v1/challenges
// endpoint
type Challenges struct {
	Success bool `json:"success"`
	Data    []struct {
		ID       int           `json:"id"`
		Type     string        `json:"type"`
		Name     string        `json:"name"`
		Value    int           `json:"value"`
		Category string        `json:"category"`
		Tags     []interface{} `json:"tags"`
		Template string        `json:"template"`
		Script   string        `json:"script"`
	} `json:"data"`
}

// Challenge describes a single challenge as returned from the
// /acpi/v1/challenges/<id> endpoint
type Challenge struct {
	Success bool `json:"success"`
	Data    struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Value       int    `json:"value"`
		Description string `json:"description"`
		Category    string `json:"category"`
		State       string `json:"state"`
		MaxAttempts int    `json:"max_attempts"`
		Type        string `json:"type"`
		TypeData    struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Templates struct {
				Create string `json:"create"`
				Update string `json:"update"`
				View   string `json:"view"`
			} `json:"templates"`
			Scripts struct {
				Create string `json:"create"`
				Update string `json:"update"`
				View   string `json:"view"`
			} `json:"scripts"`
		} `json:"type_data"`
		Solves int           `json:"solves"`
		Files  []string      `json:"files"`
		Tags   []interface{} `json:"tags"`
		Hints  []interface{} `json:"hints"`
	} `json:"data"`
}

