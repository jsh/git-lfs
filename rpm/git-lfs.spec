# keep rpm from whiningG
%define debug_package %{nil}

# %%global commit 04aa69bc94568720252181a98ce76069af21be3b
# %%global shortcommit %%(c=%%{commit}; echo ${c:0:7})
%global commit rpm-kevin1

Name:           git-lfs
Version:        0.5.1
Release:        3%{?dist}
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
rm -rf $RPM_BUILD_ROOT
%make_install -f rpm/Makefile

%files
%defattr(-,root,root)
%{_bindir}/*
%attr(644,root,root) %{_mandir}/man1/*

%doc LICENSE README.md

%changelog
* Mon Jun 01 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - 0.5.1-3.centos
- Fix most rpmlint whines

* Mon Jun 01 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - 0.5.1-2.centos
- New RPM release

* Mon May 25 2015 Jeffrey S. Haemer <jeffrey.haemer@gmail.com> - 0.5.1-1.centos
- Initial RPM release

