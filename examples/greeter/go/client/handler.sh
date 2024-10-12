function handler () {
    echo "Handler script started." >&2

    RESPONSE=$(./service)
    echo "Response: $RESPONSE" >&2
    echo $RESPONSE

    echo "Handler script finished." >&2
}
