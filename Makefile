# Fail if 'filename' is not provided
ifndef filename
$(error ❌ Please provide a filename, e.g.: make posts filename=1_Posts.csv)
endif

# Set defaults if environment variables aren't set
PGPASS ?= 123456
HOST_ADDR ?= localhost
PORT ?= 5432
USER ?= admin
DATABASE ?= postgres

tags:
# 	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -f tag.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -c "\copy tag FROM Tags.csv CSV HEADER"

users:
# 	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -f users.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -c "\copy users FROM Users.csv CSV HEADER"

badges:
# 	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -f badge.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -c "\copy badge FROM Badges.csv CSV HEADER"

votes:
# 	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -f vote.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -c "\copy vote FROM $(filename) CSV HEADER"

comments:
# 	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -f comment.sql
	PGPASSWORD=$(PGPASS) psql -h $(HOST_ADDR) -p $(PORT) -U $(DB_USER) -d $(DATABASE) -c "\copy comment FROM $(filename) CSV HEADER"
