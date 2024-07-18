docker build --rm -f Dockerfile -t devenv:ubuntu .
docker run --rm -it --network host -v `pwd`:/developer devenv:ubuntu
