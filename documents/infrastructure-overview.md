# Developer News Platform — Infrastructure & Deployment Plan

## Goal

Build a production-like deployment pipeline for a Go web application that demonstrates:

- Infrastructure as Code (Terraform)
- Kubernetes (k3s)
- Helm
- GitHub Actions
- Argo CD (GitOps)
- Containerization with Docker
- Basic cloud infrastructure management

The goal is not only to deploy a Go application, but to demonstrate the complete cloud-native delivery workflow expected from a modern backend engineer.

---

# High Level Architecture

```text
                     GitHub
          +--------------------------+
          | developer-news-app       |
          | Go application           |
          +------------+-------------+
                       |
                       | GitHub Actions
                       |
                       v
              Docker Image (GHCR)
                       |
                       v
        +-----------------------------+
        | deployment-manifests repo   |
        | Helm values / image tag     |
        +-------------+---------------+
                      |
                      v
                 Argo CD watches Git
                      |
                      v
               Kubernetes (k3s)
                      |
               Deployment
                      |
                 Go Application
```

Infrastructure:

```text
Terraform
      |
      +-- Virtual Machine
      +-- Firewall
      +-- SSH Key
      +-- (Optional) DNS
      |
      +-- cloud-init
              |
              +-- install k3s
              +-- install Helm
```

---

# Technology Stack

## Infrastructure

- Terraform
- DigitalOcean
- Ubuntu 24.04
- cloud-init

## Containerization

- Docker

## Kubernetes

- k3s (single-node cluster)

## GitOps

- Argo CD

## Package Manager

- Helm

## CI/CD

- GitHub Actions
- GitHub Container Registry (GHCR)

## Application

- Go
- PostgreSQL

---

# Why DigitalOcean

Reasons:

- inexpensive
- official Terraform Provider
- enough resources for pet project
- easy API
- widely used for personal Kubernetes clusters

Recommended droplet:

```text
2 vCPU
4 GB RAM
40 GB SSD
```

---

# Repository Structure

> Historical note: this section describes the **original** three-repo split under the project's
> old name (`developer-news-app` / `developer-news-infrastructure` / `developer-news-deployments`),
> kept here for reference. The **current** structure (project renamed to News Explorer, Helm chart
> moved into the app repo) lives in `plan.md` — read that one for anything you're about to act on.

## 1. Application Repository

```text
developer-news-app/
cmd/
internal/
web/
Dockerfile
.github/workflows/
```

Responsibilities:

- Go source code
- Docker image
- CI pipeline

---

## 2. Infrastructure Repository

```text
developer-news-infrastructure/
terraform/
    versions.tf
    main.tf
    variables.tf
    outputs.tf
cloud-init.yaml
README.md
```

Responsibilities:

- create infrastructure
- bootstrap server

---

## 3. Deployment Repository

```text
developer-news-deployments/
charts/
    developer-news/
environments/
    production/
```

Responsibilities:

- Helm charts
- Kubernetes manifests
- GitOps repository

---

# Terraform Responsibilities

Terraform should manage only infrastructure.

Creates:

- Virtual Machine
- Firewall
- SSH Key
- (Optional) DNS records

Terraform DOES NOT deploy the application.

---

# cloud-init Responsibilities

cloud-init performs initial server bootstrap.

Installs:

- curl
- k3s
- Helm

No application deployment.

---

# Kubernetes Responsibilities

Kubernetes manages application runtime.

Resources:

- Deployment
- Service
- Ingress

---

# Helm Responsibilities

Helm manages application deployment.

Contains:

- Deployment
- Service
- Ingress
- ConfigMap
- Values

---

# Argo CD Responsibilities

Argo CD continuously watches Git.

Workflow:

```text
Git changes
        ↓
Argo CD detects change
        ↓
Sync
        ↓
Deploy new version
```

---

# GitHub Actions Workflow

```text
Push
↓
Run tests
↓
Build Docker image
↓
Push image to GHCR
↓
Update image tag in deployment repository
↓
Commit
↓
Argo CD deploys new version
```

---

# Terraform Workflow

```text
terraform init
↓
terraform plan
↓
terraform apply
↓
VM created
↓
cloud-init runs
↓
k3s installed
↓
Helm installed
```

---

# Secrets

Never store secrets inside Git. Use environment variables.

Example:

```bash
export DIGITALOCEAN_TOKEN=xxxxxxxx
```

Terraform:

```hcl
provider "digitalocean" {
    token = var.do_token
}
```

Use `TF_VAR_do_token` or `DIGITALOCEAN_TOKEN`.

---

# Files NOT committed to Git

```text
terraform.tfvars
.env
terraform.tfstate
terraform.tfstate.backup
```

`.gitignore`:

```gitignore
*.tfstate
*.tfstate.*
terraform.tfvars
.env
```

Commit only:

```text
terraform.tfvars.example
.env.example
```

---

# Firewall Rules

Inbound:

```text
22
80
443
```

Outbound:

```text
Allow All
```

No PostgreSQL, Redis or Grafana ports should be exposed publicly.

---

# DNS

Optional for the first version.

Without DNS: `http://SERVER_IP`

Later: `developer-news.dev` or `app.developer-news.dev`

Initially DNS can be configured manually at the domain registrar. Terraform-managed DNS can be added later.

---

# Deployment Flow

```text
Developer
↓
Push code
↓
GitHub Actions
↓
Docker image
↓
GHCR
↓
Update Helm values
↓
Commit
↓
Argo CD
↓
k3s
↓
Go Application
```

---

# Learning Goals

By the end of the project I should be comfortable with:

- Terraform basics
- Terraform Providers
- Terraform State
- Variables
- Outputs
- cloud-init
- Infrastructure as Code
- Kubernetes fundamentals
- Helm
- Argo CD
- GitOps
- GitHub Actions
- Docker
- Kubernetes Deployments
- Services
- Ingress
- Cloud networking basics
- Firewall configuration

---

# Future Improvements

Possible additions later:

- HTTPS via Let's Encrypt
- External PostgreSQL
- Grafana
- Prometheus
- Loki
- Tempo
- OpenTelemetry
- Horizontal Pod Autoscaler
- Remote Terraform State (S3-compatible storage)
- Multiple environments (dev/stage/prod)

---

# Final Architecture

```text
                    GitHub
        developer-news-app
                 │
                 │ GitHub Actions
                 ▼
              Docker Image
                 │
                 ▼
     deployment-manifests repository
                 │
                 ▼
              Argo CD
                 │
                 ▼
          Kubernetes (k3s)
                 │
         Deployment (1 replica)
                 │
                 ▼
        Go Application + UI

Terraform
    │
    ├── VM
    ├── Firewall
    ├── SSH Key
    └── cloud-init
             │
             ├── install k3s
             └── install Helm
```

## Principle of Responsibility

- **Terraform** → Infrastructure
- **cloud-init** → Bootstrap server
- **Helm** → Kubernetes application
- **Argo CD** → GitOps deployment
- **GitHub Actions** → Build & Release
- **Kubernetes** → Runtime
