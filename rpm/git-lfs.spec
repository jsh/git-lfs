# keep rpm from whining
%define debug_package %{nil}

Name:           git-lfs
Version:        0.5.1
Release:        1%{?dist}
Summary:        Git LFS is a command line extension and specification for managing large files with Git.

License:        MIT
URL:            https://github.com/jsh/%{name}
Source0:        https://github.com/jsh/%{name}/archive/%{version}.tar.gz

BuildRequires:  bison
BuildRequires:  git
BuildRequires:  golang
BuildRequires:  make
BuildRequires:  man
BuildRequires:  ruby
BuildRequires:  ruby-devel
#Requires:       git

%description
By design, every git repository contains every version of every file.
But for some types of projects, this is not reasonable or even practical.
Multiple revisions of a large file take up space quickly,
slowing down repository operations and making fetches unwieldy.

Git LFS overcomes this limitation by storing the metadata for large files
in Git and syncing the file contents to a configurable Git LFS server

%prep
%setup -q


%build
make -f rpm/Makefile git-lfs man

%check
script/test

%install
#rm -rf $RPM_BUILD_ROOT
%make_install -f rpm/Makefile

%files
%defattr(-,root,root)
%{_bindir}/*
%{_mandir}/man1/*

%doc LICENSE README.md

%changelog
* Mon May 25 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - 0.5.1-1
- Initial RPM release

