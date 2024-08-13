function handler () {
    echo "Handler script started." >&2

    RESPONSE=$(./client --addr=$PROTO_SERVER_ADDR)  # PROTO_SERVER_ADDR is an env var
    echo "Response: $RESPONSE" >&2
    echo $RESPONSE

    echo "Handler script finished." >&2
}
