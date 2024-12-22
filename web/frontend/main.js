import htmx from 'htmx.org'
window.htmx = htmx
import './style.css'

function removeClass(elements, className) {
    for (const element of elements) {
        htmx.removeClass(element, className);
    }
}

function addClass(elements, className) {
    for (const element of elements) {
        htmx.addClass(element, className);
    }
}

document.body.addEventListener("enterEditMode", function(evt){
    htmx.addClass(htmx.find('#companies .button-append'), 'disabled');
    addClass(htmx.findAll('#table-body .button-edit'), 'disabled')
    addClass(htmx.findAll('#table-body .button-delete'), 'disabled')
})

document.body.addEventListener("exitEditMode", function(evt){
    htmx.removeClass(htmx.find('#companies .button-append'), 'disabled');
    removeClass(htmx.findAll('#table-body .button-edit'), 'disabled')
    removeClass(htmx.findAll('#table-body .button-delete'), 'disabled')
})