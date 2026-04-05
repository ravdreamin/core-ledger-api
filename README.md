# Core Ledger API

A robust, performant Go-based FinTech API leveraging `net/http`, `pgxpool`, and PASETO tokens. This project models a financial dashboard backing for managing expenses, income, category aggregation, and rigorous administrative role controls.

## Requirements
* Go 1.21+
* Docker & Docker Compose
* PostgreSQL (or rely exclusively on the `docker-compose.yml` provided)
* [golang-migrate](https://github.com/golang-migrate/migrate)

## Setup and Run

1. **Environment Variables**:
   Copy the example config to configure secrets depending on your environment.
   ```bash
   cp .env-example .env
   ```
   *Make sure `PASETO_KEY` is literally 32 characters long!*

2. **Bootstrapping the Environment**:
   A `Makefile` is configured to spin up your dockerized postgres instance, apply all migrations, and seed out admin records.
   ```bash
   make setup
   ```
   *Note: This commands runs `docker-compose up -d`, then applies migrations dynamically.*

3. **Running the API**:
   ```bash
   make run
   ```
   The API evaluates to `http://localhost:8080`.

## Endpoints and Rules

All secure endpoints utilize PASETO based Bearer `Authorization: Bearer <token>`. Users are evaluated natively via their hashed `ROLE`.

| Endpoint | Method | Role(s) Allowed | Description |
|----------|--------|-----------------|-------------|
| `/login` | `POST` | *Public* | Validates provided credentials. Returns a PASETO token resolving to specific roles. |
| `/dashboard/summary` | `GET` | `admin`, `analyst`, `viewer` | Returns `total_income`, `total_expense`, and `balance` computations dynamically calculated against your dataset. Logs 0 for unrecorded aggregates. |
| `/dashboard/categories` | `GET` | `admin`, `analyst` | Aggregates all transactions by their parent `Categories` to emit arrays of `category_name` and `total_amount`. Fully structured logic natively drops empty schemas as explicit `[]` arrays. |
| `/records` | `POST` | `admin` | Inserts a new financial record mapping to user IDs and category references respectively. |

## Seeded Data & Schema Internals
Running `make setup` triggers native seeds against:
- 5 default categorical tags: `Rent`, `Food`, `Salary`, `Utilities`, `Entertainment`.
- 3 test users (`admin`, `analyst`, `viewer`), all secured using identical standard bcrypt password hash representing identically mapping out to: **`gk123456`**.
- Native dummy data automatically cascades securely bridging testing limits immediately on project init ensuring API queries aren't returning barren `0`/null references initially.