$(document).ready(function() {
  var switchTemplate = $('#switch_template').html();
  Mustache.tags = ["[[", "]]"];
  Mustache.parse(switchTemplate);

  $('body').on("change", 'input[data-toggle="toggle"]', function() {
    var node = $(this).data("node");
    var pin = $(this).data("pin");
    var checked = $(this).prop("checked");

    $.post("switches", { node: node, pin: pin, checked: checked });
  });

  $.get("switches", function(data) {
    var rendered = "";
    for(node in data) {
      rendered += Mustache.render("Node [[ name ]]", {name: data[node]["Name"]});
      for(pin in data[node]["Pins"]) {
        var tmpl = '<div class="switch">[[ pin ]] <input [[#checked]] checked [[/checked]] data-toggle="toggle" data-node="[[ node ]]" data-pin="[[ pin ]]" type="checkbox"/></div>';
        rendered += Mustache.render(tmpl,
          { pin: pin, node: node, checked: data[node]["Pins"][pin] });
      }
    };

    $('#app').html(rendered);
    $('[data-toggle="toggle"]').bootstrapToggle();

  }, "json");
});
