rm -rf build
mkdir build
mkdir build/config
cp config/*.json build/config/
go build -o build/DoctraLisMiddleware.exe