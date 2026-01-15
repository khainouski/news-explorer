# Настройка CI для Go-сервиса: GitHub Actions + GHCR

## 1. Основные определения

### Что такое CI

**CI (Continuous Integration, непрерывная интеграция)** — это автоматическая проверка и сборка проекта после изменения кода.

В нашем случае после `push` или merge в ветку `main` автоматически выполняются:

1. загрузка исходного кода;
2. установка Go;
3. запуск тестов;
4. сборка Docker-образа;
5. публикация Docker-образа в реестр.

---

### Что такое GitHub Actions

**GitHub Actions** — встроенный в GitHub инструмент автоматизации.

Он запускает описанный нами workflow при определённом событии, например:

- `push` в ветку `main`;
- создание или обновление Pull Request;
- ручной запуск;
- создание Git-тега.

Workflow хранится прямо в репозитории в YAML-файле: `.github/workflows/ci.yml`

То есть основная настройка GitHub Actions выполняется **кодом в репозитории**, а не вручную через админку.

---

### Что такое Docker-образ

**Docker-образ** — готовый неизменяемый пакет приложения, содержащий:

- скомпилированный Go-сервис;
- необходимые системные файлы;
- сертификаты;
- инструкцию запуска приложения.

Из одного Docker-образа можно запускать одинаковые контейнеры локально, на сервере или в Kubernetes.

---

### Где хранятся Docker-образы

Docker-образы хранятся в специальном хранилище, которое называется **Container Registry — реестр контейнерных образов**.

Для проекта на GitHub удобно использовать **GitHub Container Registry (GHCR)**.

Адрес GHCR: `ghcr.io`

Адрес образа будет выглядеть так: `ghcr.io/<github-owner>/<repository>:<tag>`

Например:

```text
ghcr.io/oleg-dev/developer-news:latest
ghcr.io/oleg-dev/developer-news:sha-a1b2c3d
```

Где:

- `oleg-dev` — имя пользователя или организации GitHub;
- `developer-news` — имя репозитория;
- `latest` или `sha-a1b2c3d` — тег версии образа.

---

### Что такое workflow, job и step

#### Workflow

Полный автоматизированный процесс, описанный в YAML-файле.

#### Job

Отдельная задача внутри workflow. Например: `test`, `build-and-push`.

#### Step

Отдельный шаг внутри job. Например:

- скачать код;
- установить Go;
- выполнить `go test ./...`;
- собрать Docker-образ;
- отправить образ в GHCR.

---

## 2. Итоговая схема работы

```text
Разработчик
    |
    | git push / merge
    v
Ветка main в GitHub
    |
    v
GitHub Actions
    |
    +-- запускает тесты
    +-- собирает Docker-образ
    +-- присваивает теги
    +-- публикует образ
    |
    v
GitHub Container Registry
    |
    | docker pull или Kubernetes
    v
Сервер / k3s
```

---

## 3. Что нужно настроить

Для базового варианта нужны:

```text
Go-репозиторий
├── Dockerfile
├── .dockerignore
└── .github
    └── workflows
        └── ci.yml
```

Большая часть механизма настраивается файлами в репозитории.

Через веб-интерфейс GitHub нужно только:

1. проверить настройки GitHub Actions;
2. после первой публикации найти созданный package;
3. при необходимости изменить его видимость на public;
4. проверить успешный запуск workflow.

---

# 4. Пошаговая настройка

## Шаг 1. Создать GitHub-репозиторий

Создай обычный репозиторий, например `developer-news`, и загрузи в него Go-приложение.

Рекомендуемая ветка по умолчанию: `main`

---

## Шаг 2. Проверить структуру Go-приложения

Например:

```text
developer-news
├── cmd
│   └── app
│       └── main.go
├── internal
├── go.mod
├── go.sum
├── Dockerfile
├── .dockerignore
└── .github
    └── workflows
        └── ci.yml
```

В примерах ниже предполагается, что точка входа приложения находится здесь: `./cmd/app`

Если у тебя другой путь к `main.go`, его нужно заменить в `Dockerfile`.

---

## Шаг 3. Добавить Dockerfile

Создай в корне проекта файл `Dockerfile`:

```dockerfile
# Этап сборки приложения
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Сначала копируем файлы зависимостей,
# чтобы Docker мог использовать кэш.
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код.
COPY . .

# Собираем статический Linux-бинарник.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/service \
    ./cmd/app

# Минимальный runtime-образ.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/service ./service

USER app

EXPOSE 8080

ENTRYPOINT ["./service"]
```

