docker rm -f brcaidunitdb
docker run --name brcaidunitdb -p 3307:3306 -e MYSQL_ROOT_PASSWORD=howdy -e MYSQL_DATABASE=brcaid -d mysql

