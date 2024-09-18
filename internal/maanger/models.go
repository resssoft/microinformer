package manager

type Info struct {
	Id      string `json:"id,omitempty"`
	Bash    bool   `json:"bash,omitempty"`
	Command string `json:"command"`
	Panel   string `json:"panel,omitempty"`
	Block   string `json:"block,omitempty"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Value   string `json:"value,omitempty"`
	Error   string `json:"error,omitempty"`
	Time    string `json:"time,omitempty"`

	Once   bool `json:"once,omitempty"`
	KeepBy int  `json:"keep_by,omitempty"`
}
