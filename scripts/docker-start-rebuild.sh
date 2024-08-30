go mod tidy
docker-compose -p theoverwatchtools -f ./docker/docker-compose.yaml down --remove-orphans --volumes || true
docker-compose -p theoverwatchtools -f ./docker/docker-compose.yaml up --force-recreate --build