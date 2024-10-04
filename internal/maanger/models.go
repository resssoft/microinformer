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

	Modal  bool `json:"modal,omitempty"`
	KeepBy int  `json:"keep_by,omitempty"`
}

type ImportItems struct {
	Items []Info `json:"items"`
}

type AddedResult struct {
	Count    int    `json:"count"`
	Excluded []Info `json:"excluded"`
}

type Web struct {
	Url string `json:"url"`
}

type Command struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type Graphics struct {
}
