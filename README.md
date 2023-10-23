# golang based hotel-booking-service

## To Start the docker service:

1. Ensure you have docker compose installed
2. On the terminal from the root of this project run
   $docker compose up
3. To know container ID of hotel-booking-service
   $docker ps

4. $docker exec -it <containerId-hotel-booking-service> /bin/sh

5. $make seed
   This will give you two users (normal and admin) with their x-api-token token

6. Now open thunder client and on then experiment
   GET http://localhost:4444/api/v1/hotel/?page=1&limit=15&rating=5
   GET http://localhost:4444/api/v1/user

and so on.

## Remember if you want to use it in a normal environment without docker,

1. In config/config.go enable this line
   const projectDirName = "hotel-booking-service"

2. In .env file replace mongodb://mongodb:27017 to mongodb://localhost:27017

3. make seed and then make run
