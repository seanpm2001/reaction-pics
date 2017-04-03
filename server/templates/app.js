var timeouts = [];
function updateResults() {
  var query = $("#query").val();
  for (var x=0; x<timeouts.length; x++) {
    clearTimeout(timeouts[x]);
  }
  $.getJSON(
    "/search",
    {query: query},
    function processQueryResult(data) {
      $("#results").html("");
      for (var x=0; x<data.length; x++) {
        var postId = 'post' + x;
        var post = '<div id="' + postId + '">';
        post += '<a href="' + data[x].url + '">' + data[x].title + '</a><br />';
        post += '</div>';
        var timeout = showImage(x, postId, data[x].image);
        timeouts.push(timeout);
        $("#results").append(post);
      }
    }
  );
}

function showImage(x, postId, image) {
  var timeout = setTimeout(function() {
    var img = '<img src="' + image + '" />';
    $("#" + postId).append(img);
  }, 500 * x);
  return timeout;
}

$("#query").keyup(updateResults);
