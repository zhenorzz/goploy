<p align=center>
    <img src="https://raw.githubusercontent.com/zhenorzz/goploy/master/logo.png" alt="logo" title="logo" />
</p>

<p align="center">
  <a href="#">
      <img src="https://img.shields.io/badge/readme%20style-standard-brightgreen.svg" alt="README">
  </a>
  <a href="#">
      <img src="https://img.shields.io/badge/give%20me-a%20star-green.svg" alt="STAR">
  </a>
  <a href="../LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="LICENSE">
  </a>
</p>

Name: go + deploy

A web deployment system tool!

Support all kinds of code release and rollback, which can be done through the web with one click!

Complete installation instructions, no difficulty in getting started!

[Dockerfile](https://github.com/zhenorzz/goploy/blob/master/docker/Dockerfile)

# How to use?

Run it with all needed parameter:

```console
docker run -it \
    --name=goploy \
    --env=DB_HOST=mysql \
    --env=DB_USER=root \
    --env=DB_USER_PASSWORD=YourPasswordHere \
    --env=DB_NAME=goploy \
    -v /path/to/.ssh:/root/.ssh \
    -v /path/to/hosts:/etc/hosts \
    -v /path/to/repository:/opt/goploy/repository \
    -p 3000:80 \
    zhenorzz/goploy[:tag]
```

That's it.

## Mount

The rsync and ssh need id_rsa to connect server, so you have to mount your ssh path to container using `-v /path/to/.ssh:/root/.ssh`.

If you want to access your own domain in container, you should mount your hosts file to container using `-v /path/to/hosts:/etc/hosts`.

*Notice that the filename is case-sensitive. If your readme is called `readme.md` you have to mount the file directly, not the directory*

# Additional Information
[Windows id_rsa permission denied, click this for help.](https://stackoverflow.com/questions/9270734/ssh-permissions-are-too-open-error)
 
# Updates and updating

To update your setup simply pull the newest image version from docker hub and run it.

# License

View [license information](https://github.com/zhenorzz/goploy/blob/master/LICENSE) for the software contained in this image.

Everything in [this repository](https://github.com/zhenorzz/goploy) is published under GPLv3.

# User Feedback

## Issues

If you have any problems with or questions about this image, please contact us through a [GitHub issue](https://github.com/zhenorzz/goploy/issues).

## Contributing

You are invited to contribute new features, fixes, or updates, large or small; I'm always thrilled to receive pull requests.
