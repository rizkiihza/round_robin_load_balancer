while true; do
    curl --location 'localhost:9000/call' \
        --header 'Content-Type: application/json' \
        --data '{
        {
          "game": "Mobile Legends",
          "gamerID": "GYUTDTE",
          "points": 20
        }'
    echo ""
done
