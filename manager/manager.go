package manager

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	done    chan bool
	temp    *tempManager
	profile *acpiManager
}

const TempProveInterval = 3 // seconds

func (m *Manager) Start() {
	fmt.Println("Starting Cpu temp handling")

	ticker := time.NewTicker(TempProveInterval * time.Second)
	m.done = make(chan bool)

	m.temp = &tempManager{}

	acpi := acpiRead()
	m.profile = &acpiManager{
		Status:        acpi,
		prevStatus:    acpi,
		StatusChanged: false,
	}

	// print current status on startup
	printStatus()

	for {
		select {
		case <-m.done:
			ticker.Stop()
			return
		case <-ticker.C:

			m.temp.read()
			m.profile.read()

			if m.temp.StatusChanged || m.profile.StatusChanged {
				fmt.Printf("applying new rules... \n")
				switch acpiRead() {
				case acpiPerformance:
					switch m.temp.Status {
					case tempOk:
						ApplyPerformanceOk()
					case tempWarn:
						ApplyPerformanceWarn()
					case tempCritical:
						ApplyPerformanceCritical()
					case tempEmergency:
						ApplyPerformanceEmergency()
					}
				case acpiBalanced:
					switch m.temp.Status {
					case tempOk:
						ApplyBalanceOk()
					case tempWarn:
						ApplyBalanceWarn()
					case tempCritical:
						ApplyBalanceCritical()
					case tempEmergency:
						ApplyBalanceEmergency()
					}

				case acpiLowPower:
					// I have not experienced any thermal shutdown in low power mode
					// hence applying all default
					ApplyLowPower()
				}
			}
		}
	}

}

func (m *Manager) Stop() {
	m.done <- true
	ResetDefaults()
}

// write is a simple helper function to write a value to a path
func write(path string, v string) error {
	err := os.WriteFile(path, []byte(v), 0644)
	if err != nil {
		return err
	}
	return nil
}

func read(path string) (string, error) {
	v, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(v), "\n"), nil
}

func readInt(path string) (int, error) {
	v, err := read(path)
	if err != nil {
		return 0, err
	}
	intVar, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("unable to convert string %s, to integer", v)
	}
	return intVar, nil
}

// handleErr will handle read write errors to the /sys paths
// initially we only log the error
func handleErr(err error) {
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
}

func printStatus() {

	fmt.Printf("cpu temperature: %dº\n", tempRead())

	acpi := acpiRead()
	switch acpi {
	case acpiLowPower:
		fmt.Printf("acpi mode: low power\n")
	case acpiBalanced:
		fmt.Printf("acpi mode: balanced\n")
	case acpiPerformance:
		fmt.Printf("acpi mode: performance\n")
	}

	if cpuBoostEnabled() {
		fmt.Printf("cpu turbo boost is: on\n")
	} else {
		fmt.Printf("cpu turbo boost is: off\n")
	}

	fmt.Printf("gpu boost Mhz: %s\n", gpuBoostMhz())
	psu := psuRead()
	fmt.Printf("psu limit set to: %dW\n", psu/1000000)
}

// ===============================================================
// ACPI state
// ===============================================================
type acpiState int

const (
	acpiPerformance acpiState = iota
	acpiBalanced
	acpiLowPower
	acpiProfilePath = "/sys/firmware/acpi/platform_profile"
)

type acpiManager struct {
	Status        acpiState
	prevStatus    acpiState
	StatusChanged bool
}

func (m *acpiManager) read() {
	m.StatusChanged = false // reset the status changed on every read
	m.Status = acpiRead()
	if m.Status != m.prevStatus {
		m.StatusChanged = true
		m.prevStatus = m.Status
	}
}

func acpiRead() acpiState {

	v, err := read(acpiProfilePath)
	if err != nil {
		handleErr(err)
	}
	switch v {
	case "performance":
		return acpiPerformance
	case "balanced":
		return acpiBalanced
	case "low-power":
		return acpiLowPower
	default:
		handleErr(fmt.Errorf("unknown acpi profile: %s", v))
		return acpiBalanced
	}
}

// ===============================================================
// Intel turbo boost
// ===============================================================
const (
	cpuTurboPath = "/sys/devices/system/cpu/intel_pstate/no_turbo"
)

// set on/off state of cpu boost
func cpuBoostOn() {
	fmt.Println("enabling CPU turbo boost")
	err := write(cpuTurboPath, "0")
	handleErr(err)
}
func cpuBoostOff() {
	fmt.Println("disabling CPU turbo boost")
	err := write(cpuTurboPath, "1")
	handleErr(err)
}

func cpuBoostEnabled() bool {
	v, err := read(acpiProfilePath)
	if err != nil {
		handleErr(err)
	}
	if v == "1" {
		return false
	}
	return true
}

// ===============================================================
// GPU
// ===============================================================
const (
	gpuBoostFreq        = "/sys/class/drm/card0/gt_boost_freq_mhz"
	gpuBoostFreqDefault = 1350
	gpuBoostFreqLow     = 1000
	gpuBoostFreqLower   = 800
)

func gpuBoostMhzSet(v int) {
	fmt.Printf("setting GPU turbo freq to %d \n", v)
	err := write(cpuTurboPath, "1")
	handleErr(err)
}

func gpuBoostMhz() string {
	v, err := read(gpuBoostFreq)
	if err != nil {
		handleErr(err)
	}
	return v
}

// ===============================================================
// Power
// ===============================================================

const (
	psuWatt               = "/sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw"
	psuPerformanceDefault = 28000000
	psuPerformanceLow     = 20000000 // 20w
	psuBalancedDefault    = 15000000
	psuBalancedLow        = 12000000 // 12w
	psuLowPowerDefault    = 7500000
)

func psuRead() int {
	v, err := readInt(psuWatt)
	if err != nil {
		handleErr(err)
	}
	return v
}

func psuSet(v int) {
	fmt.Printf("setting power limit to: %dW \n", v/1000000)
	err := write(psuWatt, strconv.Itoa(v))
	handleErr(err)
}
