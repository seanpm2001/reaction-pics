var maxLength = 20;
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
      $("#results").html("");
      for (var x=0; x<data.length; x++) {
        var post = data[x];
        showImage(x, post);
        if (x >= maxLength) {
          showImage(x, {
            title: "Refine your search for more results"
          });
          break;
        }
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
  postHTML += '<h2>';
  if (postData.url) postHTML += '<a href="' + postData.internalURL() + '">';
  postHTML += postData.title;
  if (postData.url) postHTML += '</a></h2>';
  if (postData.image) postHTML += '<img data-original="' + postData.image + '" class="result-img" /><br />';
  if (postData.likes) postHTML += postData.likes + ' likes';
  postHTML += '</div>';
  $("#results").append(postHTML);
}

$("#query").on('input', updateResults);
$(function() {
  updateResults();
});
