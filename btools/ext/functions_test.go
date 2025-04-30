package ext

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	// Тест на nil-слейс: ожидание false
	var nilSlice []int
	if Contains(nilSlice, 1) {
		t.Errorf("Expected false for nil slice, got true")
	}

	// Тест, когда элемент есть в слайсе
	s := []int{1, 2, 3}
	if !Contains(s, 2) {
		t.Errorf("Expected true for element in slice, got false")
	}

	// Тест, когда элемента нет в слайсе
	if Contains(s, 4) {
		t.Errorf("Expected false for element not in slice, got true")
	}
}

func TestTernary(t *testing.T) {
	result := Ternary(true, "yes", "no")
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%v'", result)
	}

	result2 := Ternary(false, 10, 20)
	if result2 != 20 {
		t.Errorf("Expected 20, got %v", result2)
	}
}

func TestUnion(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	result := Union(a, b)
	// Ожидаемые уникальные элементы: 1, 2, 3, 4, 5
	expectedElements := []int{1, 2, 3, 4, 5}

	// Проверка длины результата
	if len(result) != len(expectedElements) {
		t.Errorf("Expected union length %d, got %d", len(expectedElements), len(result))
	}
	// Проверка наличия каждого элемента из expectedElements в результате
	for _, v := range expectedElements {
		if !Contains(result, v) {
			t.Errorf("Expected union to contain %d", v)
		}
	}
}

type Person struct {
	Name string
	Age  int
}

func TestExtractField(t *testing.T) {
	// Тест для слайса структур
	persons := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	names := ExtractField[Person, string](persons, "Name")
	expectedNames := []string{"Alice", "Bob"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf("Expected %v, got %v", expectedNames, names)
	}

	// Тест для слайса указателей на структуры
	personsPtr := []*Person{
		{Name: "Charlie", Age: 28},
		{Name: "Diana", Age: 32},
	}
	namesPtr := ExtractField[*Person, string](personsPtr, "Name")
	expectedNamesPtr := []string{"Charlie", "Diana"}
	if !reflect.DeepEqual(namesPtr, expectedNamesPtr) {
		t.Errorf("Expected %v, got %v", expectedNamesPtr, namesPtr)
	}

	// Тест для несуществующего поля: должно вернуть пустой слайс
	result := ExtractField[Person, string](persons, "NonExistent")
	if len(result) != 0 {
		t.Errorf("Expected empty slice for non-existent field, got %v", result)
	}
}

func TestDiff(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4, 5}
	result := Diff(s1, s2)
	expected := []int{4, 5}

	// Проверка длины результата
	if len(result) != len(expected) {
		t.Errorf("Expected diff length %d, got %d", len(expected), len(result))
	}
	// Проверка, что результат содержит ожидаемые элементы
	for _, v := range expected {
		if !Contains(result, v) {
			t.Errorf("Expected diff to contain %d", v)
		}
	}

	// Если все элементы из s2 уже есть в s1, ожидается пустой слайс
	s3 := []int{1, 2, 3}
	result2 := Diff(s1, s3)
	if len(result2) != 0 {
		t.Errorf("Expected empty diff, got %v", result2)
	}
}
