# User Locations App 

## Project overview

### General concept

User Locations App is a system to process user locations developed . It’s made up of two microservices with the following features:

Microservice 1 - location-mgmt: Exposes an api to do two tasks:

1. Update current user location by the username.

2. Search for users in some location within the provided radius (with pagination).

Microservice 2 - location-history-mgmt: Exposes an api to do one task:

3. Returns distance traveled by a person within some date/time range. Time range defaults to 1 day.

![alt text](./resources/locations-app-structure.png)

### Main process flows

Microservice 1 - location-mgmt:

- Save: Update current user location by the username.

1. User updates his current location
2. The current location is updated in the microservice 1 through Save endpoint
3. Microservice 1 sends the updated location to microservice 2 using Protobuf and GRPC
4. Microservice 2 prepares and persists the data

- GetUsersByLocationAndRadius: Search for users in some location within the provided radius.

1. User requires user list in some location within the provided radius
2. The user list is requested in the microservice 1 through GetUsersByLocationAndRadius REST endpoint
3. Microservice 1 sends the reference location and radius to microservice 2 using Protobuf and GRPC
4. Microservice 2 fetches user list and returns it to microservice 1 using Protobuf and GRPC
5. Microservice 1 returns user list through REST response 

Microservice 2 - location-history-mgmt:

- GetDistanceTraveled: Returns distance traveled by a person within some date/time range. Time range defaults to 1 day.

1. User requires traveled distance by a person within some date/time range
2. The traveled distance is requested in the microservice 2 through GetDistanceTraveled REST endpoint
3. Microservice 2 receives the username and date/time range as input and returns the distance traveled during that period through REST response


## Endpoints usage

### Microservice 1 - location-mgmt Endpoints

1. Update current user location by the username.

curl --location --request POST 'http://localhost:8081/location-mgmt/locations' \
--header 'Content-Type: application/json' \
--data-raw '{
"userName":"oboada",
"Latitude":35.12314,
"Longitude":27.64532
}'

2. Search for users in some location within the provided radius (with pagination).

curl --location --request POST 'http://localhost:8081/location-mgmt/locations' \
--header 'Content-Type: application/json' \
--data-raw '{
"userName":"oboada",
"Latitude":35.12314,
"Longitude":21.56878
}'

3. Returns distance traveled by a person within some date/time range. Time range defaults to 1 day.

curl --location --request GET 'http://localhost:8080/location-history-mmg/locations/distance/oboada/2022-09-02T11:26:18+00:00/2022-09-05T15:30:00+00:00' \
--header 'Content-Type: application/json'
 

## Parameter definition

username - 4-16 symbols (a-zA-Z0-9 symbols are acceptable)
coordinates - fractional part of a number should be limited by the 8 signs, latitude and longitude should be validated by the regular rules. For example:
35.12314, 27.64532
39.12355, 27.64538
dates - use ISO 8601 date format (2021-09-02T11:26:18+00:00)