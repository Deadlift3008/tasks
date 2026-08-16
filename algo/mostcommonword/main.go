package main

// Задача: Наиболее частое слово, не входящее в список запрещённых

// Дана строка `paragraph` и массив строк `banned`, содержащий запрещённые слова. Нужно вернуть наиболее часто встречающееся слово из `paragraph`, которое не является запрещённым.

// Гарантируется, что хотя бы одно слово не является запрещённым, и ответ всегда уникален.

// Слова в `paragraph` регистронезависимые, а ответ должен быть возвращён в нижнем регистре.

// ---

// Пример 1:

// - Ввод:  `paragraph = "Bob hit a ball, the hit BALL flew far after it was hit."`, `banned = ["hit"]`
// - Вывод:**  `"ball"`

// **Пояснение:**
// Слово `"hit"` встречается 3 раза, но оно запрещено.
// Слово `"ball"` встречается 2 раза и не является запрещённым, поэтому это ответ.

// Обратите внимание, что слова регистронезависимы и знаки пунктуации игнорируются.

// ---

// Пример 2:
// - Ввод: `paragraph = "a."`,  `banned = []`
// - Вывод:  `"a"`

// ---

// ### Ограничения:
// Строка `paragraph` состоит из английских букв, пробелов `' '`, или одного из символов `!?',;.`.
// `banned[i]` состоит только из строчных английских букв.

import (
	"fmt"
	"strings"
)

func mostCommonWord(paragraph string, banned []string) string {
	bannedMap := make(map[string]struct{})

	for _, bannedWord := range banned {
		bannedMap[strings.ToLower(bannedWord)] = struct{}{}
	}

	ignoredRunes := []rune{'!', '?', '\'', ',', ';', '.'}

	var isIgnored bool

	var sb strings.Builder

	for _, sym := range paragraph {
		isIgnored = false
		for _, ignoredSym := range ignoredRunes {
			if sym == ignoredSym {
				isIgnored = true
				break
			}
		}

		if !isIgnored {
			sb.WriteRune(sym)
		}
	}

	words := strings.Fields(sb.String())

	wordsMap := make(map[string]int)

	var maxCountWord string
	var maxCount int

	for _, word := range words {
		lowerWord := strings.ToLower(word)

		_, ok := bannedMap[lowerWord]

		if !ok {
			newVal := wordsMap[lowerWord] + 1
			wordsMap[lowerWord] = newVal

			if newVal > maxCount {
				maxCount = newVal
				maxCountWord = lowerWord
			}
		}
	}

	return maxCountWord
}

func main() {
	paragraph := "Bob hit a ball, the hit BALL flew far after it was hit."
	banned := []string{"hit"}
	result := mostCommonWord(paragraph, banned)
	fmt.Println("Result:", result) // Ожидаемый вывод: "ball"
}
