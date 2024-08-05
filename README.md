docker build -t test .

docker run test --url=https://www.guilhermeamorim.com --requests=300 --concurrency=10
