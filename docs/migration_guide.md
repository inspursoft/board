# Board upgrade and database migration guide

When upgrading your existing Board instance to a newer version, you may need  to migrate the data in your database. Refer to [change log](../tools/migration/changelog.md) to find out  whether there is any change in the database. If there is, you should go through the database migration process. Since the migration may alter the database schema, you should **always** backup your data before any migration.

*If your install Board for the first time, or the database version is the same as that of the latest version, you do not need any database migration.*

### Upgrading Board and migrating data
1 Log in to the host that Board runs on, stop and remove existing Board instance if it is still running:
   
   ```sh
   cd board
   docker-compose down
   ```

2 Backup Board's current files so that you can roll back to the current version when it is neccessary.
   
   ```sh
   cd ..
   mv board /my_backup_dir/board
   ```

3 Get the latest Board release package from Gogits:
   
   ```
   https://github.com/inspursoft/board
   ```

4 Before upgrading Board, perform database migration first. The migration tool is delivered as a Docker image, so you should build it yourself.
   
   ```sh
   cd board/tools/migration
   docker build -t board-migration .
   ```

5 You should start you current Board database by handy.
 
   ```sh
   docker run -d -p 3306:3306 -v /data/board/database:/var/lib/mysql -e DB_PASSWORD=root123 dev_db:dev
   ```

6 Backup database to a directory such as `/data/board-migration/backup`. You also need the IP address, port number,username and password to access the database are provided via environment variables "DB_IP", "DB_PORT", "DB_USR", "DB_PWD".
 
   ```sh
   docker run --rm -v /data/board-migration/backup:/board-migration/backup -e DB_IP=10.0.0.0 -e DB_PORT=3306 -e DB_USR=root -e DB_PWD=root123 board-migration backup
   ```
7 Upgrade database schema and migrate data.
 
   ```sh
   docker run --rm -v /data/board-migration/backup:/board-migration/backup -e DB_IP=10.0.0.0 -e DB_PORT=3306 -e DB_USR=root -e DB_PWD=root board-migration upgrade head
   ```

   **NOTE:**
   If you execute this command in a short while after starting the Board database, you may meet errors as the database is not ready for connection. Please retry it after waiting for a while.
   
### Roll back from an upgrade
For any reason, if you want to back to the previous version of Board, follow the below steps:

1 Stop and remove the current Board service if it is still running.
 
   ```sh
   cd board
   docker-compose down
   ```

2 Start stand-alone container of Board database
 
   ```sh
   docker run -d -p 3306:3306 -v /data/board/database:/var/lib/mysql -e DB_PASSWORD=root123 dev_db:dev
   ```

3 Restore database from backup file in `/data/board-migration/backup`.
 
   ```sh
   docker run --rm -v /data/board-migration/backup:/board-migration/backup -e DB_IP=10.0.0.0 -e DB_PORT=3306 -e DB_USR=root -e DB_PWD=root123 board-migration restore
   ```

4 You should use the corresponding version of Board to start with it.

### Migration tool reference
- Use `help` command to show instruction of migration tool:
ta/board-migration/backup`.
 
   ```sh
   docker run --rm -v /data/board-migration/backup:/board-migration/backup -e DB_IP=10.0.0.0 -e DB_PORT=3306 -e DB_USR=root -e DB_PWD=root123 board-migration help
   ```

