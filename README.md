# microinformer
small app - web informer with go backend for displayed by small monitor on the microPC
Be careful! Do not use without a balancer, since this project is intended for local servers, and it can run any applications on it, which will be unsafe for public servers

clone project ```git clone git@github.com:resssoft/microinformer.git```
cd project root folder (for example: ```/home/pi/apps/microinformer```)

## Build backend
```go build -o /usr/bin/microinformer ./cmd/main.go```

## App control by linux system
Copy file microinformer.service.example to ```/etc/systemd/system/microinformer.service``` and fix environments. <br />
Enabled service ```systemctl enable microinformer.service``` <br />
Start service ```systemctl start microinformer.service``` <br />
Check status ```systemctl -l status microinformer.service``` <br />
Reload systemd daemon after fixes ```systemctl daemon-reload``` <br />
Show log ```journalctl -u microinformer.service --no-pager | tail -10``` <br />

## prepare frontend
```sudo makedir /opt/microinformer/```
```sudo cp -r ./frontend/```
```sudo chmod -R 755 /opt/microinformer/```

## configure for view
for start by default page on the micro PC
``` nano $HOME/.config/lxsession/LXDE-pi/autostart```

and append ```@chromium-browser -e --kiosk -a http://127.0.0.1:8081/page/index.html```

### build after configured
```cd /home/pi/apps/microinformer && git pull && go build -o /usr/bin/microinformer ./cmd/main.go && cp -r ./frontend /opt/microinformer/ && systemctl restart microinformer.service```
