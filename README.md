# Ip address counter

[File](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip]) size ~ 120gb

## Windows 10, amd ryzen 5 1600 (6 cores, 12 Threads, 3.4 GHz), sata ssd (read 500 mbit/s), go 1.19

| Count         |  Average time  | File Size      | Ram       | Unique            |
|:------------- |----------------|:---------------|:----------|:------------------|
| 80_000        |   **62    ms** |  1.28 mb       |  30 mb    |  -                |
| 800_000       |   **180   ms** |  12.8 mb       |  30 mb    |  -                |
| 8_000_000     |   **1.3   s**  |  128  mb       |  30 mb    |  -                |
| 80_000_000    |   **13    s**  |  1.28 gb       |  30 mb    |  -                |
| 800_000_000   |   **130   s**  |  12.8 gb       |  40 mb    |  -                |
| 7_956_984_159 |   **20m20 s**  |  120  gb       |  60 mb    |  1_120_557_388    |

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