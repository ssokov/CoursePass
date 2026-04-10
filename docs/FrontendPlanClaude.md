# Frontend Implementation Plan

## Context

Go-бэкенд с JSON-RPC API (`/v1/rpc`). Нужно создать Vue 3 SPA в `frontend/`, которое работает с этим API.
Текущий `index.html` — статический макет, используется как визуальный референс.

## UX-анализ (по скриншотам)

**4 экрана:**
1. **Login / Register** — подразумевается API (`auth.Login`, `auth.RegisterStudent`)
2. **Courses** — грид карточек: название, описание, кол-во вопросов, кнопка "Pass"
3. **Exam** — прогресс-бар, счётчик вопросов, тип (single/multiple), текст, фото, варианты ответов, "Submit Answer"
4. **History** — таблица: курс, правильные ответы (x/total), статус (Passed/Failed), дата

**Навигация:** Header "Backend Courses" + кнопки [Courses] [History]. Auth guard на защищённые роуты.

**Дизайн-токены:** бежевый фон `#f4f2ec`, белые панели `#fffdf7`, teal `#0f766e`, Space Grotesk (заголовки), IBM Plex Sans (тело), скруглённые карточки, мягкие тени.

## API (источник истины: `api.ts`)

- `auth.Login({studentLogin})` -> `IToken`
- `auth.RegisterStudent({studentDraft})` -> `IToken`
- `auth.ValidateStudent` / `auth.ValidateStudentLogin` -> `IFieldError[]`
- `course.List({page, pageSize})` -> `ICourseSummary[]`
- `course.ByID({courseID})` -> `ICourse`
- `course.Me()` -> `IStudent`
- `exam.Start({courseID})` -> `IExamStart` (examId + questionIds)
- `exam.GetQuestion({examID, questionID})` -> `IQuestion`
- `exam.Answer({examID, questionID, optionIDs})` -> void
- `exam.Submit({examID})` -> `IExamResult`
- `exam.History({page, pageSize})` -> `IExamSummary[]`

Транспорт: JSON-RPC POST `/v1/rpc`. Авторизация: `Authorization: Bearer <token>`.

## Стек

- Vue 3 (Composition API, `<script setup>`)
- Vue Router (hash mode)
- VueUse (`useLocalStorage`)
- Tailwind CSS v3
- Vite
- TypeScript

## Структура проекта

```
frontend/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
├── postcss.config.js
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── api/
│   │   ├── api.ts              # сгенерированный — НЕ РЕДАКТИРОВАТЬ
│   │   └── client.ts           # JSON-RPC send + Bearer token
│   ├── composables/
│   │   └── useAuth.ts          # token, login/register/logout, isAuthenticated
│   ├── router/
│   │   └── index.ts            # routes + auth guard
│   ├── views/
│   │   ├── LoginView.vue
│   │   ├── CoursesView.vue
│   │   ├── ExamView.vue
│   │   └── HistoryView.vue
│   ├── components/
│   │   ├── AppHeader.vue
│   │   ├── CourseCard.vue
│   │   └── QuestionCard.vue
│   └── assets/
│       └── main.css            # tailwind directives + кастомные стили
```

---

## Задачи

### Задача 1: Scaffold и конфигурация проекта

**Зависимости:** нет

- Создать Vite проект с шаблоном `vue-ts` в `frontend/`
- Установить зависимости: `vue-router`, `@vueuse/core`, `tailwindcss`, `postcss`, `autoprefixer`
- Настроить Tailwind: кастомные цвета, шрифты, border-radius из дизайн-токенов
- Настроить Vite proxy: `/v1` -> `http://localhost:8075`
- Подключить Google Fonts (Space Grotesk, IBM Plex Sans) в `index.html`
- Создать `src/assets/main.css` с tailwind directives

**Результат:** `npm run dev` запускается без ошибок, пустая страница с правильным фоном и шрифтами.

---

### Задача 2: API-слой (client.ts + api.ts)

**Зависимости:** Задача 1

- Скачать `api.ts` через curl в `frontend/src/api/api.ts`
- Создать `client.ts`:
  - Функция `send(method, params)`: POST на `/v1/rpc`, JSON-RPC 2.0 формат
  - Автоматически подставляет `Authorization: Bearer <token>` из localStorage
  - Обработка ошибок JSON-RPC (error.code, error.message)
  - При 401 — очистка токена, редирект на login
- Экспорт типизированного `api` через `factory(send)`

