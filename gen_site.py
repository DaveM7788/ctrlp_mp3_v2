import os
import sys
import time

from mutagen.id3 import ID3, TALB, TPE1
from mutagen.mp3 import MP3


def check_dir():
    start = time.perf_counter()
    n = len(sys.argv)
    if n < 1:
        print("You must pass in at least 1 argument. Passed " + str(n))
        print(
            "Argument 1 = Filepath to use for generating static site. \
            Filepath should contain music files."
        )
        quit()

    directory = sys.argv[1]

    count = 0
    for root, dirs, files in os.walk(directory):
        # print(f"Current directory: {root}")
        for file_name in files:
            # print(f"  File: {os.path.join(root, file_name)}")
            if file_name.endswith(".mp3"):
                count += 1
                display_audio_metadata(os.path.join(root, file_name))
        for dir_name in dirs:
            pass
            # print(f"  Subdirectory: {os.path.join(root, dir_name)}")

    end = time.perf_counter()
    diff = end - start
    print("execution time is " + str(diff))
    print("count is  " + str(count))


def display_audio_metadata(file_path):
    audio: MP3 = MP3(file_path)
    title, album, artist = "", "", ""
    if "TIT2" in audio.tags:
        title = audio.tags["TIT2"][0]
    if "TPE1" in audio.tags:
        artist = audio.tags["TPE1"][0]
    if "TALB" in audio.tags:
        album = audio.tags["TALB"][0]
    print("meta data found is title " + title + " artist " + artist + " album " + album)


if __name__ == "__main__":
    check_dir()
