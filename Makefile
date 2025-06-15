build:
	docker build -t registry.cn-guangzhou.aliyuncs.com/scutrobot/rm-announce:latest --platform linux/amd64 .

push:
	docker push registry.cn-guangzhou.aliyuncs.com/scutrobot/rm-announce:latest
