go install github.com/pressly/goose/v3/cmd/goose@latest
git clone https://github.com/meehighlov/grats.git grats-tmp
./goose -dir=/grats-tmp/migrations sqlite3 grats.db up
./goose -dir=/grats-tmp/migrations sqlite3 grats.db status
rm -rf grats-tmp