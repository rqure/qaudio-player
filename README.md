# qaudio-player

The following steps must be performed on the machine that plays audio

Install the pulseaudio server
```bash
sudo apt install pulseaudio
```

Add a line to accept clients from docker

```bash
sudo nano /etc/pulse/default.pa

# Add the following line
load-module module-native-protocol-tcp auth-ip-acl=172.17.0.0/16 auth-anonymous=1
```

Restart the pulseaudio server

```bash
systemctl --user restart pulseaudio
```