user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

release-test: user-rpc-dev

install-server:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh