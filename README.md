# Data Transformation Platform

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB.svg)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](docker-compose.yaml)

Enterprise-grade data transformation platform for managing client data mappings and JSON transformations with real-time processing capabilities.

## Overview

Full-stack solution built with Go and React for complex data mapping and transformation scenarios. Provides robust, scalable architecture for enterprise data integration needs.

**Key Features:**
- Multi-client management with isolated mapping rules
- Advanced mapping engine with expressions and validation
- Bulk import/export operations
- Real-time transformation with streaming support (5MB+)
- JWT-based authentication
- Modern React UI with Tailwind CSS

## Technology Stack

- **Backend:** Go 1.24+, PostgreSQL 12+
- **Frontend:** React 18+, Vite, Tailwind CSS
- **Authentication:** JWT
- **Deployment:** Docker

## Quick Start

### Docker Deployment (Recommended)

```bash
git clone https://github.com/Aytaditya/datamorph.git
cd datamorph
docker-compose up --build
```

Access: Frontend at `http://localhost:5173`, Backend at `https://localhost:8080`

### Manual Setup

**Backend:**
```bash
go mod download
export DATABASE_URL="postgres://user:password@localhost:5432/data_mapping"
export JWT_SECRET="your_secure_secret_key"
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Authentication

Default credentials: `admin` / `password`  
⚠️ Change these before production deployment

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/login` | POST | User authentication |
| `/clients` | GET/POST | Client management |
| `/clients/:id/mappings` | GET/POST | Mapping rules |
| `/clients/:id/transform` | POST | Data transformation |
| `/health` | GET | Health check |

## Configuration

### Environment Variables
```bash
SERVER_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/data_mapping
JWT_SECRET=your_secure_secret_key
LOG_LEVEL=info
CERT_FILE_PATH=cert.pem
KEY_FILE_PATH=key.pem
```

### Mapping Example
```json
[
  {
    "sourcePath": "applicant.firstName",
    "destinationPath": "name.first",
    "required": true,
    "transformType": "direct"
  },
  {
    "sourcePath": "applicant.income",
    "destinationPath": "financials.annualIncome",
    "transformType": "expression",
    "transformLogic": "income * 12"
  }
]
```

## Security

- Replace default credentials in production
- Configure CORS appropriately
- Use valid SSL certificates
- Store sensitive data in environment variables