Если приложение собирается командой `go build ./cmd/server`, то в `Dockerfile` нужно указать `./cmd/server`.

---

## Шаг 4. Добавить .dockerignore

Создай в корне проекта файл `.dockerignore`:

```dockerignore
.git
.github
.idea
.vscode
.env
.env.*
bin
tmp
coverage.out
*.log
README.md
```

Этот файл не позволяет отправлять ненужные файлы в контекст Docker-сборки. Не добавляй реальные секреты в Docker-образ.

---

## Шаг 5. Создать GitHub Actions workflow

Создай файл `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main

# Не позволяет нескольким старым сборкам одной ветки
# выполняться одновременно.
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

# Минимальные разрешения встроенного GITHUB_TOKEN.
permissions:
  contents: read
  packages: write

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.24"
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run vet
        run: go vet ./...

      - name: Run tests
        run: go test -race ./...

  build-and-push:
    name: Build and push Docker image
    runs-on: ubuntu-latest
    needs:
      - test
    # Pull Request проверяется, но образ публикуется
    # только после push или merge в main.
    if: github.event_name == 'push'
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Generate Docker tags and labels
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=raw,value=latest
            type=sha,prefix=sha-

      - name: Build and push Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## Шаг 6. Почему отдельный GitHub-токен добавлять не нужно

GitHub автоматически создаёт для каждого запуска workflow временный токен `GITHUB_TOKEN`, который используется здесь:

```yaml
password: ${{ secrets.GITHUB_TOKEN }}
```

В workflow ему выдаются только необходимые разрешения:

```yaml
permissions:
  contents: read
  packages: write
```

Это означает:

- `contents: read` — разрешено скачать код репозитория;
- `packages: write` — разрешено опубликовать образ в GHCR.

Для публикации образа из этого же репозитория вручную создавать Personal Access Token не нужно.

---

## Шаг 7. Отправить файлы в GitHub

```bash
git add Dockerfile .dockerignore .github/workflows/ci.yml
git commit -m "Add Docker build and GitHub Actions CI"
git push origin main
```

После push workflow запустится автоматически.

---

## Шаг 8. Проверить GitHub Actions через веб-интерфейс

Открой репозиторий на GitHub, перейди в `Repository → Actions`. Там должен появиться workflow `CI`.

Открой запуск и проверь jobs: `Test`, `Build and push Docker image`. Зелёная галочка означает, что job выполнен успешно.

---

## Шаг 9. Где найти опубликованный образ

После первого успешного запуска GitHub создаст package — найти его можно на странице репозитория или профиля GitHub в разделе `Packages`. Название package обычно совпадает с именем репозитория.

Адрес образа: `ghcr.io/<owner>/<repository>:latest`

Например: `ghcr.io/oleg-dev/developer-news:latest`

Также workflow создаст тег конкретного коммита: `ghcr.io/oleg-dev/developer-news:sha-a1b2c3d`

---

## Шаг 10. Сделать образ публичным

Новый package может быть приватным. Чтобы сделать его публичным:

```text
GitHub profile или organization
→ Packages
→ выбрать package
→ Package settings
→ Change visibility
→ Public
```

После этого публичный образ можно скачивать без авторизации. Для учебного публичного pet-проекта это самый простой вариант.

---

# 5. Как скачать и запустить образ

## Публичный образ

Скачать:

```bash
docker pull ghcr.io/<owner>/<repository>:latest
```

Например:

```bash
docker pull ghcr.io/oleg-dev/developer-news:latest
```

Запустить:

```bash
docker run -d \
  --name developer-news \
  --restart unless-stopped \
  -p 8080:8080 \
  ghcr.io/oleg-dev/developer-news:latest
```

Проверить:

```bash
docker ps
docker logs developer-news
```

---

## Приватный образ

Для скачивания приватного образа на сервере потребуется Personal Access Token classic с правом `read:packages`.

Токен нужно хранить в переменной окружения:

```bash
export GHCR_TOKEN="your-token"
```

Авторизация:

```bash
echo "$GHCR_TOKEN" | docker login ghcr.io \
  --username YOUR_GITHUB_USERNAME \
  --password-stdin
