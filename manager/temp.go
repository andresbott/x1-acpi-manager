package manager

import (
	"fmt"
	"strconv"
)

const TempProveReadings = 5
const warnTemp = 75      // throttle gpu,  watt throttling
const criticalTemp = 85  // more watt throttling, throttle cpu
const emergencyTemp = 95 // cpu boost off extreme watt throttling

// ===============================================================
// Temp
// ===============================================================
const (
	tempPath = "/sys/class/thermal/thermal_zone0/temp"
)

type tempStat int

const (
	tempOk tempStat = iota
	tempWarn
	tempCritical
	tempEmergency
)

type tempManager struct {
	data          []int
	Status        tempStat
	prevStat      tempStat
	StatusChanged bool
}

func tempRead() int {
	tempStr, err := read(tempPath)
	if err != nil {
		handleErr(err)
	}
	tempInt, err := strconv.Atoi(tempStr)
	val := tempInt / 1000
	return val
}

// read temperature and add to the internal list
func (t *tempManager) read() {
	temp := tempRead()

	t.data = append(t.data, temp) // add to the bottom of the list

	// remove all items on the top until 10 items are achieved
	if len(t.data) > TempProveReadings {
		t.data = t.data[len(t.data)-TempProveReadings:]
	}

	t.StatusChanged = false // reset the status changed on every read
	t.Status = t.calcStatus()
	if t.Status != t.prevStat {
		t.StatusChanged = true
		t.prevStat = t.Status
	}
}

// calcStatus calculates temperature status based on temperature readings
func (t *tempManager) calcStatus() tempStat {
	warns := 0
	crits := 0
	emergencies := 0
	// check the newest 4 readings // 4 * TempProveInterval = 12s
	for i := 0; i < len(t.data); i++ {
		if t.data[i] > warnTemp {
			crits++
			warns++
			emergencies++
		} else if t.data[i] > criticalTemp {
			crits++
			emergencies++
		} else if t.data[i] > emergencyTemp {
			emergencies++
		}
	}

	curStatus := tempOk
	if emergencies >= 2 {
		curStatus = tempEmergency
	} else if crits >= 2 {
		curStatus = tempCritical
	} else if warns >= 2 {
		curStatus = tempWarn
	}
	return curStatus
}

// print all the temperature values
func (t tempManager) print() {
	fmt.Printf("%v", t.data)
}
