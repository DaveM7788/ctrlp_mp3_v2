import os
import sys
import time

from mutagen.id3 import APIC, ID3, TALB, TPE1
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
        for file_name in files:
            if file_name.endswith(".mp3"):
                count += 1
                display_audio_metadata(os.path.join(root, file_name))
        for dir_name in dirs:
            pass

    end = time.perf_counter()
    diff = end - start
    print("execution time is " + str(diff))
    print("count is  " + str(count))


def display_audio_metadata(file_path):
    audio: MP3 = MP3(file_path, ID3=ID3)
    title, album, artist = "", "", ""
    if "TIT2" in audio.tags:
        title = audio.tags["TIT2"][0]
    if "TPE1" in audio.tags:
        artist = audio.tags["TPE1"][0]
    if "TALB" in audio.tags:
        album = audio.tags["TALB"][0]
    print("meta data found is title " + title + " artist " + artist + " album " + album)

    apic_frames = audio.tags.getall("APIC")
    output_img_path = "albumart/" + artist + "/" + album + "/" + title + ".png"
    with open(output_img_path, "wb") as f:
        f.write(apic_frames[0].data)


if __name__ == "__main__":
    check_dir()
