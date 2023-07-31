#!make
include .env

cmgr:
	bin/migrate create -ext sql -dir db/migrations -seq ${name}

migup:
	bin/migrate -path db/migrations -database "${PSQL}" -verbose up

migdown:
	bin/migrate -path db/migrations -database "${PSQL}" -verbose down

setupair:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

setupmigrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz -C bin/

setup: setupair setupmigrate

run:
	bin/air