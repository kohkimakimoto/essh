Name:           essh
Version:        0.29.0
Release:        1.el%{rhel}
Summary:        essh is an extended ssh command.

Group:          Development/Tools
License:        MIT
Source0:        essh_linux_amd64.zip
Source1:        essh.config.lua
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
essh is an extended ssh command.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp essh %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/essh
cp %{SOURCE1} %{buildroot}/%{_sysconfdir}/essh/config.lua

%pre

%post

%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%dir %attr(755, root, root) %{_sysconfdir}/essh
%attr(644, root, root) %{_sysconfdir}/essh/config.lua
%attr(755, root, root) %{_bindir}/essh

%doc
