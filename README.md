<<<<<<< HEAD
docker build -t test .

docker run test --url=http://google.com --requests=300 --concurrency=10
=======
docker build -t load_tester .

docker run load_tester --url=http://google.com --requests=1000 --concurrency=10
>>>>>>> c595f71f3145ac8cbbf5672c4b80757b3788151c
