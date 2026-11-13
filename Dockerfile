# ==========================================
# TAHAP 1: Build Frontend Ionic Vue 3
# ==========================================
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# ==========================================
# TAHAP 2: Build Binary Golang + SQLite
# ==========================================
FROM golang:1.24-alpine AS backend-builder
# modernc.org/sqlite adalah Pure Go SQLite, tidak membutuhkan gcc/CGO
WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# Salin hasil build static Ionic ke folder embed backend
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Compile binary Golang statically (CGO_ENABLED=0 murni Go)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o siakad-app .

# ==========================================
# TAHAP 3: Final Production Image (Alpine)
# ==========================================
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Folder penyimpanan file database SQLite
RUN mkdir -p /app/data

# Salin binary dari backend builder
COPY --from=backend-builder /app/siakad-app .

EXPOSE 8080

ENV DB_PATH=/app/data/siakad.db
ENV PORT=8080
ENV GIN_MODE=release

CMD ["./siakad-app"]
