# map
A simple and modern modern golang network mapper!

## Installing
`map` requires Go 1.7.1 or later.
```
$ go get -u github.com/ejcx/map
```

## Example
Scan for all open redis clusters on the internet
```
  $ map redis --nets 0.0.0.0/0
```

## Usage
`map` usage is based mostly on two flags.
 - `-n`, `--net`: CIDR of network to scan. Example: `127.0.0.1/32`
 - `-p`, `--port`: Port range to scan. Examples: `1-65000`, `22`, `22-25`
