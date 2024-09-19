package valueobjects

import "time"

/**
- id: u64
- name: string
- resCode: string
- cronCode: string
- dateStart: Date
- dateEnd: Date
- deps: u64[]
*/

type CreateTask struct {
	Name      string    `json:"name"`
	ResCode   string    `json:"resCode"`
	CronCode  string    `json:"cronCode"`
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
	Deps      []uint64  `json:"deps"`
}
