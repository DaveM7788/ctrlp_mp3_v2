var posHighlight = 0;
var modal = null;

onkeydown = function (e) {
  if ((e.ctrlKey && e.key == "p") || (e.metaKey && e.key == "p")) {
    // show ctrlp modal
    e.preventDefault();
    modal = document.getElementById("ctrlpModal");
    modal.style.display = "block";
    var searchBox = document.getElementById("ctrlpInput");
    searchBox.value = ""; // clear out old values when doing new search
    searchBox.focus();
    posHighlight = 0;
    highlightResultCell(0);
  } else if (e.key == "Escape" || e.key == "Esc") {
    // close modals
    modal.style.display = "none";
    posHighlight = 0;
  } else if (e.key == "ArrowDown" || e.key == "Down") {
    if (modal.style.display == "block") {
      e.preventDefault();
      var lenCurrentResults = $("#ctrlpResultsList li").length;
      if (posHighlight < lenCurrentResults - 1) {
        posHighlight++;
        highlightResultCell(posHighlight);
      }
    }
  } else if (e.key == "ArrowUp" || e.key == "Up") {
    if (modal.style.display == "block") {
      e.preventDefault();
      if (posHighlight > 0) {
        posHighlight--;
        highlightResultCell(posHighlight);
      }
    }
  } else if (e.key == "Enter") {
    getIdAndPlay($("li.activeL").text());
  }
};

function getIdAndPlay(songToFind) {
  if (modal.style.display == "block") {
    var idxOfSong = dataSets.songs.indexOf(songToFind);
    var songIdForDB = dataSetsIds.songids[idxOfSong];

    var ctrlpPlaylist = genPlaylistFromOneId(songIdForDB, idxOfSong);
    setTrack(ctrlpPlaylist[0], ctrlpPlaylist, true);

    modal.style.display = "none";
    posHighlight = 0;
  }
}

function genPlaylistFromOneId(oneSongId, idxOfSong) {
  var returnPlaylist = [oneSongId];
  if (dataSetsIds.songids.length < 10) {
    for (i = 0; i < dataSetsIds.songids.length; i++) {
      returnPlaylist.push(dataSetsIds.songids[i]);
    }
    return returnPlaylist;
  }

  for (j = 0; j < 10; j++) {
    if (idxOfSong < dataSetsIds.songids.length - 1) {
      idxOfSong++;
    } else {
      idxOfSong = 0;
    }
    returnPlaylist.push(dataSetsIds.songids[idxOfSong]);
  }

  return returnPlaylist;
}

function highlightResultCell(indexHL) {
  var safeLength = document
    .getElementById("ctrlpResultsList")
    .getElementsByTagName("li");
  if (safeLength != null && indexHL < safeLength.length) {
    $("li").removeClass("activeL");
    safeLength[indexHL].className = "activeL";
  }
}

var patternField;
var matchFn = fuzzy_match;
var resultsList = null;
var currentDataSet = dataSets["songs"];

var asyncMatcher = null;

onload = function () {
  // Initialize document element references
  patternField = document.getElementById("ctrlpInput");
  patternField.oninput = onPatternChange;
  patternField.onpropertychange = patternField.oninput;

  resultsList = document.getElementById("ctrlpResultsList");
};

displayResults = function (results) {
  var newResultsList = resultsList.cloneNode(false);

  // Because adding too many elements is catastrophically slow because HTML is slow
  //var max_entries = 500;  // slim down max entries a bit - user shouldn't have that many songs
  var max_entries = 15;

  // Create HTML elements for results
  for (index = 0; index < results.length && index < max_entries; ++index) {
    var li = document.createElement("li");
    li.innerHTML = results[index];
    newResultsList.appendChild(li);
  }

  // Replace the old results from the DOM.
  resultsList.parentNode.replaceChild(newResultsList, resultsList);
  resultsList = newResultsList;

  highlightResultCell(0);
  $("#ctrlpResultsList li").click(function () {
    getIdAndPlay($(this).text());
  });
};

onPatternChange = function () {
  posHighlight = 0; // allow for further arrow key selection
  // Clear existing async match if it exists
  if (asyncMatcher !== null) {
    asyncMatcher.cancel();
    asyncMatcher = null;
  }

  var pattern = patternField.value;

  // Data not yet loaded
  if (currentDataSet == null) return;

  if (resultsList !== null) {
    // Clear the list
    var emptyList = resultsList.cloneNode(false);
    resultsList.parentNode.replaceChild(emptyList, resultsList);
    resultsList = emptyList;
  }

  // Early out on empty pattern (such as startup) because JS is slow
  if (pattern.length == 0) return;

  asyncMatcher = new fts_fuzzy_match_async(
    matchFn,
    pattern,
    currentDataSet,
    function (results) {
      // Scored function requires sorting
      if (matchFn == fuzzy_match) {
        results = results
          .sort(function (a, b) {
            return b[1] - a[1];
          })
          .map(function (v) {
            return v[2];
          });
      }

      displayResults(results);

      asyncMatcher = null;
    },
  );
  asyncMatcher.start();
};

onDataSetChange = function (radio) {
  var setname = radio.value;
  currentDataSet = dataSets[setname];
  onPatternChange();
};
