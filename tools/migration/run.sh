#!/bin/bash

export PYTHONPATH=$PYTHONPATH:/board-migration
if [ -z "$DB_IP" -o -z "$DB_PORT" -o -z "$DB_USR" -o -z "$DB_PWD" ]; then
    echo "DB_IP or DB_PORT or DB_USR or DB_PWD not set, exiting..."
    exit 1
fi

source ./alembic.tpl > ./alembic.ini

DBCNF="-h ${DB_IP} -P ${DB_PORT} -u ${DB_USR} -p${DB_PWD}"

#prevent shell to print insecure message
export MYSQL_PWD="${DB_PWD}"

if [[ $1 = "help" || $1 = "h" || $# = 0 ]]; then
    echo "Usage:"
    echo "backup                perform database backup"
    echo "restore               perform database restore"
    echo "up,   upgrade         perform database schema upgrade"
    echo "h,    help            usage help"
    exit 0
fi

if [[ ( $1 = "up" || $1 = "upgrade" ) && ${SKIP_CONFIRM} != "y" ]]; then
    echo "Please backup before upgrade."
    read -p "Enter y to continue updating or n to abort:" ans
    case $ans in
        [Yy]* )
            ;;
        [Nn]* )
            exit 0
            ;;
        * ) echo "illegal answer: $ans. Upgrade abort!!"
            exit 1
            ;;
    esac

fi


key="$1"
case $key in
up|upgrade)
    VERSION="$2"
    if [[ -z $VERSION ]]; then
        VERSION="head"
        echo "Version is not specified. Default version is head."
    fi
    echo "Performing upgrade ${VERSION}..."
    alembic -c ./alembic.ini current
    alembic -c ./alembic.ini upgrade ${VERSION}
    rc="$?"
    alembic -c ./alembic.ini current	
    echo "Upgrade performed."
    echo $rc
    exit $rc
    ;;
backup)
    echo "Performing backup..."
    mysqldump $DBCNF --add-drop-database --databases board > ./backup/board.sql
    echo "Backup performed."
    ;;
restore)
    echo "Performing restore..."
    mysql $DBCNF < ./backup/board.sql
    echo "Restore performed."
    ;;
*)
    echo "unknown option"
    exit 0
    ;;
esac
