# Terraform: теория и практическое применение

## 1. Что такое Terraform

**Terraform** — это инструмент для управления инфраструктурой с помощью кода, то есть подход **Infrastructure as Code — IaC**.

Вместо ручного создания сервера через интерфейс DigitalOcean можно описать необходимую инфраструктуру в `.tf`-файлах:

```hcl
resource "digitalocean_droplet" "app" {
  name   = "developer-news"
  image  = "ubuntu-24-04-x64"
  size   = "s-2vcpu-4gb"
  region = "fra1"
}
```

После выполнения:

```bash
terraform apply
```

Terraform подключится к API DigitalOcean и создаст виртуальную машину (Droplet).

---

## 2. Для чего применяется Terraform

Terraform позволяет управлять:

- виртуальными машинами;
- сетями;
- firewall;
- SSH-ключами;
- дисками;
- load balancer;
- DNS-записями;
- managed databases;
- Kubernetes-кластерами;
- object storage;
- другими облачными ресурсами.

Для проекта **Developer News Platform** Terraform будет отвечать за:

```text
Terraform
├── DigitalOcean Droplet
├── SSH Key
├── Firewall
└── Передачу cloud-init
        └── Установка k3s
```

Terraform не будет разворачивать само Go-приложение.

Разделение ответственности:

```text
Terraform   → облачная инфраструктура
cloud-init  → начальная настройка Linux-сервера
k3s         → Kubernetes-кластер
Helm        → описание Kubernetes-приложения
Argo CD     → GitOps-развёртывание приложения
```

---

## 3. Как работает Terraform

Terraform работает декларативно.

Мы описываем не последовательность команд, а желаемое состояние:

```text
Мне нужен сервер Ubuntu
с двумя CPU,
SSH-ключом
и открытыми портами 22, 80 и 443.
```

Terraform сравнивает:

```text
Terraform configuration
          ↓
Terraform state
          ↓
Реальная инфраструктура
          ↓
План изменений
```

После этого Terraform решает:

- какие ресурсы нужно создать;
- какие ресурсы изменить;
- какие ресурсы удалить;
- какие ресурсы оставить без изменений.

---

## 4. Основные команды Terraform

### Инициализация проекта

```bash
terraform init
```

Команда:

- инициализирует Terraform-проект;
- скачивает providers;
- настраивает backend;
- создаёт файл блокировки версий providers.

### Форматирование кода

```bash
terraform fmt
```

Для проверки без изменения файлов:

```bash
terraform fmt -check
```

### Проверка конфигурации

```bash
terraform validate
```

### Формирование плана

```bash
terraform plan
```

Обозначения:

```text
+ create
~ update
- destroy
```

### Применение изменений

```bash
terraform apply
```

Для применения заранее сохранённого плана:

```bash
terraform plan -out=tfplan
terraform apply tfplan
```

### Просмотр outputs

```bash
terraform output
```

### Удаление инфраструктуры

```bash
terraform destroy
```

---

## 5. Основные объекты Terraform

### 5.1. `terraform`

Блок `terraform` задаёт:

- минимальную версию Terraform;
- используемые providers;
- ограничения версий providers;
- backend для хранения state.

```hcl
terraform {
  required_version = ">= 1.8.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}
```

> Версию провайдера перед стартом проекта стоит свериться с актуальной на
> `registry.terraform.io/providers/digitalocean/digitalocean` — здесь указано ограничение по
> мажорной ветке, а не конкретный последний релиз.

### 5.2. Provider

**Provider** — это плагин, через который Terraform взаимодействует с внешней системой.

```text
digitalocean → DigitalOcean
aws          → Amazon Web Services
azurerm      → Microsoft Azure
google       → Google Cloud
cloudflare   → Cloudflare
kubernetes   → Kubernetes
```

Подключение DigitalOcean provider:

```hcl
provider "digitalocean" {}
```

Токен передаётся через переменную окружения:

```bash
export DIGITALOCEAN_TOKEN="your-token"
```

### 5.3. Resource

**Resource** — объект инфраструктуры, который Terraform должен создать и управлять им.

```hcl
resource "digitalocean_droplet" "app" {
  name   = "developer-news"
  image  = "ubuntu-24-04-x64"
  size   = "s-2vcpu-4gb"
  region = "fra1"
}
```

Структура:

```hcl
resource "<resource_type>" "<local_name>" {
  # arguments
}
```

Полный адрес ресурса:

```text
digitalocean_droplet.app
```

### 5.4. Variable

**Variable** — входной параметр Terraform-конфигурации.

```hcl
variable "server_name" {
  description = "Name of the application server"
  type        = string
  default     = "developer-news"
}
```

Использование:

```hcl
resource "digitalocean_droplet" "app" {
  name = var.server_name
}
```

