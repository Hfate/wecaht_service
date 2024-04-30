package ai

type AvgTime struct {
	ID       uint64 `json:"id"`
	SpanTime int64  `json:"spanTime"`
}

func (AvgTime) TableName() string {
	return "avg_time"
}
