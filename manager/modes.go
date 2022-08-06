package manager

import "fmt"

// ===============================================================
// Performance
// ===============================================================

func ApplyPerformanceOk() {
	fmt.Println("temperature is: OK")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqDefault)
	psuSet(psuPerformanceDefault)
}

func ApplyPerformanceWarn() {
	fmt.Println("temperature is: Warn")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqLow)
	psuSet(psuPerformanceLow)

}
func ApplyPerformanceCritical() {
	fmt.Println("temperature is: Critical")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqLow)
	psuSet(psuPerformanceLow)

}
func ApplyPerformanceEmergency() {
	fmt.Println("temperature is: Emergency")
	cpuBoostOff()
	gpuBoostMhzSet(gpuBoostFreqLower)
	psuSet(psuPerformanceLow)

}

// ===============================================================
// Balanced
// ===============================================================

func ApplyBalanceOk() {
	fmt.Println("temperature is: OK")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqDefault)
	psuSet(psuBalancedDefault)
}

func ApplyBalanceWarn() {
	fmt.Println("temperature is: Warn")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqLow)
	psuSet(psuBalancedLow)

}
func ApplyBalanceCritical() {
	fmt.Println("temperature is: Critical")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqLow)
	psuSet(psuBalancedLow)

}
func ApplyBalanceEmergency() {
	fmt.Println("temperature is: Emergency")
	cpuBoostOff()
	gpuBoostMhzSet(gpuBoostFreqLower)
	psuSet(psuBalancedLow)
}

// ===============================================================
// Low Power
// ===============================================================

func ApplyLowPower() {
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqDefault)
	psuSet(psuLowPowerDefault)
}

// ===============================================================
// Default
// ===============================================================

func ResetDefaults() {
	fmt.Println("")
	fmt.Println("Resetting default values...")
	cpuBoostOn()
	gpuBoostMhzSet(gpuBoostFreqDefault)
	switch acpiRead() {
	case acpiPerformance:
		psuSet(psuPerformanceDefault)
	case acpiBalanced:
		psuSet(psuBalancedDefault)
	case acpiLowPower:
		psuSet(psuLowPowerDefault)
	}

}
