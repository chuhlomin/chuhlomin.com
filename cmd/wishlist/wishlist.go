package main

import (
	"fmt"
	"os"

	"github.com/chuhlomin/chuhlomin.com/internal/models"
	"gopkg.in/yaml.v3"
)

// WishlistManager handles reading and writing wishlist YAML files
type WishlistManager struct {
	filePath string
}

// NewWishlistManager creates a new wishlist manager for the given file path
func NewWishlistManager(filePath string) *WishlistManager {
	return &WishlistManager{filePath: filePath}
}

// ReadItems reads all items from the wishlist YAML file
func (w *WishlistManager) ReadItems() ([]models.WishlistItem, error) {
	// If file doesn't exist, return empty list
	if _, err := os.Stat(w.filePath); os.IsNotExist(err) {
		return []models.WishlistItem{}, nil
	}

	// Read file
	data, err := os.ReadFile(w.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML
	var items []models.WishlistItem
	if err := yaml.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return items, nil
}

// AppendItem adds a new item to the end of the wishlist while preserving comments
func (w *WishlistManager) AppendItem(item models.WishlistItem) error {
	// Read the file as raw bytes
	data, err := os.ReadFile(w.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create new with single item
			return w.WriteItems([]models.WishlistItem{item})
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse as YAML node to preserve comments and structure
	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// The root node is a DocumentNode, the actual content is in Content[0]
	// which should be a SequenceNode (array)
	if len(rootNode.Content) == 0 {
		// Empty file, create new sequence
		return w.WriteItems([]models.WishlistItem{item})
	}

	sequenceNode := rootNode.Content[0]
	if sequenceNode.Kind != yaml.SequenceNode {
		return fmt.Errorf("expected YAML array at root")
	}

	// Create a new mapping node for the item
	itemNode := w.itemToNode(item)

	// Append to sequence
	sequenceNode.Content = append(sequenceNode.Content, itemNode)

	// Write back with preserved structure
	return w.writeNode(&rootNode)
}

// WriteItems writes all items to the wishlist YAML file
func (w *WishlistManager) WriteItems(items []models.WishlistItem) error {
	// Ensure directory exists
	dir := ""
	for i := len(w.filePath) - 1; i >= 0; i-- {
		if w.filePath[i] == '/' || w.filePath[i] == '\\' {
			dir = w.filePath[:i]
			break
		}
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Marshal to YAML with 2-space indent
	data, err := yaml.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(w.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// itemToNode converts a WishlistItem to a YAML mapping node
func (w *WishlistManager) itemToNode(item models.WishlistItem) *yaml.Node {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: item.Name},
			{Kind: yaml.ScalarNode, Value: "type"},
			{Kind: yaml.ScalarNode, Value: item.Type},
			{Kind: yaml.ScalarNode, Value: "url"},
			{Kind: yaml.ScalarNode, Value: item.URL},
			{Kind: yaml.ScalarNode, Value: "price"},
			{Kind: yaml.ScalarNode, Value: item.Price},
		},
	}

	// Add image field if present
	if item.Image != "" {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "image"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: item.Image},
		)
	}

	return node
}

// writeNode writes a YAML node to the wishlist file, preserving formatting
func (w *WishlistManager) writeNode(node *yaml.Node) error {
	// Ensure directory exists
	dir := ""
	for i := len(w.filePath) - 1; i >= 0; i-- {
		if w.filePath[i] == '/' || w.filePath[i] == '\\' {
			dir = w.filePath[:i]
			break
		}
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create file
	file, err := os.Create(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode with proper indentation
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	defer encoder.Close()

	if err := encoder.Encode(node); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	return nil
}
