# fc
first controller

docker volume create --name mongodb

docker run -dti -p 27017:27017 -v mongodb:/data/db mongo:3.2

export FC_CONFIG="/home/mr/Documents/work_space/fc/bin/config.json"