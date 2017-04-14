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
        var postId = 'post' + x;
        var post = '<div id="' + postId + '" class="result">';
        post += '<h2><a href="' + data[x].url + '">' + data[x].title + '</a></h2>';
        post += '</div>';
        $("#results").append(post);
        showImage(x, postId, data[x].image);
      }
      $('img.result-img').lazyload({
        effect: "fadeIn",
        threshold: 1000,
        skip_invisible: true
      });
    }
  );
}

function showImage(x, postId, image) {
  var img = '<img data-original="' + image + '" class="result-img"/>';
  $("#" + postId).append(img);
  console.log($("#" + postId));
}

$("#query").keypress(updateResults);
