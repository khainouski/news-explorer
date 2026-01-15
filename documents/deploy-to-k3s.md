# Деплой Go-приложения в k3s после создания инфраструктуры Terraform

## Цель

После выполнения Terraform у нас уже есть:

```text
DigitalOcean Droplet
└── Ubuntu
    └── k3s
        └── Kubernetes-кластер
```

Теперь необходимо:

1. Собрать Go-приложение в Docker image.
2. Опубликовать image в Container Registry.
3. Описать приложение с помощью Helm.
4. Проверить ручной деплой.
5. Установить Argo CD.
6. Настроить GitOps-деплой.
7. Настроить GitHub Actions.

Финальная цепочка:

```text
Go Application
      ↓
Docker Image
      ↓
GitHub Container Registry
      ↓
Helm Chart
      ↓
Argo CD
      ↓
k3s
      ↓
Go Application + UI
```

---

# 1. Подготовить Dockerfile

В репозитории Go-приложения должен находиться `Dockerfile`. Пример:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /app/bin/developer-news ./cmd/app

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bin/developer-news ./developer-news

USER app

EXPOSE 8080

ENTRYPOINT ["./developer-news"]
```

Для первой версии проекта проще использовать один Docker image, в котором Go отдаёт и API, и UI.

Проверка локально:

```bash
docker build -t developer-news:local .
docker run --rm -p 8080:8080 developer-news:local
```

Открыть: `http://localhost:8080`

---

# 2. Добавить health endpoints

Желательно реализовать:

```text
GET /health/live
GET /health/ready
```

Пример:

```go
mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
})

mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ready"))
})
```

---

# 3. Опубликовать image в GHCR

Название image: `ghcr.io/YOUR_GITHUB_USERNAME/developer-news:1.0.0`

Авторизация:

```bash
export GHCR_TOKEN="your-token"

echo "$GHCR_TOKEN" | docker login ghcr.io \
  -u YOUR_GITHUB_USERNAME \
  --password-stdin
```

Сборка:

```bash
docker build \
  -t ghcr.io/YOUR_GITHUB_USERNAME/developer-news:1.0.0 \
  .
```

Публикация:

```bash
docker push \
  ghcr.io/YOUR_GITHUB_USERNAME/developer-news:1.0.0
```

Для первой версии проще сделать image публичным.

---

# 4. Получить доступ к k3s

На сервере kubeconfig находится здесь: `/etc/rancher/k3s/k3s.yaml`

Проверка на сервере:

```bash
ssh root@SERVER_IP
sudo k3s kubectl get nodes
sudo k3s kubectl get pods --all-namespaces
```

Безопасный доступ с Mac через SSH tunnel:

```bash
ssh -L 6443:127.0.0.1:6443 root@SERVER_IP
```

Скопировать kubeconfig:

```bash
mkdir -p ~/.kube
scp root@SERVER_IP:/etc/rancher/k3s/k3s.yaml \
  ~/.kube/developer-news.yaml
```

Оставить в kubeconfig:

```yaml
server: https://127.0.0.1:6443
```

В другом терминале:

```bash
export KUBECONFIG=~/.kube/developer-news.yaml
kubectl get nodes
```

---

# 5. Установить kubectl и Helm на Mac

```bash
brew install kubectl
brew install helm
```

Проверка:

```bash
kubectl version --client
helm version
```

---

# 6. Создать deployment-репозиторий

Рекомендуемая структура:

```text
developer-news
└── Go application

developer-news-infrastructure
└── Terraform и cloud-init

developer-news-deployments
└── Helm и Argo CD
```

Структура deployment-репозитория:

```text
developer-news-deployments/
├── charts/
│   └── developer-news/
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── deployment.yaml
│           ├── service.yaml
│           └── ingress.yaml
└── argocd/
    └── developer-news.yaml
```

---

# 7. Создать Helm chart

## `Chart.yaml`

```yaml
apiVersion: v2
name: developer-news
description: Helm chart for Developer News Platform
type: application
version: 0.1.0
appVersion: "1.0.0"
```

## `values.yaml`

