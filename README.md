# riotgear
Utilisation of the Riot Game's League of Legend's API.

## basic usage
To run the service locally.
```
make install
make generate
sudo docker build -t riotgear .
sudo docker run -d --rm --name riotgear -p 8080:8080 -e RIOT_API_KEY=myapikey riotgear
```
To get a response from the current Riotgear server, head to:
```
http://0.0.0.0:8080/openapi-ui/
```
