package manager

type Info struct {
	Id      string `json:"id"`
	Bash    bool   `json:"bash"`
	Command string `json:"command"`
	Panel   string `json:"panel"`
	Block   string `json:"block"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Error   string `json:"error"`
	Time    string `json:"time"`
}
