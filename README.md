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

Update:

```
sudo nano /etc/systemd/system/getty@tty1.service.d/autologin.conf
```

Contents should be:

```
[Service]
ExecStart=
ExecStart=-/sbin/agetty -o '-p -f -- \\u' --noclear --autologin <username> %I $TERM
```

Enable and start the service:

```
sudo systemctl enable getty@tty1.service
sudo systemctl start getty@tty1.service
```

The user should autologin now and when the user logs in, pulseaudio daemon under the user should also start.
