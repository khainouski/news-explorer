# План действий: развёртывание News Explorer на DigitalOcean

## 0. Стратегия

Идём по принципу **walking skeleton**: сначала протаскиваем весь пайплайн от кода до кластера на
максимально простом приложении, и только когда вся цепочка работает end-to-end — наращиваем
функциональность самого приложения. Порядок фаз:

```
1. Минимальное Go-приложение                    ✅ сделано
2. GitHub Actions: сборка и публикация образа
3. Terraform: инфраструктура (Droplet + k3s)
4. Ручной деплой, проверка по IP
5. Argo CD: GitOps
6. Полноценное приложение с UI поверх готового пайплайна
```

Провайдер — **DigitalOcean** (Terraform provider `digitalocean/digitalocean`, ресурсы
`digitalocean_droplet` / `digitalocean_firewall` / `digitalocean_ssh_key`, токен
`DIGITALOCEAN_TOKEN`). `infrastructure-overview.md`, `terraform-guide.md` и `deploy-to-k3s.md` в
этой же папке — исходные обзорные доки/туториалы; репозиторий там ещё называется `developer-news`
(старое имя проекта) и Helm-chart описан как отдельный репозиторий — это **устарело**. Этот
документ (`plan.md`) — текущий источник истины: актуальное имя проекта — **News Explorer**, и
структура репозиториев изменилась (см. ниже).

---

## 1. Структура репозиториев (актуальная)

Три репозитория, разделены не по фазам пайплайна, а по **типу ответственности** — так это
масштабируется на несколько микросервисов без изменения структуры:

### 1.1. `news-explorer`

Исходный код приложения. **Helm chart живёт здесь же**, не в отдельном репозитории — он описывает,
как именно *этот* сервис запускается в Kubernetes, и меняется в том же Pull Request, что и код.

```text
news-explorer/
├── cmd/
├── internal/
├── web/
├── Dockerfile
├── .github/
│   └── workflows/
└── deploy/
    └── helm/
        ├── Chart.yaml
        ├── values.yaml
        └── templates/
```

Ответственность: Go-приложение, Web UI, Dockerfile, GitHub Actions, Helm chart.

### 1.2. `news-explorer-infra`

Infrastructure as Code.

```text
news-explorer-infra/
├── environments/
│   └── production/
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── terraform.tfvars.example
├── modules/
│   └── droplet/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── cloud-init.yaml
├── .gitignore
└── README.md
```

Ответственность: Terraform (Droplet, Firewall, Network, DNS), bootstrap Kubernetes через
cloud-init.

> Пока окружение одно (`production`), `environments/production/` — это единственный root module
> и по сути прямой перенос файлов из `terraform-guide.md` (`versions.tf`, `providers.tf`, `variables.tf`,
> `locals.tf`, `ssh.tf`, `firewall.tf`, `server.tf`, `outputs.tf`), просто переехавших под
> `environments/production/`. `modules/` пока пустой — общий код туда выносим в тот момент, когда
> появится второе окружение (`staging`) и станет что переиспользовать. Отдельной плоской папки
> `terraform/` не создаю — `environments/production/` уже и есть корневой Terraform-модуль,
> дублировать её ещё одной директорией нет смысла, пока модуль один.

### 1.3. `news-explorer-gitops`

Желаемое состояние кластера. **Chart здесь не хранится** — только ссылки на него (repo+path в
`news-explorer`) и values для каждого окружения.

```text
news-explorer-gitops/
├── applications/
│   └── news-explorer.yaml
└── environments/
    └── production/
        └── news-explorer-values.yaml
```

Готово к росту: как только появится второй сервис или окружение, структура расширяется
предсказуемо —

```text
news-explorer-gitops/
├── applications/
│   ├── news-explorer.yaml
│   └── notification-service.yaml
└── environments/
    ├── staging/
    │   ├── news-explorer-values.yaml
    │   └── notification-service-values.yaml
    └── production/
        ├── news-explorer-values.yaml
        └── notification-service-values.yaml
```

Ответственность: Argo CD `Application`-манифесты, per-environment Helm values, конфигурация
кластера.

### Итоговое разделение

