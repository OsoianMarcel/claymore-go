package claymore

import "strings"

// Stats json
type StatsJson struct {
	Id     int      `json:"error"`
	Result []string `json:"result"`
	Error  string   `json:"error"`
}

// Get result value by index
func (j StatsJson) getOneResult(index int) (string, bool) {
	if index > (len(j.Result) - 1) {
		return "", false
	}

	return j.Result[index], true
}

// Get result items by index
func (j StatsJson) getOneResultItems(index int) ([]string, bool) {
	r, ok := j.getOneResult(index)
	if !ok {
		return []string{}, false
	}

	return strings.Split(r, ";"), true
}

// Get result items by indexes
func (j StatsJson) getResultItem(resultIndex, itemIndex int) (string, bool) {
	items, ok := j.getOneResultItems(resultIndex)
	if !ok {
		return "", false
	}

	if itemIndex > (len(items) - 1) {
		return "", false
	}

	return items[itemIndex], true
}