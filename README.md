High level steps

- install dependencies
- install docker if does not have
- run docker compose (show environment variables)


### How to create [migrations](https://www.freecodecamp.org/news/database-migration-golang-migrate/)

- **Install migrate/golang-migrate**

*Windows (using [scoop](https://scoop.sh/))*
```sh
scoop install migrate
```

*Linux (using curl)*
```sh
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb*release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
```

*Mac (using hombrew)*
```sh
brew install golang-migrate
```

- **Execute create migration**

```sh
make create_migration MIGRATIONNAME=write_migration_name
```

- **Write the changes you want to make in the files generated in `database/migrations/`**

- **Execute migration up**

```sh
make migration_up
```

- **Execute migration rollback**

```sh
make migration_down
```

- **Optionally, you can override the environment used for the migrations using ENV**
- 
```sh
make migration_up ENV=.env.dev
make migration_down ENV=.env.test
```