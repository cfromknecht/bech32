# bech32
Command Line Bech32 Encode/Decode

# Install
```
$ go get -u github.com/cfromknecht/bech32
```

# Encode
```
$ bech32 encode <human-readable-prefix> <payload>
```
OR
```
$ bech32 encode --hrp=<human-readable-prefix> --payload=<payload>
```

# Decode
```
$ bech32 decode <encoding>
```
OR
```
$ bech32 decode --encoding=<encoding>
```

# Help
```
$ bech32 -h
$ bech32 encode -h
$ bech32 decode -h
```
