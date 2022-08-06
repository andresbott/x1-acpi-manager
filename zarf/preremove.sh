#!/bin/bash
# preremove  script

systemctl stop x1-carbon-gen9-acpi-manager.service
systemctl disable x1-carbon-gen9-acpi-manager.service
systemctl daemon-reload