# Daily Expense Manager API

A simple API for managing daily expenses using Go, GoFiber, GORM, and PostgreSQL. This API allows users to manage their expenses, track spending, and perform basic CRUD operations.

## Features

- User authentication with JWT
- Expense tracking with categories
- Expense entry and reporting
- Postgres as the database engine
- GORM for ORM

## Requirements

- Go 1.18+
- PostgreSQL
- Docker (for easy database setup)
- Make (optional, for managing migrations and project tasks)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/daily-expense-manager.git
cd daily-expense-manager
```
### 2. Set Up Database Configuration
```
    Host     = "localhost"
    Port     = "5432"
    User     = "youruser"
    Password = "yourpassword"
    Dname   = "expense_manager"
```
### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run Application
```bash
make api
```

The application will be accessible at http://localhost:3000.



