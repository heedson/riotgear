# riotgear
Utilisation of the Riot Game's League of Legend's API.

## basic usage
To run the service locally, it requires a PostgreSQL database to also be running.
```
$ make install
$ make generate

$ docker build -t riotgear .

$ docker network create riotgearnetwork
$ docker run -d --rm --name riotgear-db -p 5432:5432 --network riotgearnetwork -e POSTGRES_PASSWORD=mysecretpassword postgres
$ docker run -d --rm --name riotgear -p 8080:8080 --network riotgearnetwork -e RIOT_API_KEY=myapikey -e DB_URL=postgres://postgres:mysecretpassword@riotgear-db:5432/postgres riotgear
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