```

После этого:

```bash
docker pull ghcr.io/<owner>/<repository>:latest
```

Не записывай токен прямо в Dockerfile, Git-репозиторий или Kubernetes-манифест.

---

# 6. Как использовать образ в Kubernetes / k3s

В Kubernetes вручную выполнять `docker pull` обычно не нужно — Kubernetes сам скачает образ, указанный в `Deployment`.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: developer-news
  namespace: developer-news
spec:
  replicas: 1
  selector:
    matchLabels:
      app: developer-news
  template:
    metadata:
      labels:
        app: developer-news
    spec:
      containers:
        - name: developer-news
          image: ghcr.io/oleg-dev/developer-news:sha-a1b2c3d
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: 8080
```

Применить:

```bash
kubectl apply -f deployment.yaml
```

Проверить:

```bash
kubectl get pods -n developer-news
kubectl logs deployment/developer-news -n developer-news
```

---

## Приватный образ в Kubernetes

Создай токен GitHub с правом `read:packages`, затем создай Kubernetes Secret:

```bash
kubectl create secret docker-registry ghcr-secret \
  --namespace developer-news \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password="$GHCR_TOKEN"
```

Добавь в `Deployment`:

```yaml
spec:
  template:
    spec:
      imagePullSecrets:
        - name: ghcr-secret
      containers:
        - name: developer-news
          image: ghcr.io/oleg-dev/developer-news:sha-a1b2c3d
```

---

# 7. Как правильно обновлять приложение

Не рекомендуется использовать `latest` для production или GitOps.

Плохо: `image: ghcr.io/oleg-dev/developer-news:latest`

Лучше: `image: ghcr.io/oleg-dev/developer-news:sha-a1b2c3d`

Преимущества тега коммита:

- точно известно, какая версия развёрнута;
- можно быстро выполнить rollback;
- Kubernetes видит изменение тега;
- Argo CD показывает конкретное изменение;
- один и тот же тег не перезаписывается другой версией.

При GitOps-подходе процесс выглядит так:

```text
1. Код попадает в main.
2. GitHub Actions запускает тесты.
3. GitHub Actions создаёт образ sha-<commit>.
4. Образ публикуется в GHCR.
5. В GitOps-репозитории меняется image tag.
6. Argo CD замечает изменение.
7. Argo CD обновляет Deployment в k3s.
```

На первом этапе тег в GitOps-репозитории можно менять вручную.

---

# 8. Что настраивается через код, а что через админку

## Через код в репозитории

`Dockerfile`, `.dockerignore`, `.github/workflows/ci.yml`, Kubernetes `Deployment`.

Именно код определяет:

- когда запускается CI;
- какие тесты выполняются;
- как собирается приложение;
- куда отправляется образ;
- какие теги создаются.

## Через веб-интерфейс GitHub

Один раз проверяется `Repository → Actions` — что GitHub Actions разрешены.

После публикации через интерфейс настраивается `Packages → Package settings → Change visibility` — здесь выбирается публичный или приватный доступ к образу.

## Что не нужно делать вручную

Для базовой публикации в GHCR не нужно:

- устанавливать отдельный registry;
- поднимать Harbor;
- создавать пароль для GitHub Actions;
- добавлять `GITHUB_TOKEN` вручную в Secrets;
- собирать и отправлять образ со своего компьютера.

---

# 9. Финальный результат

После настройки механизм работает так:

```text
Pull Request в main
    |
    +-- go vet
    +-- go test -race
    +-- Docker-образ не публикуется

Merge или push в main
    |
    +-- go vet
    +-- go test -race
    +-- docker build
    +-- docker push
    |
    v
ghcr.io/<owner>/<repository>:latest
ghcr.io/<owner>/<repository>:sha-<commit>
```

Для pet-проекта рекомендуется:

```text
Repository: public
GHCR package: public
Runner: ubuntu-latest
Deployment tag: sha-<commit>
```

---

# 10. Официальная документация

- GitHub Actions: https://docs.github.com/actions/about-github-actions/understanding-github-actions
- Workflow syntax: https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions
- Publishing Docker images: https://docs.github.com/actions/guides/publishing-docker-images
- GitHub Container Registry: https://docs.github.com/packages/working-with-a-github-packages-registry/working-with-the-container-registry
- Package visibility and access: https://docs.github.com/packages/learn-github-packages/configuring-a-packages-access-control-and-visibility
- Using `GITHUB_TOKEN`: https://docs.github.com/actions/security-guides/automatic-token-authentication
