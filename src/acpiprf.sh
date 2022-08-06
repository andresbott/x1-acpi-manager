#!/bin/bash

# make sure whiptail is installed
command -v whiptail >/dev/null 2>&1 || {
  echo >&2 "Please make sure that package: whiptail is installed"; exit 1;
}

# check if executed as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi


## Read available acpi profiles
CURRENT=$(cat /sys/firmware/acpi/platform_profile)
for OPT in $(cat /sys/firmware/acpi/platform_profile_choices)
do
    PRINT="$PRINT $OPT"

    # note whiptail uses two strings per entry
    if [ "$CURRENT" = "$OPT" ];
    then
       PRINT="$PRINT - ON "
    else
       PRINT="$PRINT - OFF "
    fi

done

RESULT=$(whiptail --title "Select platform profile" --radiolist "Select one:" 20  70 10 $PRINT 3>&1 1>&2 2>&3)

# exit on error
exitstatus=$?
if [ $exitstatus != 0 ];
then
  exit 1
fi

if [ $exitstatus = 0 ];
then

  echo "Aplying profile: \"$RESULT\" to /sys/firmware/acpi/platform_profile"
  echo "$RESULT" > /sys/firmware/acpi/platform_profile

fi
