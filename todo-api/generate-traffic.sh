#!/usr/bin/env bash
set -euo pipefail

BASE="${TODO_API_URL:-http://todo-api-sandbox.apps-crc.testing}"
ROUNDS="${1:-3}"

create() { curl -sf -X POST "$BASE/todos" -H "Content-Type: application/json" -d "{\"title\":\"$1\"}"; }
list()   { curl -sf "$BASE/todos"; }
get()    { curl -sf "$BASE/todos/$1"; }
update() { curl -sf -X PUT "$BASE/todos/$1" -H "Content-Type: application/json" -d "$2"; }
del()    { curl -sf -X DELETE "$BASE/todos/$1"; }

echo "Target: $BASE"
echo "Rounds: $ROUNDS"
echo

for i in $(seq 1 "$ROUNDS"); do
    echo "=== Round $i ==="

    # Create a batch of todos
    id1=$(create "Buy groceries (round $i)"  | jq -r .id)
    id2=$(create "Write tests (round $i)"    | jq -r .id)
    id3=$(create "Deploy to prod (round $i)" | jq -r .id)
    echo "  created: $id1 $id2 $id3"

    # List all
    count=$(list | jq length)
    echo "  listed: $count items"

    # Get each one
    get "$id1" > /dev/null
    get "$id2" > /dev/null
    get "$id3" > /dev/null
    echo "  fetched: $id1 $id2 $id3"

    # Update some
    update "$id1" '{"completed":true}' > /dev/null
    update "$id2" '{"title":"Write MORE tests","completed":true}' > /dev/null
    echo "  updated: $id1 $id2"

    # Delete one
    del "$id3"
    echo "  deleted: $id3"

    # Generate some 404s
    curl -so /dev/null "$BASE/todos/does-not-exist"
    curl -so /dev/null -X PUT "$BASE/todos/nope" -H "Content-Type: application/json" -d '{"title":"x"}'
    curl -so /dev/null -X DELETE "$BASE/todos/gone"
    echo "  404s: 3"

    echo
done

echo "Done. $((ROUNDS * 3)) todos created, $((ROUNDS * 2)) updated, $((ROUNDS * 1)) deleted, $((ROUNDS * 3)) 404s."
