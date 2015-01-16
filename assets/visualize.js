
var width  = 500,
    height = 300,
    data;

// Initialize svg.
var svg = d3.select("#chart").append("svg")
    .attr("width", width)
    .attr("height", height);

var g = {
    x: svg.append("g").attr("class", "x axis"),
    y: svg.append("g").attr("class", "y axis"),
    data: svg.append("g").attr("class", "data"),
};

var scales = {
    x: d3.scale.ordinal(),
    y: d3.scale.linear(),
}

// Refresh the visualization data.
function update() {
    updateScales();

    g.data.selectAll("circle").data(data)
      .call(function(selection) {
        var enter = selection.enter(),
            exit  = selection.exit();

        enter.append("circle")
          .attr("r", 3)

        selection
          .attr("cx", function(d) { return scales.x(d.state) + (scales.x.rangeBand() / 2) })
          .attr("cy", function(d) { return scales.y(d.value) })
      })
}

function updateScales() {
    scales.x.domain(data.map(function(d) { return d.state; }).sort())
        .rangeBands([0, width]);

    scales.y.domain(d3.extent(data, function(d) { return d.value; }))
        .range([height, 0]);
}

function load(q) {
    // TODO: d3.json("/query?q=...", function(error, data))
    update();
}


// TEMP: Mock data.
data = [
    {state:"CO", value:10},
    {state:"CO", value:20},
    {state:"CA", value:30},
    {state:"WA", value:40},
    {state:"WA", value:50},
]

update();
