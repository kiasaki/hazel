sudo apt-get -y update
sudo apt-get -y install curl

# Setup a proper node PPA
sudo curl -sL https://deb.nodesource.com/setup | sudo bash -

# Install requirements
sudo apt-get install -y -qq \
      nodejs

# node upstart
sudo cp /ops/upstart/nodejs.conf /etc/init/nodejs.conf

# consul config
cat > /etc/consul.d/node.json <<-EOF
{
  "service": {
    "name": "{{.App.Slug}}",
    "port": 8000,
    "tags": ["web", "{{.Build.Id}}"],
    "check": {
      "id": "http",
      "name": "HTTP service on port 8080",
      "http": "http://localhost:8080/ping",
      "interval": "10s",
      "timeout": "2s"
    }
  }
}
EOF
sudo cp /ops/upstart/consul_client.conf /etc/init/consul.conf
