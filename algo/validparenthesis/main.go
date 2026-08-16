package main

// Дана строка, `s`содержащая только символы `'('`, `')'`, `'{'`, `'}'`, `'['`и `']'`, определите, является ли входная строка допустимой.

// Входная строка действительна, если:

// 1. Открытые скобки должны быть закрыты скобками того же типа.
// 2. Открытые скобки должны быть закрыты в правильном порядке.
// 3. Каждой закрывающейся скобке соответствует открывающаяся скобка того же типа.

// Пример 1:
// Ввод: s = "()"
// Вывод: true

// Пример 2:
// Ввод: s = "()[]{}"
// Вывод: true

// Пример 3:
// Ввод: s = "(]"
// Вывод: false

// Пример 4:
// Ввод: s = "([])"
// Вывод: true

// Ограничения:
// - `s`состоит только из скобок `'()[]{}'`.

import "fmt"

func isOpenParenthesis(sym rune) bool {
	openParenthesis := []rune{'(', '[', '{'}

	for _, parenthesis := range openParenthesis {
		if parenthesis == sym {
			return true
		}
	}

	return false
}

func isCorrespondingCloseParenthesis(sym rune, lastSym rune) bool {
	closeParenthesisMap := map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
	}

	return closeParenthesisMap[lastSym] == sym
}

func isValid(s string) bool {
	stack := make([]rune, 0)

	for _, sym := range s {
		if isOpenParenthesis(sym) {
			stack = append(stack, sym)
			continue
		}

		if len(stack) == 0 {
			return false
		}

		lastSymStack := stack[len(stack)-1:][0]

		if isCorrespondingCloseParenthesis(sym, lastSymStack) {
			stack = stack[:len(stack)-1]
			continue
		}

		return false
	}

	return len(stack) == 0
}

func main() {
	// Тестовые примеры
	fmt.Println(isValid("()"))     // true
	fmt.Println(isValid("()[]{}")) // true
	fmt.Println(isValid("(]"))     // false
	fmt.Println(isValid("([])"))   // true
	fmt.Println(isValid("([)]"))   // false
}
