package main

import (
	//_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

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

//TODO: step to step load by setTImeout
//TODO: dynamic change commands list
//TODO: dynamic setTImeout time from settings and changed by ip
//TODO: more often update for app settings
//TODO: save to disk items list (json file)
//TODO: one time display info show
//TODO: ps list monitor
//TODO: set custom info to api from others machines (bot statuses, cpu, home disks, pinger by the others servers)
//TODO: show images (grafics)
//TODO: parser
// DRm credit, zkh
////go:embed frontend/index.html
//var pageData string

var items []Info

var settings Settings

type request struct {
	Info   []Info `json:"info"`
	Reboot bool   `json:"reboot"`
}

type response struct {
	Info     []Info   `json:"info"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	Timeout int  `json:"timeout"`
	Reboot  bool `json:"reboot"`
}

func main() {
	fmt.Println("Start server by :8081")
	//fmt.Println(pageData)
	configure()
	http.HandleFunc("/api.json", api)
	http.HandleFunc("/update", update)
	http.HandleFunc("/settings.json", setting)
	//http.HandleFunc("/page", page)
	http.Handle("/page/", http.StripPrefix("/page", http.FileServer(http.Dir("./frontend"))))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func configure() {
	items = append(items, newItem("date", "System date", "right-foot", "main", false))
	items = append(items, newItem("hostname && hostname -I", "host", "main", "main", true))
	items = append(items, newItem("cat /proc/meminfo | grep 'MemFree:'", "meminfo", "main", "main", true))
	items = append(items, newItem("ping -c 1 8.8.8.8 | grep packet", "ping88", "main", "main", true))
	items = append(items, valItem("go version", "main", "content", runtime.Version()))
	data, err := json.Marshal(items)
	fmt.Println(string(data), err)
}

func getInfo() []Info {
	var list []Info
	list = items
	for index, item := range list {
		list[index] = item.run()
	}
	//list = append(list, Command("date", "System date", "right-foot", "main"))
	//list = append(list, Command("hostname && hostname -I", "host", "main", "main", true))
	//list = append(list, Command("hostname", "hostname", "main", "main"))
	//list = append(list, Command("hostname -I", "hostname ip", "main", "main"))
	//list = append(list, Command("cat /proc/meminfo | grep 'MemFree:'", "meminfo", "main", "main", true))
	//list = append(list, valItem("go version", "main", "content", runtime.Version()))
	//list = append(list, Command("ping -c 1 8.8.8.8 | grep packet", "ping88", "main", "main", true))
	return list
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func valItem(name, block, panel, value string) Info {
	return Info{
		Id:    strings.Replace(block+name, " ", "", -1),
		Panel: panel,
		Block: block,
		Name:  name,
		Value: value,
	}
}

func newItem(command, name, block, panel string, bash bool) Info {
	return Info{
		Id:      strings.Replace(block+name, " ", "", -1),
		Command: command,
		Panel:   panel,
		Block:   block,
		Name:    name,
		Bash:    bash,
	}
}

func (i Info) run() Info {
	switch i.Name {
	case "go version":
		i.Value = runtime.Version()
		return i
	}
	i = Command(i)
	return i
}

// func Command(command, name, block, panel string, bash bool) Info {
func Command(i Info) Info {
	var out string
	var err error
	var code int
	start := time.Now()

	if i.Bash {
		out, code, err = runRaw("bash", []string{"-c", i.Command}, []string{})
	} else {
		slittedCommand, params := SplitParams(i.Command, " ")
		out, code, err = runRaw(slittedCommand, params, []string{})
	}
	if err != nil {
		i.Error = err.Error()
	}
	i.Value = out
	if code != 0 {
		i.Error += fmt.Sprintf(" exit code: %d", code)
	}
	i.Time = duration(start)
	//out, err := exec.Command(command).Output()
	//item.Value = string(out)
	return i
}

func update(w http.ResponseWriter, r *http.Request) {
	var data response
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if data.Settings.Timeout != 0 {
		settings = data.Settings
	}
	if len(data.Info) != 0 {
		w.WriteHeader(http.StatusCreated)
		items = data.Info
		return
	}
	w.WriteHeader(http.StatusOK)
}

func setting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(settings)
	if err != nil {
		_, err := fmt.Fprintf(w, "Hi")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	data := request{Info: getInfo(), Reboot: settings.Reboot}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, err := fmt.Fprintf(w, "Hi")
		if err != nil {
			fmt.Println(err)
		}
	}
	if settings.Reboot {
		settings.Reboot = false
	}
}

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, pageData)
}

func runRaw(command string, params []string, env []string) (string, int, error) {
	var exitCode int
	cmd := exec.Command(command, params...)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		//if exitError, ok := err.(*exec.ExitError); ok {
		//	exitCode = exitError.ExitCode()
		//}
		return string(stdout), exitCode, err
	}
	return string(stdout), exitCode, nil
}

func SplitParams(command, separator string) (string, []string) {
	slitted := strings.Split(command, separator)
	if len(slitted) == 0 {
		return command, []string{}
	}
	return slitted[0], slitted[1:]
}

func duration(t time.Time) string {
	result := "?"
	dur := time.Since(t)
	switch {
	case dur < time.Microsecond:
		result = fmt.Sprintf("%dnSec", dur/time.Nanosecond)
	case dur < time.Millisecond:
		result = fmt.Sprintf("%dmcSec", dur/time.Microsecond)
	case dur < time.Second:
		result = fmt.Sprintf("%dmlSec", dur/time.Millisecond)
	case dur < time.Minute:
		result = fmt.Sprintf("%dSec", dur/time.Second)
	case dur > time.Minute:
		result = fmt.Sprintf("%v", dur)
	}
	return result
}

var pageData = ``
