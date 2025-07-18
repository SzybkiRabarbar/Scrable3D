files=(
    "repo" "db.go"
    "repo" "repository.go"
    "svc" "avchar.go"
    "svc" "field.go"
    "svc" "game.go"
    "svc" "player.go"
    "ctrl" "game.go"
    "ctrl" "words.go"
)

for ((i=0; i<${#files[@]}; i+=2)); do
    echo Generate ${files[i+1]} from ${files[i]}
    mockgen -source=internal/${files[i]}/${files[i+1]} -destination=internal/mock/mock_${files[i]}_${files[i+1]} -package=mock
done