docker build -t rootfsimage .
id=$(docker create rootfsimage true) # id was cd851ce43a403 when the image was created
mkdir -p myplugin/rootfs
rm rootfs -rf
docker export "$id" | tar -x -C myplugin/rootfs
mv myplugin/rootfs rootfs
docker rm -vf "$id"
docker rmi rootfsimage
