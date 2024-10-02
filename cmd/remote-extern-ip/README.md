# Remote ip service

## Config for nginx

```
location /api/ip/ {
                try_files $uri @proxy_to_remote_ip_service;
        }
        location @proxy_to_remote_ip_service {
              proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
              proxy_set_header X-Forwarded-Proto $scheme;
              proxy_set_header Host $http_host;
              proxy_set_header X-Real-IP $remote_addr;
              proxy_redirect off;
              proxy_pass http://127.0.0.1:8189;
        }
```

## Build backend
```go build -o /usr/bin/remote-ip ./cmd/remote-extern-ip/main.go```

## App control by linux system
Copy file remote-ip.service.example to ```/etc/systemd/system/remote-ip.service``` and fix environments. <br />
Enabled service ```sudo systemctl enable remote-ip.service``` <br />
Start service ```sudo systemctl start remote-ip.service``` <br />
reStart service ```sudo systemctl restart remote-ip.service``` <br />
Check status ```sudo systemctl -l status remote-ip.service``` <br />
Reload systemd daemon after fixes ```sudo systemctl daemon-reload``` <br />
Show log ```journalctl -u remote-ip.service --no-pager | tail -10``` <br />
