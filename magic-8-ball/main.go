package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	// answers are listed from Positive to Negative ones
	switch dice := rand.IntN(20) + 1; dice {
	case 1:
		fmt.Println("It is certain (Бесспорно)")
	case 2:
		fmt.Println("It is decidedly so (Предрешено)")
	case 3:
		fmt.Println("Without a doubt (Никаких сомнений)")
	case 4:
		fmt.Println("Yes — definitely (Определённо да)")
	case 5:
		fmt.Println("You may rely on it (Можешь быть уверен в этом)")
	case 6:
		fmt.Println("As I see it, yes (Мне кажется — «да»)")
	case 7:
		fmt.Println("Most likely (Вероятнее всего)")
	case 8:
		fmt.Println("Outlook good (Хорошие перспективы)")
	case 9:
		fmt.Println("Signs point to yes (Знаки говорят — «да»)")
	case 10:
		fmt.Println("Yes (Да)")
	case 11:
		fmt.Println("Reply hazy, try again (Пока не ясно, попробуй снова)")
	case 12:
		fmt.Println("Ask again later (Спроси позже)")
	case 13:
		fmt.Println("Better not tell you now (Лучше не рассказывать)")
	case 14:
		fmt.Println("Cannot predict now (Сейчас нельзя предсказать)")
	case 15:
		fmt.Println("Concentrate and ask again (Сконцентрируйся и спроси опять)")
	case 16:
		fmt.Println("Don’t count on it (Даже не думай)")
	case 17:
		fmt.Println("My reply is no (Мой ответ — «нет»)")
	case 18:
		fmt.Println("My sources say no (По моим данным — «нет»)")
	case 19:
		fmt.Println("Outlook not so good (Перспективы не очень хорошие)")
	case 20:
		fmt.Println("Very doubtful (Весьма сомнительно)")
	}
}
