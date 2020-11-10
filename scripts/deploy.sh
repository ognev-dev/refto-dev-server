#!/bin/bash

# this is straightforward deploy to remote that does the job for me for now.
# probably later I'll need to add a hook to deploy on merge

# how to use:
# 1. copy this file to ../bin or elsewhere you like (../bin is convenient because it is in project dir and out of git)
# 2. set variables that is just right below
# 3. run this script

# remote's requirements:
# 1. supervisor must be installed and config for supervisor to run server must be created
#     (look at ../scripts/supervisor.conf for example)

# set variables required for deploy

# server's host & port
serverHost=refto.dev
serverPort=443
# local project's dir to deploy from
projectDir=
# remote project's dir to deploy to
remoteProjectDir=
# remote's user
remoteUser=
# remote's addr
remoreAddr=
# database conf
dbAddr=localhost:5432
dbName=postgres
dbUser=postgres
dbPassword=postgres
# data repository to clone data from
dataRepo=https://github.com/refto/data.git
# This is the same secret you set on GitHub's push hook of the repo above
dataPushedHookSecret=SomeSecretToSignHooks

# necessary checks
if [ -z "$remoreAddr"  ]
then
  echo "Unable to deploy: Remote addr is not set"
  exit
fi
if [ -z "$remoteUser"  ]
then
  echo "Unable to deploy: Remote user is not set"
  exit
fi
if [ -z "$dbAddr"  ]
then
  echo "Unable to deploy: Database addr is not set"
  exit
fi
if [ -z "$dbName"  ]
then
  echo "Unable to deploy: Database name is not set"
  exit
fi
if [ -z "$dbUser"  ]
then
  echo "Unable to deploy: Database user is not set"
  exit
fi
if [ -z "$dbPassword"  ]
then
  echo "Unable to deploy: Database password is not set"
  exit
fi
if [ -z "$serverHost"  ]
then
  echo "Unable to deploy: Server's host is not set"
  exit
fi
cd $projectDir || exit

# lint
/home/vo/go/bin/golangci-lint run || { echo "Fix above errors, young apprentice"; exit; }

# test api
cd ./server/test || exit
echo "Testing API..."
go test || { echo "Unable to deploy: tests failed"; exit; }

echo "Compiling API server..."
cd $projectDir || exit
go build -ldflags "-s -w" -o ./bin/refto-server cmd/server/main.go || exit

echo "Compiling CLI..."
go build -ldflags "-s -w" -o ./bin/refto-cli cmd/cli/main.go || exit

echo "Making archive of static data..."
tar -czf refto-static.tar.gz ./web

echo "Copying files to remote (${remoteUser}@${remoreAddr})..."
scp ./bin/refto-server ./bin/refto-cli refto-static.tar.gz ${remoteUser}@${remoreAddr}:~/ || exit

echo "Setting up server on remote..."
ssh -T ${remoteUser}@${remoreAddr} << EOF
echo "Stopping supervisor..."
sudo service supervisor stop || exit

echo "Extracting static data..."
tar -xzf ~/refto-static.tar.gz -C $remoteProjectDir || exit
rm -f ~/refto-static.tar.gz || exit

echo "Moving binaries..."
mv ~/refto-server $remoteProjectDir/server || exit
mv ~/refto-cli $remoteProjectDir/cli || exit

echo "Writing config..."
/bin/cat <<EOM > $remoteProjectDir/.config.yaml
app_env: release

db:
  addr: $dbAddr
  user: $dbUser
  password: $dbPassword
  name: $dbName
  log_queries: true

server:
  host: $serverHost
  port: $serverPort
  api_base_path: api
  static:
    local_path: "./web"
    web_path: "/~/"

github:
  data_repo: $dataRepo
  data_pushed_hook_secret: $dataPushedHookSecret
  data_warden:
    app_id: 1
    install_id: 1
    pem_path: "private-key.pem"

dir:
  data: "$remoteProjectDir/data/"
  logs: ""
EOM

echo "Migrating database..."
cd $remoteProjectDir || exit
./cli migrate || exit

echo "Starting supervisor..."
sudo service supervisor start || exit
EOF

# cleanup
rm -f refto-static.tar.gz || exit
# not removing binaries they might be useful for local use

echo "Server deployed!"
