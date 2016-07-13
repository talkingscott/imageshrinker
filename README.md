# imageshrinker
Simple golang webserver that shrinks images

The IMAGEROOT environment variable sets the filesystem root where images are stored.  For a request to the server with path http://server:8080/path, the image is expected to exist at $IMAGEROOT/path in the filesystem.

The current incarnation always shrinks images to half width, maintaining the aspect ratio.  It also converts images to jpeg.  It can read images in gif, jpeg and png formats.
