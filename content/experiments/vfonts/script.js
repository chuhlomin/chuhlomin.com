function calculateDistance(elem, mouseX, mouseY) {
    return Math.floor(Math.sqrt(Math.pow(mouseX - (elem.offsetLeft+(elem.clientWidth/2)), 2) + Math.pow(mouseY - (elem.offsetTop+(elem.clientHeight/2)), 2)));
}

function wrapChars(str, tmpl) {
    return str.replace(/\w/g, tmpl || "<span>$&</span>");
}

function calculateDistanceX(elem, x) {
    return Math.floor(Math.abs(elem.offsetLeft + (elem.offsetWidth/2) - x))
}

function updateFont(x) {
    line.style.left = x + "px";

    for (var i = 0; i < symbols.length; i++) {
        distance = calculateDistanceX(symbols[i], x);

        var width = (100000/(distance+800)).toFixed(2);
        var weight = (100000/(distance+120)).toFixed(2);

        // symbols[i].style.fontVariationSettings = "\"wdth\" " + width;
        // symbols[i].style.fontVariationSettings = "\"wght\" " + weight;
        symbols[i].style.fontVariationSettings = "\"wght\" " + weight + ", \"wdth\" " + width;
        symbols[i].style.color = "hsl(" + (270+0.25*360*distance/window.innerWidth) + "deg, 100%, 50%)";
    }
}

var x = 0;
var autoplay = true;
var distance;

let degStep = 0.003 * Math.PI;
var deg = - Math.PI;

let line = document.getElementById("line")
let text = document.getElementById("var")
let symbols = document.getElementsByTagName("span")

text.innerHTML = wrapChars(text.innerHTML);

document.addEventListener('mousemove', function(e) {
    autoplay = false;
    x = e.pageX;
    updateFont(x);
    deg = Math.acos(2 * x / window.innerWidth - 1);
});

// mouse moved outside
document.addEventListener('mouseout', function(e) {
    autoplay = true;
    if (degStep > 0) {
        degStep *= -1;
    }
});

let timer = setInterval(function() {
    if (autoplay) {
        deg += degStep;
        x = ((Math.cos(deg) + 1)/2 * window.innerWidth).toFixed(2);
        updateFont(x);
    }

}, 20);
