# keep rpm from whiningG
%define debug_package %{nil}

%global commit 096b4b02cb768118ec853c175740031f8817b176
%global shortcommit %(c=%{commit}; echo ${c:0:7})

Name:           git-lfs
Version:        0.5.1
Release:        2%{?dist}
Summary:        Git command-line extension and specification for managing large files

License:        MIT
URL:            https://github.com/jsh/%{name}
Source0:        https://github.com/jsh/%{name}/archive/%{commit}/%{name}-%{commit}.tar.gz

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
%setup -qn %{name}-%{commit}


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
* Mon Jun 01 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - %{version}-%{release}
- New RPM release

* Mon May 25 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - %{version}-%{release}
- Initial RPM release

