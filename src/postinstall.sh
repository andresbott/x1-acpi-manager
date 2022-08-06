#!/bin/bash
# postinstall  script


# Declare an array of string with type
declare -a exec=(
"/usr/local/sbin/acpiprf"
"/usr/local/sbin/x1-cpu-manager.sh"
)

# Iterate the string array using for loop
for item in ${exec[@]}; do
   chown root:root "$item"
   chmod 755  "$item"
done

