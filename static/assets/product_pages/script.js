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

/* Eventlisteners related to cart and items adding */
let nextImg = 0,
	noOfItems = 0,
	clicked,
	trasitionTimer;

const minQuery = window.matchMedia("(min-width: 850px)"),
	maxQuery = window.matchMedia("(max-width: 850px)");

/*//////////////////////
 Functions 
 /////////////////////*/
/*Function to stop transition animation from triggering when page resize and reloading  */
function transitionDelay() {
	body.classList.add("preload");
	clearTimeout(trasitionTimer);
	trasitionTimer = setTimeout(() => {
		body.classList.remove("preload");
	}, 1000);
}

/* Function to get next and previous images*/
function imgBtns(btns, img, imgName) {
	btns.forEach((btn) => {
		btn.addEventListener("click", function (e) {
			if (e.target.classList.contains(`${imgName}__btnlft-img`)) {
				if (nextImg <= 0) nextImg = 3;
				else nextImg--;

				img.src = `images/image-product-${nextImg + 1}.jpg`;
			}

			if (e.target.classList.contains(`${imgName}__btnrgt-img`)) {
				if (nextImg >= 3) nextImg = 0;
				else nextImg++;

				img.src = `images/image-product-${nextImg + 1}.jpg`;
			}
		});
	});
}

imgBtns(overlayBtns, overlayImg, "item-overlay");
imgBtns(mainImgBtns, img, "img-main");

var add = document.getElementById('add')
var remove = document.getElementById('remove')
var text = document.getElementById("count")


// Function to add and remove 
add.addEventListener('click', function (e) {
	text.value = parseInt(text.value) + 1
})
remove.addEventListener('click', function (e) {
	text.value = parseInt(text.value) - 1
	if (text.value < 0) {
		text.value = 0;
	}
})
/* Function to open navigation menu */
/* Function to delete cart text 'empty cart' when cart items are > 0 */

function cartTx() {
	cartItem.classList.remove("open-cart");
	emptyCartTxt.classList.add("open-cart");
}

/* Function to delete cart text cart item  */
function emptyCart() {
	cartItem.classList.remove("open-cart");
	emptyCartTxt.classList.remove("open-cart");
}

/*//////////////////////
 Event Listeners 
 /////////////////////*/

/*  Eventlistener to close and open cart   */




/*  Eventlistener for add to cart button  */


/*  Eventlistener for delete cart item button   */


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


//Buy 


