function handler () {
    echo "Handler script started." >&2

    EVENT_DATA=$1
    echo "Event data: $EVENT_DATA" >&2

    RESPONSE=$(echo "$EVENT_DATA" | ./server)
    echo "Response: $RESPONSE" >&2
    echo $RESPONSE

    echo "Handler script finished." >&2
}
