#!/bin/bash

if (( $# < 2 )); then
    echo "Usage: sudo bash newUser.sh server-nickname port1 port2 port3..."
    echo "e.g.   sudo bash newUser.sh some-services 8000 1234"
    exit 1
fi

if (( $EUID != 0 )); then
    echo "Please run as root to spin new user"
    exit
fi

fwds=""
permitted=""
nick="$1"
user=$nick`cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1`
extip=`curl ifconfig.co || wget -qO- ifconfig.co || exit 0`

sudo useradd -m $user
sudo -Hu $user bash -c "mkdir ~/.ssh"
sudo -Hu $user bash -c "cd ~/.ssh && ssh-keygen -f id_rsa -t ed25519 -N ''"
sudo -Hu $user bash -c "mv ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys"

pub=""
priv=`cat /home/$user/.ssh/id_rsa`
privEsc=`printf %q "$priv"`

for port in "${@:2}"
do
 permitted="127.0.0.1:$port $permitted"
 fwds="\"$port:127.0.0.1:$port\", $fwds"
done

# echo "${permitted::-1}"
# echo "------"
# echo "${fwds::-2}"
# exit 0

echo "Match User "$user"
  AllowTcpForwarding yes
  X11Forwarding no
  PermitTunnel no
  GatewayPorts no
  AllowAgentForwarding no
  PermitOpen "${permitted::-1}"
  ForceCommand echo 'This account can only be used for port fwd'" >> /etc/ssh/sshd_config

sudo systemctl restart ssh

for p in /etc/ssh/*.pub; do 
    if [[ $(cat $p | wc -w) == 3 ]]; then
        pub="\"$(cat $p)\",$pub"
    fi
done

cat <<EOF > bao.conf
{
    "nickname": "$nick",
    "username": "$user",
    "addr": "$extip:22",
    "forwards": [ ${fwds::-2} ],
    "privkey" : "${privEsc:2:-1}",
    "checksums": [ ${pub::-1} ]
}
EOF

echo "all done! conf file is called bao.conf"