package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

// HexConverter это интерфейс, который описывает функцию конвертации из 16тиричной системы в 10тичную
type HexConverter interface {
	Convert(hexString string) (decString string)
}

// hexConverterImpl это структура, которая реализует интерфейс HexConverter
type hexConverterImpl struct{}

// NewHexConverter это функция, которая создает новый объект hexConverterImpl
func NewHexConverter() HexConverter {
	return &hexConverterImpl{}
}

// ConvertHexChar преобразует символ 16-ричной системы в его 10-тичное значение
func ConvertHexChar(c rune) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	return int(c-'A') + 10
}

// сложение в столбик
func add(s1, s2 string) (result string) {
	var tens int //десятки

	//приводим к одной длине
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
		s1 = strings.Repeat("0", maxLen-len(s1)) + s1
	} else {
		s2 = strings.Repeat("0", maxLen-len(s2)) + s2
	}
	//123 + 1111 = 0123
	// 1111
	//+0123

	for i := maxLen - 1; i >= 0; i-- {
		sum := int(s1[i]-'0') + int(s2[i]-'0') + tens
		tens = sum / 10
		result = string('0'+sum%10) + result
	}

	// если остались ненулевые десятки
	if tens > 0 {
		result = string('0'+tens) + result
	}

	return result
}

// Умножение в столбик
func multiply(bigint string, m int) (result string) {
	var tens int

	for i := len(bigint) - 1; i >= 0; i-- {
		prod := int(bigint[i]-'0')*m + tens
		tens = prod / 10
		result = string('0'+prod%10) + result
	}

	// если остались ненулевые десятки
	if tens > 0 {
		result = string('0'+tens) + result
	}
	return result
}

func (h *hexConverterImpl) Convert(hexString string) (decString string) {
	decString = "0"
	for _, c := range strings.ToUpper(hexString) {
		decString = multiply(decString, 16)
		decString = add(decString, fmt.Sprintf("%d", ConvertHexChar(c)))
	}
	return
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	h := NewHexConverter()       //создаем структуру
	hexed, err := randomHex(100) //создаем 16теричное число
	if err != nil {
		log.Println(err)
	}
	fmt.Println(h.Convert(hexed)) //конвертируем в 10ичное
}
