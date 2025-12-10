var currentPlaylist = [];
var shufflePlaylist = [];
var tempPlaylist = [];
var audioElement;
var mouseDown = false;
var currentIndex = 0;
var repeat = false;
var shuffle = false;
var userLoggedIn;
var timer;

$(document).click(function (click) {
  var target = $(click.target);
  if (!target.hasClass("item") && !target.hasClass("optionsButton")) {
    hideOptionsMenu();
  }

  var modal = document.getElementById("ctrlpModal");
  if (modal.style.display == "block") {
    modal.style.display = "none";
  }
});

$(window).scroll(function () {
  hideOptionsMenu();
});

$(document).on("change", "select.playlist", function () {
  var select = $(this);
  var playlistId = $(this).val();
  var songId = $(this).prev(".songId").val();
  hideOptionsMenu();
  select.val("");
});

function formatTime(seconds) {
  var time = Math.round(seconds);
  var minutes = Math.floor(time / 60); // rounds down
  var seconds = time - minutes * 60;

  var extraZero;
  if (seconds < 10) {
    extraZero = "0";
  } else {
    extraZero = "";
  }

  return minutes + ":" + extraZero + seconds;
}

function updateTimeProgressBar(audio) {
  $(".progressTime.current").text(formatTime(audio.currentTime));
  var progress = (audio.currentTime / audio.duration) * 100;
  $(".playbackBar .progress").css("width", progress + "%");
}

function updateVolumeProgressBar(audio) {
  var volume = audio.volume * 100;
  $(".volumeBar .progress").css("width", volume + "%");
}

function playFirstSong() {
  setTrack(tempPlaylist[0], tempPlaylist, true);
}

function hideOptionsMenu() {
  // hide options menu on scroll
  var menu = $(".optionsMenu");
  if (menu.css("display") != "none") {
    menu.css("display", "none");
  }
}

function showOptionsMenu(button) {
  var songId = $(button).prevAll(".songId").val();
  var menu = $(".optionsMenu");
  var menuWidth = menu.width();
  menu.find(".songId").val(songId);
  var scrollTop = $(window).scrollTop(); // distance from top of window to top of document
  var elementOffset = $(button).offset().top; // distance from top of document
  var top = elementOffset - scrollTop;
  var left = $(button).position().left;
  menu.css({
    top: top + "px",
    left: left - menuWidth + "px",
    display: "inline",
  });
}

function Audio() {
  this.currentlyPlaying;
  this.audio = document.createElement("audio");

  this.audio.addEventListener("ended", function () {
    nextSong();
  });

  this.audio.addEventListener("canplay", function () {
    // 'this' refers to the object the event was called on (audio)
    var duration = formatTime(this.duration);
    $(".progressTime.remaining").text(duration);
  });

  this.audio.addEventListener("timeupdate", function () {
    if (this.duration) {
      updateTimeProgressBar(this);
    }
  });

  this.audio.addEventListener("volumechange", function () {
    updateVolumeProgressBar(this);
  });

  this.setTrack = function (track) {
    this.currentlyPlaying = track;
    this.audio.src = track.path;
  };

  this.play = function () {
    this.audio.play();
  };

  this.pause = function () {
    this.audio.pause();
  };

  this.setTime = function (seconds) {
    this.audio.currentTime = seconds;
  };
}
