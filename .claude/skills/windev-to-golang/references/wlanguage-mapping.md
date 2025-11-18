# WLanguage to Go Function Mapping

Este documento mapeia funções comuns do WLanguage (Windev) para seus equivalentes em Go.

## Table of Contents
- [String Operations](#string-operations)
- [File Operations](#file-operations)
- [Database Operations (HyperFileSQL)](#database-operations-hyperfilesql)
- [Date and Time](#date-and-time)
- [Numeric Operations](#numeric-operations)
- [Array Operations](#array-operations)
- [HTTP/Web Services](#httpweb-services)
- [JSON/XML](#jsonxml)
- [Logging](#logging)
- [Error Handling](#error-handling)

## String Operations

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `Length(str)` | `len(str)` | builtin | |
| `Upper(str)` | `strings.ToUpper(str)` | strings | |
| `Lower(str)` | `strings.ToLower(str)` | strings | |
| `Left(str, n)` | `str[:n]` | builtin | Check bounds |
| `Right(str, n)` | `str[len(str)-n:]` | builtin | Check bounds |
| `Middle(str, pos, n)` | `str[pos:pos+n]` | builtin | Check bounds |
| `Replace(str, old, new)` | `strings.ReplaceAll(str, old, new)` | strings | |
| `StringCount(str, search)` | `strings.Count(str, search)` | strings | |
| `Position(str, search)` | `strings.Index(str, search)` | strings | Returns -1 if not found |
| `StartsWith(str, prefix)` | `strings.HasPrefix(str, prefix)` | strings | |
| `EndsWith(str, suffix)` | `strings.HasSuffix(str, suffix)` | strings | |
| `Trim(str)` | `strings.TrimSpace(str)` | strings | |
| `StringBuild(format, ...)` | `fmt.Sprintf(format, ...)` | fmt | |
| `Val(str)` | `strconv.Atoi(str)` or `strconv.ParseFloat(str, 64)` | strconv | |
| `NumToString(num)` | `strconv.Itoa(num)` or `fmt.Sprintf("%d", num)` | strconv/fmt | |

## File Operations

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `fOpen(filename, mode)` | `os.OpenFile(filename, flags, perm)` | os | Use defer for Close |
| `fClose(fileID)` | `file.Close()` | os | |
| `fReadLine(fileID)` | `bufio.Scanner` or `bufio.Reader.ReadString('\n')` | bufio | |
| `fWriteLine(fileID, str)` | `file.WriteString(str + "\n")` | os | |
| `fRead(fileID, size)` | `file.Read(buffer)` | os | |
| `fWrite(fileID, data)` | `file.Write([]byte(data))` | os | |
| `fSeek(fileID, pos)` | `file.Seek(pos, 0)` | os | |
| `fSize(filename)` | `os.Stat(filename)` then `FileInfo.Size()` | os | |
| `fFileExist(filename)` | `_, err := os.Stat(filename); !os.IsNotExist(err)` | os | |
| `fDelete(filename)` | `os.Remove(filename)` | os | |
| `fCopyFile(src, dst)` | `io.Copy` or `os.ReadFile` + `os.WriteFile` | io/os | |
| `fDir(pattern)` | `filepath.Glob(pattern)` or `os.ReadDir` | filepath/os | |
| `fExtractPath(path)` | `filepath.Dir(path)` | filepath | |
| `CompleteDir(path)` | `filepath.Join(path, "")` | filepath | |
| `fDataDir()` | Define constant or env var | - | App-specific |
| `fExeDir()` | `os.Executable()` + `filepath.Dir()` | os/filepath | |

## Database Operations (HyperFileSQL)

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `HReadFirst(file, key)` | SQL: `SELECT * FROM table ORDER BY key LIMIT 1` | database/sql | Use appropriate driver |
| `HReadLast(file, key)` | SQL: `SELECT * FROM table ORDER BY key DESC LIMIT 1` | database/sql | |
| `HReadSeek(file, key, value)` | SQL: `SELECT * FROM table WHERE key = ?` | database/sql | |
| `HReadNext(file)` | `rows.Next()` in loop | database/sql | |
| `HReadPrevious(file)` | Requires ORDER BY DESC cursor | database/sql | |
| `HAdd(file)` | SQL: `INSERT INTO table VALUES (...)` | database/sql | |
| `HModify(file)` | SQL: `UPDATE table SET ... WHERE id = ?` | database/sql | |
| `HDelete(file)` | SQL: `DELETE FROM table WHERE id = ?` | database/sql | |
| `HFound()` | Check `err != sql.ErrNoRows` | database/sql | |
| `HOut()` | Check `!rows.Next()` | database/sql | |
| `HCancelSeek(file)` | N/A - close result set | database/sql | |
| `File.Item` | `row.Scan(&var)` | database/sql | |
| `File.Item = value` | Prepare statement parameters | database/sql | |
| `HExecuteQuery(query)` | `db.Query(query, args...)` | database/sql | Use parameterized queries |
| `SQLFetch(query)` | `rows.Next()` + `rows.Scan()` | database/sql | |
| `SQLGetCol(query, col)` | `row.Scan(&var)` | database/sql | |

**Important Notes:**
- Always use parameterized queries to prevent SQL injection
- Use `database/sql` with appropriate driver (pgx for PostgreSQL, mysql for MySQL, etc.)
- Handle connection pooling properly
- Use context for timeouts and cancellation

## Date and Time

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `Today()` | `time.Now()` | time | |
| `Now()` | `time.Now()` | time | |
| `DateSys()` | `time.Now().Format("20060102")` | time | |
| `TimeSys()` | `time.Now().Format("150405")` | time | |
| `DateToString(date, format)` | `date.Format(layout)` | time | Convert format |
| `StringToDate(str, format)` | `time.Parse(layout, str)` | time | |
| `DateDifference(date1, date2)` | `date1.Sub(date2)` | time | Returns Duration |
| `DateAdd(date, days)` | `date.AddDate(0, 0, days)` | time | |
| `Year(date)` | `date.Year()` | time | |
| `Month(date)` | `int(date.Month())` | time | |
| `Day(date)` | `date.Day()` | time | |
| `Hour(time)` | `time.Hour()` | time | |
| `Minute(time)` | `time.Minute()` | time | |
| `Second(time)` | `time.Second()` | time | |

**Format Conversion:**
- Windev `YYYYMMDD` → Go `"20060102"`
- Windev `DD/MM/YYYY` → Go `"02/01/2006"`
- Windev `HHMMSS` → Go `"150405"`
- Windev `HH:MM:SS` → Go `"15:04:05"`

## Numeric Operations

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `Abs(n)` | `math.Abs(float64(n))` | math | For float; use manual for int |
| `Round(n, decimals)` | `math.Round(n*100)/100` | math | Adjust multiplier |
| `Truncate(n, decimals)` | `math.Trunc(n*100)/100` | math | |
| `Random(min, max)` | `rand.Intn(max-min+1) + min` | math/rand | Use crypto/rand for security |
| `Power(base, exp)` | `math.Pow(base, exp)` | math | |
| `SquareRoot(n)` | `math.Sqrt(n)` | math | |
| `Max(a, b)` | `math.Max(a, b)` | math | For float64 |
| `Min(a, b)` | `math.Min(a, b)` | math | For float64 |

## Array Operations

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `ArrayCount(arr)` | `len(arr)` | builtin | |
| `ArrayAdd(arr, value)` | `arr = append(arr, value)` | builtin | |
| `ArrayInsert(arr, pos, value)` | Manual slice manipulation | builtin | |
| `ArrayDelete(arr, pos)` | Manual slice manipulation | builtin | |
| `ArraySeek(arr, value)` | Loop with comparison | builtin | |
| `ArraySort(arr)` | `sort.Slice(arr, func)` | sort | |
| `ArrayDeleteAll(arr)` | `arr = make([]Type, 0)` | builtin | |
| Dynamic Array | `[]Type` (slice) | builtin | |
| Fixed Array | `[n]Type` (array) | builtin | |

## HTTP/Web Services

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `HTTPRequest()` | `http.NewRequest()` + `client.Do()` | net/http | |
| `HTTPGetResult()` | `ioutil.ReadAll(resp.Body)` | io/ioutil | |
| `HTTPSend()` | `client.Do(req)` | net/http | |
| `HTTPDestination()` | Set request URL | net/http | |
| `HTTPParameter()` | URL query or body params | net/http | |
| `JSONExecute()` | `json.Unmarshal()` | encoding/json | |
| `JSONToVariant()` | `json.Unmarshal()` into struct | encoding/json | |
| `VariantToJSON()` | `json.Marshal()` | encoding/json | |
| WEBDEV Service | Gin handler functions | github.com/gin-gonic/gin | |
| `WebserviceWriteHTTPCode()` | `c.Status(code)` | gin | |
| `WebserviceWriteMIMEType()` | `c.Header("Content-Type", mime)` | gin | |
| `WebserviceParameter()` | `c.Query()` or `c.Param()` | gin | |

**Gin API Pattern:**
```go
router := gin.Default()
router.POST("/api/endpoint", handlerFunction)
router.GET("/api/endpoint/:id", handlerFunction)

func handlerFunction(c *gin.Context) {
    // Handle request
    c.JSON(200, response)
}
```

## JSON/XML

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `JSONExecute()` | `json.Unmarshal(data, &v)` | encoding/json | |
| `JSONToVariant()` | `json.Unmarshal(data, &v)` | encoding/json | |
| `VariantToJSON()` | `json.Marshal(v)` | encoding/json | |
| `XMLDocument.Element` | Struct tags for marshaling | encoding/json | |
| `Serialize/Deserialize` | `json.Marshal/Unmarshal` | encoding/json | |

**JSON Struct Tags:**
```go
type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age,omitempty"`
    Email string `json:"email"`
}
```

## Logging

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `Trace()` | `logger.Debug()` | go.uber.org/zap | Use zap logger |
| `dbgWriteWarningAudit()` | `logger.Warn()` | go.uber.org/zap | |
| `dbgWriteErrorAudit()` | `logger.Error()` | go.uber.org/zap | |
| `LogWrite()` | `logger.Info()` | go.uber.org/zap | |
| `Info()` | `logger.Info()` | go.uber.org/zap | |
| `Warning()` | `logger.Warn()` | go.uber.org/zap | |
| `Error()` | `logger.Error()` | go.uber.org/zap | |

**Zap Logger Setup:**
```go
logger, _ := zap.NewProduction()
defer logger.Sync()
sugar := logger.Sugar()

sugar.Infow("message", "key", "value")
sugar.Errorw("error occurred", "error", err)
```

## Error Handling

| WLanguage | Go Equivalent | Package | Notes |
|-----------|---------------|---------|-------|
| `ErrorOccurred()` | Check `err != nil` | builtin | |
| `ErrorInfo()` | `err.Error()` | builtin | |
| `ExceptionEnable()` | Go doesn't use exceptions | builtin | Use error returns |
| `ExceptionThrow()` | `return err` or `panic()` | builtin | Prefer return err |
| WHEN EXCEPTION block | `if err != nil` blocks | builtin | |
| `ErrorPropagate()` | `return fmt.Errorf("context: %w", err)` | fmt | Wrap errors |

**Error Handling Pattern:**
```go
if err != nil {
    logger.Error("operation failed", zap.Error(err))
    return fmt.Errorf("failed to do something: %w", err)
}
```

## Additional Conversions

### WLanguage Control Structures to Go

| WLanguage | Go |
|-----------|-----|
| `IF ... END` | `if { }` |
| `SWITCH ... END` | `switch { }` |
| `FOR ... END` | `for { }` |
| `WHILE ... END` | `for condition { }` |
| `LOOP ... END` | `for { }` |
| `FOR EACH ... END` | `for _, item := range collection { }` |
| `BREAK` | `break` |
| `CONTINUE` | `continue` |
| `RETURN` | `return` |

### WLanguage Variable Types to Go

| WLanguage | Go |
|-----------|-----|
| `int` | `int` or `int64` |
| `real` | `float64` |
| `boolean` | `bool` |
| `string` | `string` |
| `date` | `time.Time` |
| `datetime` | `time.Time` |
| `buffer` | `[]byte` |
| `variant` | `interface{}` or `any` |
| `array of X` | `[]X` |
| `associative array` | `map[KeyType]ValueType` |

### Constants

| WLanguage | Go |
|-----------|-----|
| `True` / `False` | `true` / `false` |
| `Null` | `nil` |
| `CRLF` | `"\r\n"` or `"\n"` |
| `TAB` | `"\t"` |
| `CR` | `"\r"` |

## Package Recommendations

When converting Windev code to Go, use these packages:

- **Web Framework**: `github.com/gin-gonic/gin`
- **Logging**: `go.uber.org/zap`
- **Database**: `database/sql` + driver (e.g., `github.com/lib/pq` for PostgreSQL)
- **HTTP Client**: `net/http` (standard library)
- **JSON**: `encoding/json` (standard library)
- **Configuration**: `github.com/spf13/viper`
- **Validation**: `github.com/go-playground/validator/v10`
- **Testing**: `testing` (standard library) + `github.com/stretchr/testify`

## Notes on Conversion

1. **Error Handling**: Go uses explicit error returns instead of exceptions. Always check errors.
2. **Null Safety**: Go uses zero values. Use pointers for optional fields.
3. **Memory Management**: Go has automatic garbage collection.
4. **Concurrency**: Use goroutines and channels instead of threads.
5. **Context**: Use `context.Context` for cancellation and timeouts in APIs and database operations.
