# BBox

Beatboxer in Go

## Raspberry PI Setup

### OS

1. Download Raspian Lite: https://downloads.raspberrypi.org/raspbian_lite_latest
2. Flash `2017-07-05-raspbian-jessie-lite.zip` using Etcher
3. Remove/reinsert flash drive
4. Add `ssh` file:
```bash
touch /Volumes/boot/ssh
diskutil umount /Volumes/boot
```

### First Boot

```bash
ssh pi@raspberrypi.local
# password: raspberry

# change default password
passwd

# set quiet boot
sudo sed -i '${s/$/ quiet loglevel=1/}' /boot/cmdline.txt

# install packages
sudo apt-get update
sudo apt-get install -y git tmux vim dnsmasq hostapd

# set up wifi (note leading space to avoid bash history)
sudo tee --append /etc/wpa_supplicant/wpa_supplicant.conf > /dev/null << 'EOF'
network={
    ssid="<WIFI_SSID>"
    psk="<WIFI_PASSWORD>"
}
EOF

# set static IP address
sudo tee --append /etc/dhcpcd.conf > /dev/null << 'EOF'

# set static ip

interface eth0
static ip_address=192.168.1.141/24
static routers=192.168.1.1
static domain_name_servers=192.168.1.1

interface wlan0
static ip_address=192.168.1.142/24
static routers=192.168.1.1
static domain_name_servers=192.168.1.1
EOF

# reboot to connect over wifi
sudo shutdown -r now
```

```bash
# configure git
git config --global push.default simple
git config --global core.editor "vim"
git config --global user.email "you@example.com"
git config --global user.name "Your Name"

# disable services
sudo systemctl disable hciuart
sudo systemctl disable bluetooth
sudo systemctl disable plymouth

# remove unnecessary packages
sudo apt-get -y purge libx11-6 libgtk-3-common xkb-data lxde-icon-theme raspberrypi-artwork penguinspuzzle ntp plymouth*
sudo apt-get -y autoremove

sudo raspi-config nonint do_boot_behaviour B2 0
sudo raspi-config nonint do_boot_wait 1
sudo raspi-config nonint do_serial 1
```

## Code

```bash
wget https://dl.google.com/go/go1.10.3.linux-armv6l.tar.gz -O /tmp/go1.10.3.linux-armv6l.tar.gz
sudo tar -xzf /tmp/go1.10.3.linux-armv6l.tar.gz -C /usr/local

mkdir -p ~/code/go/src/github.com/siggy
git clone https://github.com/siggy/bbox.git ~/code/go/src/github.com/siggy/bbox
```

### portaudio

```bash
# OSX
brew install portaudio

# Raspbian
sudo apt-get install -y libasound-dev

wget http://portaudio.com/archives/pa_stable_v190600_20161030.tgz -O /tmp/pa_stable_v190600_20161030.tgz
cd /tmp
tar -xzf pa_stable_v190600_20161030.tgz
cd portaudio
./configure
make
sudo make install
sudo ldconfig
```

### rpi_ws281x

