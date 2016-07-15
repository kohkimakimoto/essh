# -*- mode: ruby -*-
# vi: set ft=ruby :

# This is a dev server to develop Essh.

Vagrant.configure(2) do |config|
  config.vm.box = "kohkimakimoto/centos-7"
  config.vm.network :private_network, type: "dhcp"
  config.vm.provision :shell, :inline => <<-EOT
    # additional yum repositories.
    yum install -y epel-release

    # install packages
    yum install -y \
      git \
      vim-enhanced \
      net-tools \
      telnet \
      ntp \
      chkconfig \
      tree \
      traceroute \
      httpd \
      httpd-devel \
      php \
      php-devel
  EOT

  config.vm.provider :virtualbox do |vb|
    vb.gui = false
    vb.customize ["modifyvm", :id, "--memory", "1024"]
  end
end
