# user-crud

To run the app first run `docker compose up` then, when all the necessary containers up and running,
`make migrate` to create db schema and as the last step run `make run`
to run Golang app.

Regarding health check of the app, I think the current response structure is good enough
to meet basic needs, if there is a need it can be improved further.

There is a swagger doc to give you more detailed information
about endpoints the app exposes.

The business logic layer is within `user-crud/internal/services` folder.

### Servers
1. HTTP - `localhost:8080`
2. GRPC - `localhost:8081`
3. HealthCheck HTTP - `localhost:8082/health`

PS, I used buf.gen tool to generate DTO objects because it's boilerplate code to map incoming request to respective data structure.