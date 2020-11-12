#!/bin/bash

if [[ "$(docker inspect --format '{{json .State.Health.Status }}' phoos-db)" != "\"healthy\"" ]]; then
    docker-compose up -d phoos-db

    while [[ "$(docker inspect --format '{{json .State.Health.Status }}' phoos-db)" != "\"healthy\"" ]]; do
        echo -n "."
        sleep 1
    done

    printf "\ncreating database\n"
    docker exec -i phoos-db mysql -h127.0.0.1 -uroot -ptesting <<< "create database phoosball;"
    docker exec -i phoos-db mysql -h127.0.0.1 -uroot -ptesting <<< "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' identified by 'testing';"
    printf "initial set up\n"
    docker exec -i phoos-db mysql -h127.0.0.1 -uroot -ptesting mysql < db/test/1.phoosball.sql
    docker exec -i phoos-db mysql -h127.0.0.1 -uroot -ptesting <<< "INSERT INTO phoosball.schema_history (version, description, name, checksum) VALUES (1, 'phoosball', '1.phoosball.sql', '6bba020d7cdd9386c3c1d59d68d5de78');"
fi

docker-compose up -d phoos-server