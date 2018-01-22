package claymore

import (
	"errors"
	"strconv"
	"fmt"
)

// Stats filler
type StatsFiller struct {
	json  *StatsJson
	model *StatsModel
}

// Create new instance of stats filler
func NewStatsFiller(json *StatsJson, model *StatsModel) StatsFiller {
	return StatsFiller{json, model}
}

// Execute fill process
func (sf *StatsFiller) Execute() error {
	// Check response error
	if sf.json.Error != "" {
		return fmt.Errorf("error response: %s", sf.json.Error)
	}

	// Check expected result length
	expLen := 9
	curLen := len(sf.json.Result)
	if expLen != curLen {
		return fmt.Errorf("current array length: %d, expected length: %d", expLen, curLen)
	}

	if err := sf.fillMinerVersion(); err != nil {
		return err
	}

	if err := sf.fillRunningMinutes(); err != nil {
		return err
	}

	if err := sf.fillTempAndFans(); err != nil {
		return err
	}

	if err := sf.fillPools(); err != nil {
		return err
	}

	if err := sf.fillCurrencyReports(); err != nil {
		return err
	}

	return nil
}

// Fill miner version
func (sf *StatsFiller) fillMinerVersion() error {
	res, ok := sf.json.getOneResult(0)
	if !ok {
		return errors.New("can not detect miner version")
	}

	sf.model.MinerVersion = res

	return nil
}

// Fill running minutes
func (sf *StatsFiller) fillRunningMinutes() error {
	res, ok := sf.json.getOneResult(1)
	if !ok {
		return errors.New("can not detect running minutes")
	}

	minutes, err := strconv.Atoi(res)
	if err != nil {
		return fmt.Errorf("can not convert mintes to int: %s", err)
	}
	sf.model.RunningMinutes = minutes

	return nil
}

// Fill pools
func (sf *StatsFiller) fillPools() error {
	items, ok := sf.json.getOneResultItems(7)
	if !ok {
		return errors.New("can not detect pools")
	}

	sf.model.Pools = items

	return nil
}

// Fill temperature and fans
func (sf *StatsFiller) fillTempAndFans() error {
	items, ok := sf.json.getOneResultItems(6)
	if !ok {
		return errors.New("can not detect temp and fans")
	}

	itemsCount := len(items)
	if itemsCount == 0 || itemsCount%2 != 0 {
		return fmt.Errorf("wrong number of items in temp and fan section: %d", itemsCount)
	}

	reports := make([]TempAndFanReport, 0, itemsCount/2)

	gpuIndex := 0
	for index := range items {
		if index%2 != 0 {
			continue
		}

		tf := TempAndFanReport{}

		temp, err := strconv.Atoi(items[index])
		if err != nil {
			return fmt.Errorf("can not convert temp to int: %s", err)
		}
		tf.Temp = temp

		fan, err := strconv.Atoi(items[index+1])
		if err != nil {
			return fmt.Errorf("can not convert fan to int: %s", err)
		}
		tf.Fan = fan

		tf.Gpu = gpuIndex
		gpuIndex = gpuIndex + 1

		reports = append(reports, tf)
	}

	sf.model.TempAndFanReports = reports

	return nil
}

// Fill currency reports
func (sf *StatsFiller) fillCurrencyReports() error {
	ethReport, err := sf.generateCurrencyReport(2, 3, 0)
	if err != nil {
		return err
	}

	sf.model.EthReport = ethReport

	altReport, err := sf.generateCurrencyReport(4, 5, 2)
	if err != nil {
		return err
	}

	sf.model.AltReport = altReport

	return nil
}

// Generate currency report by indexes
func (sf *StatsFiller) generateCurrencyReport(totalIndex, gpuIndex, invalidIndex int) (CurrencyReport, error) {
	emptyReport := CurrencyReport{}
	rep := CurrencyReport{}

	// Total items
	totalItems, ok := sf.json.getOneResultItems(totalIndex)
	if !ok || len(totalItems) != 3 {
		return emptyReport, fmt.Errorf("can not find total index: %d, or wrong number of items", totalIndex)
	}

	// Total Mhs
	totalMhs, err := strconv.Atoi(totalItems[0])
	if err != nil {
		return emptyReport, fmt.Errorf("can not convert total hash rate to int: %s", err)
	}
	rep.TotalMhs = totalMhs

	// Shares
	shares, err := strconv.Atoi(totalItems[1])
	if err != nil {
		return emptyReport, fmt.Errorf("can not convert shares to int: %s", err)
	}
	rep.Shares = shares

	// Rejected shares
	rejectedShares, err := strconv.Atoi(totalItems[2])
	if err != nil {
		return emptyReport, fmt.Errorf("can not convert rejected shares to int: %s", err)
	}
	rep.RejectedShares = rejectedShares

	// Gpu reports
	gpuItems, ok := sf.json.getOneResultItems(gpuIndex)
	if !ok || len(gpuItems) == 0 {
		return emptyReport, fmt.Errorf("can not find gpu index: %d, or no items", gpuIndex)
	}

	mhsPerGpu := make([]GpuMhs, 0, len(gpuItems))
	for i, mhs := range gpuItems {
		if mhs == "off" {
			mhsPerGpu = append(mhsPerGpu, GpuMhs{0, i})
			continue
		}

		n, err := strconv.Atoi(mhs)
		if err != nil {
			return emptyReport, fmt.Errorf("can not convert mhs to int: %s", err)
		}

		mhsPerGpu = append(mhsPerGpu, GpuMhs{n, i})
	}
	rep.MhsPerGpu = mhsPerGpu

	// Invalid shares
	invalidItems, ok := sf.json.getOneResultItems(8)
	countInvalidItems := len(invalidItems)
	if !ok || countInvalidItems != 4 || (invalidIndex+1) > (countInvalidItems-1) {
		return emptyReport, errors.New("can not detect invalid shares and pool switches")
	}

	invalidShares, err := strconv.Atoi(invalidItems[invalidIndex])
	if err != nil {
		return emptyReport, fmt.Errorf("can not convert invalid shares to int: %s", err)
	}
	rep.InvalidShares = invalidShares

	poolSwitches, err := strconv.Atoi(invalidItems[invalidIndex+1])
	if err != nil {
		return emptyReport, fmt.Errorf("can not convert invalid shares to int: %s", err)
	}
	rep.PoolSwitches = poolSwitches

	return rep, nil
}
