document.addEventListener("DOMContentLoaded", function () {
  tabs = document.querySelectorAll(".gacha__change-image");
  for (i = 0; i < tabs.length; i++) {
    tabs[i].addEventListener("click", tabSwitch, false);
  }

  function tabSwitch() {
    tabs = document.querySelectorAll(".gacha__change-image");
    var node = Array.prototype.slice.call(tabs, 0);
    node.forEach(function (element) {
      element.classList.remove("active");
    });

    this.classList.add("active");
    content = document.querySelectorAll(".gacha__container-frame-inner");
    var node = Array.prototype.slice.call(content, 0);
    node.forEach(function (element) {
      element.classList.remove("active");
    });

    content = document.querySelectorAll(".gacha__description-text");
    var node = Array.prototype.slice.call(content, 0);
    node.forEach(function (element) {
      element.classList.remove("active");
    });

    content = document.querySelectorAll(".gacha__header-text");
    var node = Array.prototype.slice.call(content, 0);
    node.forEach(function (element) {
      element.classList.remove("active");
    });

    const arrayTabs = Array.prototype.slice.call(tabs);
    const index = arrayTabs.indexOf(this);

    document
      .querySelectorAll(".gacha__container-frame-inner")
      [index].classList.add("active");
    document
      .querySelectorAll(".gacha__description-text")
      [index].classList.add("active");
    document
      .querySelectorAll(".gacha__header-text")
      [index].classList.add("active");
    
    const gachaActive = document.getElementsByClassName("gacha__change-image active");
    const gachaActiveImage = gachaActive[0].getElementsByTagName("img")
    //var gachaImgSrc = gachaActiveImage[0].getAttribute("src").replace("off","on") //現在のimgからsrcを取得し、一部を書き換える
    //gachaActiveImage[0].setAttribute("src",gachaImgSrc);
  }
});
