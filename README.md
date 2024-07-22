docker build -t load_tester .

docker run load_tester --url=http://google.com --requests=1000 --concurrency=10
