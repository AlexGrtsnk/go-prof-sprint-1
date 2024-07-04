package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// функция теста
func TestMyAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ErrCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did not panic")
		}
	}()
	analysistest.Run(t, ".../internal", ErrCheckAnalyzer, "./...")

}
