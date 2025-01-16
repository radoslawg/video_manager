# Video Manager

Idiosyncratic manager for online video consumption.

Mainly serves as training ground for learning [Go](https://go.dev/) language

# Running

```powershell
go run main.go help
```

## Starting web server
```powershell
go run main.go web
```

# Building

## Windows

```powershell
go build -o bin
```

## Cross-compiling for Linux

```powershell
$Env:GOOS = 'linux'
go build -o bin
```
