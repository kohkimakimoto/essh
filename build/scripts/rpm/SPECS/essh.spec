Name:           %{_product_name}
Version:        %{_product_version}

Release:        1.el%{_rhel_version}
Summary:        Essh is an extended ssh command.
Group:          Development/Tools
License:        MIT
Source0:        %{name}_linux_amd64.zip
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
essh is an extended ssh command.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp %{name} %{buildroot}/%{_bindir}

%pre

%post

%preun

%clean
rm -rf %{buildroot}


%files
%defattr(-,root,root,-)
%attr(755, root, root) %{_bindir}/%{name}

%doc
