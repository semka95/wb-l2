/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

package main

import (
	"fmt"
	"sort"
	"strings"
)

func unorderedEqual(str1, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}

	exists := make(map[rune]struct{})
	for _, v := range str1 {
		exists[v] = struct{}{}
	}
	for _, v := range str2 {
		if _, ok := exists[v]; !ok {
			return false
		}
	}

	return true
}

func anagram(data []string) map[string][]string {
	res := make(map[string][]string)
	exists := make(map[string]struct{})

	for i, v := range data {
		for j := i + 1; j < len(data); j++ {
			if _, ok := exists[data[j]]; ok {
				continue
			}

			if unorderedEqual(v, data[j]) {
				res[v] = append(res[v], data[j])
				exists[data[j]] = struct{}{}
				exists[v] = struct{}{}
			}
		}
	}

	return res
}

func prepareData(data []string) []string {
	res := make([]string, 0, len(data))
	exists := make(map[string]struct{})

	for _, v := range data {
		str := strings.ToLower(v)

		if _, ok := exists[str]; ok {
			continue
		}

		res = append(res, str)
		exists[str] = struct{}{}
	}

	sort.Strings(res)

	return res
}

func main() {
	data := []string{"пятка", "пятак", "тяпка", "листок", "столик", "столик", "СЛИТОК"}
	data = prepareData(data)
	res := anagram(data)

	fmt.Println(res)
}
