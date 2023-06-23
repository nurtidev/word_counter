package main

import (
	"bufio"
	"fmt"
	"os"
)

type word struct {
	value     []byte
	frequency int
}

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var words []word
	scanner := bufio.NewScanner(file)
	scanner.Split(scanWords)

	for scanner.Scan() {
		w := scanner.Bytes()
		w = bytesToLower(w)
		index := findWord(words, w)
		if index != -1 {
			words[index].frequency++
		} else {
			words = append(words, word{value: w, frequency: 1})
		}
	}

	words = sort(words)

	for i := 0; i < 20 && i < len(words); i++ {
		fmt.Printf("%d %s\n", words[i].frequency, string(words[i].value))
	}
}

// scanWords выполняет разбиение данных на слова в байтовом срезе.
// Возвращает позицию, на которой следует продолжить сканирование, токен (слово) и ошибку.
func scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0

	// Ищем начало следующего слова, пропуская символы, не являющиеся буквами.
	for ; start < len(data); start++ {
		if isLetter(data[start]) {
			break
		}
	}

	// Если достигнут конец данных, возвращаем текущую позицию и nil в качестве токена и ошибки.
	if start == len(data) {
		if atEOF {
			return start, nil, nil
		}
		return len(data), nil, nil
	}

	// Ищем конец текущего слова, чтобы вернуть его в качестве токена.
	for i := start; i < len(data); i++ {
		// Если текущий символ не является буквой, возвращаем позицию, токен (слово) и nil в качестве ошибки.
		if !isLetter(data[i]) {
			if i == start {
				start++
				continue
			} else {
				return i, data[start:i], nil
			}
		}
	}

	// Если достигнут конец данных, возвращаем текущую позицию и оставшиеся данные
	// в качестве токена, а также nil в качестве ошибки.
	if atEOF {
		return len(data), data[start:], nil
	}

	// В противном случае, возвращаем текущую позицию, nil в качестве токена и ошибки.
	return start, nil, nil
}

func findWord(words []word, value []byte) int {
	for i := range words {
		if bytesEqual(words[i].value, value) {
			return i
		}
	}
	return -1
}

func bytesToLower(b []byte) []byte {
	var result []byte
	for _, c := range b {
		result = append(result, lowerCase(c))
	}
	return result
}

func lowerCase(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func bytesEqual(b1, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

func sort(words []word) []word {
	for i := 0; i < len(words)-1; i++ {
		minIndex := i
		for j := i + 1; j < len(words); j++ {
			if words[j].frequency > words[minIndex].frequency {
				minIndex = j
			}
		}
		if minIndex != i {
			words[i], words[minIndex] = words[minIndex], words[i]
		}
	}
	return words
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}
