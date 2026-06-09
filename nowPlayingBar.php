<?php
$stmt = $con->prepare("SELECT id FROM songs ORDER BY RAND() LIMIT 10");
$stmt->execute();
$result = $stmt->get_result();
$stmt->close();
$resultArray = [];
while ($row = $result->fetch_assoc()) {
    array_push($resultArray, $row["id"]);
}
$jsonArray = json_encode($resultArray);
?>

<script>
    $(document).ready(function() {
        var newPlaylist = <?php echo $jsonArray; ?>;
        audioElement = new Audio();
        setTrack(newPlaylist[0], newPlaylist, false);
        updateVolumeProgressBar(audioElement.audio);

        $("#nowPlayingBarContainer").on("mousedown mousemove", function(e) {
            e.preventDefault();
        });

        $("#nowPlayingBarContainer").mousedown(function(e) {
            e.preventDefault();
        });

        $("#nowPlayingBarContainer").mousemove(function(e) {
            e.preventDefault();
        });



        $(".playbackBar .progressBar").mousedown(function() {
            mouseDown = true;
        });

        $(".playbackBar .progressBar").mousemove(function(e) {
            if (mouseDown) {
                // set time of song depending on position of mouse
                timeFromOffset(e, this);
            }
        });

        $(".playbackBar .progressBar").mouseup(function(e) {
            timeFromOffset(e, this);
        });

        $(".volumeBar .progressBar").mousedown(function() {
            mouseDown = true;
        });

        $(".volumeBar .progressBar").mousemove(function(e) {
            if (mouseDown) {
                // 'this' is .volumeBar .progress bar HTML element above
                var percentage = e.offsetX / $(this).width();
                if (percentage >= 0 && percentage <= 1) {
                    audioElement.audio.volume = percentage;
                }
            }
        });

        $(".volumeBar .progressBar").mouseup(function(e) {
            var percentage = e.offsetX / $(this).width();
            if (percentage >= 0 && percentage <= 1) {
                audioElement.audio.volume = percentage;
            }
        });

        $(document).mouseup(function() {
            mouseDown = false;
        });
    });

    function timeFromOffset(mouse, progressBar) {
        var percentage = mouse.offsetX / $(progressBar).width() * 100;
        var seconds = audioElement.audio.duration * (percentage / 100);
        audioElement.setTime(seconds);
    }

    function prevSong() {
        if (audioElement.audio.currentTime >= 3 || currentIndex == 0) {
            audioElement.setTime(0);
        }
        else {
            currentIndex--;
            setTrack(currentPlaylist[currentIndex], currentPlaylist, true);
        }
    }

    function nextSong() {
        if (repeat) {
            audioElement.setTime(0);
            playSong();
            return;
        }

        if (currentIndex == (currentPlaylist.length - 1)) {
            currentIndex = 0;
        }
        else {
            currentIndex++;
        }

        var trackToPlay = shuffle ? shufflePlaylist[currentIndex] : currentPlaylist[currentIndex];
        setTrack(trackToPlay, currentPlaylist, true);
    }

    function setRepeat() {
        repeat = !repeat;
        var imageName = (repeat) ? "repeat-purp.png" : "repeat-purp-inactive.png";
        $(".controlButton.repeat img").attr("src", "assets/images/icons/" + imageName);
    }

    function setMute() {
        audioElement.audio.muted = !audioElement.audio.muted;
        var imageName = (audioElement.audio.muted) ? "mute-purp.png" : "vol-purp.png";
        $(".controlButton.volume img").attr("src", "assets/images/icons/" + imageName);
    }

    function setShuffle() {
        shuffle = !shuffle;
        var imageName = (shuffle) ? "shuff-purp.png" : "shuff-purp-inactive.png";
        $(".controlButton.shuffle img").attr("src", "assets/images/icons/" + imageName);

        if (shuffle) {
            // randomize playlist
            shuffleArray(shufflePlaylist);
            // prevents next song in shuffle from being the currently playing
            currentIndex = shufflePlaylist.indexOf(audioElement.currentlyPlaying.id);
        }
        else {
            // shuffle turned off -> go back to original playlist
            currentIndex = currentPlaylist.indexOf(audioElement.currentlyPlaying.id);
        }
    }

    function shuffleArray(a) {
        var j, x, i;
        for (i = a.length; i > 0; i--) {
            j = Math.floor(Math.random() * i);
            x = a[i - 1];
            a[i - 1] = a[j];
            a[j] = x;
        }
    }

    function setTrack(trackId, newPlaylist, play) {
        if (newPlaylist != currentPlaylist) {
            currentPlaylist = newPlaylist;
            shufflePlaylist = currentPlaylist.slice();
            shuffleArray(shufflePlaylist);
        }

        if (shuffle) {
            currentIndex = shufflePlaylist.indexOf(trackId);
        } else {
            currentIndex = currentPlaylist.indexOf(trackId);
        }
        pauseSong();

        var csrfHash = document.getElementById("csrfHeader").value;

        $.post("includes/handlers/ajax/getSongJson.php", { songId: trackId, csrf: csrfHash }, function(data) {
            var track = JSON.parse(data);
            $(".trackName span").text(track.title);

            $.post("includes/handlers/ajax/getArtistJson.php", { artistId: track.artist, csrf: csrfHash }, function(data) {
                var artist = JSON.parse(data);
                $(".trackInfo .artistName span").text(artist.name);
                $(".trackInfo .artistName span").attr("onclick", "openPage('artist.php?id=" + artist.id + "')");
            });

            $.post("includes/handlers/ajax/getAlbumJson.php", { albumId: track.album, csrf: csrfHash }, function(data) {
                var album = JSON.parse(data);
                $(".content .albumLink img").attr("src", album.artworkPath);
                $(".content .albumLink img").attr("onclick", "openPage('album.php?id=" + album.id + "')");
                $(".trackInfo .trackName span").attr("onclick", "openPage('album.php?id=" + album.id + "')");
            });
            audioElement.setTrack(track);

            if (play) {
                playSong();
            }
        });
    }

    function playSong() {
        var csrfHash = document.getElementById("csrfHeader").value;
        if (audioElement.audio.currentTime == 0) {
            // update play count
            $.post("includes/handlers/ajax/updatePlays.php", { songId: audioElement.currentlyPlaying.id, csrf: csrfHash} );
        }
        $(".controlButton.play").hide();
        $(".controlButton.pause").show();
        audioElement.play();
    }

    function pauseSong() {
        $(".controlButton.play").show();
        $(".controlButton.pause").hide();
        audioElement.pause();
    }
