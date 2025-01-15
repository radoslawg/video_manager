# Video Manager

Idiosyncratic manager for online video consumption.

Mainly serves as training ground for learning [Go](https://go.dev/) language

# Running

```powershell
go run cmd\main.go
```

# Building

## Windows

```powershell
cd cmd
go build
```

## Cross-compiling for Linux

```powershell
cd cmd
$Env:GOOS = 'linux'
go build
```
