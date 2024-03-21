function updateImage(ele){
    document.querySelector("#kp-image-viewer").querySelector("img").src = ele.src;
    var images = document.querySelector(".kp-image-thumb").querySelectorAll("img");
    for(i=0;i<images.length;i++){
      images[i].classList.remove("opacity-100");
      images[i].classList.add("opacity-60");
    }
    ele.classList.remove("opacity-60");
    ele.classList.add("opacity-100");
}