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

document.addEventListener('DOMContentLoaded', (event) => {
    const themeSwitcher = document.getElementById('theme-switcher');
    const currentTheme = localStorage.getItem('theme') || 'dark';
    document.documentElement.setAttribute('data-theme', currentTheme);
    themeSwitcher.value = currentTheme;

    themeSwitcher.addEventListener('change', function() {
        document.documentElement.setAttribute('data-theme', this.value);
        localStorage.setItem('theme', this.value);
    });
});