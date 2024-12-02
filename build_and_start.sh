WORK_DIR=$(pwd)

echo "Generate certificates"
go run ${WORK_DIR}/pkg/cert/main.go dev dev ${WORK_DIR}/certs/public.crt ${WORK_DIR}/certs/private.key

echo "Init docker containers"
PRIVATE_KEY_PATH=${WORK_DIR}/certs/private.key PUBLIC_KEY_PATH=${WORK_DIR}/certs/public.crt docker-compose up -d

echo "Build server"
cd ${WORK_DIR}/cmd/server
go build -o ${WORK_DIR}/server -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'"
chmod 755 ${WORK_DIR}/server

echo "Build client"
cd ${WORK_DIR}/cmd/client
go build -o ${WORK_DIR}/client -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'"
chmod 755 ${WORK_DIR}/client

echo "Create client default config file"
cat <<EOL > "${WORK_DIR}/cfgclient.yaml"
address: "127.0.0.1:8080"
crypto_key_public_path: "${WORK_DIR}/certs/public.crt"
hash_key: "defaultHashKey"
secret_key: "strongDBKey2Ks5nM2J5JaI59PPEhL1x"
EOL
echo "File cfgclient.yaml successfully created"

echo "Start server with default params"
cd ${WORK_DIR}
${WORK_DIR}/server -d postgres://postgres:postgres@localhost:5432/postgres -a 127.0.0.1:8080 -crypto-key-public ${WORK_DIR}/certs/public.crt -crypto-key-private ${WORK_DIR}/certs/private.key -db-key strongDBKey2Ks5nM2J5JaI59PPEhL1x -j JWTKey -minio-endpoint 127.0.0.1:9000 -minio-secret minioadmin -minio-id minioadmin -k defaultHashKey

echo "Done"