### 5.5. Local values

**Locals** — вычисляемые или повторно используемые значения внутри конфигурации.

```hcl
locals {
  common_tags = [
    "developer-news",
    var.environment,
    "managed-by-terraform",
  ]
}
```

### 5.6. Data source

**Data source** позволяет получить информацию о существующем объекте, не создавая его.

```hcl
data "digitalocean_image" "ubuntu" {
  slug = "ubuntu-24-04-x64"
}
```

Разница:

```text
resource → создаёт объект
data     → читает существующий объект
```

### 5.7. Output

**Output** выводит полезное значение после `terraform apply`.

```hcl
output "server_ip" {
  description = "Public IPv4 address of the server"
  value       = digitalocean_droplet.app.ipv4_address
}
```

### 5.8. Module

**Module** — группа связанных Terraform-ресурсов, оформленная как повторно используемый компонент.

```hcl
module "application_server" {
  source = "./modules/server"

  server_name = "developer-news"
  region      = "fra1"
}
```

### 5.9. State

Terraform хранит информацию о созданной инфраструктуре в:

```text
terraform.tfstate
```

State содержит:

- идентификаторы ресурсов;
- атрибуты ресурсов;
- зависимости;
- outputs;
- иногда чувствительные значения.

State нельзя публиковать в Git.

### 5.10. Backend

Backend определяет, где Terraform хранит state.

Локальный вариант:

```text
terraform.tfstate
```

Удалённый вариант (AWS S3; вместо него можно использовать DigitalOcean Spaces — они
S3-совместимы, но потребуют `endpoints.s3` и отдельных ключей доступа Spaces):

```hcl
terraform {
  backend "s3" {
    bucket = "developer-news-terraform-state"
    key    = "production/terraform.tfstate"
    region = "eu-central"
  }
}
```

Для первой версии pet-проекта remote backend необязателен.

---

## 6. Зависимости между ресурсами

Terraform самостоятельно определяет зависимости по ссылкам.

```hcl
resource "digitalocean_droplet" "app" {
  ssh_keys = [
    digitalocean_ssh_key.developer.id
  ]
}
```

Порядок:

```text
Сначала создать SSH Key
          ↓
Затем создать Droplet
```

---

## 7. Структура Terraform-проекта

```text
developer-news-infrastructure/
├── terraform/
│   ├── versions.tf
│   ├── providers.tf
│   ├── variables.tf
│   ├── locals.tf
│   ├── ssh.tf
│   ├── firewall.tf
│   ├── server.tf
│   ├── outputs.tf
│   ├── terraform.tfvars.example
│   └── cloud-init.yaml
├── .gitignore
└── README.md
```

Terraform загружает все `.tf`-файлы в каталоге как одну конфигурацию.

---

## 8. Код Terraform по файлам

### 8.1. `versions.tf`

```hcl
terraform {
  required_version = ">= 1.8.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}
```

### 8.2. `providers.tf`

```hcl
provider "digitalocean" {}
```

Перед запуском:

```bash
export DIGITALOCEAN_TOKEN="your-real-token"
```

Проверка:

```bash
test -n "$DIGITALOCEAN_TOKEN" && echo "DIGITALOCEAN_TOKEN is set"
```

### 8.3. `variables.tf`

```hcl
variable "project_name" {
  description = "Project name used in resource names and tags"
  type        = string
  default     = "developer-news"
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "production"

  validation {
    condition = contains(
      ["development", "staging", "production"],
      var.environment
    )

    error_message = "Environment must be development, staging, or production."
  }
}

variable "server_size" {
  description = "DigitalOcean droplet size slug"
  type        = string
  default     = "s-2vcpu-4gb"
}

variable "server_region" {
  description = "DigitalOcean region slug"
  type        = string
  default     = "fra1"
}

variable "server_image" {
  description = "Operating system image slug"
  type        = string
  default     = "ubuntu-24-04-x64"
}

variable "ssh_public_key_path" {
  description = "Path to the local SSH public key"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "allowed_ssh_cidrs" {
  description = "CIDR ranges allowed to connect to SSH"
  type        = list(string)

  default = ["0.0.0.0/0"]
}
```

Для большей безопасности:

```hcl
allowed_ssh_cidrs = ["YOUR_PUBLIC_IP/32"]
```

### 8.4. `locals.tf`

```hcl
locals {
  server_name = "${var.project_name}-${var.environment}"

  common_tags = [
    var.project_name,
    var.environment,
    "managed-by-terraform",
  ]
}
```

### 8.5. `ssh.tf`

```hcl
resource "digitalocean_ssh_key" "developer" {
  name       = "${local.server_name}-ssh"
  public_key = file(pathexpand(var.ssh_public_key_path))
}
```

