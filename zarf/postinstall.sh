#!/bin/bash
# postinstall  script


# Declare an array of string with type
declare -a exec=(
"/usr/sbin/acpiprf"
"/usr/sbin/x1-acpi-manager"
)

# Iterate the string array using for loop
for item in ${exec[@]}; do
   chown root:root "$item"
   chmod 755  "$item"
done

# start the service after install
systemctl enable x1-carbon-gen9-acpi-manager.service
systemctl daemon-reload
systemctl start x1-carbon-gen9-acpi-manager.service