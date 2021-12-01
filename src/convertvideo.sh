ffmpeg -i $1 -f mp4 -vcodec libx264 -preset fast -profile:v main -acodec aac $1.mp4 -hide_banner
