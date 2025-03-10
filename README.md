# AuthSAS: Secure gRPC Authentication Microservice in Go [JWT · Redis · PostgreSQL · 2FA]

A simple authentication microservice implementing gRPC API with JWT tokens, tokens blacklist and email-based 2FA. Perfect for modern distributed systems requiring secure user management.

## 🔥 Key Features

- **JWT Authentication**
- **PostgreSQL storage**
- **Email 2FA** (TOTP codes via Yandex SMTP)
- **Password recovery**
- **Email verification**
- **Docker-ready**
- **Unit-tested core logic**

## 📦 Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL 15+
- Redis 7+
- SMTP credentials (Yandex recommended)

```bash
# 1. Clone repository
git clone https://github.com/your_username/authSAS.git && cd authSAS

# 2. Apply database migrations
edit Taskfile.yaml, then cmd-> task mg_u

# 3. Run service
edit ./config/config.yaml !!!!!!
go run cmd/app/main.go --config=./config/config.yaml
```

##⚙️ Configuration Guide ##
```yaml
app_mode: "prod"  # test/local/prod

# PostgreSQL connection
permanent_storage_path: "postgres://user:password@host:5432/dbname"

# Redis configuration
temp_storage:
  temporary_storage_path: "redis://localhost:6379/0"
  code_ttl: 10m  # 2FA/password reset/email verfy code TTL

# JWT settings
jwt_secret: "your_secure_secret_here"
jwt_token_ttl: 24h

# Email settings (Yandex SMTP)
email_sender:
  email: "your@yandex.com"
  password: "app_specific_password"
```

##Protocol Buffers Interface##
###Full API specification available in [authSASproto repository](https://github.com/BegunovDmitry/authSASproto)###
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc Logout(Empty) returns (Empty);
  rpc VerifyEmail(VerifyRequest) returns (Empty);
  rpc InitPasswordReset(ResetRequest) returns (Empty);
  rpc ConfirmPasswordReset(ConfirmResetRequest) returns (Empty);
  rpc Init2FA(Empty) returns (Init2FAResponse);
  rpc Verify2FA(Verify2FARequest) returns (AuthResponse);
}
```

##🧪 Testing Strategy##
Run unit tests for business logic:
```bash
go test ./internal/services -v -cover
```

##🐳 Docker Deployment##
```bash
# Build image
docker build -t authsas:latest .

# Run container
docker run -d \
  -p 8090:8090 \
  -v ./config:/app/config \
  --network=host \
  authsas:latest --config=/app/config/configprod.yaml
```

##📂 Project Architecture##

authSAS
├── cmd/               # Entry point
├── internal/          # Core implementation
│   ├── app/           # Application lifecycle
│   ├── config/        # Configuration parsing
│   ├── services/      # Business logic (80% test coverage)
│   ├── storages/      # PostgreSQL/Redis adapters
│   └── utils/         # Helpers (JWT, SMTP, etc)
├── migrations/        # Database schema
└── config/            # Environment configurations
