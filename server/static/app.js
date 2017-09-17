var pendingRequest = undefined;

function showPost(postID) {
  $.getJSON(
    "/postdata/" + postID,
    function processPostResult(post) {
      clearResults();
      addResults([post]);
    }
  );
}

function updateResults() {
  var query = $("#query").val();
  if (pendingRequest) {
    pendingRequest.abort();
  }
  pendingRequest = $.getJSON(
    "/search",
    {query: query},
    function processQueryResult(data) {
      clearResults();
      updateURL();
      addResults(data);
    }
  );
}

function clearResults() {
  $("#results").html("");
}

function updateURL() {
  var url = "/";
  var urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    history.pushState({}, "Reaction Pics", url);
  } else {
    history.replaceState({}, "Reaction Pics", url);
  }
}

function addResults(data) {
  for (var x=0; x<data.length; x++) {
    var post = data[x];
    addResult(post);
  }
  $('img.result-img').lazyload({
    effect: "fadeIn",
    threshold: 1000,
    skip_invisible: true
  });
}

function addResult(postData) {
  var postHTML = '<div class="result">';
  postHTML += '<h2>';
  if (postData.url) postHTML += '<a href="' + postData.internalURL + '">';
  postHTML += postData.title;
  if (postData.url) postHTML += '</a></h2>';
  if (postData.image) postHTML += '<p><img data-original="' + postData.image + '" class="result-img" /></p>';
  if (postData.likes) {
      postHTML += '<p><a href="#" id="likes" class="btn btn-success disabled">';
      postHTML += postData.likes + ' <span class="glyphicon glyphicon-thumbs-up" aria-hidden="true"></span>';
      postHTML += '</a></p>';
      postHTML += '<p><a href="' + postData.url + '">Original</a></p>';
  }
  postHTML += '</div>';
  $("#results").append(postHTML);
}

function stats() {
  $.getJSON(
    "/stats.json",
    function processStats(data) {
      var line = "Currently indexing " + data.postCount + " posts";
      $("#indexStat").text(line);
    }
  );
}

$(function() {
  $("#query").on('input', updateResults);
  var urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    showPost(urlPath[2]);
  } else {
    updateResults();
  }
  stats();
});
