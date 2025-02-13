# Level DB copy tool written in go

## Introduction

This tool helps copy only the missing data from `n` level-DB directories from the parent source
directory to the level-DB directories in the destination parent directory. Let's suppose we have 
the following directory structure:

```bash
- src             # source parent directory
   |------ A      # directory, holding files for the DB called A
   |------ B      # directory, holding files for the DB called B
   +------ C      # directory, holding files for the DB called C
   
- dest            # destination parent directory
   |------ A      # directory, holding files for the DB called A
   |------ B      # directory, holding files for the DB called B
   +------ D      # directory, holding files for the DB called D
```

The tool will first read the sub-directories from `src` and `dest` and create the set intersection.
In this case, the directories `{A, B}` will be processed, ignoring directories C and D. 
Then, it will open in an arbitrary order the `{A, B}` DBs on both source and destination locations.
For each DB from the source location, it will read all existing keys & values and will only write the 
keys & values that are missing from the destination location. Everything that exists on the destination DB
will not be re-written or replaced.

## Install
To be used, the application will first need to be compiled. Go v1.20.7 is the recommended version
of the Golang compiler. On an Ubuntu Linux distribution, it can be easily installed using the following snippet:

```bash
GO_LATEST_TESTED="go1.20.7"
ARCH=$(dpkg --print-architecture)
wget https://dl.google.com/go/${GO_LATEST_TESTED}.linux-${ARCH}.tar.gz
sudo tar -C /usr/local -xzf ${GO_LATEST_TESTED}.linux-${ARCH}.tar.gz
rm ${GO_LATEST_TESTED}.linux-${ARCH}.tar.gz

sudo apt-get update
sudo apt-get upgrade
sudo apt install build-essential
sudo apt install git rsync curl zip unzip jq gcc wget

echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile

source ~/.profile
```

Test the Golang installation by typing
```bash
go version
```
It should output a text like: `go version go1.20.7 linux/amd64`

Then, it's time to clone the repo & compile the app. The following snippet does this:
```bash
cd
git clone https://github.com/iulianpascalau/level-db-copy-go

cd level-db-copy-go/cmd/level-db-copy
go build
```

No errors should appear and a new binary should be created in the same directory with the 
`main.go` file.

## Usage
After the compiling is done, the tool is ready to be used. The application can be easily 
configured by using the binary flags. For the complete list of the available flags, just type:

```bash
cd ~/level-db-copy-go/cmd/level-db-copy
./level-db-copy -h
```

The main flags to be used are the `--source` and `--destination` that specify the parent source and
destination directories. If we want to execute the example from the **Introduction** section, 
we can do that by typing:

```bash
./level-db-copy --source /path/to/src --destination /path/to/dest
```

A text like the following section will appear:
```
INFO [2025-02-13 17:56:15.117]   Common directories between the source and destination parent paths sub-directories = A, B 
INFO [2025-02-13 17:56:15.117]   now processing sub-directory             name = A overall progress = 1/2 
INFO [2025-02-13 17:56:15.117]   successfully processed DB                name = A missing info added = 2 
INFO [2025-02-13 17:56:15.117]   now processing sub-directory             name = B overall progress = 2/2 
INFO [2025-02-13 17:56:15.117]   successfully processed DB                name = B missing info added = 2 
```
