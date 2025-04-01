# qbee-http-broker
A http broker for accessing devices using custom headers

## Build

```
docker build -t qbeeio/qbee-http-broker:1.0.0 .
```

## Run
```
docker run \
    -e QBEE_EMAIL="<username>" \
    -e QBEE_PASSWORD="<password>" \
    -e QBEE_TOKEN="<token>" \
    -p 8081:8081 qbeeio/qbee-http-broker:1.0.0
```

## Run with dev environment

```
docker run \
    -e QBEE_EMAIL="$QBEE_EMAIL" \
    -e QBEE_PASSWORD="$QBEE_PASSWORD" \
    -e QBEE_BASEURL="http://frontend-go:8080" \
    --network platform_default \
    -p 8081:8081 qbeeio/qbee-http-broker:1.0.0
```

## Relevant headers

```
X-Qbee-Device-Id
X-Qbee-Device-Port (optional, default port 80)
X-Qbee-Authorization 
```

## curl website on device

Device
```
python3 -m http.server
```

Client
```
curl -H "X-Qbee-Authorization: <token>" -H "X-Qbee-Device-Port: 8000" -H "X-Qbee-Device-Id: <device-id>" http://localhost:8081
```
