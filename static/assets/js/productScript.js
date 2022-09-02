"use strict";

const imgBtn = Array.from(document.querySelectorAll(".img-btn"));
const img = document.querySelector(".img-main");
const mainImgBtns = Array.from(document.querySelectorAll(".img-main__btn"));

const overlayCon = document.querySelector(".overlay-container");
const overlayImg = document.querySelector(".item-overlay__img");
const overlayImgBtn = Array.from(
	document.querySelectorAll(".overlay-img__btn")
);
const overlayBtnImgs = Array.from(
	document.querySelectorAll(".overlay-img__btn-img")
);
const overlayCloseBtn = document.querySelector(".item-overlay__btn ");
const overlayBtns = Array.from(document.querySelectorAll(".overlay-btn"));


const priceSingle = document.querySelector(".head-cart__price-single");
const priceTotal = document.querySelector(".head-cart__price-total");

const priceBtns = Array.from(document.querySelectorAll(".price-btn__img"));
const totalItems = document.querySelector(".price-btn__txt");

const menuOpen = document.querySelector(".head-lft__btn");
const menu = document.querySelector(".head-nav");
const menuBtnImg = document.querySelector(".head-lft__btn-img");

const bodyOverlay = document.querySelector(".body-wrapper");
const body = document.querySelector("body");

const headerCart = document.querySelector(".head-rgt");

/*//////////////////////
 Functions 
 /////////////////////*/
/*Function to stop transition animation from triggering when page resize and reloading  */
/* Function to get next and previous images*/
function rightBtn(){
	console.log(img.src)
	var number = img.src.slice(img.src.lastIndexOf("/")+1,-5);
	var numberAdd = ++number;
	const regex = /\/\d+/g;
	img.src = img.src.replace(regex, "/" + numberAdd);
}

function leftBtn(){
	var number = img.src.slice(img.src.lastIndexOf("/")+1,-5);

	const regex = /\/\d+/g;
	if(number != 1){
		var numberSubtract = --number;
		img.src = img.src.replace(regex, "/" + numberSubtract)
	}
	
}

function mainImgError(image){

	var number = img.src.slice(img.src.lastIndexOf("/")+1,-5);
	const regex = /\/\d+/g;
	var numberSubtract = --number;
	img.src = img.src.replace(regex, "/" + numberSubtract);
}

var add = document.getElementById('add')
var remove = document.getElementById('remove')
var text = document.getElementById("count")


// Function to add and remove 
add.addEventListener('click', function (e) {
	text.value = parseInt(text.value) + 1
})
remove.addEventListener('click', function (e) {
	text.value = parseInt(text.value) - 1
	if (text.value < 1) {
		text.value = 1;
	}
})
/* Function to open navigation menu */

/*//////////////////////
 Event Listeners 
 /////////////////////*/






/*  Eventlistener for  image to change when image button is clicked  */
imgBtn.forEach((btn, i) => {
	btn.addEventListener("click", function (e) {
		console.log(e.target.children[0].src)
		img.src = e.target.children[0].src;
	});
});

/*  Eventlistener to stop transition animation from triggering when page reloading  */
window.addEventListener("load", function () {
	transitionDelay();
});





function imgError(image) {
    image.onerror = nil;
	image.classList.add("hide")
    return true;
}