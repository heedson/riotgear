# riotgear
Utilisation of the Riot Game's League of Legend's API.

## basic usage
To run the service locally, it requires a PostgreSQL database to also be running.
```
$ docker run -d --rm --name riotgear-db -e POSTGRES_PASSWORD=mysecretpassword postgres
```
To then build and run the Riotgear server:
```
$ make install
$ make generate
$ docker build -t riotgear .
$ docker run -d --rm --link riotgear-db --name riotgear -p 8080:8080 -e RIOT_API_KEY=myapikey -e DB_URL=postgres://postgres:mysecretpassword@riotgear-db:5432 riotgear
```
To get a response from the current Riotgear server, head to:
```
http://0.0.0.0:8080/openapi-ui/
```
Or, alternatively, use the logs from the `riotgear` container and click the provided link.
```
$ docker logs riotgear
INFO[Jan 01 00:00:00.000] Serving gRPC-Gateway on http://0.0.0.0:8080  
INFO[Jan 01 00:00:00.001] Serving OpenAPI Documentation on http://0.0.0.0:8080/openapi-ui/ 
```
