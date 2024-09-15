goctl model mysql ddl -src="./deploy/sql/user.sql" -dir="./apps/user/models" -c

goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero
