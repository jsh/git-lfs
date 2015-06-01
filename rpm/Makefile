# A makefile for the specfile to use, when building the git-lfs-package.
#

SHELL := /bin/bash

check: git-lfs
	./script/test

git-lfs:
	./script/bootstrap
	#./bin/git-lfs init

install: git-lfs man
	strip bin/git-lfs
	install -D bin/git-lfs ${RPM_BUILD_ROOT}/usr/bin/git-lfs
	install -d ${RPM_BUILD_ROOT}/usr/share/man/man1
	install -m 644 man/*.1 ${RPM_BUILD_ROOT}/usr/share/man/man1

man:
	./script/man

.PHONY: check git-lfs install
