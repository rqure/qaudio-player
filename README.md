# qaudio-player

The following steps must be performed on the machine that plays audio

Install the pulseaudio server
```bash
sudo apt install pulseaudio
```

Add a line to accept clients from docker

```bash
sudo nano /etc/pulse/default.pa

# Add the following lines
load-module module-native-protocol-tcp auth-ip-acl=172.17.0.0/16 auth-anonymous=1
load-module module-combine-sink sink_name=combined
set-default-sink combined
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

# Bluetooth Setup

Install bluetooth support for pulseaudio

```
sudo apt update
sudo apt install pulseaudio-module-bluetooth
```

Launch the bluetooth CLI

```
bluetoothctl
```

Turn on the agent, which will allow us to search for and pair with other bluetooth devices (in this case, our bluetooth speaker)

```
agent on
```

Start scanning for bluetooth devices

```
scan on
```

Once you have found the MAC address of the device you want to connect to, you can now proceed to pair with it

```
pair [XX:XX:XX:XX:XX:XX]
```

When you first pair a device, you will be immediately connected to it.

However, once you have gone out of range of the Raspberry Pi’s Bluetooth, you will need to re-connect the device by using the following command

```
connect [XX:XX:XX:XX:XX:XX]
```

If you don’t want to have to re-pair your device, then you can make use of the trust command.

```
trust [XX:XX:XX:XX:XX:XX]
```

