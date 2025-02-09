<H1> best team <H1>

<h2>


docker-compose down
docker-compose up --build

docker-compose down
docker-compose build --no-cache
docker-compose up







Microservice architecture:
docker-compose build --no-cache
docker-compose up

curl -X POST "http://localhost:8083/game/submit" -H "Content-Type: application/json" -d "{\"user_id\": 1, \"score\": 100}" //test game service
docker-compose down game-service                                                                                           //purposely shut down game service
test risk assessment still works