```yaml
replicaCount: 1

image:
  repository: ghcr.io/YOUR_GITHUB_USERNAME/developer-news
  tag: "1.0.0"
  pullPolicy: IfNotPresent

containerPort: 8080

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: traefik
  host: app.example.com

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi

env:
  APP_ENV: production
  HTTP_PORT: "8080"
```

## `templates/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}
  labels:
    app.kubernetes.io/name: {{ .Chart.Name }}
    app.kubernetes.io/instance: {{ .Release.Name }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app.kubernetes.io/name: {{ .Chart.Name }}
      app.kubernetes.io/instance: {{ .Release.Name }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: {{ .Chart.Name }}
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: {{ .Values.containerPort }}
          env:
            - name: APP_ENV
              value: {{ .Values.env.APP_ENV | quote }}
            - name: HTTP_PORT
              value: {{ .Values.env.HTTP_PORT | quote }}
          readinessProbe:
            httpGet:
              path: /health/ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health/live
              port: http
            initialDelaySeconds: 15
            periodSeconds: 15
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

## `templates/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}
spec:
  type: {{ .Values.service.type }}
  selector:
    app.kubernetes.io/name: {{ .Chart.Name }}
    app.kubernetes.io/instance: {{ .Release.Name }}
  ports:
    - name: http
      port: {{ .Values.service.port }}
      targetPort: http
```

## `templates/ingress.yaml`

```yaml
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}
spec:
  ingressClassName: {{ .Values.ingress.className }}
  rules:
    - host: {{ .Values.ingress.host }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {{ .Release.Name }}
                port:
                  number: {{ .Values.service.port }}
{{- end }}
```

---

# 8. Проверить Helm chart

```bash
cd developer-news-deployments
helm lint charts/developer-news
helm template developer-news charts/developer-news
```

---

# 9. Первый ручной деплой

```bash
helm upgrade --install developer-news \
  charts/developer-news \
  --namespace developer-news \
  --create-namespace
```

Проверка:

```bash
kubectl get pods -n developer-news
kubectl get deployment -n developer-news
kubectl get service -n developer-news
kubectl get ingress -n developer-news
```

Логи:

```bash
kubectl logs -n developer-news deployment/developer-news --follow
```

---

# 10. Проверить приложение через port-forward

```bash
kubectl port-forward -n developer-news service/developer-news 8080:80
```

Открыть: `http://localhost:8080`

---

# 11. Настроить DNS

Создать A-запись: `app.example.com → SERVER_IP`

После обновления DNS приложение будет доступно: `http://app.example.com`

Позже можно добавить HTTPS через cert-manager и Let's Encrypt.

---

# 12. Установить Argo CD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Проверка:

```bash
kubectl get pods -n argocd
```

Доступ к UI:

```bash
kubectl port-forward service/argocd-server -n argocd 8081:443
```

Открыть: `https://localhost:8081`

Получить пароль:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode
echo
```

---

# 13. Создать Argo CD Application

Файл: `argocd/developer-news.yaml`

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: developer-news
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/YOUR_USERNAME/developer-news-deployments.git
    targetRevision: main
    path: charts/developer-news
    helm:
      valueFiles:
        - values.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: developer-news
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

Применить:

```bash
kubectl apply -f argocd/developer-news.yaml
```

---

# 14. Настроить GitHub Actions

Файл: `.github/workflows/build.yml`

```yaml
name: Build and publish image

on:
  push:
    branches:
      - main

permissions:
  contents: read
  packages: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.24"

      - name: Run tests
        run: go test ./...

      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ghcr.io/YOUR_USERNAME/developer-news:${{ github.sha }}
            ghcr.io/YOUR_USERNAME/developer-news:latest
```

---

# 15. Обновление image tag

После сборки нового image нужно изменить в deployment-репозитории:

```yaml
image:
  repository: ghcr.io/YOUR_USERNAME/developer-news
  tag: "NEW_GIT_SHA"
```

После commit Argo CD выполнит деплой новой версии.

Полный процесс:

```text
Push Go code
      ↓
GitHub Actions
      ↓
go test
      ↓
docker build
      ↓
docker push
      ↓
Update Helm image tag
      ↓
Commit deployment repository
      ↓
Argo CD sync
      ↓
Kubernetes rollout
```

---

# 16. Проверка rollout

```bash
kubectl rollout status deployment/developer-news -n developer-news
```

История:

```bash
kubectl rollout history deployment/developer-news -n developer-news
```

Откат:

```bash
kubectl rollout undo deployment/developer-news -n developer-news
```

В GitOps-подходе правильнее откатывать commit в deployment-репозитории.

---

# 17. Если GHCR image приватный

```bash
kubectl create secret docker-registry ghcr-secret \
  --namespace developer-news \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password="$GHCR_TOKEN"
```

Добавить в Deployment:

```yaml
spec:
  imagePullSecrets:
    - name: ghcr-secret
```

---

# 18. PostgreSQL и секреты

Для pet-проекта PostgreSQL можно запустить внутри k3s либо использовать managed PostgreSQL.

Несекретные настройки хранятся в `ConfigMap`, секретные — в `Secret`.

Создание секрета вручную:

```bash
kubectl create secret generic developer-news \
  --namespace developer-news \
  --from-literal=DATABASE_URL="$DATABASE_URL"
```

Реальные пароли нельзя коммитить в Git.

---

# 19. Минимальный рабочий путь

Сначала достаточно выполнить:

```bash
docker build -t ghcr.io/YOUR_USERNAME/developer-news:1.0.0 .
docker push ghcr.io/YOUR_USERNAME/developer-news:1.0.0

helm lint charts/developer-news
helm upgrade --install developer-news \
  charts/developer-news \
  --namespace developer-news \
  --create-namespace

kubectl get pods -n developer-news

kubectl port-forward -n developer-news service/developer-news 8080:80
```

После этого приложение должно открываться: `http://localhost:8080`

---

# 20. Рекомендуемый порядок реализации

```text
Шаг 1. Docker image работает локально
Шаг 2. Image опубликован в GHCR
Шаг 3. Helm chart проходит helm lint
Шаг 4. Ручной helm install работает
Шаг 5. Приложение открывается через port-forward
Шаг 6. Настроены DNS и Ingress
Шаг 7. Установлен Argo CD
Шаг 8. Argo CD управляет приложением
Шаг 9. GitHub Actions собирает image
Шаг 10. CI обновляет image tag
```

---

# 21. Финальная архитектура

```text
                      GitHub
        developer-news repository
                    │
                    │ push
                    ▼
             GitHub Actions
                    │
                    ├── go test
                    ├── docker build
                    └── docker push
                            │
                            ▼
                           GHCR
                            │
                            ▼
       developer-news-deployments
                    │
                    │ image tag
                    ▼
                  Argo CD
                    │
                    ▼
             Helm templates
                    │
                    ▼
                    k3s
                    │
                    ├── Deployment
                    ├── Service
                    └── Ingress
                            │
                            ▼
               Go Application + UI
```

---

# 22. Разделение ответственности

| Инструмент | Ответственность |
|---|---|
| Terraform | Создание VM, firewall, SSH и инфраструктуры |
| cloud-init | Установка и запуск k3s |
| Docker | Упаковка Go-приложения |
| GHCR | Хранение Docker images |
| Helm | Описание Kubernetes-ресурсов приложения |
| Argo CD | Синхронизация Kubernetes с Git |
| GitHub Actions | Тестирование, сборка и публикация image |
| k3s | Запуск приложения |
| Traefik | Входящий HTTP/HTTPS-трафик |
| DNS | Привязка домена к IP сервера |

---

# Итог

Terraform уже подготовил инфраструктуру: `VM + Firewall + SSH + k3s`

Дальше приложение разворачивается не Terraform, а через:

```text
Docker
   ↓
GHCR
   ↓
Helm
   ↓
Argo CD
   ↓
k3s
```

Сначала следует выполнить ручной деплой через Helm. После проверки приложения управление деплоем можно передать Argo CD и подключить GitHub Actions.
