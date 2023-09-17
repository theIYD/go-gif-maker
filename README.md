# gif-maker

Using a simple command, create a GIF from a video in no time. :rocket:

### How does it work ?

The CLI creates a GIF from a video using `ffmpeg`. The video can be a remote URL resource OR an absolute path to a video file on your local machine.

**Options**

- `-start: (HH:MM:SS)` - The time at which the trim should start cutting the video.
- `-end: (HH:MM:SS)` - The time at which the trim should end cutting the video. 
- `-path:` - Remote URL / Absolute path to a video file.
- `-out:` - Output path for the generation of the GIF.
- `-h:` - See usage of the options.

### Installation
```bash
$ brew tap theIYD/tools
$ brew install gif-maker
```

### Usage
```bash
$ gif-maker -start=00:00:23 -end=00:00:30 -path=input.mp4

// use a youtube video
$ gif-maker -start=00:00:23 -end=00:00:30 -path=https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Development

- Clone the repository
```bash
$ git clone https://github.com/theIYD/go-gif-maker 
```

- Build
```bash
$ make build
```

- Run
```bash
$ ./bin/gif-maker -start=HH:MM:SS -end=HH:MM:SS -path=(remote url / absolute path)
```

- Compile for other distributions
```bash
$ make compile
```

### License
 The project is licensed under <a href="https://github.com/theIYD/go-gif-maker/LICENSE">MIT</a>