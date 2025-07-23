package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func TestGenericItem(t *testing.T) {
	item := GenericItem{
		Title:       "Test Title",
		Description: "Test Description",
		Value:       "test-value",
	}
	
	if item.FilterValue() != "Test Title Test Description" {
		t.Errorf("Expected FilterValue to return 'Test Title Test Description', got %q", item.FilterValue())
	}
}

func TestNewList(t *testing.T) {
	items := []list.Item{
		GenericItem{
			Title:       "Item 1",
			Description: "Description 1",
			Value:       1,
		},
		GenericItem{
			Title:       "Item 2", 
			Description: "Description 2",
			Value:       2,
		},
	}
	
	listModel := NewList(items, "Test List", 80, 20)
	
	// Check that list was created
	if listModel.Items() == nil {
		t.Error("Expected list items to be set")
	}
	
	if len(listModel.Items()) != 2 {
		t.Errorf("Expected 2 items, got %d", len(listModel.Items()))
	}
}

func TestNewListWithEmptyItems(t *testing.T) {
	var items []list.Item
	
	listModel := NewList(items, "Empty List", 80, 20)
	
	// Should handle empty items gracefully - list.Items() returns empty slice, not nil
	if len(listModel.Items()) != 0 {
		t.Errorf("Expected 0 items, got %d", len(listModel.Items()))
	}
}

func TestGenericItemInterface(t *testing.T) {
	item := GenericItem{
		Title:       "Test",
		Description: "Test Description",
		Value:       42,
	}
	
	// Test that GenericItem implements list.Item interface
	var listItem list.Item = item
	
	if listItem.FilterValue() != "Test Test Description" {
		t.Error("GenericItem does not properly implement list.Item interface")
	}
}

func TestGenericItemWithComplexValue(t *testing.T) {
	type ComplexValue struct {
		ID   int
		Name string
	}
	
	complexVal := ComplexValue{ID: 123, Name: "Complex"}
	
	item := GenericItem{
		Title:       "Complex Item",
		Description: "Item with complex value",  
		Value:       complexVal,
	}
	
	// Test that complex values are stored correctly
	storedValue, ok := item.Value.(ComplexValue)
	if !ok {
		t.Error("Expected value to be of type ComplexValue")
	}
	
	if storedValue.ID != 123 || storedValue.Name != "Complex" {
		t.Error("Complex value not stored correctly")
	}
}