// shake function
var shakingElements = [];

var shake = function (element, magnitude = 16) {
  //First set the initial tilt angle to the right (+1) 

  //A counter to count the number of shakes
  var counter = 1;

  //The total number of shakes (there will be 1 shake per frame)
  var numberOfShakes = 15;

  //Capture the element's position and angle so you can
  //restore them after the shaking has finished
  var startX = 0,
    startY = 0;

  // Divide the magnitude into 10 units so that you can 
  // reduce the amount of shake by 10 percent each frame
  var magnitudeUnit = magnitude / numberOfShakes;

  //The `randomInt` helper function
  var randomInt = (min, max) => {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  };

  //Add the element to the `shakingElements` array if it
  //isn't already there
  if (shakingElements.indexOf(element) === -1) {
    //console.log("added")
    shakingElements.push(element);
    upAndDownShake();
  }

  //The `upAndDownShake` function
  function upAndDownShake() {

    //Shake the element while the `counter` is less than 
    //the `numberOfShakes`
    if (counter < numberOfShakes) {

      //Reset the element's position at the start of each shake
      element.style.transform = 'translate(' + startX + 'px, ' + startY + 'px)';

      //Reduce the magnitude
      magnitude -= magnitudeUnit;

      //Randomly change the element's position
      var randomX = randomInt(-magnitude, magnitude);
      var randomY = randomInt(-magnitude, magnitude);

      element.style.transform = 'translate(' + randomX + 'px, ' + randomY + 'px)';

      //Add 1 to the counter
      counter += 1;

      requestAnimationFrame(upAndDownShake);
    }

    //When the shaking is finished, restore the element to its original 
    //position and remove it from the `shakingElements` array
    if (counter >= numberOfShakes) {
      element.style.transform = 'translate(' + startX + ', ' + startY + ')';
      shakingElements.splice(shakingElements.indexOf(element), 1);
    }
  }

};

var goButton = document.getElementById('go');
var copyButton = document.getElementById('copy');
var resetButton = document.getElementById('reset');
var urlBox = document.getElementById('url');
var clipboard = new ClipboardJS('#copy');

var cmtc = "Click me to copy!";
var copied = "Copied!";

function go() {
  var longURL = urlBox.value;
  if (longURL == "") {
    return
  }
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.open("POST", "/");
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4) {
      var r = xmlhttp.responseText;
      if (xmlhttp.status != 200) {
        var e = "error: " + r;
        console.error(e);
        urlBox.title = e;
        shake(urlBox);
      } else {
        console.debug(r);
        urlBox.value = r;
        urlBox.readOnly = true;
        goButton.hidden = true;
        copyButton.hidden = false;
        resetButton.hidden = false;
      }
    }
  };
  console.debug("POST: " + longURL);
  xmlhttp.send(longURL);
};

goButton.onclick = go;

function reset() {
  console.debug("reset");
  urlBox.value = "";
  urlBox.title = "";
  urlBox.readOnly = false;
  goButton.hidden = false;
  copyButton.hidden = true;
  copyButton.innerHTML = "Copy";
  resetButton.hidden = true;
};

resetButton.onclick = reset;

function copy() {
  console.debug("copy");
  urlBox.select();
  copyButton.innerHTML = copied.bold();
};

copyButton.onclick = copy;

urlBox.onclick = function () {
  this.select();
};

urlBox.onkeydown = function (ev) {
  if (ev.keyCode != 13) {
    return
  }
  // go
  if (!goButton.hidden && urlBox.value != "") {
    go();
  } else {
    // reset
    if (urlBox.value == "") {
      reset();
    } else {
      if (copyButton.innerText == copied) {
        return
      }
      urlBox.select();
      copyButton.innerText = cmtc;
      shake(copyButton);
    }
  }
};