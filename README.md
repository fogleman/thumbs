# Thumbs

Go binary that watches a folder for images and generates thumbnails of them.

### Install

    go get github.com/fogleman/thumbs

### Run

The app will run forever, watching the `src` folder for images.

    thumbs

Probably, you would use something like [supervisor](http://supervisord.org/) to
launch and monitor the `thumbs` process.

### Arguments

    thumbs -src IMAGE_FOLDER -dst THUMB_FOLDER -w MAX_WIDTH -h MAX_HEIGHT -q JPG_QUALITY

All arguments are optional. See the defaults below.

| Flag | Default | Description |
| --- | --- | --- |
| -src | `.` | directory to watch for images |
| -dst | `thumbs` | directory to place thumbnails |
| -w | 1024 | max thumbnail width |
| -h | 1024 | max thumbnail height |
| -q | 95 | thumbnail jpeg quality |
