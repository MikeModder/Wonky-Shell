# ensure go-prompt and govvv are installed
echo Installing govv...
go get -u github.com/JoshuaDoes/govvv
echo now go-prompt and deps...
go get -u github.com/c-bata/go-prompt
go get -u github.com/mattn/go-tty
go get -u github.com/mattn/go-colorable

rm -rf out
mkdir out
echo win32
GOOS=windows GOARCH=amd64 govvv build -o out/Wonky-Shell-win32.exe
echo win64
GOOS=windows GOARCH=386   govvv build -o out/Wonky-Shell-win64.exe
echo mac64
GOOS=darwin  GOARCH=amd64 govvv build -o out/Wonky-Shell-mac64
echo linux64
GOOS=linux   GOARCH=amd64 govvv build -o out/Wonky-Shell-lin64
echo linux32
GOOS=linux   GOARCH=386   govvv build -o out/Wonky-Shell-lin32