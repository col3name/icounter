# Ip address counter

Console application for count unique ip address on file. File handled parallel using hyperLogLog.
Optimized only for ssd, on hdd need use `-n 1` parameter.

Time complexity O(1)

[Download testing file](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip]). Packed Size ~ 12gb, unpacked ~ 120gb
## Variants
- HyperLogLog
- Count every time the bloom filter doesn't contain IP, but its difficult to parallelize

## Windows 10 amd ryzen 5 1600 (6 cores, 12 Threads, 3.4 GHz), sata ssd (read 220 mbit/s*), go 1.19 

*on nvme m2 ssd not enough space

| Count         |  Average time  | File Size      | Ram        | Unique            |
|:------------- |----------------|:---------------|:-----------|:------------------|
| 80_000        |   **10   ms**  |  1.28 mb       |  9.5 mb    |  -                |
| 800_000       |   **50   ms**  |  12.8 mb       |  9.5 mb    |  -                |
| 8_000_000     |   **340  ms**  |  128  mb       |  9.5 mb    |  -                |
| 80_000_000    |   **3.4  s**   |  1.28 gb       |  9.5 mb    |  -                |
| 800_000_000   |   **34   s**   |  12.8 gb       |  9.5 mb    |  -                |
| 7_956_984_159 |   **8m15 s**   |  120  gb       |  9.5 mb    |  1_122_438_962    |

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

### Usage
`unique -f path -n countParallelTask`

### How optimize
- faster ssd
- faster processor

### Format of the file with ip addresses

```
145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
```