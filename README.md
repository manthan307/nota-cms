# Nota CMS

I building a Headless cms in golang which use postgresql as database.
The project is not complete yet i'm still developing it just to learn more about CMS.

## Run this porject

To run this project you need go@1.25.1 and docker. Then run following command in your terminal

make a .env file with this variables

```env
POSTGRES_USER=admin
POSTGRES_PASSWORD=password
POSTGRES_DB=cms
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

PORT=8000
JWT_SECRET_KEY=your_jwt_secret_key

MINIO_USE_SSL=false
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET_NAME=cms-bucket
MINIO_ENDPOINT=localhost:9000
MINIO_REGION=us-east-1
```

and then start the server

```bash
docker compose up
```

REMEMBER: make .env file in root with variable shown in .env.example without it will not run

# Nota CMS API Routes

Base URL: `/api/v1`

---

## Authentication

| Method | Endpoint         | Description         |
| ------ | ---------------- | ------------------- |
| POST   | `/auth/register` | Register a new user |
| POST   | `/auth/login`    | Login user          |
| POST   | `/auth/verify`   | Verify JWT token    |

---

## Schemas

| Method | Endpoint                     | Role   | Description         |
| ------ | ---------------------------- | ------ | ------------------- |
| POST   | `/schemas/create`            | editor | Create a new schema |
| GET    | `/schemas/get_by_id/:id`     | viewer | Get schema by ID    |
| GET    | `/schemas/get_by_name/:name` | viewer | Get schema by name  |
| GET    | `/schemas/list`              | viewer | List all schemas    |
| DELETE | `/schemas/delete/:id`        | editor | Delete schema by ID |

---

## Content

| Method | Endpoint                        | Role   | Description                          |
| ------ | ------------------------------- | ------ | ------------------------------------ |
| POST   | `/content/create`               | editor | Create a new content item            |
| DELETE | `/content/delete/:id`           | editor | Delete content by ID                 |
| GET    | `/content/get/:id`              | all    | Get content by ID                    |
| GET    | `/content/get_all/:schema_name` | all    | Get all content for a schema         |
| POST   | `/content/update`               | editor | Update content item (data/published) |

---

## Media

| Method | Endpoint        | Role   | Description                                 |
| ------ | --------------- | ------ | ------------------------------------------- |
| POST   | `/media/upload` | editor | Upload a new media file                     |
| DELETE | `/media/delete` | editor | Delete media file (pass `file_url` in body) |
