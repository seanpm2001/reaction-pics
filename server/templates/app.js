function updateResults() {
  var query = $("#query").val();
  $.getJSON(
    "/search",
    {query: query},
    function processQueryResult(data) {
      $("#results").html("");
      for (var x=0; x<data.length; x++) {
        var post = '<a href="' + data[x].url + '">' + data[x].title + '</a><br />';
        $("#results").append(post);
      }
    }
  );
}

$("#query").keyup(updateResults);
