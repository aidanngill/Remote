echo "[%] Building for all available systems..."

dpkg -s mingw-w64 > /dev/null
if [ $? -eq 1 ]
then
    echo "[!] MinGW was not found on your computer, installing now..."
    apt install mingw-w64
fi

if [ -d "bin" ]
then
    rm -rf bin/*
else
    mkdir bin
fi

for p in windows linux darwin; do
    export CC=gcc
    export CXX=g++

    export GOOS=$p
    export GOARCH=amd64

    export CGO_ENABLED=1

    LDFLAGS="-s -w"
    OUTFILE="bin/remote-$p"

    if [ $p == windows ]
    then
        export CC=x86_64-w64-mingw32-gcc
        export CXX=x86_64-w64-mingw32-g++

        LDFLAGS="$LDFLAGS -H=windowsgui"
        OUTFILE="$OUTFILE.exe"
    fi

    go build -ldflags "$LDFLAGS" -o $OUTFILE
    echo "[%] Finished build for '$p' (exit code $?)"
done