- **`news-explorer`** → код + Dockerfile + Helm chart (**что** и **как** запускать)
- **`news-explorer-infra`** → Terraform (**где** это работает)
- **`news-explorer-gitops`** → Argo CD + per-environment values (**какая версия** и **в каком
  окружении** сейчас развёрнута)

Каждый следующий микросервис (`notification-service`, `user-service`, ...) владеет своим Helm
chart-ом в своём репозитории — `news-explorer-gitops` только перечисляет, какие сервисы
развёрнуты, какая версия/чарт используется и с какими values по каждому окружению. Сам он не
содержит бизнес-логики шаблонов Kubernetes-ресурсов.

---

## Фаза 1 — Минимальное Go-приложение ✅

Сделано: `news-explorer/cmd/app` + `pkg/logger` (`chi` router, `zerolog`), `GET /`, `GET /live`,
`GET /ready`, Dockerfile (multi-stage `golang:1.25.6-alpine` → `alpine:3.22`, non-root),
`golangci-lint` чист. Go 1.25.6.

---

## Фаза 2 — GitHub Actions: сборка и публикация образа

Цель фазы: `push` в `main` сам собирает и публикует образ в GHCR, без ручного `docker push`.

- [ ] `.github/workflows/ci.yml` в `news-explorer`:
  - job `test`: `actions/checkout`, `actions/setup-go` (1.24), `go vet ./...`,
    `go test -race ./...`;
  - job `build-and-push` (только на `push` в `main`, `needs: test`): логин в GHCR через встроенный
    `GITHUB_TOKEN` (права `contents: read`, `packages: write` — отдельный PAT не нужен),
    `docker/build-push-action` с тегами `latest` и `sha-<commit>`.
- [ ] `concurrency` группа по ветке, чтобы старые запуски не мешали новым.
- [ ] Push → проверить `Repository → Actions`, дождаться зелёной сборки.
- [ ] Сделать package в GHCR публичным (`Packages → Package settings → Change visibility →
  Public`) — проще для pet-проекта, не нужен `imagePullSecret` при деплое.

На выходе фазы: `ghcr.io/khainouski/news-explorer:sha-<commit>` доступен `docker pull` без
авторизации.

---

## Фаза 3 — Terraform: инфраструктура

Репозиторий `news-explorer-infra/`. Не зависит от фаз 1–2, можно делать параллельно.

- [ ] Аккаунт DigitalOcean, API Token (`API → Tokens/Keys`), SSH-ключ в аккаунте.
  ```bash
  export DIGITALOCEAN_TOKEN="dop_v1_xxxxxxxx"
  ```
- [ ] `environments/production/variables.tf` (`do_token`, `region`, `droplet_size`).
- [ ] `environments/production/main.tf`:
  ```hcl
  provider "digitalocean" {
    token = var.do_token
  }

  resource "digitalocean_ssh_key" "default" {
    name       = "news-explorer"
    public_key = file("~/.ssh/id_ed25519.pub")
  }

  resource "digitalocean_droplet" "app" {
    name      = "news-explorer"
    image     = "ubuntu-24-04-x64"
    region    = var.region          # например fra1
    size      = "s-2vcpu-4gb"
    ssh_keys  = [digitalocean_ssh_key.default.id]
    user_data = file("${path.module}/../../cloud-init.yaml")
  }

  resource "digitalocean_firewall" "app" {
    name        = "news-explorer-fw"
    droplet_ids = [digitalocean_droplet.app.id]

    inbound_rule { protocol = "tcp" port_range = "22"  source_addresses = ["0.0.0.0/0", "::/0"] }
    inbound_rule { protocol = "tcp" port_range = "80"  source_addresses = ["0.0.0.0/0", "::/0"] }
    inbound_rule { protocol = "tcp" port_range = "443" source_addresses = ["0.0.0.0/0", "::/0"] }
    outbound_rule { protocol = "tcp" port_range = "1-65535" destination_addresses = ["0.0.0.0/0", "::/0"] }
  }
  ```
- [ ] `environments/production/outputs.tf` — публичный IP droplet-а.
- [ ] `.gitignore`: `*.tfstate*`, `terraform.tfvars`, `.env`. Коммитить только
  `terraform.tfvars.example`.
