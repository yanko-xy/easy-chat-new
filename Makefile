user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user-api.mk release-test

social-rpc-dev:
	@make -f deploy/make/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/make/social-api.mk release-test

release-test: user-rpc-dev social-rpc-dev user-api-dev social-api-dev

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh
