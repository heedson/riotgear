# riotgear
Utilisation of the Riot Game's League of Legend's API.

## basic usage
To run the service locally.
```
sudo docker build -t riotgear .
sudo docker run -p 8080:8080 --name riotgear -e RIOT_API_KEY=myapikey riotgear
```
To get a response from the current Echo test server.
```
curl -d '{"value":"hello world"}' -H "Content-Type: application/json" -X POST http://$(sudo docker inspect --format '{{.NetworkSettings.IPAddress}}' riotgear):8080/api/v1/echo
```
