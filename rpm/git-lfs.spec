Name:           git-lfs
Version:        0.5.1
Release:        1%{?dist}
Summary:        Git LFS is a command line extension and specification for managing large files with Git.

License:        MIT
URL:            https://github.com/github/git-lfs
Source0:        https://github.com/github/git-lfs/archive/master.zip

BuildRequires:  bison
BuildRequires:  git
BuildRequires:  golang
BuildRequires:  make
BuildRequires:  man
BuildRequires:  ruby
BuildRequires:  ruby-devel
Requires:       

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
%configure
make %{?_smp_mflags}

%check
script/test

%install
rm -rf $RPM_BUILD_ROOT
%make_install


%files
%doc



%changelog
