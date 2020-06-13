cd ..
CALL docker-compose down -v --rmi all --remove-orphans
CALL docker image prune