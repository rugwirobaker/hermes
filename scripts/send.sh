#/bin/bash


export ADDRESS=localhost:8000

export ID="$(curl -s -d '@hermes.json' -H 'Content-Type: application/json' -X POST $ADDRESS/api/send | jq -r '.id')"

echo $ID

sleep 5

curl -s $ADDRESS/api/events/$ID/status | jq -r '.'
