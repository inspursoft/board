
board_dir=$1
service_cov_name="service_cov.out"
collecter_cov_name="collecter_cov.out"

# check the go path

if [ ! -d $go ]; then
  mkdir -p $go
fi

#export GOPATH=

# if the first clone

if [ ! -d $board_dir ]; then
   #mkdir -p $board_dir
   #cd $board_dir/../
   #git clone http://10.110.18.40:10080/inspursoft/board.git
   cd $board_dir
   make prepare
   make compile_ui
   make -e start DEVFLAG=dev
   
# if have cloned 

else
   cd $board_dir/make/dev
   docker-compose down -v
   docker-compose up -d
   cd $board_dir
   git pull
fi

