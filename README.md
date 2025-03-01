# Web Calculator

Проект представляет собой распределенную систему для вычисления математических выражений. Система состоит из двух основных компонентов:
1. **Оркестратор** — принимает выражения, разбивает их на задачи и распределяет между воркерами.
2. **Агент** — управляет воркерами и выполняет задачи.

---

## Как это работает?

### Схема работы системы

```plaintext
+-------------------+       +-------------------+    
|                   |       |                   |     
|   Оркестратор     | <---> |      Агент        |
|                   |       |                   |     
+-------------------+       +-------------------+     
      ^                         ^                         
      |                         |                         
      |                         |                         
      v                         v                         
+-------------------+       +-------------------+    
|                   |       |                   |  
|  Пользователь     |       |  Вычисления       |  
|  (HTTP-запросы)   |       |  (арифметика)     |      
|                   |       |                   |     
+-------------------+       +-------------------+   
``` 
1. Пользователь отправляет математическое выражение на оркестратор через HTTP-запрос.

2. Оркестратор разбивает выражение на задачи и отправляет их агенту.

3. Агент распределяет задачи между воркерами.

4. Воркеры выполняют вычисления и возвращают результаты агенту.

5. Агент отправляет результаты обратно оркестратору.

6. Оркестратор сохраняет результаты и возвращает их пользователю. 

1. Запуск с помощью Docker Compose
   Для удобства запуска системы используется docker-compose. Убедитесь, что у вас установлены Docker и Docker Compose.

## Запуск системы

Клонируйте репозиторий:

```bash
git clone https://github.com/InsafMin/web-calculator.git
cd web-calculator
```
Запустите систему:

```bash
docker-compose up --build
```
Это запустит:

 - Оркестратор на http://localhost:8080

 - Агент и воркеры

## Примеры запросов
### 1. Отправка выражения
   Отправьте математическое выражение на оркестратор.

 - Метод: POST

 - URL: http://localhost:8080/api/v1/calculate

 - Тело запроса:

```json
{
"expression": "2 + 2 * 2"
}
```
 - Пример с curl:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "2 + 2 * 2"}'
```
 - Ответ:

```json
{
"id": "12345"
}
```
### 2. Получение списка выражений
   Получите список всех выражений и их статусов.

 - Метод: GET

 - URL: http://localhost:8080/api/v1/expressions

 - Пример с curl:

```bash
curl http://localhost:8080/api/v1/expressions
```
 - Ответ:

```json
{
   "expressions": [
      {
         "id": "12345",
         "expression": "2 + 2 * 2",
         "status": "done",
         "result": 6
      }
   ]
}
```
### 3. Получение результата конкретного выражения
   Получите результат конкретного выражения по его ID.

 - Метод: GET

 - URL: http://localhost:8080/api/v1/expressions/{id}

 - Пример с curl:

```bash
curl http://localhost:8080/api/v1/expressions/12345
```
 - Ответ:

```json
{
   "expression": {
      "id": "12345",
      "expression": "2 + 2 * 2",
      "status": "done",
      "result": 6
   }
}
```
## Документация API
### 1. Отправка выражения
 - Метод: POST

 - URL: /api/v1/calculate

 - Тело запроса: {"expression": "математическое выражение"}

### 2. Получение списка выражений
 - Метод: GET

 - URL: /api/v1/expressions

### 3. Получение результата конкретного выражения

 - Метод: GET

 - URL: /api/v1/expressions/{id}

## Контакты
Если у вас есть вопросы или предложения, свяжитесь с автором проекта:

 - Имя: Insaf Mingazov

 - GitHub: [InsafMin](https://github.com/InsafMin)
 - Email: [insaf.min.in@yandex.ru](mailto:insaf.min.in@yandex.ru)