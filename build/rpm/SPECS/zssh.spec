Name:           zssh-cli
Version:        0.10.0
Release:        1.el%{rhel}
Summary:        zssh is an extended ssh command.

Group:          Development/Tools
License:        MIT
Source0:        zssh_linux_amd64.zip
Source1:        zssh.config.lua
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
zssh is an extended ssh command.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp zssh %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/zssh
cp %{SOURCE1} %{buildroot}/%{_sysconfdir}/zssh/config.lua

%pre

%post

%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%dir %attr(755, root, root) %{_sysconfdir}/zssh
%attr(644, root, root) %{_sysconfdir}/zssh/config.lua
%attr(755, root, root) %{_bindir}/zssh

%doc
