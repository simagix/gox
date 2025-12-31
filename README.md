# gox

[![Go Report Card](https://goreportcard.com/badge/github.com/simagix/gox)](https://goreportcard.com/report/github.com/simagix/gox)
[![GoDoc](https://godoc.org/github.com/simagix/gox?status.svg)](https://godoc.org/github.com/simagix/gox)

Go Extra Functions and Utilities - a collection of reusable Go utilities for MongoDB tools and general-purpose applications.

## Installation

```bash
go get github.com/simagix/gox
```

## Features

### Obfuscation (`obfuscate.go`)

Deterministic PII obfuscation with consistent mappings - the same input always produces the same output.

```go
import "github.com/simagix/gox"

o := gox.NewObfuscator()

// Obfuscate various PII types
o.ObfuscateIP("192.168.1.100")       // → "192.X.X.100" or "10.X.X.X"
o.ObfuscateHostname("server1.com")   // → "tulip.atlanta.local"
o.ObfuscateEmail("user@example.com") // → "begonia@chicago.com"
o.ObfuscateSSN("123-45-6789")        // → "XXX-XX-XXXX" (shuffled)
o.ObfuscatePhoneNo("555-123-4567")   // → "555-12X-XXXX"
o.ObfuscateMAC("AA:BB:CC:11:22:33")  // → "AA:BB:CC:XX:XX:XX"
o.ObfuscateCreditCardNo("4532...")   // → "************1234"
o.ObfuscateDate("2024-06-15")        // → shifted by DateOffset days
```

**Configuration Options:**

```go
o := gox.NewObfuscator()

// IP obfuscation style
o.IPStyle = gox.IPStyleKeepEnds  // 192.168.1.100 → 192.X.X.100 (default)
o.IPStyle = gox.IPStylePrivate   // 192.168.1.100 → 10.X.X.X

// Name obfuscation style  
o.NameStyle = gox.NameStyleReadable  // city/flower names (default)
o.NameStyle = gox.NameStyleHash      // host-abc123.local

// Numeric obfuscation
o.Coefficient = 0.917  // multiplier for numbers (default)
o.DateOffset = -42     // days to shift dates (default)
```

**PII Detection:**

```go
gox.ContainsIP("192.168.1.1")           // true
gox.ContainsEmail("user@example.com")   // true
gox.ContainsSSN("123-45-6789")          // true
gox.ContainsPhoneNo("555-123-4567")     // true
gox.ContainsCreditCardNo("4532...")     // true (with Luhn check)
gox.ContainsFQDN("server.example.com")  // true
gox.IsNamespace("mydb.mycollection")    // true
```

**Traversal:**

```go
// Obfuscate nested maps
doc := map[string]interface{}{
    "ip": "192.168.1.1",
    "user": map[string]interface{}{
        "email": "user@example.com",
    },
}
obfuscated := o.ObfuscateMap(doc)
```

### I/O Utilities (`ioutil.go`)

Read files with automatic decompression (gzip, zstd, snappy).

```go
// Auto-detect and decompress
reader, _ := gox.NewReader(file)        // from io.Reader
reader, _ := gox.NewFileReader(path)    // from file path

// Compression detection
gox.IsGzip(data)   // check gzip magic bytes
gox.IsZstd(data)   // check zstd magic bytes
gox.IsSnappy(data) // check snappy magic bytes
```

### Map Walker (`map_walker.go`)

Traverse nested maps with callbacks.

```go
walker := gox.NewMapWalker(func(v interface{}) interface{} {
    if s, ok := v.(string); ok {
        return strings.ToUpper(s)
    }
    return v
})

result := walker.Walk(myMap)
level := walker.GetNestedLevel()
maxArrayLen := walker.GetMaxArrayLength()
```

### Logger (`logger.go`)

Simple logging utilities.

```go
gox.GetLogger().Info("message")
gox.GetLogger().Error("error occurred")
```

### Numbers (`numbers.go`)

Type conversion utilities.

```go
gox.ToInt(value)            // convert any type to int
gox.ToInt64(value)          // convert to int64 with error
gox.ToFloat64(value)        // convert to float64 with error
```

### Random Strings (`random_string.go`)

Generate random strings.

```go
gox.GetRandomDigitString(10)  // random digits
gox.GetRandomHexString(16)    // random hex string
gox.GetRandomUUIDString()     // UUID format string
```

### Measurements (`measurements.go`)

Human-readable size and duration formatting.

```go
gox.FormatBytes(1024)        // "1 KB"
gox.FormatDuration(time.Hour) // "1h"
```

### Ordered Map (`ordered_map.go`)

Map that preserves insertion order.

```go
om := gox.NewOrderedMap()
om.Set("key1", "value1")
om.Set("key2", "value2")
for _, key := range om.Keys() {
    value := om.Get(key)
}
```

### Stringify (`stringify.go`)

Convert values to strings with formatting.

```go
gox.Stringify(value)
```

### HTTP Utilities (`http_util.go`, `http_digest.go`)

HTTP client utilities including digest authentication.

```go
gox.HTTPGet(url, headers)
gox.HTTPDigestAuth(...)
```

### Web Server (`web_server.go`)

Simple web server utilities.

### Wait Group (`wait_group.go`)

Enhanced wait group utilities.

## Used By

- [hatchet](https://github.com/simagix/hatchet) - MongoDB JSON Log Analyzer
- [mongo-ftdc](https://github.com/simagix/mongo-ftdc) - MongoDB FTDC Metrics Viewer

## License

[Apache 2.0](LICENSE)
