cd ../frontend
CALL ng build --prod
cd ..
CALL docker-compose up -d --build