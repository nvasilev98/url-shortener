# URL-Shortener

## Set up environment for local run:
### Load Firestore
1. Local

    Navigate to root folder of the project `url-shortener` and execute `docker compose up -d`

    Export the following environment variable:

    ```export FIRESTORE_EMULATOR_HOST=0.0.0.0:8342```

2. GCP Service

    Unset emulator if set and export the following environment variable
    
    ```
    export FIRESTORE_EMULATOR_HOST=0.0.0.0:8342
    export GOOGLE_APPLICATION_CREDENTIALS="<path-to-service-account>"
    ```

### Set up application config
1. Set the following environment variables

    ```
    export HOST=<host>
    export PORT=<port>
    ```
    Note: If app config is not set default one will be used and the application will be availabe on `localhost:8080`

### Start application

1. Execute from root folder of the project: `go run cmd/urlshortener/main.go`

## Run unit tests
1. Start Firestore emulator:

    In order to prepare environment for the tests we need to start a Firestore emulator. Navigate root folder of the project `url-shortener` and execute `docker compose up`

    Export the following environment variable:

    ```export FIRESTORE_EMULATOR_HOST=0.0.0.0:8342```

2. Run: 

   Execute: `make unit-tests`