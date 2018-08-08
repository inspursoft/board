for ip in 10.110.25.{1..253}; do 
{ 
ping -c1 $ip >/dev/null 2>&1 ; 
}&
done

basepath=$(cd `dirname $0`; pwd)


$basepath/getIp.py $1

