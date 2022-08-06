#!/bin/bash -e
# partially based on: https://raw.githubusercontent.com/m0ppers/x1-carbon-gen9-hammer/main/hammer.sh

# check if executed as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi


# ===============================================================
# Intel turbo boost
# ===============================================================
disable_turbo_boost () {
  echo "disabling CPU turbo boost"
  echo 1 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo
}
enable_turbo_boost () {
  echo "enabling CPU turbo boost"
  echo 0 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo
}

# ===============================================================
# GPU
# ===============================================================
# limit gpu to 1000MHz (instead of 1350) all the time. power reduced from 13W to ~7.5W
disable_gpu_boost () {
  echo "setting GPU turbo freq to 1000"
  echo 1000 > /sys/class/drm/card0/gt_boost_freq_mhz
}
enable_gpu_boost () {
  echo "setting GPU turbo freq to 1350"
  echo 1350 > /sys/class/drm/card0/gt_boost_freq_mhz
}

# ===============================================================
# Power
# ===============================================================
power_state_perf_full () {
    echo "setting power limit to performance default: 28W"
    echo 28000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
}

power_state_perf_throttled () {
      echo "setting throttled power limit 20W"
      echo 20000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
}

power_state_bal_full () {
    echo "setting power limit to balanced default: 15W"
    echo 15000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
}

power_state_bal_throttled () {
    echo "setting throttled power limit 12W"
    echo 12000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
}

# ===============================================================
# Reset default on exit
# ===============================================================

function reset_defaults {
  echo ""
  echo "resetting default values..."
  enable_gpu_boost
  enable_turbo_boost

  local PROFILE
  PROFILE=$(cat /sys/firmware/acpi/platform_profile)
  if [ "$PROFILE" == "performance" ]; then
    power_state_full
  elif [ "$PROFILE" == "balanced" ]; then
    power_state_bal_full
  fi

}
trap reset_defaults EXIT


# ===============================================================
# Start management
# ===============================================================

function handle_performance () {
    echo "performance"

    echo "bla"
#        if [ "$POWER_LIMIT" -lt "20000000" ]; then
#          if [ "$TEMP" -lt "75000" ]; then
#            echo "Throttling to $POWER_LIMIT detected in performance profile. Resetting to full wattage (28W) because temp is okish at $TEMP"
#            echo 28000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
#          else
#            echo "Throttling to $POWER_LIMIT detected. Temperature is $TEMP so reducing to an hopefully sane level of 20W"
#            # play safe
#            echo 20000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
#          fi
#        elif [ "$POWER_LIMIT" != "28000000" ] && [ "$TEMP" -lt "80000" ]; then
#          echo "Resetting performance power limit to 28W"
#          echo 28000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
#        fi
}

function handle_balanced () {
    echo "Handling Balanced mode"
    echo "${POWER_LIMIT}"


#        if [ "$POWER_LIMIT" -lt "15000000" ]; then
#          if [ "$TEMP" -lt "75000" ]; then
#            echo "Throttling to $POWER_LIMIT detected in balanced profile. Resetting to full wattage (15W) because temp is okish at $TEMP"
#            # default
#            echo 15000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
#          elif [ "$POWER_LIMIT" -lt "12000000" ]; then
#            echo "Throttling to $POWER_LIMIT detected. Temperature is $TEMP so reducing to an hopefully sane level of 12W"
#            # lower by 3 watt and pray that is enough
#            echo 12000000 > /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw
#          fi
#        fi
}

disable_gpu_boost

while true; do
  PROFILE=$(cat /sys/firmware/acpi/platform_profile)
  POWER_LIMIT=$(cat /sys/devices/virtual/powercap/intel-rapl-mmio/intel-rapl-mmio:0/constraint_0_power_limit_uw)
  TEMP=$(cat /sys/class/thermal/thermal_zone0/temp)
  MODE=0


  if [ "$PROFILE" == "performance" ]; then
    handle_performance

  elif [ "$PROFILE" == "balanced" ]; then
    handle_balanced
  fi
  # it seems it take quite a while for throttling to happen
  sleep 1
done
