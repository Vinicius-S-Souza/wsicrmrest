package main

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Example handler structure - adapt based on your Windev procedures

// GetItem retrieves an item by ID
// Converted from Windev procedure: ws_GetItem
func GetItem(c *gin.Context) {
	// Parse ID from URL parameter
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("invalid item ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	// Query database
	var item Item
	query := "SELECT id, name, description, price FROM items WHERE id = $1"
	err = db.QueryRow(query, itemID).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
	)

	if err == sql.ErrNoRows {
		logger.Warn("item not found", zap.Int("itemID", itemID))
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	if err != nil {
		logger.Error("database error", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	logger.Info("item retrieved successfully", zap.Int("itemID", itemID))
	c.JSON(200, item)
}

// CreateItem creates a new item
// Converted from Windev procedure: ws_CreateItem
func CreateItem(c *gin.Context) {
	var req CreateItemRequest

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request body", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if req.Name == "" {
		logger.Warn("validation failed: name is required")
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	// Insert into database
	var itemID int
	query := `INSERT INTO items (name, description, price) 
	          VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, req.Name, req.Description, req.Price).Scan(&itemID)

	if err != nil {
		logger.Error("failed to create item", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to create item"})
		return
	}

	logger.Info("item created successfully", zap.Int("itemID", itemID))
	c.JSON(201, CreateItemResponse{ItemID: itemID})
}

// UpdateItem updates an existing item
// Converted from Windev procedure: ws_UpdateItem
func UpdateItem(c *gin.Context) {
	// Parse ID from URL parameter
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("invalid item ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	var req UpdateItemRequest

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request body", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update database
	query := `UPDATE items 
	          SET name = $1, description = $2, price = $3 
	          WHERE id = $4`
	result, err := db.Exec(query, req.Name, req.Description, req.Price, itemID)

	if err != nil {
		logger.Error("failed to update item", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to update item"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("item not found for update", zap.Int("itemID", itemID))
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	logger.Info("item updated successfully", zap.Int("itemID", itemID))
	c.JSON(200, gin.H{"message": "Item updated successfully"})
}

// DeleteItem deletes an item by ID
// Converted from Windev procedure: ws_DeleteItem
func DeleteItem(c *gin.Context) {
	// Parse ID from URL parameter
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("invalid item ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	// Delete from database
	query := "DELETE FROM items WHERE id = $1"
	result, err := db.Exec(query, itemID)

	if err != nil {
		logger.Error("failed to delete item", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to delete item"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("item not found for deletion", zap.Int("itemID", itemID))
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	logger.Info("item deleted successfully", zap.Int("itemID", itemID))
	c.JSON(200, gin.H{"message": "Item deleted successfully"})
}

// ListItems returns a list of items with optional filtering
// Converted from Windev procedure: ws_ListItems
func ListItems(c *gin.Context) {
	// Get query parameters
	category := c.Query("category")
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")

	// Convert limit and offset to integers
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	// Build query
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT id, name, description, price 
		         FROM items WHERE category = $1 
		         ORDER BY id LIMIT $2 OFFSET $3`
		args = []interface{}{category, limitInt, offsetInt}
	} else {
		query = `SELECT id, name, description, price 
		         FROM items 
		         ORDER BY id LIMIT $1 OFFSET $2`
		args = []interface{}{limitInt, offsetInt}
	}

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Error("query failed", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to retrieve items"})
		return
	}
	defer rows.Close()

	// Scan results
	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to scan items"})
			return
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		logger.Error("rows iteration error", zap.Error(err))
		c.JSON(500, gin.H{"error": "Error iterating items"})
		return
	}

	logger.Info("items listed", zap.Int("count", len(items)))
	c.JSON(200, ListItemsResponse{
		Items: items,
		Total: len(items),
	})
}

// Models

// Item represents an item entity
type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

// CreateItemResponse represents the response after creating an item
type CreateItemResponse struct {
	ItemID int `json:"item_id"`
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

// ListItemsResponse represents the response for listing items
type ListItemsResponse struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}
