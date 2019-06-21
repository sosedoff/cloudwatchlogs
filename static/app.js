var textRenderer = new AnsiUp();

function parseTimestamps() {
  $("span.ts").each(function(idx, el) {
    if ($(el).data("loaded")) return;

    var unixTime = parseInt($(el).text()) / 1000;

    $(el).
      data("loaded", true).
      text(moment.unix(unixTime).format("MM/DD/YYYY, h:mm:ss a"));
  });
}

function resetLogLines() {
  $(".lines").html("");
}


function setLines(events, nextToken) {
  var lines = [];

  for (var i = 0; i < events.length; i++) {
    var event = events[i];
    var item = "<div class='line'>\
      <span class='ts'>" + event.Timestamp  + "</span>\
      <span class='src'>" + event.LogStreamName + "</span>\
      <span class='data'>" + textRenderer.ansi_to_html(event.Message) + "</span>\
      </div>";

    lines.push(item);
  }

  $("#results").data("next", nextToken);
  lines.push("<div class='paginator'>load more</div>");

  if ($("#results").data("append")) {
    $(lines.join("\n")).appendTo(".lines");
  } else {
    $(".lines").html(lines.join("\n"));
  }
  
  parseTimestamps();
}

function fetch() {
  $("form").submit();
}

function loadInitialSearch() {
  var q = window.location.search.split("?")[1] || "";
  var chunks = q.split("&");
  var params = {};

  for (var i = 0; i < chunks.length; i++) {
    var kv = chunks[i].split("=");
    params[kv[0]] = unescape(kv[1]);
  }

  if (params.filter)     $("#filter").val(params.filter);
  if (params.start_time) $("#start_time").val(params.start_time);

  if (params.group) {
    $("#log_group").val(params.group);

    loadGroupStreams(params.group, function() {
      if (params.stream) {
        $("#log_stream").val(params.stream);
      }
      fetch();
    });
  }
}

function loadGroupStreams(group, cb) {
  $.get("/streams", { group: group }, function (resp) {
    var opts = ["<option value=''>All Streams</options>"];
    for (var i = 0; i < resp.length; i++) {
      opts.push("<option value='" + resp[i].LogStreamName + "'>" + resp[i].LogStreamName + "</option>");
    }
    $("#log_stream").html(opts.join("\n"));

    if (cb) {
      cb();
    }
  });
}

function fetchLogEvents(form) {
  var formData = form.serializeArray();
  var params = [];

  for (var i = 0; i < formData.length; i++) {
    var field = formData[i];
    if (field.name == "next_token") continue;

    params.push(field.name + "=" + field.value);
  }

  var url = "/?" + params.join("&");
  history.pushState(params, "search", url);

  $.ajax({
    method: "POST",
    url: "/logs",
    data: form.serialize(),
    success: function (data) {
      $("#results").data("loading", false);
      $("#error").text("").hide();
      setLines(data.Events, data.NextToken);
    },
    error: function (xhr) {
      $("#results").data("loading", false);
      $("#results .lines").html("");
      $("#error").text(xhr.responseJSON.error).show();
    }
  });
}

$(function() {
  loadInitialSearch();

  $("form").on("submit", function(e) {
    e.preventDefault();
    fetchLogEvents($(this));
  });

  $("#log_group").on("change", function() {
    $("#results").data("append", false);
    $("#log_stream").val("");
    $("#start_time").val("1h");
    $("#next_token").val("");
    $("#filter").val("");
    $("#results").scrollTop(0);
    loadGroupStreams($(this).val());
    fetch();
  });

  $("#log_stream").on("change", function () {
    $("#results").data("append", false);
    $("#next_token").val("");
    $("#results").scrollTop(0);
    fetch();
  });

  $("#start_time").on("change", function() {
    $("#results").data("append", false);
    $("#next_token").val("");
    $("#results").scrollTop(0);
    fetch();
  });

  $("#filter").on("change", function() {
    $("#results").data("append", false);
    $("#next_token").val("");
    $("#results").scrollTop(0);
  });

  $("body").on("click", "span.src", function(e) {
    var stream = $(this).text();

    $("#log_stream").val(stream);
    fetch();
  });

  $("#results").on("scroll", function(e) {
    var height          = $("#results .lines").height();
    var containerHeight = $("#results").height();
    var offset          = Math.abs($("#results .lines").offset().top);
    var scrollPercent   = ((containerHeight + offset) * 100) / height;

    if (scrollPercent < 95) {
      return;
    }

    var loading = $("#results").data("loading");
    if (loading) return;

    var token = $("#results").data("next");
    if (!token) return;

    $("#next_token").val(token);
    $("#results").data("loading", true);
    $("#results").data("append", true);

    fetch();
  });
});