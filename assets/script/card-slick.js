$(document).ready(function(){
    $('.slick-carousel').slick({
        autoplay: true,
        dots: true,
        speed: 300,
        slidesToShow: 3,
        slidesToScroll: 1,
        centerMode: true,
        variableWidth: true,
        centerPadding: '0px',
        arrows: true,
        prevArrow: '<button type="button" class="absolute top-[40%] left-0 z-[10]"> <div class="h-[40px] w-[40px]"><img src="/assets/images/arrowSlider.svg" alt="arrowSlider" width="100%" height="auto" /></div></button>',
        nextArrow: '<button type="button" class="absolute top-[40%] right-0 z-[10]"> <div class="h-[40px] w-[40px] rotate-180"><img src="/assets/images/arrowSlider.svg" alt="arrowSlider" width="100%" height="auto" /></div></button>',
        responsive: [
        {
            breakpoint: 768,
            settings: {
            arrows: true,
            centerMode: true,
            centerPadding: '40px',
            slidesToShow: 2,
            slidesToScroll: 1,
            }
        },
        {
            breakpoint: 480,
            settings: {
            arrows: true,
            centerMode: true,
            centerPadding: '40px',
            slidesToShow: 1
            }
        }
        ]
        // Add more configuration options as needed
    });
    });


    $('.slider-for').slick({
        slidesToShow: 1,
        slidesToScroll: 1,
        arrows: false,
        fade: true,
        asNavFor: '.slider-nav'
      });
      $('.slider-nav').slick({
        slidesToShow: 3,
        slidesToScroll: 1,
        asNavFor: '.slider-for',
        dots: true,
        centerMode: true,
        focusOnSelect: true
      });