</script>

<div id="nowPlayingBarContainer">
    <div id="nowPlayingBar">
        <div id="nowPlayingLeft">
            <div class="content">
                <span class="albumLink">
                    <img
                        src="assets/images/controlp_sq_jpg.jpg"
                        role="link"
                        tabindex="0"
                        class="albumArtwork"
                    />
                </span>
                <div class="trackInfo">
                    <span class="trackName">
                        <span role="link" tabindex="0">Track Name</span>
                    </span>
                    <span class="artistName">
                        <span role="link" tabindex="0">Artist Name</span>
                    </span>
                </div>
            </div>
        </div>
        <div id="nowPlayingCenter">
            <div class="content playerControls">
                <div class="buttons">
                    <button
                        class="controlButton shuffle"
                        title="Shuffle Button"
                        onclick="setShuffle()"
                    >
                        <img
                            src="assets/images/icons/shuff-purp-inactive.png"
                            alt="Shuffle"
                        />
                    </button>

                    <button
                        class="controlButton previous"
                        title="Previous Button"
                        onclick="prevSong()"
                    >
                        <img
                            src="assets/images/icons/prev-purp.png"
                            alt="Previous"
                        />
                    </button>

                    <button
                        class="controlButton play"
                        title="Play Button"
                        onclick="playSong()"
                    >
                        <img
                            src="assets/images/icons/play-purp.png"
                            alt="Play"
                        />
                    </button>

                    <button
                        class="controlButton pause"
                        title="Pause Button"
                        style="display: none"
                        onclick="pauseSong()"
                    >
                        <img
                            src="assets/images/icons/pause-purp.png"
                            alt="Pause"
                        />
                    </button>

                    <button
                        class="controlButton next"
                        title="Next Button"
                        onclick="nextSong()"
                    >
                        <img
                            src="assets/images/icons/next-purp.png"
                            alt="Next"
                        />
                    </button>

                    <button
                        class="controlButton repeat"
                        title="Repeat Button"
                        onclick="setRepeat()"
                    >
                        <img
                            src="assets/images/icons/repeat-purp-inactive.png"
                            alt="Repeat"
                        />
                    </button>
                </div>

                <div class="playbackBar">
                    <span class="progressTime current">0.00</span>

                    <div class="progressBar">
                        <div class="progressBarBg">
                            <div class="progress"></div>
                        </div>
                    </div>

                    <span class="progressTime remaining">0.00</span>
                </div>
            </div>
        </div>

        <div id="nowPlayingRight">
            <div class="volumeBar">
                <button
                    class="controlButton volume"
                    title="Volume Button"
                    onclick="setMute()"
                >
                    <img src="assets/images/icons/vol-purp.png" alt="Volume" />
                </button>

                <div class="progressBar">
                    <div class="progressBarBg">
                        <div class="progress"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