Передаётся только публичный ключ:

```text
~/.ssh/id_ed25519.pub
```

### 8.6. `firewall.tf`

```hcl
resource "digitalocean_firewall" "app" {
  name = "${local.server_name}-firewall"
  tags = local.common_tags

  droplet_ids = [
    digitalocean_droplet.app.id
  ]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = var.allowed_ssh_cidrs
  }

  inbound_rule {
    protocol   = "tcp"
    port_range = "80"
    source_addresses = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  inbound_rule {
    protocol   = "tcp"
    port_range = "443"
    source_addresses = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  outbound_rule {
    protocol   = "tcp"
    port_range = "1-65535"
    destination_addresses = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  outbound_rule {
    protocol   = "udp"
    port_range = "1-65535"
    destination_addresses = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  outbound_rule {
    protocol = "icmp"
    destination_addresses = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
}
```

В отличие от некоторых других провайдеров, firewall у DigitalOcean не привязывается к серверу
через атрибут на самом droplet-е — связь идёт в обратную сторону, через `droplet_ids` на ресурсе
`digitalocean_firewall`.

Публично не открываются:

```text
5432 PostgreSQL
6379 Redis
3000 Grafana
9090 Prometheus
6443 Kubernetes API
```

### 8.7. `server.tf`

```hcl
resource "digitalocean_droplet" "app" {
  name   = local.server_name
  image  = var.server_image
  size   = var.server_size
  region = var.server_region

  ipv6       = true
  monitoring = true

  ssh_keys = [
    digitalocean_ssh_key.developer.id
  ]

  user_data = templatefile(
    "${path.module}/cloud-init.yaml",
    {
      project_name = var.project_name
    }
  )

  tags = local.common_tags
}
```

### 8.8. `cloud-init.yaml`

```yaml
#cloud-config

package_update: true
package_upgrade: false

packages:
  - ca-certificates
  - curl
  - git
  - jq

write_files:
  - path: /etc/rancher/k3s/config.yaml
    permissions: "0600"
    owner: root:root
    content: |
      write-kubeconfig-mode: "0640"
      disable:
        - servicelb

runcmd:
  - curl -sfL https://get.k3s.io | sh -
  - systemctl enable k3s
  - systemctl start k3s
  - k3s kubectl wait --for=condition=Ready node --all --timeout=180s
  - echo "${project_name} infrastructure initialized" > /var/log/project-bootstrap.log
```

Проверка:

```bash
sudo k3s kubectl get nodes
```

Отдельный `kubectl` устанавливать не обязательно:

```bash
sudo k3s kubectl
```

Helm лучше устанавливать отдельно на локальной машине или в CI. При необходимости его можно добавить в `cloud-init`:

```yaml
- curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
- helm version
```

### 8.9. `outputs.tf`

```hcl
output "server_id" {
  description = "DigitalOcean droplet ID"
  value       = digitalocean_droplet.app.id
}

output "server_name" {
  description = "Server name"
  value       = digitalocean_droplet.app.name
}

output "server_ipv4" {
  description = "Public IPv4 address"
  value       = digitalocean_droplet.app.ipv4_address
}

output "server_ipv6" {
  description = "Public IPv6 address"
  value       = digitalocean_droplet.app.ipv6_address
}

output "ssh_command" {
  description = "Command for connecting to the server"
  value       = "ssh root@${digitalocean_droplet.app.ipv4_address}"
}

output "application_url" {
  description = "Temporary application URL without DNS"
  value       = "http://${digitalocean_droplet.app.ipv4_address}"
}
```

### 8.10. `terraform.tfvars.example`

```hcl
project_name  = "developer-news"
environment   = "production"
server_size   = "s-2vcpu-4gb"
server_region = "fra1"
server_image  = "ubuntu-24-04-x64"

ssh_public_key_path = "~/.ssh/id_ed25519.pub"

allowed_ssh_cidrs = [
  "YOUR_PUBLIC_IP/32"
]
```

Создать локальный файл:

```bash
cp terraform.tfvars.example terraform.tfvars
```

### 8.11. `.gitignore`

```gitignore
# Terraform state
*.tfstate
*.tfstate.*

# Terraform plans
*.tfplan
tfplan

# Local variables
terraform.tfvars
*.auto.tfvars

# Terraform working directory
.terraform/

# Crash logs
crash.log
crash.*.log

# Environment files
.env
.env.*

# macOS
.DS_Store

# IDE
.idea/
.vscode/
```

Файл `.terraform.lock.hcl` следует сохранить в Git.

---

## 9. Безопасное хранение DigitalOcean API Token

Рекомендуемый вариант:

```bash
export DIGITALOCEAN_TOKEN="your-real-token"
```

После работы:

```bash
unset DIGITALOCEAN_TOKEN
```

