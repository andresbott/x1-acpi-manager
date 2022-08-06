# x1-acpi-manager

Simple service that makes some runtime modifications on a X1 carbon Gen 9 on top of the provided 
acpi profiles to try to prevent excessive cpu throttling as well as prevent thermal shutdowns

(MS teams i'm looking at you)

## Usage

run as root, if using systemd you can see the logs:

```
 systemctl status x1-carbon-gen9-acpi-manager.service 
```

## How it behaves:

the service will monitor the cpu temperature every 3 seconds, if the temperature is too high over certain 
amount of readings, it will switch between different criticality levels, and apply some changes.

* level OK 
  * all values are reset to default
* level warn
  * the boost frequency of the GPU is reduced
  * the power cap is reduced
* level critical
  * nothing special at the moment, to be seen how the system behaves
* level emergency
  * the cpu turbo boost is disabled