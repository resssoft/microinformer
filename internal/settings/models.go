package settings

import (
	"sync"

	"microinformer/internal/repository"
)

const (
	Version = "1.0.0.14"
)

type Public struct {
	Timeout int    `json:"timeout"`
	Reboot  bool   `json:"reboot"`
	Panel   *Panel `json:"panel"`
	Version string `json:"version"`
}

type Panel struct {
	ShowCount bool  `json:"show_count"`
	Rows      []Row `json:"rows"`
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

func defaultPanel() *Panel {
	return &Panel{
		Rows: []Row{
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
		},
	}
}
