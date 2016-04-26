# -*- mode: ruby -*-
# vi: set ft=ruby :

# This is a dev server to develop Essh.

Vagrant.configure(2) do |config|

  config.vm.box = "https://github.com/kohkimakimoto/vagrantboxes/releases/download/centos7.1/centos-7.1.box"
  config.vm.hostname = "essh-dev-server"
  config.vm.network :private_network, ip:"192.168.56.13"
  config.vm.provision :shell, :inline => <<-EOT
    # add user
    curl -sL https://raw.githubusercontent.com/kohkimakimoto/scripts/master/useradd.sh | bash -s -- --username kohkimakimoto --uid 1001 --sudoers

    # additional yum repositories.
    yum install -y \
      epel-release
    rpm -ivh http://rpms.famillecollet.com/enterprise/remi-release-7.rpm

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

    yum update  -y

    # selinux off
    setenforce 0
    sed -i -e "s/SELINUX=enforcing/SELINUX=disabled/g" /etc/selinux/config

    # configure sshd config
    sed -i -e "s/^PasswordAuthentication yes/PasswordAuthentication no/g" /etc/ssh/sshd_config

    # reboot
    reboot
  EOT

  config.vm.provider :virtualbox do |vb|
    vb.gui = false
    vb.customize ["modifyvm", :id, "--memory", "1024"]
  end

end
