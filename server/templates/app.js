var pendingRequest = undefined;
function updateResults(e) {
  if (e.which === 0) {
      return;
  }
  var query = $("#query").val();
  if (pendingRequest) {
    pendingRequest.abort();
  }
  pendingRequest = $.getJSON(
    "/search",
    {query: query},
    function processQueryResult(data) {
      $("#results").html("");
      for (var x=0; x<data.length; x++) {
        var post = data[x];
        showImage(x, post);
      }
      $('img.result-img').lazyload({
        effect: "fadeIn",
        threshold: 1000,
        skip_invisible: true
      });
    }
  );
}

function showImage(x, postData) {
  var postHTML = '<div id="post' + x + '" class="result">';
  postHTML += '<h2><a href="' + postData.url + '">' + postData.title + '</a></h2>';
  postHTML += '<img data-original="' + postData.image + '" class="result-img" />';
  postHTML += '</div>';
  $("#results").append(postHTML);
}

$("#query").keypress(updateResults);
