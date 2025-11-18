# Common Windev Patterns and Go Equivalents

Este documento descreve padrões comuns encontrados em código Windev e como convertê-los para Go idiomático.

## Table of Contents
- [Procedure Declarations](#procedure-declarations)
- [Web Services](#web-services)
- [Database Access Patterns](#database-access-patterns)
- [File Processing](#file-processing)
- [Initialization and Cleanup](#initialization-and-cleanup)
- [Error Handling Patterns](#error-handling-patterns)
- [Configuration and Constants](#configuration-and-constants)

## Procedure Declarations

### Windev Pattern
```wlanguage
PROCEDURE MyProcedure(param1 is string, param2 is int)
nResult is int
sMessage is string

// Procedure logic here
sMessage = "Hello " + param1
nResult = param2 * 2

RESULT nResult
```

### Go Equivalent
```go
// MyProcedure performs specific business logic
// param1: description of param1
// param2: description of param2
// Returns: description of return value
func MyProcedure(param1 string, param2 int) int {
    var result int
    var message string
    
    // Procedure logic here
    message = "Hello " + param1
    result = param2 * 2
    
    return result
}
```

### Procedure with Multiple Returns (Windev using global or by reference)

**Windev:**
```wlanguage
PROCEDURE Calculate(nValue is int, sResult is string by reference)
sResult = "Calculated"
RESULT nValue * 10
```

**Go:**
```go
func Calculate(value int) (int, string) {
    result := "Calculated"
    return value * 10, result
}
```

## Web Services

### REST API Endpoint (WEBDEV Service)

**Windev:**
```wlanguage
PROCEDURE ws_GetUser(nUserID is int)
stUser is ST_User
sJSON is string

// Read from database
HReadSeek(User, UserID, nUserID)
IF HFound() THEN
    stUser.Name = User.Name
    stUser.Email = User.Email
    sJSON = VariantToJSON(stUser)
    WebserviceWriteHTTPCode(200)
    RESULT sJSON
ELSE
    WebserviceWriteHTTPCode(404)
    RESULT "User not found"
END
```

**Go with Gin:**
```go
// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        logger.Error("invalid user ID", zap.Error(err))
        c.JSON(400, gin.H{"error": "Invalid user ID"})
        return
    }
    
    var user User
    
    // Read from database
    query := "SELECT id, name, email FROM users WHERE id = $1"
    err = db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
    
    if err == sql.ErrNoRows {
        logger.Warn("user not found", zap.Int("userID", userID))
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    
    if err != nil {
        logger.Error("database error", zap.Error(err))
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    logger.Info("user retrieved successfully", zap.Int("userID", userID))
    c.JSON(200, user)
}
```

### POST Endpoint with JSON Body

**Windev:**
```wlanguage
PROCEDURE ws_CreateUser()
stUser is ST_User
sJSON is string
nUserID is int

// Get JSON from request
sJSON = WebserviceParameter("json")
JSONToVariant(stUser, sJSON)

// Validate
IF stUser.Name = "" THEN
    WebserviceWriteHTTPCode(400)
    RESULT "Name is required"
END

// Insert into database
User.Name = stUser.Name
User.Email = stUser.Email
HAdd(User)
nUserID = User.UserID

WebserviceWriteHTTPCode(201)
RESULT VariantToJSON(nUserID)
```

**Go with Gin:**
```go
// CreateUser creates a new user
// @Summary Create a new user
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // Parse JSON body
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Error("invalid request body", zap.Error(err))
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Validate
    if req.Name == "" {
        logger.Warn("validation failed: name is required")
        c.JSON(400, gin.H{"error": "Name is required"})
        return
    }
    
    // Insert into database
    var userID int
    query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
    err := db.QueryRow(query, req.Name, req.Email).Scan(&userID)
    
    if err != nil {
        logger.Error("failed to create user", zap.Error(err))
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }
    
    logger.Info("user created successfully", zap.Int("userID", userID))
    c.JSON(201, CreateUserResponse{UserID: userID})
}
```

## Database Access Patterns

### Simple Read

**Windev:**
```wlanguage
PROCEDURE GetProduct(nProductID is int)
stProduct is ST_Product

HReadSeek(Product, ProductID, nProductID)
IF HFound() THEN
    stProduct.Name = Product.Name
    stProduct.Price = Product.Price
    RESULT stProduct
ELSE
    RESULT Null
END
```

**Go:**
```go
func GetProduct(productID int) (*Product, error) {
    var product Product
    
    query := "SELECT id, name, price FROM products WHERE id = $1"
    err := db.QueryRow(query, productID).Scan(
        &product.ID,
        &product.Name,
        &product.Price,
    )
    
    if err == sql.ErrNoRows {
        logger.Debug("product not found", zap.Int("productID", productID))
        return nil, nil
    }
    
    if err != nil {
        logger.Error("database error", zap.Error(err))
        return nil, fmt.Errorf("failed to get product: %w", err)
    }
    
    return &product, nil
}
```

### List with Filter

**Windev:**
```wlanguage
PROCEDURE ListProducts(sCategory is string)
arrProducts is array of ST_Product
stProduct is ST_Product

FOR EACH Product WHERE Category = sCategory
    stProduct.ProductID = Product.ProductID
    stProduct.Name = Product.Name
    stProduct.Price = Product.Price
    ArrayAdd(arrProducts, stProduct)
END

RESULT arrProducts
```

**Go:**
```go
func ListProducts(category string) ([]Product, error) {
    var products []Product
    
    query := "SELECT id, name, price FROM products WHERE category = $1"
    rows, err := db.Query(query, category)
    if err != nil {
        logger.Error("query failed", zap.Error(err))
        return nil, fmt.Errorf("failed to list products: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var product Product
        err := rows.Scan(&product.ID, &product.Name, &product.Price)
        if err != nil {
            logger.Error("scan error", zap.Error(err))
            return nil, fmt.Errorf("failed to scan product: %w", err)
        }
        products = append(products, product)
    }
    
    if err = rows.Err(); err != nil {
        logger.Error("rows iteration error", zap.Error(err))
        return nil, fmt.Errorf("error iterating products: %w", err)
    }
    
    logger.Info("products listed", zap.Int("count", len(products)))
    return products, nil
}
```

### Insert/Update

**Windev:**
```wlanguage
PROCEDURE SaveProduct(stProduct is ST_Product)
bSuccess is boolean

IF stProduct.ProductID = 0 THEN
    // Insert
    Product.Name = stProduct.Name
    Product.Price = stProduct.Price
    Product.Category = stProduct.Category
    bSuccess = HAdd(Product)
ELSE
    // Update
    HReadSeek(Product, ProductID, stProduct.ProductID)
    IF HFound() THEN
        Product.Name = stProduct.Name
        Product.Price = stProduct.Price
        Product.Category = stProduct.Category
        bSuccess = HModify(Product)
    END
END

RESULT bSuccess
```

**Go:**
```go
func SaveProduct(product *Product) error {
    if product.ID == 0 {
        // Insert
        query := `INSERT INTO products (name, price, category) 
                  VALUES ($1, $2, $3) RETURNING id`
        err := db.QueryRow(query, product.Name, product.Price, product.Category).Scan(&product.ID)
        if err != nil {
            logger.Error("insert failed", zap.Error(err))
            return fmt.Errorf("failed to insert product: %w", err)
        }
        logger.Info("product inserted", zap.Int("productID", product.ID))
    } else {
        // Update
        query := `UPDATE products 
                  SET name = $1, price = $2, category = $3 
                  WHERE id = $4`
        result, err := db.Exec(query, product.Name, product.Price, product.Category, product.ID)
        if err != nil {
            logger.Error("update failed", zap.Error(err))
            return fmt.Errorf("failed to update product: %w", err)
        }
        
        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            logger.Warn("no rows updated", zap.Int("productID", product.ID))
            return fmt.Errorf("product not found")
        }
        logger.Info("product updated", zap.Int("productID", product.ID))
    }
    
    return nil
}
```

## File Processing

### Read File Line by Line

**Windev:**
```wlanguage
PROCEDURE ProcessFile(sFilePath is string)
nFileID is int
sLine is string
nCount is int = 0

nFileID = fOpen(sFilePath, foRead)
IF nFileID = -1 THEN
    Error("Cannot open file: " + sFilePath)
    RESULT False
END

WHILE NOT fEndOfFile(nFileID)
    sLine = fReadLine(nFileID)
    // Process line
    nCount++
END

fClose(nFileID)
RESULT nCount
```

**Go:**
```go
func ProcessFile(filePath string) (int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        logger.Error("cannot open file", zap.String("path", filePath), zap.Error(err))
        return 0, fmt.Errorf("cannot open file: %w", err)
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    count := 0
    
    for scanner.Scan() {
        line := scanner.Text()
        // Process line
        count++
    }
    
    if err := scanner.Err(); err != nil {
        logger.Error("error reading file", zap.Error(err))
        return count, fmt.Errorf("error reading file: %w", err)
    }
    
    logger.Info("file processed", zap.String("path", filePath), zap.Int("lines", count))
    return count, nil
}
```

### Write to File

**Windev:**
```wlanguage
PROCEDURE WriteLog(sMessage is string)
nFileID is int
sFilePath is string = fDataDir() + ["\"] + "app.log"

nFileID = fOpen(sFilePath, foCreateIfNotExist + foWrite + foAppend)
IF nFileID <> -1 THEN
    fWriteLine(nFileID, DateSys() + " " + TimeSys() + " - " + sMessage)
    fClose(nFileID)
END
```

**Go:**
```go
func WriteLog(message string) error {
    filePath := filepath.Join(dataDir, "app.log")
    
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        logger.Error("cannot open log file", zap.Error(err))
        return fmt.Errorf("cannot open log file: %w", err)
    }
    defer file.Close()
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logLine := fmt.Sprintf("%s - %s\n", timestamp, message)
    
    if _, err := file.WriteString(logLine); err != nil {
        logger.Error("cannot write to log", zap.Error(err))
        return fmt.Errorf("cannot write to log: %w", err)
    }
    
    return nil
}
```

## Initialization and Cleanup

### Windev Global Initialization

**Windev:**
```wlanguage
// Project initialization code
PROCEDURE InitProject()
gsDataPath is string = fDataDir()
gbDebugMode is boolean = InTestMode()

// Initialize database connection
IF NOT HOpenConnection("MyConnection") THEN
    Error("Cannot connect to database")
    RESULT False
END

// Load configuration
LoadConfiguration()

RESULT True
```

**Go Main Function:**
```go
var (
    dataPath   string
    debugMode  bool
    db         *sql.DB
    logger     *zap.Logger
    config     Config
)

func main() {
    // Initialize logger
    var err error
    if debugMode {
        logger, err = zap.NewDevelopment()
    } else {
        logger, err = zap.NewProduction()
    }
    if err != nil {
        log.Fatal("failed to initialize logger:", err)
    }
    defer logger.Sync()
    
    // Load configuration
    config, err = LoadConfiguration()
    if err != nil {
        logger.Fatal("failed to load configuration", zap.Error(err))
    }
    
    // Initialize database
    db, err = sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        logger.Fatal("failed to connect to database", zap.Error(err))
    }
    defer db.Close()
    
    // Verify connection
    if err = db.Ping(); err != nil {
        logger.Fatal("failed to ping database", zap.Error(err))
    }
    
    logger.Info("application initialized successfully")
    
    // Start server
    router := setupRouter()
    if err := router.Run(":8080"); err != nil {
        logger.Fatal("failed to start server", zap.Error(err))
    }
}

func setupRouter() *gin.Engine {
    router := gin.Default()
    
    // Middleware
    router.Use(LoggingMiddleware())
    router.Use(gin.Recovery())
    
    // Routes
    api := router.Group("/api")
    {
        api.GET("/users/:id", GetUser)
        api.POST("/users", CreateUser)
        // ... more routes
    }
    
    return router
}
```

### Service Initialization (Windev Service)

**Windev:**
```wlanguage
// Service initialization
PROCEDURE Service_Init()
// Open connections, load resources
HOpenConnection("MainDB")
LoadCache()
```

**Go:**
```go
func init() {
    // This runs automatically before main()
    // Use for package-level initialization only
}

// Better approach: explicit initialization
func InitializeService() error {
    // Open database connection
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    // Load cache
    if err := LoadCache(); err != nil {
        return fmt.Errorf("failed to load cache: %w", err)
    }
    
    logger.Info("service initialized successfully")
    return nil
}
```

## Error Handling Patterns

### Windev Exception Handling

**Windev:**
```wlanguage
PROCEDURE ProcessData()
WHEN EXCEPTION
    Error("An error occurred: " + ExceptionInfo())
    RESULT False
END

// Code that might fail
nResult is int = DangerousOperation()
RESULT True
```

**Go Error Handling:**
```go
func ProcessData() error {
    // Code that might fail
    result, err := DangerousOperation()
    if err != nil {
        logger.Error("operation failed", zap.Error(err))
        return fmt.Errorf("failed to process data: %w", err)
    }
    
    // Continue processing
    logger.Info("operation successful", zap.Int("result", result))
    return nil
}
```

### Windev Error Propagation

**Windev:**
```wlanguage
PROCEDURE OuterFunction()
IF NOT InnerFunction() THEN
    ErrorPropagate()
    RESULT False
END
RESULT True
```

**Go Error Wrapping:**
```go
func OuterFunction() error {
    if err := InnerFunction(); err != nil {
        // Wrap error with context
        return fmt.Errorf("outer function failed: %w", err)
    }
    return nil
}

func InnerFunction() error {
    // Some operation
    if someCondition {
        return fmt.Errorf("inner function failed: %v", reason)
    }
    return nil
}
```

## Configuration and Constants

### Windev Constants

**Windev:**
```wlanguage
// Constants
CONSTANT
    DEFAULT_TIMEOUT = 30
    MAX_RETRIES = 3
    API_BASE_URL = "https://api.example.com"
END
```

**Go Constants:**
```go
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
    APIBaseURL     = "https://api.example.com"
)
```

### Windev Configuration File

**Windev:**
```wlanguage
PROCEDURE LoadConfiguration()
sIniFile is string = fExeDir() + ["\"] + "config.ini"
gsAPIKey = INIRead("API", "Key", "", sIniFile)
gnTimeout = Val(INIRead("API", "Timeout", "30", sIniFile))
```

**Go Configuration (using environment variables):**
```go
type Config struct {
    APIKey     string
    Timeout    time.Duration
    DatabaseURL string
}

func LoadConfiguration() (Config, error) {
    config := Config{
        APIKey:     os.Getenv("API_KEY"),
        Timeout:    30 * time.Second,
        DatabaseURL: os.Getenv("DATABASE_URL"),
    }
    
    if config.APIKey == "" {
        return config, fmt.Errorf("API_KEY is required")
    }
    
    if timeoutStr := os.Getenv("API_TIMEOUT"); timeoutStr != "" {
        timeout, err := time.ParseDuration(timeoutStr)
        if err != nil {
            return config, fmt.Errorf("invalid timeout: %w", err)
        }
        config.Timeout = timeout
    }
    
    logger.Info("configuration loaded successfully")
    return config, nil
}
```

## Additional Patterns

### Retry Logic

**Go:**
```go
func RetryOperation(operation func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = operation()
        if err == nil {
            return nil
        }
        
        logger.Warn("operation failed, retrying",
            zap.Int("attempt", i+1),
            zap.Int("maxRetries", maxRetries),
            zap.Error(err))
        
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}
```

### Middleware Pattern (Logging)

**Go with Gin:**
```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // Process request
        c.Next()
        
        // Log after request
        duration := time.Since(start)
        logger.Info("request processed",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
        )
    }
}
```

### Context Usage for Timeouts

**Go:**
```go
func ProcessWithTimeout(ctx context.Context, data string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // Create a channel for result
    done := make(chan error, 1)
    
    go func() {
        done <- performLongOperation(data)
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        logger.Error("operation timed out")
        return fmt.Errorf("operation timed out: %w", ctx.Err())
    }
}
```
