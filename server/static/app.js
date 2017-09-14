var pendingRequest = undefined;

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
  );
}

function clearResults() {
  $("#results").html("");
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

$("#query").on('input', updateResults);
$(function() {
  var results = $("#results").html();
  if ($.trim(results).length === 0) {
    updateResults();
  }
  stats();
});
