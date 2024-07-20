while true; do
    curl --location 'localhost:9000' \
        --header 'Content-Type: application/json' \
        --data '{
    "Hello": "world"
}'
    echo ""
done
