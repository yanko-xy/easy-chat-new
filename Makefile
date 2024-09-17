user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user-api.mk release-test

release-test: user-rpc-dev user-api-dev

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh