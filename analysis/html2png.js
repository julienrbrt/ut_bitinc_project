//uses phamtomjs to generate an html to a png
var page = require('webpage').create();

//viewportSize being the actual size of the headless browser
page.viewportSize = { width: 1900, height: 900 };

page.open('{{.}}.html', function() {
    page.render('{{.}}.png');
    phantom.exit();
});