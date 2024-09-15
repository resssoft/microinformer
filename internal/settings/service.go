package settings

import (
	"sync"
)

type Public struct {
	Timeout int  `json:"timeout"`
	Reboot  bool `json:"reboot"`
}

type Service struct {
	Data  *Public
	mutex sync.Mutex
}

func NewService() *Service {
	return &Service{
		Data: &Public{
			Reboot: true,
		},
	}
}

func (s *Service) Taught() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data.Reboot = false
}

func (s *Service) Get() *Public {
	return s.Data
}

func (s *Service) Set(p *Public) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Data = p
}
