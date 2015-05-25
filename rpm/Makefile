SHELL := /bin/bash

check: git-lfs
	./script/test

git-lfs:
	./script/bootstrap
	./git-lfs init

install: git-lfs man
	install -D bin/git-lfs ${RPM_BUILD_ROOT}/usr/bin/git-lfs
	install -d ${RPM_BUILD_ROOT}/usr/share/man/man1
	install man/*.1 ${RPM_BUILD_ROOT}/usr/share/man/man1

man:
	./script/man

.PHONY: check git-lfs install