- [ ] `cloud-init.yaml` (в корне репозитория, передаётся как `user_data`): установка k3s
  (`curl -sfL https://get.k3s.io | sh -`) и Helm. Приложение здесь **не** разворачивается.
- [ ] Применить:
  ```bash
  cd environments/production
  terraform init
  terraform plan
  terraform apply
  ```
- [ ] Проверить кластер с Mac:
  ```bash
  scp root@DROPLET_IP:/etc/rancher/k3s/k3s.yaml ~/.kube/news-explorer.yaml
  # заменить server: https://127.0.0.1:6443 на https://DROPLET_IP:6443,
  # либо держать SSH-туннель: ssh -L 6443:127.0.0.1:6443 root@DROPLET_IP
  export KUBECONFIG=~/.kube/news-explorer.yaml
  kubectl get nodes
  ```

На выходе фазы: живой k3s-кластер на DigitalOcean, доступный с Mac через `kubectl`, но без
приложения внутри.

---

## Фаза 4 — Ручной деплой, проверка по IP

Цель фазы: образ из фазы 2 крутится в кластере из фазы 3, и `curl http://DROPLET_IP` отвечает —
без Argo CD, без DNS, минимальными средствами.

- [ ] Самый простой путь — голые манифесты (`kubectl apply -f`), без Helm, чтобы не тратить время
  на шаблонизацию раньше времени:
  - `Deployment` с `image: ghcr.io/khainouski/news-explorer:sha-<commit>`, портом 8080,
    `readinessProbe`/`livenessProbe` на `/ready` и `/live`;
  - `Service` — **важно для проверки по голому IP без DNS**: либо `type: NodePort`, либо
    `ClusterIP` + `Ingress` **без указания `host`** (Traefik, встроенный в k3s, уже слушает 80/443
    на дроплете — host-based Ingress с конкретным доменом тут не сработает, т.к. домена ещё нет).
  ```bash
  kubectl create namespace news-explorer
  kubectl apply -n news-explorer -f deployment.yaml -f service.yaml
  kubectl get pods,svc -n news-explorer
  ```
- [ ] Проверка:
  ```bash
  curl http://DROPLET_IP            # если Ingress без host / NodePort на 80
  # или временно:
  kubectl port-forward -n news-explorer service/news-explorer 8080:80
  curl http://localhost:8080
  ```
- [ ] Когда голый деплой подтверждён — оформить то же самое как **Helm chart прямо в
  `news-explorer/deploy/helm/`** (`Chart.yaml`, `values.yaml`, `templates/deployment.yaml`,
  `templates/service.yaml`, `templates/ingress.yaml`) и закоммитить **в том же репозитории и той
  же PR**, что и код приложения — chart живёт с кодом, не в отдельном репозитории.
  ```bash
  helm lint deploy/helm
  helm upgrade --install news-explorer deploy/helm \
    --namespace news-explorer --create-namespace
  ```
  `values.yaml` в чарте задаёт безопасный дефолт (`image.tag` можно оставить как
  `"latest"`/placeholder) — конкретную версию для конкретного окружения переопределяет отдельный
  values-файл из `news-explorer-gitops` (см. Фазу 5).

На выходе фазы: приложение реально работает на кластере и доступно по IP; в `news-explorer`
лежит Helm chart, соответствующий этому рабочему состоянию.

---

## Фаза 5 — Argo CD: GitOps

Цель фазы: дальнейшие изменения версии/конфигурации применяются к кластеру автоматически, без
ручного `helm upgrade`.

- [ ] Установка:
  ```bash
  kubectl create namespace argocd
  kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
  kubectl port-forward service/argocd-server -n argocd 8081:443
  kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode
  ```
- [ ] `news-explorer-gitops/environments/production/news-explorer-values.yaml` — конкретная
  версия и конфигурация для production:
  ```yaml
  image:
    tag: "sha-<commit>"

  ingress:
    host: app.example.com     # позже, когда появится DNS
  ```
- [ ] `news-explorer-gitops/applications/news-explorer.yaml` — Argo CD `Application` с **двумя
  источниками** (chart из `news-explorer`, values из `news-explorer-gitops`):
  ```yaml
  apiVersion: argoproj.io/v1alpha1
  kind: Application

  metadata:
    name: news-explorer
    namespace: argocd

  spec:
    project: default

    sources:
      - repoURL: https://github.com/khainouski/news-explorer.git
        targetRevision: main
        path: deploy/helm
        helm:
          valueFiles:
            - $values/environments/production/news-explorer-values.yaml

      - repoURL: https://github.com/khainouski/news-explorer-gitops.git
        targetRevision: main
        ref: values

    destination:
      server: https://kubernetes.default.svc
      namespace: news-explorer

    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
        - CreateNamespace=true
  ```
  ```bash
  kubectl apply -f applications/news-explorer.yaml
  ```
- [ ] Проверка GitOps-цикла: поменять `image.tag` в
  `news-explorer-gitops/environments/production/news-explorer-values.yaml` на другой существующий
  `sha-<commit>` из GHCR → закоммитить → убедиться, что Argo CD сам подхватил и обновил
  Deployment (`kubectl rollout status deployment/news-explorer -n news-explorer`).
- [ ] `latest` в качестве тега **не использовать** в values окружения — только `sha-<commit>`,
  иначе Argo CD и Kubernetes не увидят изменение тега.
- [ ] Откат делать через git-revert коммита в `news-explorer-gitops` (в values-файле окружения),
  а не через `kubectl rollout undo` — так остаётся GitOps-консистентность, а сам chart в
  `news-explorer` не трогается.

На выходе фазы: полный GitOps-контур работает — push в `news-explorer` (фаза 2) собирает образ,
после обновления тега в values-файле `news-explorer-gitops` Argo CD раскатывает новую версию сам.

---

## Фаза 6 — Полноценное приложение с UI

Теперь, когда пайплайн полностью проверен, можно спокойно наращивать сам продукт поверх готовой
инфраструктуры, ничего в CI/CD/деплое принципиально не меняя:

- [ ] Бизнес-логика, роуты, шаблоны/статика для UI (`web/` в `news-explorer`).
- [ ] PostgreSQL — либо в самом k3s, либо managed DB от DigitalOcean
  (`digitalocean_database_cluster` в `news-explorer-infra`).
- [ ] Секреты (`DATABASE_URL` и т.п.) — через Kubernetes `Secret`, не в `ConfigMap` и не в коде:
  ```bash
  kubectl create secret generic news-explorer \
    --namespace news-explorer \
    --from-literal=DATABASE_URL="$DATABASE_URL"
  ```
- [ ] DNS: A-запись `app.example.com → DROPLET_IP`, обновить `ingress.host` в values-файле
  `news-explorer-gitops`.
- [ ] HTTPS через cert-manager + Let's Encrypt.
- [ ] По желанию — автоматизировать обновление `image.tag` в
  `news-explorer-gitops/environments/production/news-explorer-values.yaml` отдельным Action/ботом
  после успешной сборки, вместо ручного коммита (убирает последний ручной шаг в GitOps-цикле).
- [ ] Если появится второй сервис (например `notification-service`) — он получает свой `deploy/helm`
  в своём репозитории и свою запись в `news-explorer-gitops/applications/` + values по каждому
  окружению; структура `news-explorer-gitops` уже это предполагает (см. §1.3).

---

## Безопасность — сквозной чеклист по всем фазам

- [ ] `DIGITALOCEAN_TOKEN`, `GHCR_TOKEN`, `DATABASE_URL` — только через переменные окружения /
  `TF_VAR_*` / Kubernetes `Secret`, никогда в коде или коммитах.
- [ ] `.gitignore` во всех трёх репозиториях покрывает `*.tfstate*`, `.env`, `terraform.tfvars`.
- [ ] Firewall: inbound только 22/80/443, всё остальное закрыто (PostgreSQL/Redis наружу не
  торчат).
- [ ] Образ в values окружения (`news-explorer-gitops`) всегда указывается конкретным тегом
  (`sha-<commit>`), не `latest`.
