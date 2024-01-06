document.addEventListener("htmx:responseError", (e) => {
    document.getElementById("error").innerHTML = e.detail.xhr.responseText
})
