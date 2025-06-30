Write-Host " 正在启动 RPC 服务..."

Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\article\rpc'; go run article.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\chat\rpc'; go run chat.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\follow\rpc'; go run follow.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\like\rpc'; go run like.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\message\rpc'; go run message.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\qa\rpc'; go run qa.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\reply\rpc'; go run reply.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\tag\rpc'; go run tag.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\user\rpc'; go run user.go" -WindowStyle Normal
