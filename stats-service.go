package claymore

import (
	"fmt"
	"encoding/json"
)

// Stats service
type StatsService struct {
	conn Connection
}

// Create new stats service instance
func NewStatsService(conn Connection) StatsService {
	return StatsService{conn}
}

// Request json
func (s StatsService) requestJson() (StatsJson, error) {
	jsonReq := `{"id":0,"jsonrpc":"2.0","method":"miner_getstat1"}`

	resp, err := s.conn.Request([]byte(jsonReq))
	if err != nil {
		return StatsJson{}, err
	}

	var sj StatsJson
	err = json.Unmarshal(resp, &sj)
	if err != nil {
		return StatsJson{}, fmt.Errorf("can not unmarshal stats: %s", err.Error())
	}

	return sj, nil
}

// Execute
func (s StatsService) Execute() (StatsModel, error) {
	sj, err := s.requestJson()
	if err != nil {
		return StatsModel{}, err
	}

	sm := StatsModel{}

	sf := NewStatsFiller(&sj, &sm)
	err = sf.Execute()
	if err != nil {
		return StatsModel{}, err
	}

	return sm, nil
}