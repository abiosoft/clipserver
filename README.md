# clipserver
Barebones clipboard server.

## Installing

Requires [Go](https://go.dev/dl).

```
go install github.com/abiosoft/clipserver@latest
```

## Usage

```
$ clipserver
starting server on 8080
```

Visit `localhost:8080` to use the clipboard.

On other devices on the local network, visit `<device-ip>:8080` to use the clipboard.

## ⚠️  Warning

This is written for convenience and not for security. Clipboard contents are being sent in plain text over the network.


