Name:           zssh
Version:        0.8.0
Release:        1.el%{rhel}
Summary:        zssh is an extended ssh command.

Group:          Development/Tools
License:        MIT
Source0:        %{name}_linux_amd64.zip
Source1:        %{name}.config.lua
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
zssh is an extended ssh command.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp zssh %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/%{name}
cp %{SOURCE1} %{buildroot}/%{_sysconfdir}/%{name}/config.lua

%pre

%post

%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%dir %attr(755, root, root) %{_sysconfdir}/%{name}
%attr(644, root, root) %{_sysconfdir}/%{name}/config.lua
%attr(755, root, root) %{_bindir}/zssh

%doc
