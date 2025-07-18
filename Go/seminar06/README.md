# CRUD

---

## **1. Введение в CRUD**  
- **Что такое CRUD?**  
  - Базовая концепция: Create (создание), Read (чтение), Update (обновление), Delete (удаление).  

- **Разбор примеров**  
  - Как CRUD реализован в классическом REST API.  
    - GET /users - получить список пользователей  
    - GET /users/1 - получить пользователя с id = 1  
    - POST /users - создать пользователя  
    - PUT /users/1 - обновить пользователя с id = 1  
    - DELETE /users/1 - удалить пользователя с id = 1  

---

## **2. Реализация CRUD в базе данных**  
- **Create (INSERT)**  
  - Добавление новых записей в БД.  
  - Использование `RETURNING` (PostgreSQL) для получения ID новой записи.  
  - Автоинкрементные поля и UUID как идентификаторы.  

- **Read (SELECT)**  
  - Простые и сложные `SELECT`-запросы.  
  - Использование индексов (`PRIMARY KEY`, `UNIQUE`, `INDEX`) для ускорения поиска.  
  - Фильтрация (`WHERE`, `ORDER BY`, `LIMIT`, `OFFSET`).  

- **Update (UPDATE)**  
  - Обновление записей с `SET`.  
  - Использование `RETURNING` для получения измененных данных.  
  - Влияние `UPDATE` на производительность (неизменяемые столбцы, обновление индексов).  

- **Delete (DELETE)**  
  - Физическое (`DELETE`) vs логическое удаление (`is_deleted = true`).  
  - Каскадное удаление (`CASCADE`).  

---

## **3. CRUD в API и микросервисах**  
- **RESTful CRUD**  
  - Соответствие CRUD и HTTP-методов:  
    - `POST /users` → создание.  
    - `GET /users/{id}` → чтение.  
    - `PUT /users/{id}` → обновление.  
    - `DELETE /users/{id}` → удаление.  

- **Генерация API-документации**  
  - OpenAPI (Swagger) для описания CRUD API.  
  - Postman для тестирования запросов.  

---

## **4. Оптимизация CRUD-операций**  
- **Индексы и их влияние**  
  - Как индексы ускоряют `SELECT`, но замедляют `INSERT/UPDATE/DELETE`.  
  - Виды индексов (`B-TREE`, `HASH`, `GIN`).  

- **Пагинация в CRUD**  
  - **OFFSET-LIMIT** (простая, но медленная на больших данных).  
  - **Keyset pagination** (по `id` или `created_at`, более производительная).  
  - **Оконные функции (ROW_NUMBER)** для гибкой пагинации.  

- **Кеширование CRUD-запросов**  
  - Redis / Memcached для хранения часто запрашиваемых данных.  
  - Когда стоит кешировать: только `READ` или также `WRITE`?  

- **Шардирование и репликация в больших CRUD-системах**  
  - Разделение базы по пользователям (sharding).  
  - Использование read-replic для разгрузки `SELECT`-запросов.  

---

## **5. CRUD и современные технологии**  
- **CRUD в микросервисной архитектуре**  
  - Разделение CRUD между сервисами (например, отдельный сервис для пользователей, товаров, заказов).  
  - Проблемы согласованности данных (eventual consistency, SAGA-паттерн).  

- **NoSQL vs SQL в CRUD**  
  - MongoDB: работа с документами вместо строк таблиц.  
  - Redis как key-value хранилище для ускорения операций.  

- **Event-driven подход в CRUD**  
  - CQRS (разделение команд и запросов).  
  - Event Sourcing (CRUD через события вместо изменения БД).  

## 6. Пагинация

### 1. **OFFSET-LIMIT (Классическая пагинация)**
**Описание**: Использует `OFFSET` для пропуска записей и `LIMIT` для ограничения выборки.  
**Пример SQL-запроса**:  
```sql
SELECT * FROM items ORDER BY created_at DESC LIMIT 10 OFFSET 20;
```
#### **Плюсы**:
- Простая реализация.
- Подходит для небольших и средних таблиц.
- Хорошо работает при небольшом количестве страниц.  

#### **Минусы**:
- **Неэффективно на больших данных**: Чем больше `OFFSET`, тем дольше выполнение запроса.
- **Проблема пропуска записей**: Если данные изменяются в процессе пагинации, возможны дубликаты или пропуски.
- **Нагрузка на БД**: БД сначала выбирает все `OFFSET + LIMIT` записей, а потом отбрасывает `OFFSET`.  

---

### 2. **Пагинация по курсору (Keyset Pagination)**
**Описание**: Использует `WHERE` для фильтрации на основе ключа последней записи предыдущей страницы.  
**Пример SQL-запроса**:  
```sql
SELECT * FROM items WHERE created_at < '2025-03-10 12:00:00' ORDER BY created_at DESC LIMIT 10;
```
#### **Плюсы**:
- **Высокая производительность**: Нет необходимости пропускать записи (`OFFSET`), сразу выбираются нужные.
- **Стабильность данных**: Нет проблемы дубликатов и пропущенных записей при изменении данных.
- **Эффективно на больших объемах данных**.  

#### **Минусы**:
- **Не подходит для произвольного доступа (например, "перейти на страницу 50")**.
- **Зависит от уникального поля (например, `id` или `created_at`)**.
- **Сложность реализации**: Клиент должен передавать последний загруженный курсор.  

---

### 3. **Пагинация по оконным функциям (ROW_NUMBER)**
**Описание**: Использует оконную функцию `ROW_NUMBER()` для нумерации записей.  
**Пример SQL-запроса**:
```sql
WITH numbered AS (
    SELECT id, name, ROW_NUMBER() OVER (ORDER BY created_at DESC) AS row_num
    FROM items
)
SELECT * FROM numbered WHERE row_num BETWEEN 21 AND 30;
```
#### **Плюсы**:
- **Позволяет быстро получать произвольную страницу**, как в `OFFSET-LIMIT`, но без лишнего перебора строк.
- **Гибкость**: Можно комбинировать с фильтрацией.  

#### **Минусы**:
- **Тяжелая операция на больших данных**: `ROW_NUMBER()` сначала вычисляется для всей выборки.
- **Не во всех СУБД одинаково эффективно**.  
