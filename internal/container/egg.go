package container

type Egg struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Startup     string            `json:"startup"`
	Env         map[string]string `json:"env"`
	Image       string            `json:"image"`
	Ports       []int             `json:"ports"`
	Volumes     []string          `json:"volumes"`
}
