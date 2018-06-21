var goButton = document.getElementById('go');
var urlBox = document.getElementById('url');

goButton.onclick = function () {
    var longURL = urlBox.value;
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST", "/");
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4) {
            var r = xmlhttp.responseText;
            console.debug(r);
            if (xmlhttp.status != 200) {
                alert(r);
            }else{
                urlBox.value = r;
            }
        }
    };
    console.debug("POST: " + longURL);
    xmlhttp.send(longURL);
};

urlBox.onclick = function(){
    this.select();
};