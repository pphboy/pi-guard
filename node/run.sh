# 运行的命令

# Node运行
go-node -port 5703 -name nd3 -grpcPort 5704 -center nd1.pi.g:9981
# Center运行
go run . -ctrlDomain=ctrl.pi.g:7431
