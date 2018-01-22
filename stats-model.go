package claymore

import (
	"fmt"
	"errors"
)

// Currency report
type GpuMhs struct {
	Mhs int `json:"mhs"`
	Gpu int `json:"gpu"`
}

// Currency report
type CurrencyReport struct {
	TotalMhs       int   `json:"total_mhs"`
	Shares         int   `json:"shares"`
	RejectedShares int   `json:"rejected_shares"`
	InvalidShares  int   `json:"invalid_shares"`
	PoolSwitches   int   `json:"pool_switches"`
	MhsPerGpu      []GpuMhs `json:"mhs_per_gpu"`
}

// Temp and fan
type TempAndFanReport struct {
	Temp int `json:"temp"`
	Fan  int `json:"fan"`
	Gpu  int `json:"gpu"`
}

// String transformation
func (tfr TempAndFanReport) String() string {
	return fmt.Sprintf("GPU%d | Temp: +%d C | Fan: %d%% RPM", tfr.Gpu, tfr.Temp, tfr.Fan)
}

// Stats model
type StatsModel struct {
	MinerVersion      string             `json:"miner_version"`
	RunningMinutes    int                `json:"running_minutes"`
	EthReport         CurrencyReport     `json:"eth_report"`
	AltReport         CurrencyReport     `json:"alt_report"`
	TempAndFanReports []TempAndFanReport `json:"temp_and_fan_reports"`
	Pools             []string           `json:"pools"`
}

// Get highest temp report
func (sm StatsModel) GetHighestTemp() (TempAndFanReport, error) {
	if len(sm.TempAndFanReports) == 0 {
		return TempAndFanReport{}, errors.New("empty temp and fans report list")
	}

	highest := struct {
		index int
		temp  int
	}{
		-1,
		-1,
	}

	for i, v := range sm.TempAndFanReports {
		if v.Temp > highest.temp {
			highest.temp = v.Temp
			highest.index = i
		}
	}

	return sm.TempAndFanReports[highest.index], nil
}
