##
主要是监控服务，用于监控所有节点的状态。包括节点的运行状态，流量。

```
{
	"nodeinfo":{
		"name":"normal",
		"connects":10,
		"bind":"0.0.0.0",
		"port":"50000",
		"hash":""
	},
	"packages":90,
	"nextnode":[
		{
			"name":"extra1",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"56001",
			"hash":""
		},
		{
			"name":"extra2",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"56002",
			"hash":""
		}
	],
	"genesis": true
}
```
```
{
	"nodeinfo":{
		"name":"extra1",
		"connects":10,
		"bind":"0.0.0.0",
		"port":"56001",
		"hash":""
	},
	"packages":90,
	"nextnode":[
		{
			"name":"extra1",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"57001",
			"hash":""
		},
		{
			"name":"extra2",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"57002",
			"hash":""
		}
	],
	"genesis": false
}
```

```
{
	"nodeinfo":{
		"name":"extra2",
		"connects":10,
		"bind":"0.0.0.0",
		"port":"56002",
		"hash":""
	},
	"packages":90,
	"nextnode":[
		{
			"name":"extra1",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"58001",
			"hash":""
		},
		{
			"name":"extra2",
			"connects":10,
			"bind":"0.0.0.0",
			"port":"58002",
			"hash":""
		}
	],
	"genesis": false
}
```

### GET
```
{
    "code": 200,
    "data": {
        "2dc7b18859cbf00d0b6712c1006b9895": {
            "nodeinfo": {
                "name": "extra1",
                "connects": 10,
                "bind": "0.0.0.0",
                "port": "56002",
                "hash": "2dc7b18859cbf00d0b6712c1006b9895"
            },
            "packages": 90,
            "nextnode": [
                {
                    "name": "extra1",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "58001",
                    "hash": "49f2f8229af20a4903e2f5fc3230ed8d"
                },
                {
                    "name": "extra2",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "58002",
                    "hash": "32875f9a1db4f8b7002035949716dd77"
                }
            ]
        },
        "c1a773193f60c622c17e435109af6f8b": {
            "nodeinfo": {
                "name": "extra2",
                "connects": 10,
                "bind": "0.0.0.0",
                "port": "56001",
                "hash": "c1a773193f60c622c17e435109af6f8b"
            },
            "packages": 90,
            "nextnode": [
                {
                    "name": "extra1",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "57001",
                    "hash": "b00c57f1047fd0999fc6dec9121d45f1"
                },
                {
                    "name": "extra2",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "57002",
                    "hash": "791925fabd8831ca214bc8a9b6ec691f"
                }
            ]
        },
        "genesis": {
            "nodeinfo": {
                "name": "normal",
                "connects": 10,
                "bind": "0.0.0.0",
                "port": "50000",
                "hash": "3a37c58b6fba167efd8c8b518a757090"
            },
            "packages": 90,
            "nextnode": [
                {
                    "name": "extra1",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "56001",
                    "hash": "c1a773193f60c622c17e435109af6f8b"
                },
                {
                    "name": "extra2",
                    "connects": 10,
                    "bind": "0.0.0.0",
                    "port": "56002",
                    "hash": "2dc7b18859cbf00d0b6712c1006b9895"
                }
            ],
            "genesis": true
        }
    },
    "msg": "Success"
}
```
前端通过API请求数据，将其展示出来，以便查看状态。