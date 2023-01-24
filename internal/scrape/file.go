package scrape

type File struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Timestamp uint64 `json:"timestamp"`
	Cdn       string `json:"cdn"`
	I         string `json:"i"`
}
