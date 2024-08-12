function handler () {
    echo "Handler script started." >&2

    EVENT_DATA=$1
    echo "event data: $EVENT_DATA" >&2

    RESPONSE=$(echo "$EVENT_DATA" | ./server)
    echo $RESPONSE

    echo "Handler script finished." >&2
}
