cd ..
CALL docker-compose down -v --rmi all --remove-orphans
CALL docker image prune
CALL FOR /f "tokens=*" %i IN ('docker ps -a -q') DO docker stop %i
CALL FOR /f "tokens=*" %i IN ('docker ps -a -q') DO docker rm %i