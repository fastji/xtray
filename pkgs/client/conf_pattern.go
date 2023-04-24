package client

var StreamStr string = `{
	"network": "ws",
	"security": "tls",
	"tlsSettings": {
		"disableSystemRoot": false
	},
	"wsSettings": {
		"path": "/current_time"
	},
	"xtlsSettings": {
		"disableSystemRoot": false
	}
}`

var XrayConfStr string = `{
    "api": {
        "services": [
            "ReflectionService",
            "HandlerService",
            "LoggerService",
            "StatsService"
        ],
        "tag": "QV2RAY_API"
    },
    "dns": {
        "servers": [
            "1.1.1.1",
            "8.8.8.8",
            "8.8.4.4"
        ]
    },
    "fakedns": {
        "ipPool": "198.18.0.0/15",
        "poolSize": 65535
    },
    "inbounds": [
        {
			"port": 1080, 
			"listen": "127.0.0.1",
			"protocol": "socks",
			"settings": {
			  "udp": true
			}
		}
    ],
    "log": {
        "loglevel": "warning"
    },
    "outbounds": [
        {
            "protocol": "vmess",
            "sendThrough": "0.0.0.0",
            "settings": %s,
            "streamSettings": %s,
            "tag": "PROXY"
        },
        {
            "protocol": "freedom",
            "sendThrough": "0.0.0.0",
            "settings": {
                "domainStrategy": "AsIs",
                "redirect": ":0"
            },
            "streamSettings": {
            },
            "tag": "DIRECT"
        },
        {
            "protocol": "blackhole",
            "sendThrough": "0.0.0.0",
            "settings": {
                "response": {
                    "type": "none"
                }
            },
            "streamSettings": {
            },
            "tag": "BLACKHOLE"
        }
    ],
    "policy": {
        "system": {
            "statsInboundDownlink": true,
            "statsInboundUplink": true,
            "statsOutboundDownlink": true,
            "statsOutboundUplink": true
        }
    },
    "routing": {
        "domainMatcher": "mph",
        "domainStrategy": "AsIs",
        "rules": [
            {
                "inboundTag": [
                    "QV2RAY_API_INBOUND"
                ],
                "outboundTag": "QV2RAY_API",
                "type": "field"
            },
            {
                "ip": [
                    "geoip:private"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            },
            {
                "ip": [
                    "geoip:cn"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            },
            {
                "domain": [
                    "geosite:cn"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            }
        ]
    },
    "stats": {
    }
}`
