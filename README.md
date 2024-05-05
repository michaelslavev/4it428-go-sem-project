# 4it428-go-sem-project
Semestral project - STRV GO Lecture (4IT428)
- **DB & Auth hosted by:** supabase.com
- **Hosted on:** http://109.205.180.225:9069
- **Postman collection:** [Download colletion here](https://we.tl/t-szgWD8cDL4)


## Project structure
Project is divided into multiple services, each service is in its own directory. 
Each service has its own **.env** file, where you can set environment variables for the service.

- **api-gateway** 
  - is the main entry point working as reverse proxy between docker containers which are running services strictly without exposing them to the outside world
  - it is also responsible for routing requests to the correct service with jwt authentication if service is not public
- **auth-service**
  - handles register, login and refresh token requests 
- **newsletter-management-service**
  - handles CRUD operations for newsletters
- **publishing-service**
  - handles operations for posts
- **subscription-service**
  - handles subscribing and unsubscribing to newsletters
```shell
.
└── services
    ├── api-gateway
    │   ├── .env
    │   └── main.go
    ├── auth-service
    │   ├── .env
    │   └── main.go
    ├── newsletter-management-service
    │   ├── .env
    │   └── main.go
    ├── publishing-service
    │   ├── .env
    │   └── main.go
    └── subscription-service
        ├── .env
        └── main.go
```

## How to run project with docker-compose?
- Don't forget to fill **.env** for each service


### Running docker-compose without detached mode (for debugging):
```shell
..\> docker-compose up --build
```

### Running docker-compose with detached mode:
```shell
..\> docker-compose up --build -d
```