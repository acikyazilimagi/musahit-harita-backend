musahit-harita-backend

## Proje Hakkinda

![architecture](/static/structure.png)

Projenin localde kaldirilmasi icin asagidaki adimlar izlenmelidir.

### Gereksinimler

- Docker
- Docker Compose
- Go 1.20

### Calistirma

1. Docker compose ile uygulama ayaga kaldirilir. Her seferinde temiz bir ortam olusturmaktan emin olmak icin `--build --remove-orphans --force-recreate` parametreleri asagidaki komuta eklenebilir (container'lara ait varolan veriler silinmez).
```bash
docker compose up
```

2. env dosyasi olusturulur sisteme gore degiskenler eklenir.
   DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/postgres
   REDIS_HOST=localhost:6379
   REDIS_PASSWORD=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

3. api klasorune girilir ve asagidaki komutlar calistirilir.
```bash
go run . 
```
