# Set defaults if environment variables aren't set
PGPASS ?= 123456
HOST_ADDR ?= localhost
PORT ?= 5432
USER ?= admin
DATABASE ?= postgres

tags:
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -f tag.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -c "\copy tag FROM Tags.csv CSV HEADER"

users:
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -f users.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -c "\copy users FROM Users.csv CSV HEADER"
