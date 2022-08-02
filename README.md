

```md
Welcome to DECForum
===================

This is a forum where the conversation is backed by a git repository. Each message is a separate commit which is referenced by it’s hash. You can use markdown to write posts, and they will be rendered back to you as HTML. There is no login, posters are identified by the base32 address of the client which they use to connect to the server.
```

Build Dependencies:
-------------------

Node, NPM, Gulp.

```sh
sudo apt-get install nodejs npm gulp
```

### Docker Setup:

```sh
docker build -t eyedeekay/decforum https://github.com/eyedeekay/decforum.git
docker run \
    --detach \
    --name=decforum \
    --net=host \
    --restart=always \
    --volume=decforum:/var/lib/decforum \
    eyedeekay/decforum
echo "https://$(docker exec decforum cat gitforum.i2p.public.txt)"
```