**Результат:** `api.course.list({page: 1, pageSize: 10})` возвращает данные от бэкенда.

---

### Задача 3: Auth (composable + LoginView + router guard)

**Зависимости:** Задача 2

- `useAuth()` composable:
  - `token` — `useLocalStorage('token', '')`
  - `isAuthenticated` — computed от наличия токена
  - `login(login, password)` — вызывает `api.auth.login`, сохраняет токен
  - `register(draft)` — вызывает `api.auth.registerStudent`, сохраняет токен
  - `logout()` — очищает токен
- `LoginView.vue`:
  - Переключатель Login / Register
  - Форма логина: login, password
  - Форма регистрации: login, email, password, firstName, lastName
  - Валидация через `auth.ValidateStudent` / `auth.ValidateStudentLogin`
  - Отображение `IFieldError[]`
- Router:
  - `/login` — LoginView (только для неавторизованных)
  - `/` — redirect на `/courses`
  - `/courses` — CoursesView (auth required)
  - `/exam/:courseId` — ExamView (auth required)
  - `/history` — HistoryView (auth required)
  - `beforeEach` guard: нет токена → `/login`

**Результат:** можно залогиниться/зарегистрироваться, токен сохраняется, неавторизованные редиректятся на логин.

---

### Задача 4: AppHeader + App.vue layout

**Зависимости:** Задача 3

- `AppHeader.vue`:
  - Заголовок "Backend Courses"
  - Кнопки [Courses] [History] — `router-link`
  - Информация о пользователе (из `course.Me()`) + кнопка Logout
  - Скрыт на странице логина
- `App.vue`:
  - AppHeader + `<router-view>`
  - Общий layout: `.wrap` контейнер с max-width 1100px

**Результат:** навигация между экранами работает, header отображает имя пользователя.

---

### Задача 5: CoursesView + CourseCard

**Зависимости:** Задача 4

- `CoursesView.vue`:
  - На mount: `course.List({page: 1, pageSize: 50})`
  - Loading state
  - Грид карточек: 3 колонки десктоп, 1 колонка мобайл
- `CourseCard.vue`:
  - Props: `ICourseSummary`
  - Заголовок, описание (через `course.ByID` или из summary), кнопка "Pass"
  - "Pass" → `router.push('/exam/' + courseId)`

**Результат:** отображаются курсы из API, кнопка "Pass" ведёт на экзамен.

---

### Задача 6: ExamView + QuestionCard

**Зависимости:** Задача 5

- `ExamView.vue`:
  - На mount: `exam.Start({courseID})` → получаем `examId`, `questionIds`
  - Заголовок "Pass Course: {courseName}" (через `course.ByID`)
  - Подзаголовок "Answer all N questions."
  - Прогресс-бар (текущий вопрос / всего)
  - Последовательная загрузка вопросов: `exam.GetQuestion({examID, questionID})`
  - По "Submit Answer": `exam.Answer({examID, questionID, optionIDs})` → следующий вопрос
  - На последнем вопросе: "Finish Course" → `exam.Submit({examID})` → показать результат или перейти в History
- `QuestionCard.vue`:
  - Props: `IQuestion`, номер, всего
  - Счётчик "Question X of Y"
  - Тип: "Single choice (radio)" / "Multiple choice (checkboxes)"
  - Текст вопроса
  - Фото (если `photoUrl` есть)
  - Варианты: radio для single, checkbox для multiple
  - Emit выбранных `optionIds`

**Результат:** полный флоу экзамена работает от начала до конца.

---

### Задача 7: HistoryView

**Зависимости:** Задача 4

- `HistoryView.vue`:
  - На mount: `exam.History({page: 1, pageSize: 50})`
  - Таблица: Course, Correct Answers, Result, Date
  - Course name: resolve `courseId` → title (через отдельный запрос или кэш)
  - Correct Answers: `correctAnswers/totalQuestions` (из `IExamSummary` — нужно проверить наличие полей, возможно из `finalScore`)
  - Result: badge "Passed" (зелёный) / "Failed" (красный) по `status`
  - Date: `finishedAt` отформатированный

**Результат:** история экзаменов отображается корректно.

---

## Проверка

1. `cd frontend && npm run dev` — приложение стартует
2. Логин/регистрация работают, токен сохраняется
3. Страница курсов показывает карточки из API
4. "Pass" запускает экзамен, вопросы рендерятся (включая фото)
5. После всех ответов — результат / переход в историю
6. Таблица истории отображает прошлые результаты
7. Адаптивная вёрстка на мобильных
