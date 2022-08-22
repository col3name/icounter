# Ip address counter

[File](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip]) size ~ 120gb

Windows 10, AMD Ryzen 5 1600, sata ssd, read 500 mbit/s

    result:
        count all ip: 212_031_415
        time taken = 39s
	    size ~ 3.18 gb
	result:
        count all ip = 7_956_984_159
        uniq = 1_120_557_388
        maxRam = 100mb
        average = 50mb
	    Time taken - 23m3.6871014s

## Build
### Build for macos for apple silicon
`make buildMacos`
### Build for linux
`make build`
### Build for windows
`make buildWin`

### Run 
`./bin/unique -f filepath`

### Lint 
`make lint`