# archpi-zerow-sensors
ArchLinuxARM Raspberry PI Zero send sensors (DHT22, DS18B20) to influxdb v1.8 using GO


2 programs: 1 for DHT22 and one for DS18B20

TODO:
* merge both programs into 1 and share single influxdb client
* use config file for all settings
* make non-root user read DHT22 sensor

Tested on:
* Linux archlinuxarm 5.4.83-1-ARCH #1 SMP PREEMPT Mon Dec 14 15:44:54 UTC 2020 armv6l GNU/Linux
* go version go1.15.6 linux/arm
* /boot/config.txt:
	gpu_mem=16
	initramfs initramfs-linux.img followkernel

	device_tree=bcm2835-rpi-zero-w.dtb

	# does not work if below is enabled..
	#dtoverlay=dht11,gpiopin=17

	dtoverlay=w1-gpio,gpiopin=24
* Connections:
	- DS18B20 data line connected to GPIO24
 	- DHT22 data line connected to GPIO17

* Compile:
    - go build dht22.go
    - go buold ds18b20.go

* go dependencies to be picked up in common manner

* InfluxDB 1.8.3

Thanks to all examples provided online
