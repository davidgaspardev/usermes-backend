package valueobjects

type CreateTask struct {
	Name      string `json:"name"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
}