Beatboxer depends on a fork of (https://github.com/jgarff/rpi_ws281x). See that
repo for complete instructions.

```bash
cd ~/code/go/src/github.com/siggy/bbox
sudo cp rpi_ws281x/rpihw.h  /usr/local/include/
sudo cp rpi_ws281x/ws2811.h /usr/local/include/
sudo cp rpi_ws281x/pwm.h    /usr/local/include/

sudo cp rpi_ws281x/libws2811.a /usr/local/lib/

# osx
export CGO_CFLAGS="$CGO_CFLAGS -I/usr/local/include"
export CGO_LDFLAGS="$CGO_LDFLAGS -L/usr/local/lib"
```

## Env / bootup

```bash
# set bootup and shell env
cd ~/code/go/src/github.com/siggy/bbox
cp rpi/.local.bash ~/
source ~/.local.bash

cp rpi/bboxgo.sh ~/
sudo cp rpi/bbox.service /etc/systemd/system/bbox.service
sudo systemctl enable bbox

echo "[[ -s ${HOME}/.local.bash ]] && source ${HOME}/.local.bash" >> ~/.bashrc

# audio setup

# external sound card
sudo cp rpi/asound.conf /etc/

# *output of raspi-config after forcing audio to hdmi*
numid=3,iface=MIXER,name='Mic Playback Switch'
  ; type=BOOLEAN,access=rw------,values=1
  : values=on
# *also this might work*
amixer cset numid=3 2
# OR:
sudo raspi-config nonint do_audio 2

echo "blacklist snd_bcm2835" | sudo tee --append /etc/modprobe.d/snd-blacklist.conf

echo "hdmi_force_hotplug=1" | sudo tee --append /boot/config.txt
echo "hdmi_force_edid_audio=1" | sudo tee --append /boot/config.txt

# make usb audio card #0
sudo vi /lib/modprobe.d/aliases.conf
#options snd-usb-audio index=-2

# reboot

aplay -l
# ... should match the contents of asound.conf, and also:
sudo vi /usr/share/alsa/alsa.conf
defaults.ctl.card 0
defaults.pcm.card 0
```

## Build

```bash
go build cmd/beatboxer_noleds.go && \
  go build cmd/beatboxer_leds.go && \
  go build cmd/baux.go &&      \
  go build cmd/clear.go &&     \
  go build cmd/fishweb.go &&   \
  go build cmd/human.go &&      \
  go build cmd/leds.go &&      \

  go build cmd/amplitude.go && \
  go build cmd/aud.go &&       \
  go build cmd/crane.go &&     \
  go build cmd/crawler.go &&   \
  go build cmd/fish.go &&      \
  go build cmd/keys.go &&      \
  go build cmd/noleds.go &&    \
  go build cmd/record.go
```

## Run

All programs that use LEDs must be run with `sudo`.

```bash
sudo ./beatboxer # main program
sudo ./leds # led testing
sudo ./clear # clear LEDs
./noleds # beatboxer without LEDs (for testing without pi)
./aud # audio testing
./keys # keyboard test
```

## Stop bbox process

```bash
# the systemd way
sudo systemctl stop bbox

# send SIGINT to turn off LEDs
sudo kill -2 <PID>
```

## Check for voltage drop

```bash
vcgencmd get_throttled
```

## Editing SD card

Launch Ubuntu in VirtualBox

```bash
sudo mount /dev/sdb7 ~/usb
sudo umount /dev/sdb7
```

## Wifi access point

Based on:
https://frillip.com/using-your-raspberry-pi-3-as-a-wifi-access-point-with-hostapd/

```bash
sudo tee --append /etc/dhcpcd.conf > /dev/null <<'EOF'

# this must go above any `interface` line
denyinterfaces wlan0

# this must go below `interface wlan0`
nohook wpa_supplicant
EOF

sudo tee --append /etc/network/interfaces > /dev/null <<'EOF'

allow-hotplug wlan0
iface wlan0 inet static
    address 192.168.4.1
    netmask 255.255.255.0
    network 192.168.4.0
    broadcast 192.168.1.255
#    wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
EOF

sudo tee /etc/hostapd/hostapd.conf > /dev/null <<'EOF'
interface=wlan0
driver=nl80211
ssid=sigpi
hw_mode=g
channel=6
ieee80211n=1
wmm_enabled=0
ht_capab=[HT40][SHORT-GI-20][DSSS_CCK-40]
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_passphrase=showmethepi
# wpa_pairwise=TKIP
rsn_pairwise=CCMP
EOF

sudo tee /etc/default/hostapd > /dev/null <<'EOF'
DAEMON_CONF="/etc/hostapd/hostapd.conf"
EOF

sudo tee --append /etc/dnsmasq.conf > /dev/null <<'EOF'

interface=wlan0
listen-address=192.168.4.1
bind-interfaces
domain-needed
dhcp-range=192.168.4.2,192.168.4.100,255.255.255.0,24h
EOF

sudo tee --append /etc/sysctl.conf > /dev/null <<'EOF'

net.ipv4.ip_forward=1
EOF

sudo service dhcpcd restart
sudo systemctl start hostapd
sudo systemctl start dnsmasq

# reboot to connect over wifi
sudo shutdown -r now
```

### To re-enable internet wifi

Comment out from `/etc/dhcpcd.conf`:
```
# denyinterfaces wlan0
# nohook wpa_supplicant
```

Re-enable in `/etc/network/interfaces`:
```
allow-hotplug wlan0
iface wlan0 inet manual
    wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
```

sudo service dhcpcd restart
sudo ifdown wlan0; sudo ifup wlan0
sudo systemctl stop hostapd
sudo systemctl stop dnsmasq

### Mounting / syncing pi

```bash
# mount pi volume locally
alias pifs='umount /Volumes/pi; sudo rmdir /Volumes/pi; sudo mkdir /Volumes/pi; sudo chown sig:staff /Volumes/pi && sshfs pi@raspberrypi.local:/ /Volumes/pi -f'

# rsync local repo to pi
rsync -vr ~/code/go/src/github.com/siggy/bbox/.git /Volumes/pi/home/pi/code/go/src/github.com/siggy/bbox

# remove volume mount
alias pifsrm='umount /Volumes/pi; sudo rmdir /Volumes/pi; sudo mkdir /Volumes/pi; sudo chown sig:staff /Volumes/pi'
```

### Set restart crontab

```bash
sudo crontab -e
```

```
*/10 * *   *   *     sudo /sbin/shutdown -r now
```

```
sudo crontab -l
```

## Docs

```bash
jekyll serve -s docs
open http://127.0.0.1:4000/bbox
```

## Credits

- [wavs](wavs) courtesy of (http://99sounds.org/drum-samples/)
- [rpi_ws281x](rpi_ws281x) courtesy of (https://github.com/jgarff/rpi_ws281x)
