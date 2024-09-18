package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"microinformer/internal/repository"
	"microinformer/internal/settings"
)

type Service struct {
	Items    []Info
	Settings *settings.Service
	mutex    sync.Mutex
	repo     *repository.Service
}

func NewService(
	settingsService *settings.Service) *Service {
	return &Service{
		Settings: settingsService,
		repo:     repository.NewFileRepo("items.json"),
	}
}

func (s *Service) load() {
	rawData, err := s.repo.Load()
	if err != nil {
		fmt.Println("manager load from disk err", err)
	} else {
		var items []Info
		err = json.Unmarshal(rawData, &items)
		if err != nil {
			fmt.Println("manager Unmarshal items err", err)
		} else {
			s.Items = items
		}
	}
}

func (s *Service) save() {
	data, err := json.MarshalIndent(s.Items, "", "  ")
	if err != nil {
		fmt.Println("manager marshalling items err", err)
	}
	err = s.repo.Save(data)
	if err != nil {
		fmt.Println("manager save to disk err", err)
	}
}

func (s *Service) Configure() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.load()
	if len(s.Items) == 0 {
		s.Items = append(s.Items, s.NewItem("date", "System date", "right-foot", "main", false))
		s.Items = append(s.Items, s.NewItem("hostname && hostname -I", "host", "main", "main", true))
		s.Items = append(s.Items, s.NewItem("cat /proc/meminfo | grep 'MemFree:'", "meminfo", "main", "main", true))
		s.Items = append(s.Items, s.NewItem("ping -c 1 8.8.8.8 | grep packet", "ping88", "main", "main", true))
		s.Items = append(s.Items, s.NewItem("", "go version", "main", "content", true))
	}
	s.save()
}

func (s *Service) GetInfo() []Info {
	var list []Info
	list = s.Items
	for index, item := range list {
		list[index] = s.run(item)
		if item.Once {
			_ = s.DelItem(item)
		}
	}
	return list
}

func (s *Service) ValItem(name, block, panel, value string) Info {
	return Info{
		Id:    strings.Replace(block+name, " ", "", -1),
		Panel: panel,
		Block: block,
		Name:  name,
		Value: value,
	}
}

func (s *Service) NewItem(command, name, block, panel string, bash bool) Info {
	return Info{
		Id:      strings.Replace(block+name, " ", "", -1),
		Command: command,
		Panel:   panel,
		Block:   block,
		Name:    name,
		Bash:    bash,
	}
}

func (s *Service) run(i Info) Info {
	switch i.Name {
	case "go version":
		i.Value = runtime.Version()
		return i
	}
	i = s.Command(i)
	return i
}

func (s *Service) ListItem() []Info {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Items
}

func (s *Service) AddItem(i Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if i.Command == "" || i.Name == "" {
		return errors.New("invalid command")
	}
	s.Items = append(s.Items, i)
	s.save()
	return nil
}

func (s *Service) Update(list []Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Items = list
	s.save()
	return nil
}

func (s *Service) DelItem(i Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newList := make([]Info, 0)
	for _, item := range s.Items {
		if item.Name != i.Name && item.Command != i.Command {
			newList = append(newList, item)
		}
	}
	s.Items = newList
	s.save()
	return nil
}

func (s *Service) RunItem(i Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newList := make([]Info, 0)
	for _, item := range s.Items {
		if item.Name != i.Name && item.Command != i.Command {
			newList = append(newList, item)
		}
	}
	s.Items = newList
	s.save()
	return nil
}

// func Command(command, name, block, panel string, bash bool) Info {
func (s *Service) Command(i Info) Info {
	var out string
	var err error
	var code int
	start := time.Now()

	if i.Bash {
		out, code, err = s.RunRaw("bash", []string{"-c", i.Command}, []string{})
	} else {
		slittedCommand, params := SplitParams(i.Command, " ")
		out, code, err = s.RunRaw(slittedCommand, params, []string{})
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

func (s *Service) RunRaw(command string, params []string, env []string) (string, int, error) {
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
