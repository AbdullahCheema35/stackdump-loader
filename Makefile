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

badges:
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -f badge.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -c "\copy badge FROM Badges.csv CSV HEADER"

votes:
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -f vote.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -c "\copy vote FROM Votes.csv CSV HEADER"

comments:
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -f comment.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(USER) -d $(DATABASE) -c "\copy comment FROM Comments.csv CSV HEADER"
