var collums = document.getElementsByClassName("collapsible");
for(var i = 0 ; i < collums.length; i++) {
    var collum = collums[i];
    collum.addEventListener("click" , function() {
        this.classList.toggle("active");
        var dropdownContent = this.nextElementSibling;
        if(dropdownContent.style.display === "block") {
            dropdownContent.style.display = "none";
        } else {
            dropdownContent.style.display = "block";
        }
    });
}