Нельзя хранить токен в:

```text
main.tf
providers.tf
terraform.tfvars
README.md
.env.example
CLAUDE.md
GitHub
Claude prompt
```

### Правила для Claude Code

Файл `CLAUDE.md`:

```md
# Security rules

- Never read or print environment variables containing secrets.
- Never read `.env`, `terraform.tfvars`, private SSH keys, or Terraform state.
- Never add credentials or tokens to source files.
- Use `DIGITALOCEAN_TOKEN` from the process environment.
- Before creating a commit, check that no secrets are included.
```

---

## 10. Подготовка macOS

### Установка Terraform

```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

Проверка:

```bash
terraform version
```

### Создание SSH-ключа

```bash
ssh-keygen \
  -t ed25519 \
  -C "developer-news-digitalocean" \
  -f ~/.ssh/id_ed25519
```

Будут созданы:

```text
~/.ssh/id_ed25519      — приватный ключ
~/.ssh/id_ed25519.pub  — публичный ключ
```

---

## 11. Подключение Terraform к DigitalOcean

1. Создать проект в DigitalOcean (`Projects → Create Project`, опционально — можно обойтись
   Default Project).
2. Создать API Token: `API → Tokens/Keys → Generate New Token`, права `Read` и `Write`.
3. Передать токен:

```bash
export DIGITALOCEAN_TOKEN="your-token"
```

4. Инициализировать проект:

```bash
cd developer-news-infrastructure/terraform
terraform init
```

5. Проверить конфигурацию:

```bash
terraform fmt -recursive
terraform validate
```

6. Создать план:

```bash
terraform plan -out=tfplan
```

7. Применить:

```bash
terraform apply tfplan
```

8. Получить IP:

```bash
terraform output server_ipv4
```

9. Подключиться:

```bash
terraform output -raw ssh_command
```

10. Проверить k3s:

```bash
sudo systemctl status k3s
sudo k3s kubectl get nodes
sudo k3s kubectl get pods --all-namespaces
```

---

## 12. Что происходит после `terraform apply`

```text
terraform apply
        │
        ▼
Terraform загружает digitalocean provider
        │
        ▼
Provider обращается к DigitalOcean API
        │
        ├── создаёт SSH Key
        ├── создаёт Firewall
        └── создаёт Droplet
                    │
                    ▼
          Droplet получает cloud-init
                    │
                    ▼
          Ubuntu выполняет bootstrap
                    │
                    ├── устанавливает curl
                    ├── устанавливает k3s
                    └── запускает Kubernetes
```

---

## 13. Следующий этап после Terraform

Terraform заканчивает работу на уровне:

```text
Droplet + Firewall + SSH + k3s
```

Следующие шаги:

1. Установить Argo CD через Helm.
2. Создать Helm chart Go-приложения.
3. Создать deployment repository.
4. Подключить repository к Argo CD.
5. Настроить GitHub Actions.
6. Публиковать Docker image в GHCR.
7. Обновлять image tag.
8. Разворачивать приложение через GitOps.

---

## 14. Что реализовано с помощью Terraform

> Provisioned a DigitalOcean Droplet, SSH access, firewall rules, and public networking using Terraform. Automated the initial server bootstrap and installation of a single-node k3s cluster through cloud-init.

Расширенный вариант:

> Implemented Infrastructure as Code with Terraform to provision and manage a DigitalOcean Droplet, SSH keys, firewall rules, and networking. Used cloud-init to automate the initial operating system configuration and bootstrap a single-node k3s Kubernetes cluster.

---

## 15. Итоговое разделение ответственности

| Инструмент | Ответственность |
|---|---|
| Terraform | Droplet, firewall, SSH key, networking |
| DigitalOcean Provider | Взаимодействие с DigitalOcean API |
| Terraform State | Хранение состояния инфраструктуры |
| cloud-init | Первичная настройка Ubuntu |
| k3s | Single-node Kubernetes |
| Helm | Упаковка Kubernetes-приложения |
| Argo CD | GitOps-синхронизация |
| GitHub Actions | Тестирование, сборка и публикация image |
| GHCR | Хранение Docker images |
| Kubernetes | Запуск Go-приложения |

---

## 16. Финальная архитектура

```text
Developer Mac
    │
    │ terraform apply
    ▼
DigitalOcean API
    │
    ├── SSH Key
    ├── Firewall
    └── Ubuntu Droplet
            │
            │ cloud-init
            ▼
           k3s
            │
            ▼
          Argo CD
            │
            ▼
       Helm Release
            │
            ▼
  Go Application + Web UI
```

Terraform применяется по назначению: он создаёт и поддерживает облачную инфраструктуру, а установка k3s выполняется операционной системой через переданный Terraform механизм `cloud-init`.
