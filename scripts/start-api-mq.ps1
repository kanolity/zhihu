Write-Host " 正在启动 API 服务..."

Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\article\api'; go run article.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\chat\api'; go run chat.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\follow\api'; go run follow.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\like\api'; go run like.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\message\api'; go run message.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\qa\api'; go run qa.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\reply\api'; go run reply.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\tag\api'; go run tag.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\user\api'; go run user.go" -WindowStyle Normal

Write-Host " 正在启动 MQ 服务..."

Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\article\mq'; go run main.go" -WindowStyle Normal
Start-Process powershell -ArgumentList "cd 'C:\go_work\src\go_code\zhihu\application\like\mq'; go run main.go" -WindowStyle Normal