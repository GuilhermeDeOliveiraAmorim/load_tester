docker build -t test .

docker run test --url=http://google.com --requests=300 --concurrency=10
