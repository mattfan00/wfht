document.addEventListener("htmx:beforeSwap", (e) => {
    if (
        e.detail.xhr.status === 500
        || e.detail.xhr.status === 400
    ) {
        e.detail.shouldSwap = true
        e.detail.isError = false
    }
})
