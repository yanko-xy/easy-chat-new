user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user-api.mk release-test

social-rpc-dev:
	@make -f deploy/make/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/make/social-api.mk release-test

release-test: user-rpc-dev social-rpc-dev user-api-dev social-api-dev

install-user-rpc:
	cd ./deploy/script && chmod +x ./user-rpc-test.sh && ./user-rpc-test.sh

install-user-api:
	cd ./deploy/script && chmod +x ./user-api-test.sh && ./user-api-test.sh

install-social-rpc:
	cd ./deploy/script && chmod +x ./social-rpc-test.sh && ./social-rpc-test.sh

install-social-api:
	cd ./deploy/script && chmod +x ./social-api-test.sh && ./social-api-test.sh

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh

