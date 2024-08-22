function handler () {
    EVENT_DATA=$1
    echo "event data: $EVENT_DATA" >&2

    BODY=$(echo $EVENT_DATA | /opt/jq/jq -r '.body')
    HEADERS=$(echo $EVENT_DATA | /opt/jq/jq -r '.headers')
    METHOD=$(echo $EVENT_DATA | /opt/jq/jq -r '.requestContext.http.method')
    PATH=$(echo $EVENT_DATA | /opt/jq/jq -r '.requestContext.http.path')

    HOST="localhost"
    PORT=8080

    /opt/envoy/envoy -c ./envoy.yaml &
    ./server &
    sleep 2

    RESPONSE=$(curl -s -X $METHOD \
        -H "$(echo $HEADERS | /opt/jq/jq -r 'to_entries[] | "\(.key): \(.value)"')" \
        -d "$BODY" \
        http://$HOST:$PORT$PATH)

    echo $RESPONSE
}
