Реализация функции tr для интернационализации

## Ограничения и требования
Доступное подмножество языка: функции, ассоциативный массив (map), struct, пакет для работы с переменными окружения

Минимально поддерживаемое количество локалей - две.

## Интерфейс
```go
tr("text") // получение заданного текста в локали ОС (если перевода нет, то возвращается текст переданный в функцию tr)
```

## Пример
```go
tr("Hello") // "Привет" (ru_RU локаль, есть перевод)
tr("Wow") // "Wow" (ru_RU локаль, нет перевода)
```

## Демо

Example 1
```sh
LANG=C.UTF-8 make run
```
```
Apple
Something went wrong
wtf?
bla-bl
```
Example 2
```sh
LANG=ru_RU.UTF-8 make run
```
```
Яблоко
Что-то пошло не так
wtf?
bla-bla
```

## Ресурсы:

- [LC_ALL and LC_*](https://unix.stackexchange.com/a/87763)
- [ISO 3166-1 Codes for the representation of names of countries and their subdivision - Part 1: Country code](https://en.wikipedia.org/wiki/ISO_3166-1#Naming_and_code_construction)