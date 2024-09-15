package settings

import (
	"encoding/json"
	"fmt"
	"sync"

	"microinformer/internal/repository"
)

type Public struct {
	Timeout int    `json:"timeout"`
	Reboot  bool   `json:"reboot"`
	Panel   *Panel `json:"panel"`
}

type Panel struct {
	Rows []Row `json:"rows"`
}

type Row struct {
	Id     string  `json:"id"`
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Id     string `json:"id"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type Service struct {
	Data  *Public
	repo  *repository.Service
	mutex sync.Mutex
}

func NewService() *Service {
	s := &Service{
		Data: &Public{
			Timeout: 5000,
		},
		repo: repository.NewFileRepo("settings.json"),
	}
	s.load()
	s.Data.Reboot = true
	if s.Data.Panel == nil {
		s.Data.Panel = &Panel{}
		s.Data.Panel.Rows = []Row{
			{
				Id: "header",
				Blocks: []Block{
					{
						Id: "error",
					},
				},
			},
			{
				Id: "middle",
				Blocks: []Block{
					{
						Id:    "left-side",
						Width: "15%",
					},
					{
						Id:    "main",
						Width: "60%",
					},
					{
						Id:    "right-side",
						Width: "15%",
					},
				},
			},
			{
				Id: "footer",
				Blocks: []Block{
					{
						Id: "left-foot",
					},
					{
						Id: "right-foot",
					},
				},
			},
		}
	}
	s.save()
	return s
}

func (s *Service) NoReboot() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data.Reboot = false
}

func (s *Service) SetReboot() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data.Reboot = true
}

func (s *Service) Get() *Public {
	return s.Data
}

func (s *Service) Set(p *Public) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data = p
	s.save()
}

func (s *Service) load() {
	rawData, err := s.repo.Load()
	if err != nil {
		fmt.Println("manager load from disk err", err)
	} else {
		var items Public
		err = json.Unmarshal(rawData, &items)
		if err != nil {
			fmt.Println("manager Unmarshal items err", err)
		} else {
			s.Data = &items
		}
	}
}

func (s *Service) save() {
	data, err := json.MarshalIndent(&s.Data, "", "  ")
	if err != nil {
		fmt.Println("manager marshalling items err", err)
	}
	err = s.repo.Save(data)
	if err != nil {
		fmt.Println("manager save to disk err", err)
	}
}
