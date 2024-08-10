function handler () {
    EVENT_DATA=$1
    echo "event data: $($EVENT_DATA)" >&2

    RESPONSE="{\"statusCode\": 200, \"body\": \"Hello from Lambda!\"}"
    echo $RESPONSE
}
