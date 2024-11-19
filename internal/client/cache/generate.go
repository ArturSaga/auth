package cache

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate /Users/artursagataev/GolandProjects/auth/bin/minimock -i RedisClient -o ./mocks/ -s "_minimock.go"
