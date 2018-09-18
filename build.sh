echo Windows x86
GOOS=windows GOARCH=amd64 govvv build -o out/Wonky-Shell-win32.exe
echo Windows x86_64
GOOS=windows GOARCH=386   govvv build -o out/Wonky-Shell-win64.exe
echo OSX x86_64
GOOS=darwin  GOARCH=amd64 govvv build -o out/Wonky-Shell-mac64
echo Linux x86_64
GOOS=linux   GOARCH=amd64 govvv build -o out/Wonky-Shell-lin64
echo Linux x86
GOOS=linux   GOARCH=386   govvv build -o out/Wonky-Shell-lin32