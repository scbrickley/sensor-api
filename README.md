# Sensor API

This program runs a service for storing, retrieving, and manipulating metadata for various sensors.

## Running the service

### Set up the database

The sensor service uses PostgreSQL for the data layer. If you have the Docker CLI installed, you can run the following
command to start a suitable postgres instance on your machine:

```
docker run --rm --name sensor-db \
	-e POSTGRES_HOST_AUTH_METHOD=trust \
	-e POSTGRES_DB=sensor-db \
	-p 5432:5432 postgres:15
```

If you have the `psql` command line tool installed, you may optionally connect to the postgres instance with the following command,
if you wish to manually inspect the database:

```
psql --host=localhost -U postgres sensor-db
```

Once the database is running, you can start the actual service.

### Start the Sensor API

With the postgres instance running in the background, run the following command from the root of the `sensor-api` repo to start the
service:

```
go run main.go
```

On startup, the service will connect to the postgres instance, create a suitable table for the data it wants to store, and start listening
for user requests on port `8000` of your machine.

To provide some insight into the schema of the sensor database, here is the SQL command that the service uses to create the table it uses:

```
CREATE TABLE IF NOT EXISTS sensors (
    name VARCHAR NOT NULL UNIQUE,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    tags VARCHAR[] NOT NULL
);
```

## Using the service

The Sensor API exposes six endpoints that allow users to:

- Add new sensors to the database
- List all sensors in the database
- View metadata for a specific sensor
- Update metadata for a specific sensor
- Delete a specific sensor

### List all sensors

Run the following cURL command to get a list of all the sensors stored in the database:

```
curl -X GET localhost:8000/sensors
```

Here is a sample response:

```
{
    "success": true,
    "sensor": [
        {
            "name": "sensor1",
            "latitude": 1,
            "longitude": -1,
            "tags": [
                "theFirstOne"
            ]
        },
        {
            "name": "sensor2",
            "latitude": -1,
            "longitude": 2,
            "tags": [
                "theSecondOne",
            "letsDoItAgain"
            ]
        },
        {
            "name": "sensor3",
            "latitude": -10,
            "longitude": 10,
            "tags": [
                "outlier"
            ]
        }
    ],
    "error_msg": ""
}
```

### Add a new sensor

You can add a new sensor to the database by making a `POST` request with an appropriately constructed JSON payload

```
curl -X POST -d '{
	"name": "sensor4",
	"longitude": 80.123,
	"latitude": 35.456,
	"tags": ["barometer", "requiresMaintenance"]
}' localhost:8000/sensors
```

The response will be the JSON data from the row created as the result of your request, which should match exactly with the JSON payload
provided as part of the original `POST` request, along with some extra metadata that will show whether or not the request was successful,
and an error message if relevant.

Here's what a successful response might look like for the above request:

```
{
    "success": true,
    "sensor": {
        "name": "sensor4",
        "latitude": 35.456,
        "longitude": 80.123,
        "tags": [
            "barometer",
            "requiresMaintenance"
        ]
    },
    "error_msg": ""
}
```

Users will receive an error if the JSON payload is malformed, or if the provided sensor name is not unique.

### View a sensor

To view metadata for a particular sensor, run a `GET` request similar to the following:

```
curl -X GET localhost:8000/sensors/<name>
```

...replacing `<name>` with the unique name of the sensor in quesiton.

### Update a sensor

To overwrite the metadata for a particular sensor, you can make a `PUT` request to the same endpoint as above, replacing <name> with the
unique name of the sensor you want to update, and providing a JSON payload that includes the updated name, latitude, longitude, and tags
for the specified sensor.

Here's a sample cURL request:

```
curl -X PUT -d '{
	"name": "sensor4",
	"longitude": 80.123,
	"latitude": 35.456,
	"tags": ["barometer"]
}' localhost:8000/sensors/sensor4
```

### Delete a sensor

You can make a `DELETE` request to the same endpoint to remove a particular sensor from the database. Here is a sample cURL request:

```
curl -X DELETE localhost:8000/sensors/sensor4
```

### Find the nearest sensor

You can ask the service to identify the sensor that's nearest to a provided coordinate. To do so, make a `GET` request to the `nearest`
endpoint, like so:

```
curl -X GET -d '{
    "longitude": 50.002,
    "latitude": 23.108
}' localhost:8000/sensors/nearest
```

Given this request, as well as the following data:

```
sensor-db=# select * from sensors;
  name   | latitude | longitude |              tags
---------+----------+-----------+---------------------------------
 sensor1 |        1 |        -1 | {theFirstOne}
 sensor2 |       -1 |         2 | {theSecondOne,letsDoItAgain}
 sensor3 |      -10 |        10 | {}
 sensor4 |   35.456 |    80.123 | {barometer}
```

We would expect the following response:

```
{
    "success": true,
    "sensor": {
        "name": "sensor4",
        "latitude": 35.456,
        "longitude": 80.123,
        "tags": [
            "barometer",
        ]
    },
    "error_msg": ""
}
```
