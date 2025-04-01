# CourseProject
This is backend with microservice architecture for the educational courses platform.

Stack: Golang, C++, Docker, Swagger, PostgreSQL, JWT.

Description: course platform where user can authorize, register for the courses, track his activity.

Create .env file in root directory and add your secret values (there is an example):
```
JWT_SECRET_KEY=12345
REDIS_REFRESH_PASSWORD=12345
REDIS_VERIFY_PASSWORD=12345
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_PASSWORD=your_host_pass
FROM_EMAIL=your_host_email

DB_USER=admin
DB_PASSWORD=admin
DB_SOURCE_NAME=postgres://admin:admin@postgres_users:5432/usersdb?sslmode=disable
DB_PORT=5432

```


