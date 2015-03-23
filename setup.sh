#cloud-config
users:
  - name: hazel
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh-authorized-keys:
      - asd
packages:
  - git
  - bzr
  - mercurial
write_files:
  - path: /hazel/env.sh
    owner: hazel
    content: |
      export HAZEL_API_URL=http://api.hazel.merd.io
  - path: /home/hazel/.bash_aliases
    owner: hazel
    content: |
      red='\[\e[0;31m\]' # Red
      grn='\[\e[0;32m\]' # Green
      ylw='\[\e[0;33m\]' # Yellow
      pur='\[\e[0;35m\]' # Purple
      wht='\[\e[0;37m\]' # White
      rst='\[\e[0m\]'    # Text Reset
      export PS1="$pur\h $wht- $ylw\W\n\$$rst "
      export PATH="$PATH:/hazel/go/bin"
      export GOROOT=/hazel/go
      export GOPATH=/hazel/gopath
      export GOBIN=/home/hazel/bin
runcmd:
  - sed -i -e '/^Port/s/^.*$/Port 2222/' /etc/ssh/sshd_config
  - sed -i -e '/^PermitRootLogin/s/^.*$/PermitRootLogin no/' /etc/ssh/sshd_config
  - sed -i -e '$aAllowUsers hazel' /etc/ssh/sshd_config
  - restart ssh
  - sudo mkdir /hazel && sudo chown -R hazel /hazel && sudo chgrp -R hazel /hazel
  - sudo su hazel -c "mkdir /home/hazel/bin"
  - sudo su hazel -c "mkdir /hazel/gopath"
  - curl -o /hazel/go.tar.gz https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
  - cd /hazel && tar xzf go.tar.gz && rm /hazel/go.tar.gz
