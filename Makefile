user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user-api.mk release-test

social-rpc-dev:
	@make -f deploy/make/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/make/social-api.mk release-test

im-rpc-dev:
	@make -f deploy/make/im-rpc.mk release-test

im-api-dev:
	@make -f deploy/make/im-api.mk release-test

im-ws-dev:
	@make -f deploy/make/im-ws.mk release-test

task-mq-dev:
	@make -f deploy/make/task-mq.mk release-test

release-test: user-rpc-dev social-rpc-dev user-api-dev social-api-dev im-rpc-dev im-api-dev im-ws-dev task-mq-dev

install-user-rpc:
	cd ./deploy/script && chmod +x ./user-rpc-test.sh && ./user-rpc-test.sh

install-user-api:
	cd ./deploy/script && chmod +x ./user-api-test.sh && ./user-api-test.sh

install-social-rpc:
	cd ./deploy/script && chmod +x ./social-rpc-test.sh && ./social-rpc-test.sh

install-social-api:
	cd ./deploy/script && chmod +x ./social-api-test.sh && ./social-api-test.sh

install-im-rpc:
	cd ./deploy/script && chmod +x ./im-rpc-test.sh && ./im-rpc-test.sh

install-im-api:
	cd ./deploy/script && chmod +x ./im-api-test.sh && ./im-api-test.sh

install-im-ws:
	cd ./deploy/script && chmod +x ./im-ws-test.sh && ./im-ws-test.sh

install-task-mq:
	cd ./deploy/script && chmod +x ./task-mq-test.sh && ./task-mq-test.sh

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh

