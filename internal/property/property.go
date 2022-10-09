package property

type property struct {
	ccid         string `json:"ccid"`
	slaveId      int    `json:"slaveId"`
	propertyType int    `json:"propertyType"`
	dataType     int    `json:"dataType"`
}

type Interface interface {
}
